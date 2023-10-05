package simgo

import (
	"reflect"
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
	convertToInterface(reflect.ValueOf(&v), vars)
}

func TestConvertWithReflect(t *testing.T) {
	some := func(v interface{}) {
		vars := convertToSimSimVar(reflect.ValueOf(v))
		assert.Greater(t, len(vars), 0, "vars is zero")

		convertToInterface(reflect.ValueOf(v), vars)
	}

	some(&Report{})
}
