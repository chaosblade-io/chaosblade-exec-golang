package action

import (
	"context"
	"fmt"
	"sync"

	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/model"
	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/model/response"
)

var CommonActionFlagMap = map[string]model.ExperimentFlag{
	Index.Name:         Index,
	EffectCount.Name:   EffectCount,
	EffectPercent.Name: EffectPercent,
}

var Index = model.ExperimentFlag{
	Name:     "index",
	Desc:     "",
	NoArgs:   false,
	Required: false,
}
var EffectCount = model.ExperimentFlag{
	Name:     "effect-count",
	Desc:     "",
	NoArgs:   false,
	Required: false,
}
var EffectPercent = model.ExperimentFlag{
	Name:     "effect-percent",
	Desc:     "",
	NoArgs:   false,
	Required: false,
}

type Action interface {
	// Action name
	Name() string
	// Flags return all of the action flags contain matcher flags
	Flags() map[string]model.ExperimentFlag
	// Execute experiment rule
	Execute(ctx context.Context, rule *model.Experiment) response.Response
}

type CloseAction interface {
	Action
	Close(ctx context.Context, rule *model.ExperimentRule) error
}

var neededClosedActions = make(map[string]CloseAction, 0)
var registerActionLock sync.RWMutex

func RegisterNeededClosedActions(action CloseAction) error {
	registerActionLock.Lock()
	defer registerActionLock.Unlock()
	if _, ok := neededClosedActions[action.Name()]; ok {
		return fmt.Errorf("the %s action has been registed", action.Name())
	}
	neededClosedActions[action.Name()] = action
	return nil
}

func GetNeededClosedActions() map[string]CloseAction {
	registerActionLock.RLock()
	defer registerActionLock.RUnlock()
	return neededClosedActions
}

func GetAllActions() []Action {
	return []Action{
		NewDelayAction(), NewErrorAction(), NewModifyAction(nil), NewOomAction(), NewPanicAction(),
	}
}
