package transport

import (
	"context"
	"net/http"

	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/api"
	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/model/response"
)

type enableScenarioHandler struct {
	scenarioFunc func(context.Context) response.Response
}

func (e *enableScenarioHandler) Name() string {
	return "/chaos/scenario/enable"
}

func (e *enableScenarioHandler) Execute(ctx context.Context, request *http.Request) response.Response {
	return e.scenarioFunc(ctx)
}
func newEnableScenarioHandler(api api.ChaosScenarioApi) Handler {
	return &enableScenarioHandler{scenarioFunc: api.EnableScenario}
}

type disableScenarioHandler struct {
	scenarioFunc func(context.Context) response.Response
}

func (d *disableScenarioHandler) Name() string {
	return "/chaos/scenario/disable"
}

func (d *disableScenarioHandler) Execute(ctx context.Context, request *http.Request) response.Response {
	return d.scenarioFunc(ctx)
}
func newDisableScenarioHandler(api api.ChaosScenarioApi) Handler {
	return &disableScenarioHandler{scenarioFunc: api.DisableScenario}
}

type listScenarioHandler struct {
	scenarioFunc func(context.Context) response.Response
}

func (l *listScenarioHandler) Name() string {
	return "/chaos/scenario/list"
}

func (l *listScenarioHandler) Execute(ctx context.Context, request *http.Request) response.Response {
	return l.scenarioFunc(ctx)
}
func newListScenarioHandler(scenarioApi api.ChaosScenarioApi) Handler {
	return &listScenarioHandler{scenarioFunc: scenarioApi.ListScenarios}
}
