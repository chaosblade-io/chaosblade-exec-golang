package action

import (
	"context"
	"math/rand"
	"time"

	"go.uber.org/zap"

	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/log"
	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/model"
	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/model/response"
)

var Time = model.ExperimentFlag{
	Name:     "time",
	Desc:     "delay time, unit is duration",
	NoArgs:   false,
	Required: true,
}

var Offset = model.ExperimentFlag{
	Name:     "offset",
	Desc:     "time offset",
	NoArgs:   false,
	Required: false,
}

type delayAction struct {
}

//NewDelayAction
func NewDelayAction() *delayAction {
	return &delayAction{
	}
}

func (d *delayAction) Name() string {
	return "delay"
}

func (d *delayAction) Flags() map[string]model.ExperimentFlag {
	return map[string]model.ExperimentFlag{
		Time.Name:   Time,
		Offset.Name: Offset,
	}
}

func (d *delayAction) Execute(ctx context.Context, rule *model.Experiment) response.Response {
	delayTimeDuration, err := time.ParseDuration(rule.Flags[Time.Name])
	if err != nil {
		return response.ReturnIllegalParameters(err.Error())
	}
	offsetDuration := time.Duration(0)
	offsetValue := rule.Flags[Offset.Name]
	if offsetValue != "" {
		offsetDuration, err = time.ParseDuration(offsetValue)
		if err != nil {
			return response.ReturnIllegalParameters(err.Error())
		}
	}
	offsetMillis := offsetDuration.Milliseconds()
	if offsetMillis == 0 {
		time.Sleep(delayTimeDuration)
		return response.ReturnOK("success")
	}
	timeMillis := delayTimeDuration.Milliseconds()
	if offsetMillis >= timeMillis {
		return response.ReturnIllegalParameters("the offset value must be less the time value")
	}
	rand.Seed(time.Now().UnixNano())
	offset := rand.Int63n(offsetMillis*2+1) - offsetMillis
	sleepTimeInMillis := timeMillis - offset
	duration := time.Duration(sleepTimeInMillis) * time.Millisecond
	log.Zap.Debug("start to delay",
		zap.String("target", rule.Target),
		zap.String("action", rule.Action),
		zap.String("index", rule.Index),
		zap.Duration("time", duration),
		zap.Duration("offset", offsetDuration),
		zap.Duration("delay", duration))
	time.Sleep(duration)
	return response.ReturnOK("success")
}
