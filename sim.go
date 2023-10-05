package simgo

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/flysim-apps/simgo/websockets"
	sim "github.com/micmonay/simconnect"
	"github.com/op/go-logging"
)

type SimGo struct {
	Error      error
	State      chan int
	Connection <-chan bool
	TrackEvent chan interface{}
	Logger     *logging.Logger
	Socket     *websockets.Websocket
	Context    context.Context
}

// creates new simgo instance
func NewSimGo(logger *logging.Logger) *SimGo {
	return &SimGo{State: make(chan int, 1), TrackEvent: make(chan interface{}, 1), Logger: logger}
}

// starts web socket server on given host and port
func (s *SimGo) Start(httpListen string) chan bool {
	c := make(chan bool)
	s.Socket = websockets.New()
	http.HandleFunc("/socket.io", s.Socket.Serve)
	go func() {
		s.Logger.Debugf("Socket starting on port %s", httpListen)
		go func() {
			c <- true
		}()
		if err := http.ListenAndServe(httpListen, nil); err != nil {
			go func() {
				s.Logger.Errorf("Server could not be started! Reason: %s", err.Error())
				s.Error = err
				s.State <- STATE_FATAL
			}()
		}
	}()
	return c
}

// connects to MSFS
func (s *SimGo) connectToMsfs(name string) (*sim.EasySimConnect, error) {
	// dllPath := filepath.Join(filepath.Dir(exePath), "SimConnect.dll")
	// if _, err = os.Stat(dllPath); os.IsNotExist(err) {
	// 	buf := MustAsset("../simgo/SimConnect.dll")

	// 	if err := ioutil.WriteFile(dllPath, buf, 0644); err != nil {
	// 		return nil, err
	// 	}
	// }
	s.Logger.Info("Connecting to MSFS...")
	sc, err := sim.NewEasySimConnect()
	if err != nil {
		s.Logger.Errorf("Failed to NewEasySimConnect: %v", err.Error())
		return nil, err
	}

	sc.SetDelay(1 * time.Second)
	sc.SetLoggerLevel(sim.LogInfo)

	go func() {
		for {
			c, err := sc.Connect(name)
			if err != nil {
				go func() {
					s.Error = err
					s.State <- STATE_CONNECTION_FAILED
				}()
				time.Sleep(30 * time.Second)
				continue
			}
			<-c
			go func() {
				s.Connection = c
				s.Error = nil
				s.State <- STATE_CONNECTION_READY
			}()
			return
		}
	}()

	return sc, nil
}

var connectToMsfsInProgress = false
var lastMessageReceived time.Time

func (s *SimGo) TrackWithRecover(name string, report interface{}, maxTries int, trackID int) {
	go recoverer(maxTries, trackID, func() {
		checker := time.NewTicker(30 * time.Second)

		ctx, cancel := context.WithCancel(context.Background())
		s.Context = ctx
		wait := &sync.WaitGroup{}
		wait.Add(2)

		if err := s.Track(name, report, wait); err != nil {
			s.Logger.Errorf("Tracking error: %v", err.Error())
			wait.Done()
		}

		go func() {
			defer wait.Done()
			defer fmt.Println(`Exiting `)

			for {
				select {
				case <-checker.C:
					timeNow := time.Now().Add(-5 * time.Second)
					s.Logger.Info("Timeout checker")
					if !connectToMsfsInProgress && !lastMessageReceived.IsZero() && lastMessageReceived.Before(timeNow) {
						s.Logger.Info("Last received message was received 5 sec ago. Cancel tracking")
						cancel()
						return
					}
				}
			}
		}()

		wait.Wait()

		s.Logger.Warning("Exiting from tracker routine...")

		return
	})
}

