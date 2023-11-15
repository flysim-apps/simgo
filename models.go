package simgo

import (
	"fmt"
	"reflect"
	"strings"

	sim "github.com/flysim-apps/simgo/simconnect"
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

type Report struct {
	AltitudeAboveGround  int     `json:"" unit:"feet"`
	PlaneAltitude        int     `json:"" unit:"feet"`
	Altitude             int     `json:"" unit:"feet"`
	RadioHeight          int     `json:"" unit:"feet"`
	Latitude             float64 `json:"" unit:"degrees"`
	Longitude            float64 `json:"" unit:"degrees"`
	Heading              float64 `json:"" unit:"degrees"`
	HeadingMagnetic      float64 `json:"" unit:"degrees"`
	Airspeed             int     `json:"" unit:"knots"`
	AirspeedTrue         int     `json:"" unit:"knots"`
	AirspeedMach         float64 `json:"" unit:"mach"`
	VerticalSpeed        int     `json:"" unit:"ft/s"`
	FlapsLeft            int     `json:"" unit:"position"`
	FlapsRight           int     `json:"" unit:"position"`
	Trim                 int     `json:"" unit:"position"`
	RudderTrim           int     `json:"" unit:"position"`
	AleronTrim           int     `json:"" unit:"position"`
	AmbientWindDirection float64 `json:"" unit:"degrees"`
	AmbientWindVelocity  int     `json:"" unit:"knots"`
	AmbientTemperature   int     `json:"" unit:"celcium"`
	SurfaceType          int     `json:""`
	SurfaceCondition     int     `json:""`
	GroundVelocity       int     `json:"" unit:"knots"`
	Pitch                float64 `json:"" unit:"degrees"`
	Bank                 float64 `json:"" unit:"degrees"`
	Title                string  `json:"" type:"string" size:"256"`
	OnGround             bool    `json:"" unit:"bool"`
	APUSwitch            bool    `json:"" unit:"bool"`
	BatterySwitch        bool    `json:"" unit:"bool"`
	ExtPowerOn           bool    `address:"0x07AB" type:"uint" size:"1" unit:"bool"`
	DoorsClosed          bool    `address:"0x2A70" type:"uint" size:"8" unit:"bool"`
	LightsNav            bool    `address:"0x0D0C" type:"bits" size:"2" index:"0" unit:"bool"`
	LightsBeacon         bool    `address:"0x0D0C" type:"bits" size:"2" index:"1" unit:"bool"`
	LightsLanding        bool    `address:"0x0D0C" type:"bits" size:"2" index:"2" unit:"bool"`
	LightsTaxi           bool    `address:"0x0D0C" type:"bits" size:"2" index:"3" unit:"bool"`
	LightsStrobe         bool    `address:"0x0D0C" type:"bits" size:"2" index:"4" unit:"bool"`
	LightsInstruments    bool    `address:"0x0D0C" type:"bits" size:"2" index:"5" unit:"bool"`
	LightsRecognition    bool    `address:"0x0D0C" type:"bits" size:"2" index:"6" unit:"bool"`
	LightsWing           bool    `address:"0x0D0C" type:"bits" size:"2" index:"7" unit:"bool"`
	LightsLogo           bool    `address:"0x0D0C" type:"bits" size:"2" index:"8" unit:"bool"`
	LightsCabin          bool    `address:"0x0D0C" type:"bits" size:"2" index:"9" unit:"bool"`
	FastenSeatBealts     bool    `address:"0x341D" type:"int" size:"1" unit:"bool"`
	NoSmoking            bool    `address:"0x341C" type:"int" size:"1" unit:"bool"`
	StallWarning         bool    `address:"0x036C" type:"int" size:"1" unit:"bool"`
	OverspeedWarning     bool    `address:"0x036D" type:"int" size:"1" unit:"bool"`
	InParkingState       bool    `address:"0x062B" type:"int" size:"1" unit:"bool"`
	PushbackAngle        float64 `address:"0x0334" type:"unit" size:"4" unit:"degrees"`
	PushbackState        int     `address:"0x31F4" `
	GearHandlePosition   int     `address:"0x0BE8" unit:"position"`
	GForce               float64 `address:"0x1140" type:"uint" size:"8" unit:"GForce"`
	GForceMin            float64 `address:"0x34D8" type:"uint" size:"8" unit:"GForce"`
	GForceMax            float64 `address:"0x34D0" type:"uint" size:"8" unit:"GForce"`
	NumberOfEngines      int     `address:"0x0AEC"`
	Engine1Combustion    bool    `address:"0x0894" unit:"bool"`
	Engine2Combustion    bool    `address:"0x092C" unit:"bool"`
	Engine3Combustion    bool    `address:"0x09C4" unit:"bool"`
	Engine4Combustion    bool    `address:"0x0A5C" unit:"bool"`
	Engine1Failed        bool    `address:"0x0B6B" type:"bits" size:"1" index:"0" unit:"bool"`
	Engine2Failed        bool    `address:"0x0B6B" type:"bits" size:"1" index:"1" unit:"bool"`
	Engine3Failed        bool    `address:"0x0B6B" type:"bits" size:"1" index:"2" unit:"bool"`
	Engine4Failed        bool    `address:"0x0B6B" type:"bits" size:"1" index:"3" unit:"bool"`
	Engine1TurbN1        float64 `json:"" unit:"percent"`
	Engine2TurbN1        float64 `json:"" unit:"percent"`
	Engine3TurbN1        float64 `json:"" unit:"percent"`
	Engine4TurbN1        float64 `json:"" unit:"percent"`
	Engine1TurbN2        float64 `json:"" unit:"percent"`
	Engine2TurbN2        float64 `json:"" unit:"percent"`
	Engine3TurbN2        float64 `json:"" unit:"percent"`
	Engine4TurbN2        float64 `json:"" unit:"percent"`
	LocalTime            int     `json:"localTime"`
	ZuluHour             int     `json:"zuluHour"`
	ZuluMinute           int     `json:"zuluMinute"`
	ZuluDayOfWeek        int     `json:"zuluDayOfWeek"`
	ZuluDayOfMonth       int     `json:"zuluDayOfMonth"`
	ZuluMonthOfYear      int     `json:"zuluMonthOfYear"`
	ZuluDayOfYear        int     `json:"zuluDayOfYear"`
	ZuluYear             int     `json:"zuluYear"`
	BrakeParkingPosition int     `json:"breakeParkingPosition"`
	BrakeIndicator       bool    `json:"brakeIndicator" unit:"bool"`
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

type FSUIPC_Offset_Linkage struct {
	AltitudeAboveGround  int64        `address:"0x0020" type:"int" size:"8" unit:"feet"`
	PlaneAltitude        int64        `address:"0x0570" type:"int" size:"8" unit:"feet"`
	Altitude             int64        `address:"0x0590" type:"int" size:"8" unit:"feet"`
	RadioHeight          int64        `address:"0x31E4" type:"int" size:"8" unit:"feet"`
	Latitude             float64      `address:"0x0560" type:"uint" size:"8" unit:"degrees"`
	Longitude            float64      `address:"0x0568" type:"uint" size:"8" unit:"degrees"`
	Heading              float64      `address:"0x0580" type:"uint" size:"4" unit:"degrees"`
	HeadingMagnetic      float64      `address:"0x2B00" type:"uint" size:"8" unit:"degrees"`
	Airspeed             int          `address:"0x02BC" type:"int" size:"4" unit:"knots"`
	AirspeedTrue         int          `address:"0x02B8" type:"int" size:"4" unit:"knots"`
	AirspeedMach         float64      `address:"0x11C6" type:"uint" size:"2" unit:"mach"`
	VerticalSpeed        int          `address:"0x02C8" type:"int" size:"4" unit:"ft/s"`
	FlapsLeft            int          `address:"0x0BE0" type:"int" size:"4" unit:"position"`
	FlapsRight           int          `address:"0x0BE4" type:"int" size:"4" unit:"position"`
	Trim                 int          `address:"0x0BC0" type:"int" size:"4" unit:"position"`
	RudderTrim           int          `address:"0x0C04" type:"int" size:"4" unit:"position"`
	AleronTrim           int          `address:"0x0C02" type:"int" size:"4" unit:"position"`
	AmbientWindDirection float64      `address:"0x0E92" type:"uint" size:"8" unit:"wind"`
	AmbientWindVelocity  int          `address:"0x0E90" type:"int" size:"2" unit:"knots"`
	AmbientTemperature   int          `address:"0x0E8C" type:"int" size:"2" unit:"celcium"`
	SurfaceType          int          `address:"0x31E8" type:"int" size:"4"`
	SurfaceCondition     int          `address:"0x31EC" type:"int" size:"4"`
	GroundVelocity       int          `address:"0x02B4" type:"int" size:"4" unit:"knots"`
	Pitch                float64      `address:"0x0578" type:"uint" size:"8" unit:"degrees"`
	Bank                 float64      `address:"0x057C" type:"uint" size:"8" unit:"degrees"`
	Title                string       `address:"0x3D00" type:"string" size:"256"`
	OnGround             int          `address:"0x0366" type:"uint" size:"2" unit:"bool"`
	APUSwitch            int          `address:"0x029D" type:"uint" size:"1" unit:"bool"`
	BatterySwitch        int          `address:"0x281C" type:"uint" size:"4" unit:"bool"`
	ExtPowerOn           int          `address:"0x07AB" type:"uint" size:"1" unit:"bool"`
	DoorsClosed          int          `address:"0x2A70" type:"uint" size:"8" unit:"bool"`
	Lights               map[int]bool `address:"0x0D0C" type:"bits" size:"2"`
	LightsNav            bool         `address:"0x0D0C" type:"bits" size:"2" index:"0" unit:"bool"`
	LightsBeacon         bool         `address:"0x0D0C" type:"bits" size:"2" index:"1" unit:"bool"`
	LightsLanding        bool         `address:"0x0D0C" type:"bits" size:"2" index:"2" unit:"bool"`
	LightsTaxi           bool         `address:"0x0D0C" type:"bits" size:"2" index:"3" unit:"bool"`
	LightsStrobe         bool         `address:"0x0D0C" type:"bits" size:"2" index:"4" unit:"bool"`
	LightsInstruments    bool         `address:"0x0D0C" type:"bits" size:"2" index:"5" unit:"bool"`
	LightsRecognition    bool         `address:"0x0D0C" type:"bits" size:"2" index:"6" unit:"bool"`
	LightsWing           bool         `address:"0x0D0C" type:"bits" size:"2" index:"7" unit:"bool"`
	LightsLogo           bool         `address:"0x0D0C" type:"bits" size:"2" index:"8" unit:"bool"`
	LightsCabin          bool         `address:"0x0D0C" type:"bits" size:"2" index:"9" unit:"bool"`
	FastenSeatBealts     int          `address:"0x341D" type:"int" size:"1" unit:"bool"`
	NoSmoking            int          `address:"0x341C" type:"int" size:"1" unit:"bool"`
	StallWarning         int          `address:"0x036C" type:"int" size:"1" unit:"bool"`
	OverspeedWarning     int          `address:"0x036D" type:"int" size:"1" unit:"bool"`
	InParkingState       int          `address:"0x062B" type:"int" size:"1" unit:"bool"`
	PushbackAngle        float64      `address:"0x0334" type:"uint" size:"4" unit:"degrees"`
	PushbackState        int          `address:"0x31F4" type:"int" size:"4"`
	GearHandlePosition   int          `address:"0x0BE8" type:"int" size:"4" unit:"position"`
	GForce               float64      `address:"0x1140" type:"uint" size:"8" unit:"GForce"`
	GForceMin            float64      `address:"0x34D8" type:"uint" size:"8" unit:"GForce"`
	GForceMax            float64      `address:"0x34D0" type:"uint" size:"8" unit:"GForce"`
	NumberOfEngines      int          `address:"0x0AEC" type:"int" size:"2"`
	Engine1Combustion    int          `address:"0x0894" unit:"bool"`
	Engine2Combustion    int          `address:"0x092C" unit:"bool"`
	Engine3Combustion    int          `address:"0x09C4" unit:"bool"`
	Engine4Combustion    int          `address:"0x0A5C" unit:"bool"`
	EngineFailed         map[int]bool `address:"0x0B6B" type:"bits" size:"1"`
	Engine1Failed        bool         `address:"0x0B6B" type:"bits" size:"1" index:"0" unit:"bool"`
	Engine2Failed        bool         `address:"0x0B6B" type:"bits" size:"1" index:"1" unit:"bool"`
	Engine3Failed        bool         `address:"0x0B6B" type:"bits" size:"1" index:"2" unit:"bool"`
	Engine4Failed        bool         `address:"0x0B6B" type:"bits" size:"1" index:"3" unit:"bool"`
	Engine1TurbN1        float64      `address:"0x0898" unit:"percent"`
	Engine2TurbN1        float64      `address:"0x0930" unit:"percent"`
	Engine3TurbN1        float64      `address:"0x09C8" unit:"percent"`
	Engine4TurbN1        float64      `address:"0x0A60" unit:"percent"`
	Engine1TurbN2        float64      `address:"0x0896" unit:"percent"`
	Engine2TurbN2        float64      `address:"0x092E" unit:"percent"`
	Engine3TurbN2        float64      `address:"0x09C6" unit:"percent"`
	Engine4TurbN2        float64      `address:"0x0A5E" unit:"percent"`
	LocalTime            int          `address:"0x023A" type:"int" size:"1"`
	ZuluHour             int          `address:"0x023B" type:"int" size:"1"`
	ZuluMinute           int          `address:"0x023C" type:"int" size:"1"`
	ZuluDayOfWeek        int          `address:"0x0243" type:"int" size:"1"`
	ZuluDayOfMonth       int          `address:"0x023D" type:"int" size:"1"`
	ZuluMonthOfYear      int          `address:"0x0242" type:"int" size:"1"`
	ZuluDayOfYear        int          `address:"0x023E"`
	ZuluYear             int          `address:"0x0240"`
	BrakeParkingPosition int          `address:"0x0BC8"`
	BrakeIndicator       int          `address:"0x0BCA" unit:"bool"`
}

type FSUIPC_Offset_Payload struct {
	WeightUnit              string                   `json:"weightUnit"`
	VolumeUnit              string                   `json:"volumeUnit"`
	LengthUnit              string                   `json:"lengthUnit"`
	GrossWeight             float64                  `json:"grossWeight"`
	MaxGrossWeight          float64                  `json:"maxGrossWeight"`
	EmptyWeight             float64                  `json:"emptyWeight"`
	TotalPayloadWeight      float64                  `json:"totalPayloadWeight"`
	TotalFuelWeight         float64                  `json:"totalFuelWeight"`
	TotalFuelVolume         float64                  `json:"totalFuelVolume"`
	TotalFuelCapacityWeight float64                  `json:"totalFuelCapacityWeight"`
	TotalFuelCapacityVolume float64                  `json:"totalFuelCapacityVolume"`
	TotalFuelPercent        float64                  `json:"totalFuelPercent"`
	FuelTanks               []FSUIPC_FuelTank        `json:"fuelTanks"`
	PayloadStations         []FSUIPC_PayloadStations `json:"payloadStations"`
}

type FSUIPC_FuelTank struct {
	Index          int     `json:"index"`
	Name           string  `json:"name"`
	IsPresent      bool    `json:"isPresent"`
	Weight         float64 `json:"weight"`
	Volume         float64 `json:"volume"`
	Percent        float64 `json:"percent"`
	CapacityWeight float64 `json:"capacityWeight"`
	CapacityVolume float64 `json:"capacityVolume"`
}

type FSUIPC_PayloadStations struct {
	Index    int                `json:"index"`
	Name     string             `json:"name"`
	Weight   float64            `json:"weight"`
	Position map[string]float64 `json:"position"`
}

type FSUIPC_Command struct {
	Command  string          `json:"command"`
	Name     string          `json:"name"`
	Offsets  []FSUIPC_Offset `json:"offsets,omitempty"`
	Interval int             `json:"interval,omitempty"`
}

type FSUIPC_Command_Payload struct {
	Command    string `json:"command"`
	WeightUnit string `json:"weightUnit"`
	VolumeUnit string `json:"volumeUnit"`
	LengthUnit string `json:"lengthUnit"`
	Interval   int    `json:"interval,omitempty"`
}

type FSUIPC_Offset struct {
	Name    string `json:"name"`
	Address uint32 `json:"address"`
	Type    string `json:"type"`
	Size    int    `json:"size"`
}

type FSUIPC_Response struct {
	Success      bool                   `json:"success"`
	Command      string                 `json:"command"`
	Name         string                 `json:"name"`
	ErrorCode    string                 `json:"errorCode"`
	ErrorMessage string                 `json:"errorMessage"`
	Data         map[string]interface{} `json:"data"`
}

type Provider string

const (
	SimConnect Provider = "SimConnect"
	FSUIPC     Provider = "FSUIPC"
)
