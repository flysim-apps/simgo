package main

import (
	"context"

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

	if err := sim.OffsetPolling("event", simgo.FSUIPC_Offset_Linkage{}, 1000); err != nil {
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
			logger.Debugf("Event: %+v", result)
		case result := <-payloadChan:
			logger.Debugf("Payload: %+v", result)
		case <-sim.TrackFailed:
			logger.Errorf("Track failed: %s", sim.Error.Error())
		}
	}

}
