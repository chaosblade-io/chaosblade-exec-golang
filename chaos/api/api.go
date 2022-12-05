package api

import (
	"context"

	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/action"
	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/model"
	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/model/response"
)

type ChaosClientApi interface {
	// 设置
	Set(ctx context.Context, action action.Action) response.Response
	SetWithIndex(ctx context.Context, actionIndex int, act action.Action) response.Response
}

type ChaosServiceApi interface {
	Inject(ctx context.Context, model model.ExperimentRule) response.Response
	Recover(ctx context.Context, model model.ExperimentRule) response.Response
	Check(ctx context.Context, model model.ExperimentRule) response.Response
	Metric(ctx context.Context, model model.ExperimentRule) response.Response
	Circuit(ctx context.Context) response.Response
	CancelCircuit(ctx context.Context) response.Response
}

type ChaosScenarioApi interface {
	EnableScenario(ctx context.Context) response.Response
	DisableScenario(ctx context.Context) response.Response
	ListScenarios(ctx context.Context) response.Response
}
