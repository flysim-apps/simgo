package simgo

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"

	"nhooyr.io/websocket/wsjson"
)

func (s *SimGo) OffsetPolling(name string, report interface{}, interval int) error {
	cmd := FSUIPC_Command{
		Command: "offsets.declare",
		Name:    name,
		Offsets: s.FSUIPC_Remap(reflect.ValueOf(report)),
	}

	b, _ := json.Marshal(cmd)

	s.Logger.Debugf("%s", b)

	if err := wsjson.Write(s.Context, s.WS, cmd); err != nil {
		return err
	}

	cmd = FSUIPC_Command{
		Command:  "offsets.read",
		Name:     name,
		Interval: interval,
	}

	s.Logger.Debugf("CMD: %+v", cmd)

	if err := wsjson.Write(s.Context, s.WS, cmd); err != nil {
		return err
	}

	return nil
}

func (s *SimGo) Payload(interval int) error {
	cmd := FSUIPC_Command_Payload{
		Command:    "payload.read",
		WeightUnit: "Lbs",
		VolumeUnit: "gal",
		LengthUnit: "ft",
	}

	s.Logger.Debugf("CMD: %+v", cmd)

	if err := wsjson.Write(s.Context, s.WS, cmd); err != nil {
		return err
	}

	return nil
}

func (s *SimGo) ReadData(name string, v interface{}, eventChan chan interface{}, payloadChan chan interface{}) {
	done := make(chan bool)
	go func() {
		defer s.FSUIPC_Close()
		for {
			var msg FSUIPC_Response
			if err := wsjson.Read(s.Context, s.WS, &msg); err != nil {
				s.Logger.Errorf("Unable to read from socket: %s", err.Error())
				return
			}

			if msg.Success {
				switch msg.Command {
				case "offsets.declare":
					s.Logger.Debugf("Offsets %s have been declared", msg.Name)
					break
				case "offsets.read":
					if name == msg.Name {
						linkage := FSUIPC_Offset_Linkage{}
						linkage.LoadFromMap(msg.Data)
						eventChan <- s.FSUIPC_ToInterface(linkage, reflect.ValueOf(v))
					}
					break
				case "payload.read":
					payload := FSUIPC_Offset_Payload{}
					payload.LoadFromMap(msg.Data)
					payloadChan <- payload
					break
				}
			} else {
				// The request failed. Handle the errors here.
				// In this example we just display errors to the webpage
				var error = fmt.Sprintf("Error for %s (%s): ", msg.Name, msg.Command)
				error += fmt.Sprintf("%v - %s", msg.ErrorCode, msg.ErrorMessage)
				s.Logger.Error(error)
				if msg.ErrorCode == "NoFlightSim" {
					s.Error = errors.New(fmt.Sprintf("%v - %s", msg.ErrorCode, msg.ErrorMessage))
					s.TrackFailed <- true
					close(done)
				}
			}
		}
	}()
}

func (s *SimGo) FSUIPC_Remap(val reflect.Value) []FSUIPC_Offset {
	vars := make([]FSUIPC_Offset, 0)

	for i := 0; i < val.Type().NumField(); i++ {
		addressTag, _ := val.Type().Field(i).Tag.Lookup("address")
		typeTag, _ := val.Type().Field(i).Tag.Lookup("type")
		sizeTag, _ := val.Type().Field(i).Tag.Lookup("size")
		indexTag, _ := val.Type().Field(i).Tag.Lookup("index")

		if addressTag == "" || typeTag == "" || sizeTag == "" {
			continue
		}

		// skip bits sub index fields
		if typeTag == "bits" && indexTag != "" {
			continue
		}

		sizeInt, _ := strconv.Atoi(sizeTag)
		uin, _ := strconv.ParseUint(addressTag, 0, 32)
		addr32 := (uint32)(uintptr(uin))

		simv := FSUIPC_Offset{
			Name:    val.Type().Field(i).Name,
			Address: addr32,
			Type:    typeTag,
			Size:    sizeInt,
		}

		vars = append(vars, simv)
	}

	return vars
}

