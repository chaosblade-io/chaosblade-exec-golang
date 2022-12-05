package action

import (
	"context"
	"runtime"
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/log"
	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/model"
	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/model/response"
)

var MemRate = model.ExperimentFlag{
	Name:     "rate",
	Desc:     "burn memory rate, unit is M/s, default value is 100",
	NoArgs:   false,
	Required: false,
}

// 128K，为了扩容时申请粒度问题
type Block [32 * 1024]int32

var oomInstance *oomAction

func init() {
	oomInstance = &oomAction{}
	RegisterNeededClosedActions(oomInstance)
}

type oomAction struct {
}

func NewOomAction() *oomAction {
	return oomInstance
}

func (o *oomAction) Name() string {
	return "oom"
}

func (o *oomAction) Flags() map[string]model.ExperimentFlag {
	return map[string]model.ExperimentFlag{
		MemRate.Name: MemRate,
	}
}

var cache map[int][]Block
var quit = make(chan struct{})
var lock sync.Mutex

func (o *oomAction) Execute(ctx context.Context, rule *model.Experiment) response.Response {
	lock.Lock()
	defer lock.Unlock()
	if cache != nil {
		log.Zap.Info("oom experiment has started")
		return response.ReturnOK("oom experiment has started")
	}
	rate := 100
	rateValue := rule.Flags[MemRate.Name]
	if rateValue != "" {
		if r, err := strconv.Atoi(rateValue); err != nil {
			return response.ReturnIllegalParameters(err.Error())
		} else if r > 0 {
			rate = r
		}
	}
	cache = make(map[int][]Block, 0)
	go startOom(rate)
	return response.ReturnOK("start oom")
}

func startOom(rate int) {
	tick := time.Tick(time.Second)
	var count = 1
	cache[count] = make([]Block, 0)
	for range tick {
		select {
		case <-quit:
			cache = nil
			runtime.GC()
			return
		default:
			blockSize := rate * 8
			count += 1
			cache[count] = make([]Block, blockSize)
			log.Zap.Debug("malloc memory",
				zap.Int("count", count),
				zap.Int("blockSize", blockSize))
		}
	}
}

func (o *oomAction) Close(ctx context.Context, rule *model.ExperimentRule) error {
	go func() {
		quit <- struct{}{}
	}()
	return nil
}
