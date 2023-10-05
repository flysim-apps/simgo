package simgo

import (
	"math"
)

func EventNotification(event string) map[string]interface{} {
	return map[string]interface{}{
		"event": event,
	}
}

func systemNotification(r *System) map[string]interface{} {
	return map[string]interface{}{
		"event": SYSTEM_STATE,
		"data": map[string]interface{}{
			"zuluTime": r.ZuluTime,
		},
	}
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
