package simgo

import (
	"fmt"
	"reflect"
	"strings"

	sim "github.com/micmonay/simconnect"
)

const (
	MSFS_CONNECTION_READY = "msfs_connection_ready"
	AIRCARFT_POSITION     = "aircraft_position"
	AIRCARFT_STATE        = "aircraft_state"
	SYSTEM_STATE          = "system_state"
	FLIGHT_DATA           = "flight_data"
	LANDING_REPORT        = "passcargo_landing"
	FAILURE_REPORT        = "passcargo_failures"
)

const (
	STATE_OK                = 0
	STATE_CONNECTION_READY  = 1
	STATE_FATAL             = -1
	STATE_CONNECTION_ERROR  = -2
	STATE_CONNECTION_FAILED = -3
)

type Report struct {
	AltitudeAboveGround   float64 `name:"PLANE ALT ABOVE GROUND" unit:"feet"`
	PlaneAltitude         float64 `name:"PLANE ALTITUDE" unit:"Radians"`
	Altitude              float64 `name:"INDICATED ALTITUDE" unit:"feet"`
	RadioHeight           float64 `name:"RADIO HEIGHT" unit:"feet"`
	Latitude              float64 `name:"PLANE LATITUDE" unit:"degrees"`
	Longitude             float64 `name:"PLANE LONGITUDE" unit:"degrees"`
	Heading               float64 `name:"PLANE HEADING DEGREES TRUE" unit:"degrees"`
	HeadingMagnetic       float64 `name:"PLANE HEADING DEGREES MAGNETIC" unit:"degrees"`
	Airspeed              float64 `name:"AIRSPEED INDICATED" unit:"knot"`
	AirspeedTrue          float64 `name:"AIRSPEED TRUE" unit:"knot"`
	AirspeedMach          float64 `name:"AIRSPEED MACH" unit:"Mach"`
	VerticalSpeed         float64 `name:"VERTICAL SPEED" unit:"ft/min"`
	Flaps                 float64 `name:"TRAILING EDGE FLAPS LEFT ANGLE" unit:"degrees"`
	Trim                  float64 `name:"ELEVATOR TRIM PCT" unit:"percent"`
	RudderTrim            float64 `name:"RUDDER TRIM PCT" unit:"percent"`
	WindDirection         float64 `name:"AMBIENT WIND DIRECTION" unit:"degrees"`
	WindVelocity          float64 `name:"AMBIENT WIND VELOCITY" unit:"knots"`
	SurfaceType           float64 `name:"SURFACE TYPE" unit:"Enum"`
	OnAnyRunway           float64 `name:"ON ANY RUNWAY" unit:"Bool"`
	GroundVelocity        float64 `name:"GROUND VELOCITY" unit:"knots"`
	TouchdownLat          float64 `name:"PLANE TOUCHDOWN LATITUDE" unit:"radians"`
	TouchdownLon          float64 `name:"PLANE TOUCHDOWN LONGITUDE" unit:"radians"`
	TouchdownBank         float64 `name:"PLANE TOUCHDOWN BANK DEGREES" unit:"degrees"`
	TouchdownPitch        float64 `name:"PLANE TOUCHDOWN PITCH DEGREES" unit:"degrees"`
	Bank                  float64 `name:"PLANE BANK DEGREES" unit:"Radians"`
	Title                 string  `name:"TITLE" unit:"String"`
	OnGround              float64 `name:"SIM ON GROUND" unit:"Bool"`
	APUSwitch             float64 `name:"APU SWITCH" unit:"Bool"`
	BatterySwitch         float64 `name:"ELECTRICAL MASTER BATTERY" unit:"Bool"`
	ExtPowerOn            float64 `name:"EXTERNAL POWER ON" unit:"Bool"`
	DoorsClosed           float64 `name:"CANOPY OPEN" unit:"Bool"`
	NavLights             float64 `name:"LIGHT NAV" unit:"Bool"`
	BeaconLights          float64 `name:"LIGHT BEACON" unit:"Bool"`
	FastenSeatBealts      float64 `name:"CABIN SEATBELTS ALERT SWITCH" unit:"Bool"`
	StallWarning          float64 `name:"STALL WARNING" unit:"Bool"`
	OverspeedWarning      float64 `name:"OVERSPEED WARNING" unit:"Bool"`
	TaxiLights            float64 `name:"LIGHT TAXI" unit:"Bool"`
	LandingLights         float64 `name:"LIGHT LANDING" unit:"Bool"`
	StrobeLights          float64 `name:"LIGHT STROBE" unit:"Bool"`
	LogoLights            float64 `name:"LIGHT LOGO" unit:"Bool"`
	InParkingState        float64 `name:"PLANE IN PARKING STATE" unit:"Bool"`
	NoSmokingOn           float64 `name:"CABIN NO SMOKING ALERT SWITCH" unit:"Bool"`
	PushbackAngle         float64 `name:"PUSHBACK ANGLE" unit:"radians"`
	PushbackState         float64 `name:"PUSHBACK STATE" unit:"Enum"`
	PushbackAttached      float64 `name:"PUSHBACK ATTACHED" unit:"Bool"`
	PushbackAvailable     float64 `name:"PUSHBACK AVAILABLE" unit:"Bool"`
	BrakeParkingIndicator float64 `name:"BRAKE PARKING INDICATOR" unit:"Bool"`
	GearIsOnGround        float64 `name:"GEAR IS ON GROUND" unit:"Bool"`
	GForce                float64 `name:"G FORCE" unit:"GForce"`
	GForceMin             float64 `name:"MIN G FORCE" unit:"GForce"`
	GForceMax             float64 `name:"MAX G FORCE" unit:"GForce"`
	NumberOfEngines       float64 `name:"NUMBER OF ENGINES" unit:"number"`
	Engine1Combustion     float64 `name:"ENG COMBUSTION:index" index:"1" unit:"bool"`
	Engine2Combustion     float64 `name:"ENG COMBUSTION:index" index:"2" unit:"bool"`
	Engine3Combustion     float64 `name:"ENG COMBUSTION:index" index:"3" unit:"bool"`
	Engine4Combustion     float64 `name:"ENG COMBUSTION:index" index:"4" unit:"bool"`
	Engine1Failed         float64 `name:"ENG FAILED:index" index:"1" unit:"bool"`
	Engine2Failed         float64 `name:"ENG FAILED:index" index:"2" unit:"bool"`
	Engine3Failed         float64 `name:"ENG FAILED:index" index:"3" unit:"bool"`
	Engine4Failed         float64 `name:"ENG FAILED:index" index:"4" unit:"bool"`
	Engine1N1Rpm          float64 `name:"ENG N1 RPM:index" index:"1" unit:"rpm"`
	Engine2N1Rpm          float64 `name:"ENG N1 RPM:index" index:"2" unit:"rpm"`
	Engine3N1Rpm          float64 `name:"ENG N1 RPM:index" index:"3" unit:"rpm"`
	Engine4N1Rpm          float64 `name:"ENG N1 RPM:index" index:"4" unit:"rpm"`
	Engine1N2Rpm          float64 `name:"ENG N2 RPM:index" index:"1" unit:"rpm"`
	Engine2N2Rpm          float64 `name:"ENG N2 RPM:index" index:"2" unit:"rpm"`
	Engine3N2Rpm          float64 `name:"ENG N2 RPM:index" index:"3" unit:"rpm"`
	Engine4N2Rpm          float64 `name:"ENG N2 RPM:index" index:"4" unit:"rpm"`
	Engine1TurbN1         float64 `name:"TURB ENG N1:index" index:"1" unit:"percent"`
	Engine2TurbN1         float64 `name:"TURB ENG N1:index" index:"2" unit:"percent"`
	Engine3TurbN1         float64 `name:"TURB ENG N1:index" index:"3" unit:"percent"`
	Engine4TurbN1         float64 `name:"TURB ENG N1:index" index:"4" unit:"percent"`
	Engine1TurbN2         float64 `name:"TURB ENG N2:index" index:"1" unit:"percent"`
	Engine2TurbN2         float64 `name:"TURB ENG N2:index" index:"2" unit:"percent"`
	Engine3TurbN2         float64 `name:"TURB ENG N2:index" index:"3" unit:"percent"`
	Engine4TurbN2         float64 `name:"TURB ENG N2:index" index:"4" unit:"percent"`
	ZuluTime              float64 `name:"ZULU TIME" unit:"Seconds"`
	ZuluDayOfWeek         float64 `name:"ZULU DAY OF WEEK" unit:"number"`
	ZuluDayOfMonth        float64 `name:"ZULU DAY OF MONTH" unit:"number"`
	ZuluMonthOfYear       float64 `name:"ZULU MONTH OF YEAR" unit:"number"`
	ZuluDayOfYear         float64 `name:"ZULU DAY OF YEAR" unit:"number"`
	ZuluYear              float64 `name:"ZULU YEAR" unit:"number"`
}

type FlightEntry struct {
	Payload string `json:"payload"`
}

func getValue(field reflect.Value, simVar sim.SimVar) {
	if strings.Contains(string(simVar.Unit), "String") {
		field.SetString(simVar.GetString())
	} else if simVar.Unit == "SIMCONNECT_DATA_LATLONALT" {
		data, _ := simVar.GetDataLatLonAlt()
		field.SetString(fmt.Sprintf("%#v", data))
	} else if simVar.Unit == "SIMCONNECT_DATA_XYZ" {
		data, _ := simVar.GetDataXYZ()
		field.SetString(fmt.Sprintf("%#v", data))
	} else if simVar.Unit == "SIMCONNECT_DATA_WAYPOINT" {
		data, _ := simVar.GetDataWaypoint()
		field.SetString(fmt.Sprintf("%#v", data))
	} else {
		f, _ := simVar.GetFloat64()
		field.SetFloat(f)
	}
}
