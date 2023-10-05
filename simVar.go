package simgo

import (
	sim "github.com/micmonay/simconnect"
)

var SIMVAR_APU_SWITCH = SimVar{
	Index:    0,
	Name:     "APU SWITCH",
	Unit:     sim.UnitBool,
	Settable: false,
}

var SIMVAR_ON_ANY_RUNWAY = SimVar{
	Index:    0,
	Name:     "ON ANY RUNWAY",
	Unit:     sim.UnitBool,
	Settable: false,
}

var SIMVAR_LIGHT_LANDING_ON = SimVar{
	Index:    0,
	Name:     "LIGHT LANDING ON",
	Unit:     sim.UnitBool,
	Settable: false,
}

var SIMVAR_PLANE_IN_PARKING_STATE = SimVar{
	Index:    0,
	Name:     "PLANE IN PARKING STATE",
	Unit:     sim.UnitBool,
	Settable: false,
}

var SIMVAR_EXT_PWR_CON_ON = SimVar{
	Index:    0,
	Name:     "EXTERNAL POWER CONNECTION ON",
	Unit:     sim.UnitBool,
	Settable: false,
}

var SIMVAR_EXT_PWR_ON = SimVar{
	Index:    1,
	Name:     "EXTERNAL POWER ON",
	Unit:     sim.UnitBool,
	Settable: false,
}

var SIMVAR_PUSHBACK_ATTACHED = SimVar{
	Index:    0,
	Name:     "PUSHBACK ATTACHED",
	Unit:     sim.UnitBool,
	Settable: false,
}

var SIMVAR_PUSHBACK_AVAILABLE = SimVar{
	Index:    0,
	Name:     "PUSHBACK AVAILABLE",
	Unit:     sim.UnitBool,
	Settable: false,
}

var SIMVAR_GEAR_IS_ON_GROUND = SimVar{
	Index:    0,
	Name:     "GEAR IS ON GROUND",
	Unit:     sim.UnitBool,
	Settable: false,
}

var SIMVAR_N1_RPM_IDX_1 = SimVar{
	Index:    1,
	Name:     "ENG N1 RPM:index",
	Unit:     sim.UnitRpm,
	Settable: false,
}
var SIMVAR_N1_RPM_IDX_2 = SimVar{
	Index:    2,
	Name:     "ENG N1 RPM:index",
	Unit:     sim.UnitRpm,
	Settable: false,
}
var SIMVAR_N1_RPM_IDX_3 = SimVar{
	Index:    3,
	Name:     "ENG N1 RPM:index",
	Unit:     sim.UnitRpm,
	Settable: false,
}

var SIMVAR_N1_RPM_IDX_4 = SimVar{
	Index:    4,
	Name:     "ENG N1 RPM:index",
	Unit:     sim.UnitRpm,
	Settable: false,
}

var SIMVAR_ENG_N2_RPM_IDX_1 = SimVar{
	Index:    1,
	Name:     "ENG N2 RPM:index",
	Unit:     sim.UnitRpm,
	Settable: false,
}

var SIMVAR_ENG_N2_RPM_IDX_2 = SimVar{
	Index:    2,
	Name:     "ENG N2 RPM:index",
	Unit:     sim.UnitRpm,
	Settable: false,
}

var SIMVAR_ENG_N2_RPM_IDX_3 = SimVar{
	Index:    3,
	Name:     "ENG N2 RPM:index",
	Unit:     sim.UnitRpm,
	Settable: false,
}

var SIMVAR_ENG_N2_RPM_IDX_4 = SimVar{
	Index:    4,
	Name:     "ENG N2 RPM:index",
	Unit:     sim.UnitRpm,
	Settable: false,
}

var SIMVAR_TURB_ENG_N1_IDX_1 = SimVar{
	Index:    1,
	Name:     "TURB ENG N1:index",
	Unit:     sim.UnitPercent,
	Settable: false,
}
var SIMVAR_TURB_ENG_N1_IDX_2 = SimVar{
	Index:    2,
	Name:     "TURB ENG N1:index",
	Unit:     sim.UnitPercent,
	Settable: false,
}
var SIMVAR_TURB_ENG_N1_IDX_3 = SimVar{
	Index:    3,
	Name:     "TURB ENG N1:index",
	Unit:     sim.UnitPercent,
	Settable: false,
}
var SIMVAR_TURB_ENG_N1_IDX_4 = SimVar{
	Index:    4,
	Name:     "TURB ENG N1:index",
	Unit:     sim.UnitPercent,
	Settable: false,
}

var SIMVAR_TURB_ENG_N2_IDX_1 = SimVar{
	Index:    1,
	Name:     "TURB ENG N2:index",
	Unit:     sim.UnitPercent,
	Settable: false,
}
var SIMVAR_TURB_ENG_N2_IDX_2 = SimVar{
	Index:    2,
	Name:     "TURB ENG N2:index",
	Unit:     sim.UnitPercent,
	Settable: false,
}
var SIMVAR_TURB_ENG_N2_IDX_3 = SimVar{
	Index:    3,
	Name:     "TURB ENG N2:index",
	Unit:     sim.UnitPercent,
	Settable: false,
}
var SIMVAR_TURB_ENG_N2_IDX_4 = SimVar{
	Index:    4,
	Name:     "TURB ENG N2:index",
	Unit:     sim.UnitPercent,
	Settable: false,
}

