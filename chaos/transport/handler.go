package transport

import (
	"context"
	"sync"

	"go.uber.org/zap"

	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/log"
	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/manager"
	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/model"
	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/model/response"
	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/scenario"
)

func init() {
	serviceHandler := &ServiceHandler{lock: sync.RWMutex{}}
	RegisterHandler(newInjectHandler(serviceHandler))
	RegisterHandler(newRecoverHandler(serviceHandler))
	RegisterHandler(newCheckHandler(serviceHandler))
	RegisterHandler(newMetricHandler(serviceHandler))
	RegisterHandler(newCircuitHandler(serviceHandler))
	RegisterHandler(newCancelCircuitHandler(serviceHandler))

	RegisterHandler(newEnableScenarioHandler(serviceHandler))
	RegisterHandler(newDisableScenarioHandler(serviceHandler))
	RegisterHandler(newListScenarioHandler(serviceHandler))

	RegisterHandler(&logHandler{})
}

type ServiceHandler struct {
	lock sync.RWMutex
}

func (h *ServiceHandler) Inject(ctx context.Context, experimentRule model.ExperimentRule) response.Response {
	h.lock.Lock()
	defer h.lock.Unlock()
	if manager.CircuitEnabled() {
		return response.ReturnExperimentInCircuit()
	}
	err := manager.PutExperiment(&experimentRule)
	if err != nil {
		return response.ReturnFail(response.IllegalParameters, err.Error())
	}
	return response.ReturnOK("success")
}

func (h *ServiceHandler) Recover(ctx context.Context, experiment model.ExperimentRule) response.Response {
	h.lock.Lock()
	defer h.lock.Unlock()
	manager.RemoveExperiment(&experiment)
	return response.ReturnOK("success")
}

// 是否生效，如何验证，只有真正触发时才能验证
func (h *ServiceHandler) Check(ctx context.Context, experiment model.ExperimentRule) response.Response {
	h.lock.RLock()
	defer h.lock.RUnlock()
	panic("implement me")
}

func (h *ServiceHandler) Metric(ctx context.Context, rule model.ExperimentRule) response.Response {
	h.lock.RLock()
	defer h.lock.RUnlock()
	metric := manager.GetExperimentMetric(&rule)
	return response.ReturnOK(metric)
}

func (h *ServiceHandler) Circuit(ctx context.Context) response.Response {
	h.lock.Lock()
	defer h.lock.Unlock()
	manager.OpenCircuit()
	for key, experiment := range manager.Experiments() {
		err := manager.RemoveExperiment(&experiment.ExperimentRule)
		if err == nil {
			delete(manager.Experiments(), key)
			continue
		}
		log.Zap.Error("recover experiment failed under circuit",
			zap.String("target", experiment.Target),
			zap.String("action", experiment.Action),
			zap.String("index", experiment.Index),
			zap.Any("flags", experiment.Flags),
			zap.String("error", err.Error()))
	}
	return response.ReturnOK("success")
}

func (h *ServiceHandler) CancelCircuit(ctx context.Context) response.Response {
	h.lock.Lock()
	defer h.lock.Unlock()
	manager.CancelCircuit()
	return response.ReturnOK("success")
}

func (h *ServiceHandler) SetLogLevel(logLevel string) response.Response {
	err := log.SetLogLevel(logLevel)
	if err != nil {
		return response.ReturnIllegalParameters(err.Error())
	}
	return response.ReturnOK("success")
}

func (h *ServiceHandler) EnableScenario(ctx context.Context) response.Response {
	scenario.Enable()
	return response.ReturnOK("success")
}

func (h *ServiceHandler) DisableScenario(ctx context.Context) response.Response {
	scenario.Disable()
	return response.ReturnOK("success")
}

func (h *ServiceHandler) ListScenarios(ctx context.Context) response.Response {
	return response.ReturnOK(scenario.ListScenarios())
}
