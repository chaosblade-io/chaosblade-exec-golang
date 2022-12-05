package action

import (
	"context"
	"encoding/json"
	"reflect"
	"strconv"

	"go.uber.org/zap"

	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/log"
	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/model"
	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/model/response"
)

var value = model.ExperimentFlag{
	Name:     "value",
	Desc:     "json string",
	NoArgs:   false,
	Required: true,
}

// modifyAction contains modify parameters and returned obj.
type modifyAction struct {
	// 原有对象
	OriginalValue interface{}
}

func NewModifyAction(value interface{}) *modifyAction {
	return &modifyAction{
		OriginalValue: value,
	}
}

func (m *modifyAction) Name() string {
	return "modify"
}

func (m *modifyAction) Flags() map[string]model.ExperimentFlag {
	return map[string]model.ExperimentFlag{
		value.Name: value,
	}
}

func (m *modifyAction) Execute(ctx context.Context, rule *model.Experiment) response.Response {
	value := rule.Flags[value.Name]
	if value == "" {
		return response.ReturnIllegalParameters("the value flag is missing")
	}
	log.Zap.Debug("modify value",
		zap.String("target", rule.Target),
		zap.String("action", rule.Action),
		zap.String("index", rule.Index),
		zap.Any("originalValue", m.OriginalValue),
		zap.Any("newValue", value),
	)
	v := reflect.ValueOf(m.OriginalValue)
	kind := reflect.TypeOf(m.OriginalValue).Kind()
	switch kind {
	case reflect.Ptr:
		if err := json.Unmarshal([]byte(value), m.OriginalValue); err != nil {
			return response.ReturnIllegalParameters(err.Error())
		}
		return response.ReturnOK(m.OriginalValue)
	case reflect.String:
		return response.ReturnOK(value)
	case reflect.Bool:
		parseBool, err := strconv.ParseBool(value)
		if err != nil {
			return response.ReturnIllegalParameters(err.Error())
		}
		return response.ReturnOK(parseBool)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return response.ReturnIllegalParameters(err.Error())
		}
		if v.OverflowInt(n) {
			return response.ReturnIllegalParameters("overflow type scope")
		}
		return response.ReturnOK(reflect.ValueOf(n).Convert(v.Type()).Interface())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		n, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return response.ReturnIllegalParameters(err.Error())
		}
		if v.OverflowUint(n) {
			return response.ReturnIllegalParameters("overflow type scope")
		}
		return response.ReturnOK(reflect.ValueOf(n).Convert(v.Type()).Interface())
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(value, v.Type().Bits())
		if err != nil {
			return response.ReturnIllegalParameters(err.Error())
		}
		if v.OverflowFloat(n) {
			return response.ReturnIllegalParameters("overflow type scope")
		}
		return response.ReturnOK(reflect.ValueOf(n).Convert(v.Type()).Interface())
	case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array:
		newObj := reflect.New(reflect.TypeOf(m.OriginalValue))
		if err := json.Unmarshal([]byte(value), newObj.Interface()); err != nil {
			return response.ReturnIllegalParameters(err.Error())
		}
		return response.ReturnOK(newObj.Elem().Interface())
	default:
		log.Zap.Error("not support the type",
			zap.String("target", rule.Target),
			zap.String("action", rule.Action),
			zap.String("index", rule.Index),
			zap.Any("type", kind),
			zap.Any("originalValue", m.OriginalValue),
			zap.Any("newValue", value),
		)
	}
	return response.ReturnOK(m.OriginalValue)
}
