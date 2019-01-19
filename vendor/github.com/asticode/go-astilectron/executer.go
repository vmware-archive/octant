package astilectron

import (
	"os/exec"
	"strings"

	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

// Executer represents an object capable of executing Astilectron run command
type Executer func(a *Astilectron, cmd *exec.Cmd) (err error)

// DefaultExecuter represents the default executer
func DefaultExecuter(a *Astilectron, cmd *exec.Cmd) (err error) {
	// Start command
	astilog.Debugf("Starting cmd %s", strings.Join(cmd.Args, " "))
	if err = cmd.Start(); err != nil {
		err = errors.Wrapf(err, "starting cmd %s failed", strings.Join(cmd.Args, " "))
		return
	}

	// Watch command
	go a.watchCmd(cmd)
	return
}
