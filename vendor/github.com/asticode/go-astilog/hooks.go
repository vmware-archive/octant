package astilog

import "github.com/sirupsen/logrus"

type withFieldHook struct {
	k, v string
}

func newWithFieldHook(k, v string) *withFieldHook {
	return &withFieldHook{
		k: k,
		v: v,
	}
}

func (h *withFieldHook) Fire(e *logrus.Entry) error {
	if len(h.v) > 0 {
		e.Data[h.k] = h.v
	}
	return nil
}

func (h *withFieldHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
