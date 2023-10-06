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

var connectToMsfsInProgress = false
var lastMessageReceived time.Time

// creates new simgo instance
func NewSimGo(logger *logging.Logger) *SimGo {
	return &SimGo{State: make(chan int, 1), TrackEvent: make(chan interface{}, 1), Logger: logger}
}

// starts web socket server on given host and port
func (s *SimGo) StartWebSocket(httpListen string) error {
	s.Socket = websockets.New()
	http.HandleFunc("/socket.io", s.Socket.Serve)
	s.Logger.Debugf("Socket starting on port %s", httpListen)
	go func() {
		if err := http.ListenAndServe(httpListen, nil); err != nil {
			s.Logger.Errorf("Server could not be started! Reason: %s", err.Error())
			return
		}
	}()
	return nil
}

// connects to MSFS
func (s *SimGo) connect(name string) (*sim.EasySimConnect, error) {
	s.Logger.Info("Connecting to MSFS...")
	sc, err := sim.NewEasySimConnect()
	if err != nil {
		return nil, err
	}

	sc.SetDelay(1 * time.Second)
	sc.SetLoggerLevel(sim.LogInfo)

	c, err := sc.Connect(name)
	if err != nil {
		return nil, err
	}

	<-c // wait connection confirmation

	for {
		if <-sc.ConnectSysEventSim() {
			break // wait sim start
		}
	}

	return sc, nil
}

func (s *SimGo) TrackWithRecover(name string, report interface{}, maxTries int, trackID int) {
	go recoverer(maxTries, trackID, func() {
		checker := time.NewTicker(15 * time.Second)
		ctx, cancel := context.WithCancel(context.Background())
		wait := sync.WaitGroup{}

		defer checker.Stop()
		wait.Add(2)

		go s.track(name, report, ctx, &wait)

		go func() {
			defer wait.Done()

			for {
				select {
				case <-ctx.Done():
					s.Logger.Warning("Checking routine will exit")
					return
				case <-checker.C:
					timeNow := time.Now().Add(-5 * time.Second)
					s.Logger.Info("Timeout checker")
					if !connectToMsfsInProgress && !lastMessageReceived.IsZero() && lastMessageReceived.Before(timeNow) {
						s.Logger.Info("Last received message was received 5 sec ago. Cancel tracking")
						cancel()
					}
				}
			}
		}()

		wait.Wait()

		panic("Exiting from tracker routine...")
	})
}

func (s *SimGo) track(name string, report interface{}, ctx context.Context, wg *sync.WaitGroup) {
	sc, err := s.connect(name)
	defer wg.Done()
	defer sc.Close()

	if err != nil {
		panic("connection to MSFS has been failed. Reason: %s" + err.Error())
	}

	cSimVar, err := sc.ConnectToSimVar(convertToSimSimVar(reflect.ValueOf(report))...)
	if err != nil {
		panic("connection to MSFS has been failed. Reason: %s" + err.Error())
	}

	crashed := sc.ConnectSysEventCrashed()

	for {
		select {
		case <-ctx.Done():
			s.Logger.Warning("Tracking routine will exit")
			connectToMsfsInProgress = true
			return
		case sv := <-cSimVar:
			lastMessageReceived = time.Now()
			connectToMsfsInProgress = false
			s.TrackEvent <- convertToInterface(reflect.ValueOf(report), sv)
		case <-crashed:
			s.Logger.Error("Your are crashed !!")
			<-sc.Close() // Wait close confirmation
			return
		}
	}
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
