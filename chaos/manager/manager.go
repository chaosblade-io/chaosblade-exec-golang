package manager

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"

	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/action"
	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/log"
	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/model"
)

var circuitEnabled bool

var mgr = manager{
	experiments: make(map[string]*model.Experiment, 0),
}

type manager struct {
	experiments map[string]*model.Experiment
}

func Experiments() map[string]*model.Experiment {
	return mgr.experiments
}

func PutExperiment(experimentRule *model.ExperimentRule) error {
	if experimentRule.Target == "" || experimentRule.Action == "" {
		return fmt.Errorf("less target or action parameter")
	}
	experimentKey := GetExperimentKey(experimentRule.Target, experimentRule.Action, experimentRule.Index)
	if _, ok := mgr.experiments[experimentKey]; ok {
		return fmt.Errorf("%s target and %s action at %s position exits",
			experimentRule.Target, experimentRule.Action, experimentRule.Index)
	}
	mgr.experiments[experimentKey] = &model.Experiment{
		Identifier:       experimentKey,
		ExperimentRule:   *experimentRule,
		ExperimentMetric: model.ExperimentMetric{},
		Lock:             sync.RWMutex{},
	}
	log.Zap.Info("inject chaos experimentRule",
		zap.String("target", experimentRule.Target),
		zap.String("action", experimentRule.Action),
		zap.String("index", experimentRule.Index),
		zap.String("flags", fmt.Sprintf("%+v", experimentRule.Flags)),
	)
	return nil
}

func RemoveExperiment(experiment *model.ExperimentRule) error {
	actionIndex := experiment.Flags[action.Index.Name]
	experimentKey := GetExperimentKey(experiment.Target, experiment.Action, actionIndex)
	if _, ok := mgr.experiments[experimentKey]; ok {
		if action, ok := action.GetNeededClosedActions()[experiment.Action]; ok {
			if err := action.Close(context.Background(), experiment); err != nil {
				return err
			}
		}
		delete(mgr.experiments, experimentKey)
		log.Zap.Info("remove chaos experiment",
			zap.String("target", experiment.Target),
			zap.String("action", experiment.Action),
			zap.String("index", actionIndex),
			zap.String("flags", fmt.Sprintf("%+v", experiment.Flags)),
		)
	}
	return nil
}

// key=target+actionName
func GetExperiment(target, actionName, actionIndex string) *model.Experiment {
	experiment, ok := mgr.experiments[GetExperimentKey(target, actionName, actionIndex)]
	if ok {
		return experiment
	}
	return nil
}

func GetExperimentKey(target, actionName, actionIndex string) string {
	if actionIndex == "" {
		actionIndex = "0"
	}
	return fmt.Sprintf("%s-%s-%s", target, actionName, actionIndex)
}

func CircuitEnabled() bool {
	return circuitEnabled
}

func OpenCircuit() {
	circuitEnabled = true
	log.Zap.Info("open chaos circuit")
}

func CancelCircuit() {
	circuitEnabled = false
	log.Zap.Info("cancel chaos circuit")
}

func GetExperimentMetric(rule *model.ExperimentRule) model.ExperimentMetric {
	experiment := GetExperiment(rule.Target, rule.Action, rule.Index)
	if experiment == nil {
		return model.ExperimentMetric{}
	}
	return experiment.ExperimentMetric
}
