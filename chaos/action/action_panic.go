package action

import (
	"context"

	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/model"
	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/model/response"
)

var PanicMessage = model.ExperimentFlag{
	Name:     "message",
	Desc:     "panic message",
	NoArgs:   false,
	Required: false,
}

type panicAction struct {
}

func NewPanicAction() *panicAction {
	return &panicAction{}
}

func (p *panicAction) Flags() map[string]model.ExperimentFlag {
	return map[string]model.ExperimentFlag{PanicMessage.Name: PanicMessage}
}

func (p *panicAction) Name() string {
	return "panic"
}

func (p *panicAction) Execute(ctx context.Context, rule *model.Experiment) response.Response {
	panic(rule.Flags["message"])
}