var SIMVAR_ENG_COMBUSTION_IDX_1 = SimVar{
	Index:    1,
	Name:     "ENG COMBUSTION:index",
	Unit:     sim.UnitBool,
	Settable: false,
}
var SIMVAR_ENG_COMBUSTION_IDX_2 = SimVar{
	Index:    2,
	Name:     "ENG COMBUSTION:index",
	Unit:     sim.UnitBool,
	Settable: false,
}
var SIMVAR_ENG_COMBUSTION_IDX_3 = SimVar{
	Index:    3,
	Name:     "ENG COMBUSTION:index",
	Unit:     sim.UnitBool,
	Settable: false,
}
var SIMVAR_ENG_COMBUSTION_IDX_4 = SimVar{
	Index:    4,
	Name:     "ENG COMBUSTION:index",
	Unit:     sim.UnitBool,
	Settable: false,
}

var SIMVAR_ELECTRICAL_MASTER_BATTERY_IDX_1 = SimVar{
	Index:    1,
	Name:     "ELECTRICAL MASTER BATTERY:index",
	Unit:     sim.UnitBool,
	Settable: true,
}

var SIMVAR_ELECTRICAL_MASTER_BATTERY_IDX_2 = SimVar{
	Index:    2,
	Name:     "ELECTRICAL MASTER BATTERY:index",
	Unit:     sim.UnitBool,
	Settable: true,
}

var SIMVAR_PLANE_TD_LAT = SimVar{
	Index:    0,
	Name:     "PLANE TOUCHDOWN LATITUDE",
	Unit:     sim.UnitRadians,
	Settable: false,
}

var SIMVAR_PLANE_TD_LON = SimVar{
	Index:    0,
	Name:     "PLANE TOUCHDOWN LONGITUDE",
	Unit:     sim.UnitRadians,
	Settable: false,
}

var SIMVAR_PLANE_TD_NORMAL_VELOCITY = SimVar{
	Index:    0,
	Name:     "PLANE TOUCHDOWN NORMAL VELOCITY",
	Unit:     sim.UnitFeetpersecond,
	Settable: false,
}

var SIMVAR_PLANE_TD_BANK_DEGREES = SimVar{
	Index:    0,
	Name:     "PLANE TOUCHDOWN BANK DEGREES",
	Unit:     sim.UnitDegrees,
	Settable: false,
}

var SIMVAR_PLANE_TD_PITCH_DEGREES = SimVar{
	Index:    0,
	Name:     "PLANE TOUCHDOWN PITCH DEGREES",
	Unit:     sim.UnitDegrees,
	Settable: false,
}

var SIMVAR_LIGHT_LANDING = SimVar{
	Index:    0,
	Name:     "LIGHT LANDING",
	Unit:     sim.UnitBool,
	Settable: false,
}

func BuiltInSimVars() []SimVar {
	return []SimVar{
		SIMVAR_APU_SWITCH,
		SIMVAR_ON_ANY_RUNWAY,
		SIMVAR_LIGHT_LANDING_ON,
		SIMVAR_PLANE_IN_PARKING_STATE,
		SIMVAR_EXT_PWR_CON_ON,
		SIMVAR_EXT_PWR_ON,
		SIMVAR_PUSHBACK_ATTACHED,
		SIMVAR_PUSHBACK_AVAILABLE,
		SIMVAR_GEAR_IS_ON_GROUND,
		SIMVAR_N1_RPM_IDX_1,
		SIMVAR_N1_RPM_IDX_2,
		SIMVAR_N1_RPM_IDX_3,
		SIMVAR_N1_RPM_IDX_4,
		SIMVAR_ENG_N2_RPM_IDX_1,
		SIMVAR_ENG_N2_RPM_IDX_2,
		SIMVAR_ENG_N2_RPM_IDX_3,
		SIMVAR_ENG_N2_RPM_IDX_4,
		SIMVAR_ENG_N2_RPM_IDX_4,
		SIMVAR_TURB_ENG_N1_IDX_1,
		SIMVAR_TURB_ENG_N1_IDX_2,
		SIMVAR_TURB_ENG_N1_IDX_3,
		SIMVAR_TURB_ENG_N1_IDX_4,
		SIMVAR_TURB_ENG_N2_IDX_1,
		SIMVAR_TURB_ENG_N2_IDX_2,
		SIMVAR_TURB_ENG_N2_IDX_3,
		SIMVAR_TURB_ENG_N2_IDX_4,
		SIMVAR_ENG_COMBUSTION_IDX_1,
		SIMVAR_ENG_COMBUSTION_IDX_2,
		SIMVAR_ENG_COMBUSTION_IDX_3,
		SIMVAR_ENG_COMBUSTION_IDX_4,
		SIMVAR_ELECTRICAL_MASTER_BATTERY_IDX_1,
		SIMVAR_ELECTRICAL_MASTER_BATTERY_IDX_2,
		SIMVAR_PLANE_TD_LAT,
		SIMVAR_PLANE_TD_LON,
		SIMVAR_PLANE_TD_NORMAL_VELOCITY,
		SIMVAR_PLANE_TD_BANK_DEGREES,
		SIMVAR_PLANE_TD_PITCH_DEGREES,
		SIMVAR_LIGHT_LANDING,
	}
}
