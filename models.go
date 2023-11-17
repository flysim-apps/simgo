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

type Offsets struct {
	Agl                  float64      `address:"0x6020" type:"float" size:"8" fsuipc:"feet"`
	Alt                  int64        `address:"0x0570" type:"int" size:"8" fsuipc:"fractional"`
	AltRadio             int64        `address:"0x31E4" type:"int" size:"4" fsuipc:"radio"`
	Lat                  float64      `address:"0x0560" type:"int" size:"8" fsuipc:"lat"`
	Lon                  float64      `address:"0x0568" type:"int" size:"8" fsuipc:"lng"`
	Heading              float64      `address:"0x0580" type:"uint" size:"4" fsuipc:"degrees"`
	MagVar               float64      `address:"0x02A0" type:"int" size:"2" fsuipc:"magvar"`
	Airspeed             int          `address:"0x02BC" type:"int" size:"4" fsuipc:"knots"`
	AirspeedTrue         int          `address:"0x02B8" type:"int" size:"4" fsuipc:"knots"`
	AirspeedMach         float64      `address:"0x11C6" type:"uint" size:"2" fsuipc:"mach"`
	VerticalSpeed        int          `address:"0x02C8" type:"int" size:"4" fsuipc:"ftm"`
	FlapsLeft            int          `address:"0x0BE0" type:"int" size:"4" fsuipc:"position"`
	FlapsRight           int          `address:"0x0BE4" type:"int" size:"4" fsuipc:"position"`
	ElevatorTrim         float64      `address:"0x0BC0" type:"int" size:"2" fsuipc:"percent"`
	RudderTrim           float64      `address:"0x0C04" type:"int" size:"2" fsuipc:"percent"`
	AleronTrim           float64      `address:"0x0C02" type:"int" size:"2" fsuipc:"percent"`
	Pitch                float64      `address:"0x0578" type:"int" size:"4" fsuipc:"degrees"`
	Bank                 float64      `address:"0x057C" type:"int" size:"4" fsuipc:"degrees"`
	AmbientWindDirection float64      `address:"0x3490" type:"float" size:"8" fsuipc:"wind"`
	AmbientWindVelocity  float64      `address:"0x0E90" type:"int" size:"2"`
	AmbientTemperature   float64      `address:"0x34A8" type:"float" size:"8"`
	SurfaceType          int          `address:"0x31E8" type:"int" size:"4"`
	SurfaceCondition     int          `address:"0x31EC" type:"int" size:"4"`
	GroundVelocity       int          `address:"0x02B4" type:"int" size:"4" fsuipc:"velocity"`
	Title                string       `address:"0x3D00" type:"string" size:"256"`
	OnGround             int          `address:"0x0366" type:"uint" size:"2" fsuipc:"bool"`
	APUSwitch            int          `address:"0x029D" type:"uint" size:"1" fsuipc:"bool"`
	BatterySwitch        int          `address:"0x281C" type:"uint" size:"4" fsuipc:"bool"`
	ExtPowerOn           int          `address:"0x07AB" type:"uint" size:"1" fsuipc:"bool"`
	IsDoorsOpen          int          `address:"0x2A70" type:"uint" size:"8" fsuipc:"bool"`
	Lights               map[int]bool `address:"0x0D0C" type:"bits" size:"2"`
	FastenSeatBealts     int          `address:"0x341D" type:"int" size:"1" fsuipc:"bool"`
	NoSmoking            int          `address:"0x341C" type:"int" size:"1" fsuipc:"bool"`
	StallWarning         int          `address:"0x036C" type:"int" size:"1" fsuipc:"bool"`
	OverspeedWarning     int          `address:"0x036D" type:"int" size:"1" fsuipc:"bool"`
	InParkingState       int          `address:"0x062B" type:"int" size:"1" fsuipc:"bool"`
	BrakeParkingPosition int          `address:"0x0BC8" type:"int" size:"2"`
	BrakeIndicator       int          `address:"0x0BCA" type:"int" size:"2"`
	PushbackAngle        float64      `address:"0x0334" type:"int" size:"4" fsuipc:"degrees"`
	PushbackStatus       int          `address:"0x31F0" type:"int" size:"4"`
	GearHandlePosition   int          `address:"0x0BE8" type:"int" size:"4" fsuipc:"position"`
	GForce               float64      `address:"0x11BA" type:"int" size:"2" fsuipc:"GForce"`
	NumberOfEngines      int          `address:"0x0AEC" type:"int" size:"2"`
	Engine1Combustion    int          `address:"0x0894" type:"int" size:"2" fsuipc:"bool"`
	Engine2Combustion    int          `address:"0x092C" type:"int" size:"2" fsuipc:"bool"`
	Engine3Combustion    int          `address:"0x09C4" type:"int" size:"2" fsuipc:"bool"`
	Engine4Combustion    int          `address:"0x0A5C" type:"int" size:"2" fsuipc:"bool"`
	EngineFailed         map[int]bool `address:"0x0B6B" type:"bits" size:"1"`
	Engine1TurbN1        int          `address:"0x0898" type:"int" size:"2" fsuipc:"percent"`
	Engine2TurbN1        int          `address:"0x0930" type:"int" size:"2" fsuipc:"percent"`
	Engine3TurbN1        int          `address:"0x09C8" type:"int" size:"2" fsuipc:"percent"`
	Engine4TurbN1        int          `address:"0x0A60" type:"int" size:"2" fsuipc:"percent"`
	Engine1TurbN2        int          `address:"0x0896" type:"int" size:"2" fsuipc:"percent"`
	Engine2TurbN2        int          `address:"0x092E" type:"int" size:"2" fsuipc:"percent"`
	Engine3TurbN2        int          `address:"0x09C6" type:"int" size:"2" fsuipc:"percent"`
	Engine4TurbN2        int          `address:"0x0A5E" type:"int" size:"2" fsuipc:"percent"`
	LocalTime            int          `address:"0x023A" type:"int" size:"1"`
	ZuluHour             int          `address:"0x023B" type:"int" size:"1"`
	ZuluMinute           int          `address:"0x023C" type:"int" size:"1"`
	ZuluDayOfWeek        int          `address:"0x0243" type:"int" size:"1"`
	ZuluDayOfMonth       int          `address:"0x023D" type:"int" size:"1"`
	ZuluMonthOfYear      int          `address:"0x0242" type:"int" size:"1"`
	ZuluDayOfYear        int          `address:"0x023E" type:"int" size:"2"`
	ZuluYear             int          `address:"0x0240" type:"int" size:"2"`
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
