package astios

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/asticode/go-astilog"
)

type SignalsFunc func(s os.Signal)

func HandleSignals(fn SignalsFunc) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch)
	for s := range ch {
		astilog.Debugf("astios: received signal %s", s)
		fn(s)
	}
}

func ContextSignalsFunc(c context.CancelFunc) SignalsFunc {
	return func(s os.Signal) {
		if s == syscall.SIGABRT || s == syscall.SIGINT || s == syscall.SIGQUIT || s == syscall.SIGTERM {
			c()
		}
	}
}
