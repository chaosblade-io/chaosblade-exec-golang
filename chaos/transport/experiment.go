package transport

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/action"
	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/api"
	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/log"
	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/model"
	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/model/response"
)

const (
	InjectUrl        = "/chaos/inject"
	RecoverUrl       = "/chaos/recover"
	CheckUrl         = "/chaos/check"
	MetricUrl        = "/chaos/metric"
	CircuitUrl       = "/chaos/circuit"
	CancelCircuitUrl = "/chaos/cancelcircuit"
	LogUrl           = "/chaos/log"
)

type ExperimentHandler struct {
	serviceFunc func(ctx context.Context, experimentRule model.ExperimentRule) response.Response
}

func (e *ExperimentHandler) Execute(ctx context.Context, request *http.Request) response.Response {
	experimentRule, err := convertRequestToExperimentRule(request)
	if err != nil {
		return response.ReturnIllegalParameters(err.Error())
	}
	return e.serviceFunc(ctx, *experimentRule)
}

func convertRequestToExperimentRule(request *http.Request) (*model.ExperimentRule, error) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}
	rule := &model.ExperimentRule{}
	err = json.Unmarshal(body, rule)
	if err != nil {
		return nil, err
	}
	indexValue := rule.Flags[action.Index.Name]
	if indexValue == "" {
		rule.Index = "0"
	} else {
		rule.Index = indexValue
	}
	return rule, nil
}

type injectHandler struct {
	ExperimentHandler
}

func (i *injectHandler) Name() string {
	return InjectUrl
}
func newInjectHandler(serviceApi api.ChaosServiceApi) Handler {
	return &injectHandler{ExperimentHandler{serviceFunc: serviceApi.Inject}}
}

type recoverHandler struct {
	ExperimentHandler
}

func (r recoverHandler) Name() string {
	return RecoverUrl
}
func newRecoverHandler(serviceApi api.ChaosServiceApi) Handler {
	return &recoverHandler{ExperimentHandler{serviceFunc: serviceApi.Recover}}
}

type checkHandler struct {
	ExperimentHandler
}

func (c *checkHandler) Name() string {
	return CheckUrl
}
func newCheckHandler(serviceApi api.ChaosServiceApi) Handler {
	return &checkHandler{ExperimentHandler{serviceFunc: serviceApi.Check}}
}

type metricHandler struct {
	ExperimentHandler
}

func (m *metricHandler) Name() string {
	return MetricUrl
}
func newMetricHandler(serviceApi api.ChaosServiceApi) Handler {
	return &metricHandler{ExperimentHandler{serviceFunc: serviceApi.Metric}}
}

type circuitHandler struct {
	circuitFunc func(ctx context.Context) response.Response
}

func (c *circuitHandler) Name() string {
	return CircuitUrl
}

func (c *circuitHandler) Execute(ctx context.Context, request *http.Request) response.Response {
	return c.circuitFunc(ctx)
}
func newCircuitHandler(serviceApi api.ChaosServiceApi) Handler {
	return &circuitHandler{circuitFunc: serviceApi.Circuit}
}

type cancelCircuitHandler struct {
	circuitFunc func(ctx context.Context) response.Response
}

func (c *cancelCircuitHandler) Name() string {
	return CancelCircuitUrl
}

func (c *cancelCircuitHandler) Execute(ctx context.Context, request *http.Request) response.Response {
	return c.circuitFunc(ctx)
}
func newCancelCircuitHandler(serviceApi api.ChaosServiceApi) Handler {
	return &cancelCircuitHandler{circuitFunc: serviceApi.CancelCircuit}
}

type logHandler struct {
}

func (l *logHandler) Name() string {
	return LogUrl
}

func (l *logHandler) Execute(ctx context.Context, request *http.Request) response.Response {
	level := request.FormValue("level")
	if level == "" {
		return response.ReturnIllegalParameters("less level parameter")
	}
	if err := log.SetLogLevel(level); err != nil {
		return response.ReturnIllegalParameters(err.Error())
	}
	return response.ReturnOK("success")
}
