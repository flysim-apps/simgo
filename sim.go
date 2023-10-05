package simgo

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/flysim-apps/simgo/websockets"
	sim "github.com/micmonay/simconnect"
	"github.com/op/go-logging"
)

type SimGo struct {
	Error         error
	State         chan int
	Connection    <-chan bool
	TrackEvent    chan interface{}
	Logger        *logging.Logger
	Socket        *websockets.Websocket
	Context       context.Context
	ContextCancel context.CancelFunc
}

// creates new simgo instance
func NewSimGo(logger *logging.Logger) *SimGo {
	ctx, cancel := context.WithCancel(context.Background())
	return &SimGo{State: make(chan int, 1), TrackEvent: make(chan interface{}, 1), Logger: logger, Context: ctx, ContextCancel: cancel}
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

func (s *SimGo) Track(name string, report interface{}) error {
	sc, err := s.connectToMsfs(name)
	if err != nil {
		return errors.New(fmt.Sprintf("connection to MSFS has been failed. Reason: %s", err.Error()))
	}

	checker := time.NewTicker(30 * time.Second)

	go func() {
		for {
			select {
			case <-checker.C:
				timeNow := time.Now().Add(-5 * time.Second)
				if !connectToMsfsInProgress && !lastMessageReceived.IsZero() && lastMessageReceived.Before(timeNow) {
					s.Logger.Infof("Last received message was received 5 sec ago. Restart tracking")
					sc.Close()
					//sc, err = sim.ConnectToMsfs()
					connectToMsfsInProgress = true
					s.ContextCancel()
					return
					// if err != nil {
					// 	globals.Logger.Errorf("Connection to MSFS has been failed. Reason: %s", err.Error())
					// 	return
					// }
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-s.Context.Done():
				s.Logger.Warning("Tracking routine will exit")
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
					s.trackSimVars(sc, report)
				default:
					s.Logger.Warningf("Received simVar: %v", state)
				}
			}
		}
	}()

	return nil
}

func (s *SimGo) trackSimVars(sc *sim.EasySimConnect, report interface{}) error {
	if err := s.ConnectToSimVar(sc, convertToSimSimVar(reflect.ValueOf(report)), report); err != nil {
		return errors.New(fmt.Sprintf("failed to connect to SimVar: %v ", err.Error()))
	}
	return nil
}

func (s *SimGo) ConnectToSimVar(sc *sim.EasySimConnect, listSimVar []sim.SimVar, result interface{}) error {
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
			s.TrackEvent <- convertToInterface(reflect.ValueOf(result).Addr(), sv)
		case <-crashed:
			s.Logger.Error("Your are crashed !!")
			<-sc.Close() // Wait close confirmation
			return nil
		}

	}
}

func convertToSimSimVar(v reflect.Value) []sim.SimVar {
	vars := make([]sim.SimVar, 0)

	for i := 0; i < v.Type().NumField(); i++ {
		nameTag, _ := v.Type().Field(i).Tag.Lookup("name")
		indexTag, _ := v.Type().Field(i).Tag.Lookup("index")
		unitTag, _ := v.Type().Field(i).Tag.Lookup("unit")
		settableTag, _ := v.Type().Field(i).Tag.Lookup("settable")

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
	if val.Kind() == reflect.Interface && !val.IsNil() {
		elm := val.Elem()
		if elm.Kind() == reflect.Ptr && !elm.IsNil() && elm.Elem().Kind() == reflect.Ptr {
			val = elm
		}
	}
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	found := make([]string, 0)
	for _, simVar := range vars {
		fmt.Printf("iterateSimVars(): Name: %s                                               Index: %b    Unit: %s\n", simVar.Name, simVar.Index, simVar.Unit)
		for j := 0; j < val.NumField(); j++ {
			nameTag, _ := val.Type().Field(j).Tag.Lookup("name")
			indexTag, _ := val.Type().Field(j).Tag.Lookup("index")
			if indexTag == "" {
				indexTag = "0"
			}

			idx, _ := strconv.Atoi(indexTag)

			if simVar.Index == idx && simVar.Name == nameTag {
				found = append(found, fmt.Sprintf("Name: %s                   Index: %b    Unit: %s\n", simVar.Name, simVar.Index, simVar.Unit))
				getValue(val.Field(j), simVar)
			}
		}
	}
	return val
}
