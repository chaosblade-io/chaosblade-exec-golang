package model

import (
	"sync"
)

type Empty struct{}

// ExperimentRule is an experiment with flags
type ExperimentRule struct {
	// ExperimentRule target
	Target string `json:"target"`
	// Action name
	Action string `json:"action"`
	// Flags contains all matcher rules
	Flags map[string]string `json:"flags"`
	// 由于一个 target 可以设置多个相同的场景，所以需要设置索引，是否单独设置还是放到 flags
	Index string `json:"index"`
}

// ExperimentFlag
type ExperimentFlag struct {
	// Name returns the flag FlagName
	Name string `yaml:"name"`
	// Desc returns the flag description
	Desc string `yaml:"desc"`
	// NoArgs means no arguments
	NoArgs bool `yaml:"noArgs"`
	// Required means necessary or not
	Required bool `yaml:"required"`
}

// ExperimentMetric
type ExperimentMetric struct {
	Count int64 `json:"count"`
}

// Experiment
type Experiment struct {
	Identifier string `json:"identifier"`
	ExperimentRule
	ExperimentMetric

	Lock sync.RWMutex
}

func (e *Experiment) Inc() {
	e.Count++
}

func (e *Experiment) Dec() {
	e.Count--
}

func (e *Experiment) Get() int64 {
	return e.Count
}

type ExperimentModel struct {
	Target   string           `json:"target"`
	Action   string           `json:"action"`
	Matchers []ExperimentFlag `json:"matchers"`
	Flags    []ExperimentFlag `json:"flags"`
}
