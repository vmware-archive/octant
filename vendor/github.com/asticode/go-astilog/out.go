// +build windows

package astilog

import (
	"io"
	"log"

	colorable "github.com/mattn/go-colorable"
)

func stdOut() io.Writer {
	return colorable.NewColorableStdout()
}

func syslogOut(c Configuration) io.Writer {
	log.Println("astilog: syslog is not implemented on this os, using stdout instead")
	return stdOut()
}
