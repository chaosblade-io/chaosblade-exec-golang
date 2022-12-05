package matcher

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/log"
)

type EqualMatcher struct {
	LocalMatcher
}

func NewEqualMatcher(name string, actualValue interface{}) Matcher {
	return &EqualMatcher{
		LocalMatcher{
			name:        name,
			actualValue: actualValue,
		}}
}

func (e *EqualMatcher) Match(expectedValue string) (bool, error) {
	if e.actualValue == nil {
		return false, fmt.Errorf("%s value is nil", e.name)
	}
	log.Zap.Debug("contains matcher",
		zap.String("expectedValue", expectedValue),
		zap.String("actualValue", e.actualValue.(string)))
	return e.actualValue.(string) == expectedValue, nil
}
