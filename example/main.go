package main

import (
	"context"
	"reflect"

	"github.com/flysim-apps/simgo"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("simgo")

func main() {
	sim := simgo.NewSimGo(logger, simgo.FSUIPC)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := sim.FSUIPC(ctx, "ws://localhost:2048/fsuipc/"); err != nil {
		panic("Unable to establish websocket connection: " + err.Error())
	}

	if err := sim.OffsetPolling("event", simgo.Offsets{}, 1000); err != nil {
		sim.Logger.Errorf("Failed to obtain polling: %s", err.Error())
	}
	if err := sim.Payload(10000); err != nil {
		sim.Logger.Errorf("Failed to obtain payload: %s", err.Error())
	}

	eventChan := make(chan interface{})
	payloadChan := make(chan interface{})

	sim.ReadData("event", simgo.Report{}, eventChan, payloadChan)

	for {
		select {
		case result := <-eventChan:
			logger.Debugf("===================================================================================")
			val := reflect.ValueOf(result)
			for i := 0; i < val.Type().NumField(); i++ {
				logger.Debugf("%s = %+v", val.Type().Field(i).Name, val.Field(i))
			}
		case result := <-payloadChan:
			logger.Debugf("Payload: %+v", result)
		case <-sim.TrackFailed:
			logger.Errorf("Track failed: %s", sim.Error.Error())
		}
	}

}
