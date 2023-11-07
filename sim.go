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

	sim "github.com/flysim-apps/simgo/simconnect"
	"github.com/flysim-apps/simgo/websockets"
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

var maxTriesInitial int
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
func (s *SimGo) connect(ctx context.Context, name string) (*sim.EasySimConnect, error) {
	s.Logger.Info("Connecting to MSFS...")
	sc, err := sim.NewEasySimConnect(ctx)
	if err != nil {
		return nil, err
	}

	sc.SetDelay(1 * time.Second)
	sc.SetLoggerLevel(sim.LogInfo)

	//time.Sleep(5 * time.Second)

	c, err := sc.Connect(name)
	if err != nil {
		return nil, err
	}

	<-c // wait connection confirmation

	s.Logger.Info("Still working on connection to MSFS...")

	for {
		if <-sc.ConnectSysEventSim() {
			break // wait sim start
		}
	}

	s.Logger.Info("Connected to MSFS!")

	connectToMsfsInProgress = false
	lastMessageReceived = time.Now()

	return sc, nil
}

func (s *SimGo) TrackWithRecover(name string, report interface{}, maxTries int, trackID int) {
	maxTriesInitial = maxTries

	go s.recoverer(maxTries, trackID, func() {
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
					timeOut5s := time.Now().Add(-5 * time.Second)
					timeOutConnection := time.Now().Add(5 * time.Minute)
					//s.Logger.Infof("Timeout checker: %v %v", connectToMsfsInProgress, lastMessageReceived)
					if connectToMsfsInProgress && !lastMessageReceived.IsZero() && lastMessageReceived.Before(timeOutConnection) {
						panic("Connection was not confirmed for 5m. Cancel tracking")
					}
					if !connectToMsfsInProgress && !lastMessageReceived.IsZero() && lastMessageReceived.Before(timeOut5s) {
						s.Logger.Info("Last received message was received 5 sec ago. Cancel tracking")
						connectToMsfsInProgress = true
						cancel()
					}
				}
			}
		}()

		wait.Wait()

		s.Logger.Debug(ctx.Err())

		if ctx.Err() != nil {
			panic(ctx.Err().Error())
		}

		panic("Exiting from tracker routine...")
	})
}

func (s *SimGo) track(name string, report interface{}, ctx context.Context, wg *sync.WaitGroup) {
	sc, err := s.connect(ctx, name)
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
	paused := sc.ConnectSysEventPause()
	airloaded := sc.ConnectSysEventAircraftLoaded()

	for {
		select {
		case <-ctx.Done():
			s.Logger.Warning("Tracking routine will exit")
			return
		case sv := <-cSimVar:
			lastMessageReceived = time.Now()
			s.TrackEvent <- convertToInterface(reflect.ValueOf(report), sv)
		case r := <-paused:
			s.Logger.Debugf("Sim is paused: %v", r)
		case r := <-airloaded:
			s.Logger.Debugf("Aircraft: %v", r)
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

func (s *SimGo) recoverer(maxPanics, id int, f func()) {
	defer func() {
		if err := recover(); err != nil {
			s.Logger.Error(err)
			if maxPanics == 0 {
				panic("SimGo exceeded max tries. Exiting...")
			} else {
				if err.(string) != "context canceled" {
					maxPanics -= 1
					s.Logger.Info("Panic caused by error")
				} else {
					maxPanics = maxTriesInitial
				}
				s.Logger.Info("Recovering...")
				go s.recoverer(maxPanics, id, f)
			}
		}
	}()
	f()
}
