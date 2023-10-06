# SimGo Connect

Current version: v1.0.0

Library provides abilitu to connect Microsoft Flight Simulator 2020 using SimConnect (SimConnect.dll) and get/set data from/to it. 

Written in Go.

## Features:
- built-in web sockets
- re-connection in case of any stucks
- predefined set of variables which will be enough for the first use, but this list is fully customizable

## ToDo:
- work with SimConnect events
- work w/o web sockets

## Getting Started

Make sure that you have copy of SimConnect.dll on root of your go project (you can copy dll provided in this repo). Otherwise you will get an error.

```
func main() {
    
    sim := simgo.NewSimGo(logging.MustGetLogger("simgo"))
    
    if err := Sim.StartWebSocket(":4050"); err != nil {
		panic(err.Error())
	}

    sim.TrackWithRecover("simgo", sim.Report{}, 5, 1)

    ...
```

You will set connection to Microsoft Flight Simulator 2020 and get notifications like this:

```
    ...
    go func() {
        for {
            select {
            case result := <-sim.TrackEvent:
                report := result.(simgo.Report)
                ...
                // do your stuff here
                ...
            }
        }
    }()
}    
```



## Known Projects

- PassCargo Stream Overlay - [https://passcargo.app](https://passcargo.app)