func (s *SimGo) FSUIPC_ToInterface(data FSUIPC_Offset_Linkage, dst reflect.Value) interface{} {
	val := reflect.ValueOf(data)
	r := reflect.New(reflect.TypeOf(dst.Interface())).Elem()
	for i := 0; i < val.Type().NumField(); i++ {
		for j := 0; j < dst.NumField(); j++ {
			dstField := dst.Type().Field(j)
			if val.Type().Field(i).Name == dstField.Name {
				typeTag, _ := val.Type().Field(i).Tag.Lookup("type")
				indexTag, _ := val.Type().Field(i).Tag.Lookup("index")
				if typeTag == "bits" && indexTag != "" {
					continue
				}

				if err := setValueForField(val.Type().Field(i), val.Field(i), r.Field(j)); err != nil {
					s.Logger.Warningf("Failed set value for %s", err.Error())
				}
			}
		}
	}
	return r.Interface()
}

func setValueForField(srcType reflect.StructField, src reflect.Value, dst reflect.Value) error {
	unitTag, _ := srcType.Tag.Lookup("unit")
	switch unitTag {
	case "knots":
		if src.CanInt() {
			dst.SetInt(src.Int() / 128)
		} else {
			return errors.New(fmt.Sprintf("[knots  ] %s = %s", srcType.Name, dst.String()))
		}
	case "mach":
		if src.CanFloat() {
			dst.SetFloat(src.Float() / 2048 / 10)
		} else {
			return errors.New(fmt.Sprintf("[mach   ] %s = %s", srcType.Name, dst.String()))
		}
	case "degrees":
		dst.SetFloat(src.Float() * 360 / (65536 * 65536))
	case "raddeg":
		dst.SetFloat(src.Float() * 180 / math.Pi)
	case "GForce":
		dst.SetFloat(src.Float() / 624)
	case "radio":
		dst.SetInt(int64(float64(src.Int()/65536) * 3.28084))
	case "floating":
		dst.SetInt(int64(float64(src.Int() * 90.0 / (10001750.0 * 65536.0 * 65536.0))))
	case "ftm":
		dst.SetInt(int64(float64(src.Int()*60.0) * 3.28084 / 256))
	case "velocity":
		dst.SetInt(int64(float64(src.Int()/65536) * 1.944))
	case "feet":
		if src.CanFloat() {
			dst.SetInt(int64(math.Round(src.Float() * 3.28084)))
		} else {
			return errors.New(fmt.Sprintf("%s = %s", srcType.Name, dst.String()))
		}
	case "bool":
		if src.CanInt() {
			dst.SetBool(src.Int() > 0)
		} else {
			return errors.New(fmt.Sprintf("%s = %s", srcType.Name, dst.String()))
		}
	case "position":
		if src.CanInt() {
			dst.SetInt(src.Int())
		} else {
			return errors.New(fmt.Sprintf("%s = %s", srcType.Name, dst.String()))
		}
	case "percent":
		total := float64(16384)
		if src.Kind().String() == "float64" {
			dst.SetFloat(float64(src.Float()) / total * 100)
		} else {
			fmt.Printf("%s (%s) = %s\n", srcType.Name, src.Kind(), dst.String())
			dst.SetFloat(float64(src.Int()) / total * 100)
		}
	case "bits":
	default:
		//fmt.Printf("%s (%s) = %s\n", srcType.Name, src.Kind(), dst.String())
		if src.CanFloat() {
			dst.SetFloat(src.Float())
		} else if src.CanInt() {
			//fmt.Printf("Kind: %s\n", dst.Kind().String())
			if strings.Contains(dst.Kind().String(), "float") {
				dst.SetFloat(float64(src.Int()))
			} else {
				dst.SetInt(src.Int())
			}
		} else if src.Kind().String() != "map" {
			dst.SetString(src.String())
		}
	}

	return nil
}

func (c *FSUIPC_Offset_Linkage) LoadFromMap(m map[string]interface{}) error {
	data, err := json.Marshal(m)
	if err == nil {
		err = json.Unmarshal(data, c)
	}
	return err
}

func (c *FSUIPC_Offset_Payload) LoadFromMap(m map[string]interface{}) error {
	data, err := json.Marshal(m)
	if err == nil {
		err = json.Unmarshal(data, c)
	}
	return err
}
