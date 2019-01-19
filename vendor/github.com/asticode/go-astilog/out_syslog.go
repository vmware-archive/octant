// +build !windows

package astilog

import (
	"io"
	"log"
	"log/syslog"
	"os"

	"github.com/pkg/errors"
)


func stdOut() io.Writer {
	return os.Stdout
}

func syslogOut(c Configuration) (w io.Writer) {
	var err error
	if w, err = syslog.New(syslog.LOG_INFO|syslog.LOG_USER, c.AppName); err != nil {
		log.Println(errors.Wrap(err, "astilog: new syslog failed"))
		return os.Stdout
	}
	return
}
