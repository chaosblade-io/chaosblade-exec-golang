package scenario

import (
	"sync"

	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/action"
	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/matcher"
	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/model"
)

var sce = &scenario{
	enable:    false,
	scenarios: make(map[string]*model.ExperimentModel, 0),
	lock:      sync.RWMutex{},
}

type scenario struct {
	enable    bool
	scenarios map[string]*model.ExperimentModel
	lock      sync.RWMutex
}

func Enabled() bool {
	return sce.enable
}

func Enable() {
	sce.enable = true
}

func Disable() {
	sce.enable = false
}

func Register(target string, matchers map[string]matcher.Matcher, act action.Action) {
	if !Enabled() {
		return
	}
	sce.lock.Lock()
	defer sce.lock.Unlock()
	key := target + "-" + act.Name()
	if _, ok := sce.scenarios[key]; ok {
		return
	}
	matcherFlags := make([]model.ExperimentFlag, 0)
	for key := range matchers {
		matcherFlags = append(matcherFlags, model.ExperimentFlag{
			Name:     key,
			Desc:     "custom define key",
			NoArgs:   false,
			Required: false,
		})
	}
	flags := make([]model.ExperimentFlag, 0)
	for _, value := range act.Flags() {
		flags = append(flags, value)
	}
	for _, value := range action.CommonActionFlagMap {
		flags = append(flags, value)
	}
	sce.scenarios[key] = &model.ExperimentModel{
		Target:   target,
		Action:   act.Name(),
		Matchers: matcherFlags,
		Flags:    flags,
	}
}

func ListScenarios() []model.ExperimentModel {
	sce.lock.RLock()
	defer sce.lock.RUnlock()
	scenarios := make([]model.ExperimentModel, 0)
	for _, scenario := range sce.scenarios {
		scenarios = append(scenarios, *scenario)
	}
	return scenarios
}