func (s *SimGo) Track(name string, report interface{}, wait *sync.WaitGroup) error {
	sc, err := s.connectToMsfs(name)
	if err != nil {
		return errors.New(fmt.Sprintf("connection to MSFS has been failed. Reason: %s", err.Error()))
	}

	go func() {
		defer wait.Done()
		defer fmt.Println(`Exiting Track`)
		defer sc.Close()

		for {
			select {
			case <-s.Context.Done():
				s.Logger.Warning("Tracking routine will exit")
				connectToMsfsInProgress = true
				return
			case open := <-s.Connection:
				s.Logger.Debugf("Open: %b", open)
			case state := <-s.State:
				switch state {
				case STATE_CONNECTION_FAILED:
					s.Logger.Debugf("Waiting for MSFS... %v", s.Error)
				case STATE_CONNECTION_READY:
					s.Logger.Infof("Connection to MSFS has been established")
					s.Socket.Broadcast(EventNotification(MSFS_CONNECTION_READY))
					time.Sleep(20 * time.Second)
					connectToMsfsInProgress = false
					s.trackSimVars(sc, reflect.ValueOf(report))
				default:
					s.Logger.Warningf("Received simVar: %v", state)
				}
			}
		}
	}()

	return nil
}

func (s *SimGo) trackSimVars(sc *sim.EasySimConnect, report reflect.Value) error {
	if err := s.ConnectToSimVar(sc, convertToSimSimVar(report), report); err != nil {
		return errors.New(fmt.Sprintf("failed to connect to SimVar: %v ", err.Error()))
	}
	return nil
}

func (s *SimGo) ConnectToSimVar(sc *sim.EasySimConnect, listSimVar []sim.SimVar, returnType reflect.Value) error {
	if sc == nil {
		return errors.New("sim connect is nil")
	}
	cSimVar, err := sc.ConnectToSimVar(listSimVar...)
	if err != nil {
		return err
	}

	cSimStatus := sc.ConnectSysEventSim()
	//wait sim start
	for {
		if <-cSimStatus {
			break
		}
	}

	crashed := sc.ConnectSysEventCrashed()
	for {
		select {
		case sv := <-cSimVar:
			lastMessageReceived = time.Now()
			s.TrackEvent <- convertToInterface(returnType, sv)
		case <-crashed:
			s.Logger.Error("Your are crashed !!")
			<-sc.Close() // Wait close confirmation
			return nil
		}

	}
}

func convertToSimSimVar(val reflect.Value) []sim.SimVar {
	vars := make([]sim.SimVar, 0)

	for i := 0; i < val.Type().NumField(); i++ {
		nameTag, _ := val.Type().Field(i).Tag.Lookup("name")
		indexTag, _ := val.Type().Field(i).Tag.Lookup("index")
		unitTag, _ := val.Type().Field(i).Tag.Lookup("unit")
		settableTag, _ := val.Type().Field(i).Tag.Lookup("settable")

		if nameTag == "" || unitTag == "" {
			continue
		}

		simv := sim.SimVar{
			Name: nameTag,
			Unit: sim.SimVarUnit(unitTag),
		}

		if indexTag != "" {
			idx, _ := strconv.Atoi(indexTag)
			simv.Index = idx
		}

		if settableTag != "" {
			simv.Settable = settableTag == "1" || strings.ToLower(settableTag) == "true"
		}

		vars = append(vars, simv)
	}

	return vars
}

func convertToInterface(val reflect.Value, vars []sim.SimVar) interface{} {
	found := make([]string, 0)
	r := reflect.New(reflect.TypeOf(val.Interface())).Elem()
	for _, simVar := range vars {
		//fmt.Printf("iterateSimVars(): Name: %s                                               Index: %b    Unit: %s\n", simVar.Name, simVar.Index, simVar.Unit)
		for j := 0; j < val.NumField(); j++ {
			nameTag, _ := val.Type().Field(j).Tag.Lookup("name")
			indexTag, _ := val.Type().Field(j).Tag.Lookup("index")
			if indexTag == "" {
				indexTag = "0"
			}

			idx, _ := strconv.Atoi(indexTag)

			if simVar.Index == idx && simVar.Name == nameTag {
				found = append(found, fmt.Sprintf("Name: %s                   Index: %b    Unit: %s\n", simVar.Name, simVar.Index, simVar.Unit))
				getValue(r.Field(j), simVar)
			}
		}
	}
	return r.Interface()
}
