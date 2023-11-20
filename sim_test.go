package simgo

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/flysim-apps/simgo/simconnect"
	"github.com/stretchr/testify/assert"
)

type ReportTest struct {
	fieldAddress        uintptr
	AltitudeAboveGround float64 `name:"PLANE ALT ABOVE GROUND" unit:"feet"`
	PlaneAltitude       float64 `name:"PLANE ALTITUDE" unit:"feet"`
	Altitude            float64 `name:"INDICATED ALTITUDE" unit:"feet"`
	RadioHeight         float64 `name:"RADIO HEIGHT" unit:"feet"`
	Latitude            float64 `name:"PLANE LATITUDE" unit:"degrees"`
	Longitude           float64 `name:"PLANE LONGITUDE" unit:"degrees"`
	Heading             float64 `name:"PLANE HEADING DEGREES TRUE" unit:"degrees"`
	HeadingMagnetic     float64 `name:"PLANE HEADING DEGREES MAGNETIC" unit:"degrees"`
}

func TestConvertToSimSimVar(t *testing.T) {
	vars := convertToSimSimVar(reflect.ValueOf(ReportTest{}))
	assert.Greater(t, len(vars), 0, "vars is zero")
}

func TestConvertToInterface(t *testing.T) {
	vars := []simconnect.SimVar{
		simconnect.SimVarSimOnGround(),
	}
	v := ReportTest{}
	convertToInterface(reflect.ValueOf(v), vars)
}

func TestConvertWithReflect(t *testing.T) {
	some := func(v interface{}) {
		vars := convertToSimSimVar(reflect.ValueOf(v))
		assert.Greater(t, len(vars), 0, "vars is zero")

		val := reflect.ValueOf(v)
		r := reflect.New(reflect.TypeOf(v))

		for _, simVar := range vars {
			t.Logf("iterateSimVars(): Name: %s  Index: %b    Unit: %s\n", simVar.Name, simVar.Index, simVar.Unit)
			for j := 0; j < val.NumField(); j++ {
				nameTag, _ := val.Type().Field(j).Tag.Lookup("name")
				indexTag, _ := val.Type().Field(j).Tag.Lookup("index")
				if indexTag == "" {
					indexTag = "0"
				}

				idx, _ := strconv.Atoi(indexTag)

				if simVar.Index == idx && simVar.Name == nameTag {
					getValue(r.Elem().Field(j), simVar)
				}
			}
		}
	}

	some(Report{})
}
