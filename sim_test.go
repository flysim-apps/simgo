package simgo

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/micmonay/simconnect"
	"github.com/stretchr/testify/assert"
)

func TestConvertToSimSimVar(t *testing.T) {
	vars := convertToSimSimVar(reflect.ValueOf(Report{}))
	assert.Greater(t, len(vars), 0, "vars is zero")
}

func TestConvertToInterface(t *testing.T) {
	vars := []simconnect.SimVar{
		simconnect.SimVarSimOnGround(),
	}
	v := Report{}
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
