package chaos

import (
	"context"
	"flag"
	"math/rand"
	"runtime"
	"strconv"
	"time"

	"github.com/spf13/pflag"
	"go.uber.org/zap"

	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/action"
	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/log"
	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/manager"
	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/matcher"
	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/model"
	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/model/response"
	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/scenario"
	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/transport"
)

var (
	chaosFlagSet   *pflag.FlagSet
	ip             string
	port           int
	scenarioEnable bool
)

func init() {
	chaosFlagSet = pflag.NewFlagSet("chaos", pflag.ContinueOnError)
	flag.StringVar(&ip, "chaos-ip", "127.0.0.1", "chaos service host ip")
	flag.IntVar(&port, "chaos-port", 9526, "chaos service port")
	flag.BoolVar(&scenarioEnable, "chaos-scenario-enable", false, "chaos scenario enable or not")
}

// Must be invoked
func Init() {
	transport.Run(ip, strconv.Itoa(port))
}

type Chaos struct {
	name     string
	matchers map[string]matcher.Matcher
}

func New() *Chaos {
	name := getTargetName()
	log.Zap.Debug("new chaos name: " + name)
	return &Chaos{
		name:     name,
		matchers: make(map[string]matcher.Matcher, 0),
	}
}

func (c *Chaos) Name() string {
	return c.name
}

func (c *Chaos) Matchers() map[string]matcher.Matcher {
	return c.matchers
}

func (c *Chaos) AddMatchers(matchers ...matcher.Matcher) {
	for _, m := range matchers {
		if _, ok := c.matchers[m.Name()]; ok {
			log.Zap.Warn("the matcher exists, so skip it", zap.String("name", m.Name()))
			continue
		}
		c.matchers[m.Name()] = m
	}
}

func (c *Chaos) Set(ctx context.Context, action action.Action) response.Response {
	return c.SetWithIndex(ctx, 0, action)
}

// 是否
func (c *Chaos) SetWithIndex(ctx context.Context, actionIndex int, act action.Action) response.Response {
	if manager.CircuitEnabled() {
		log.Zap.Debug("in circuit, skip to execute chaos")
		return response.ReturnExperimentInCircuit()
	}
	// for experiments exporting
	scenario.Register(c.Name(), c.matchers, act)
	experiment := manager.GetExperiment(c.Name(), act.Name(), strconv.Itoa(actionIndex))
	if experiment == nil {
		return response.ReturnExperimentNotFound()
	}
	for key, value := range experiment.Flags {
		if _, ok := action.CommonActionFlagMap[key]; ok {
			continue
		}
		if _, ok := act.Flags()[key]; ok {
			continue
		}
		matcher, ok := c.Matchers()[key]
		if !ok {
			log.Zap.Debug("can not find the experiment matcher",
				zap.String("target", experiment.Target),
				zap.String("action", experiment.Action),
				zap.Int("index", actionIndex),
				zap.String("matcher", key))
			return response.ReturnExperimentMatcherNotFound(key)
		}
		if matched, err := matcher.Match(value); err != nil || !matched {
			log.Zap.Debug("the matcher does not match",
				zap.String("target", experiment.Target),
				zap.String("action", experiment.Action),
				zap.Int("index", actionIndex),
				zap.String("matcher", key),
				zap.String("expectedValue", value),
			)
			return response.ReturnExperimentNotMatched(key)
		}
	}
	less, err := c.lessEffectedLimit(experiment)
	if err != nil {
		return response.ReturnFail(response.IllegalParameters, err.Error())
	}
	if !less {
		return response.ReturnExperimentLimited()
	}
	experiment.Lock.Lock()
	defer experiment.Lock.Unlock()
	experiment.Inc()
	log.Zap.Info("matched experiment",
		zap.String("target", experiment.Target),
		zap.String("action", experiment.Action),
		zap.Int("index", actionIndex),
		zap.Any("flags", experiment.Flags),
	)
	response := act.Execute(ctx, experiment)
	if !response.Success {
		log.Zap.Warn("execute experiment failed",
			zap.String("target", experiment.Target),
			zap.String("action", experiment.Action),
			zap.Int("index", actionIndex),
			zap.Any("flags", experiment.Flags),
			zap.String("error", response.Error))
		experiment.Dec()
		return response
	}
	return response
}

// lessEffectedLimit
func (c *Chaos) lessEffectedLimit(experiment *model.Experiment) (bool, error) {
	effectCountValue := experiment.Flags[action.EffectCount.Name]
	effectCountPercent := experiment.Flags[action.EffectPercent.Name]
	if effectCountValue == "" && effectCountPercent == "" {
		return true, nil
	}
	if effectCountValue != "" {
		effectCount, err := strconv.Atoi(effectCountValue)
		if err != nil {
			return false, err
		}
		if experiment.Count >= int64(effectCount) {
			return false, nil
		}
	}
	if effectCountPercent != "" {
		effectPercent, err := strconv.Atoi(effectCountPercent)
		if err != nil {
			return false, err
		}
		if effectPercent != 100 {
			rand.Seed(time.Now().UnixNano())
			randValue := rand.Intn(100) + 1
			if randValue > effectPercent {
				return false, nil
			}
		}
	}
	return true, nil
}

// examples.(*Business).Execute
// github.com/chaosblade-io/chaosblade-exec-golang/demo.Log
// examples.init
func getTargetName() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return frame.Function
}
