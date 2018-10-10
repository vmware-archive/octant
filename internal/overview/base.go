package overview

import "k8s.io/apimachinery/pkg/util/clock"

// Describer creates content.
type Describer interface {
	Describe(prefix, namespace string, cache Cache, fields map[string]string) ([]Content, error)
}

type baseDescriber struct{}

func newBaseDescriber() *baseDescriber {
	return &baseDescriber{}
}

func (d *baseDescriber) clock() clock.Clock {
	return &clock.RealClock{}
}
