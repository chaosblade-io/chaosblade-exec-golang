package matcher

type Matcher interface {
	Name() string
	Match(expectedValue string) (bool, error)
}

type LocalMatcher struct {
	name        string
	actualValue interface{}
}

func (l *LocalMatcher) Name() string {
	return l.name
}

func (l *LocalMatcher) ActualValue() interface{} {
	return l.actualValue
}
