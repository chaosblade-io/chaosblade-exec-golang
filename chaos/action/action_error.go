package action

import (
	"context"
	"fmt"

	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/model"
	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/model/response"
)

var message = model.ExperimentFlag{
	Name:     "message",
	Desc:     "error message",
	NoArgs:   false,
	Required: true,
}

type errorAction struct {
}

func NewErrorAction() *errorAction {
	return &errorAction{}
}

func (e *errorAction) Name() string {
	return "error"
}

func (e *errorAction) Flags() map[string]model.ExperimentFlag {
	return map[string]model.ExperimentFlag{message.Name: message}
}

func (e *errorAction) Execute(ctx context.Context, rule *model.Experiment) response.Response {
	message := rule.Flags[message.Name]
	return response.ReturnOK(fmt.Errorf(message))
}
