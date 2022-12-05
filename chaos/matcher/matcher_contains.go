package matcher

import (
	"fmt"
	"strings"

	"go.uber.org/zap"

	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/log"
)

type ContainsMatcher struct {
	LocalMatcher
}

func NewContainsMatcher(name string, actualValue interface{}) Matcher {
	return &ContainsMatcher{LocalMatcher{
		name:        name,
		actualValue: actualValue,
	}}
}

func (c *ContainsMatcher) Match(expectedValue string) (bool, error) {
	if c.actualValue == nil {
		return false, fmt.Errorf("%s value is nil", c.name)
	}
	if expectedValue == "" {
		return false, fmt.Errorf("%s expected value is empty", c.name)
	}
	log.Zap.Debug("contains matcher",
		zap.String("expectedValue", expectedValue),
		zap.String("actualValue", c.actualValue.(string)))
	return strings.Contains(c.actualValue.(string), expectedValue), nil
}
