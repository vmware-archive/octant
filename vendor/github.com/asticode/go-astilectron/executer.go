package astilectron

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/asticode/go-astikit"
)

// Executer represents an object capable of executing Astilectron run command
type Executer func(l astikit.SeverityLogger, a *Astilectron, cmd *exec.Cmd) (err error)

// DefaultExecuter represents the default executer
func DefaultExecuter(l astikit.SeverityLogger, a *Astilectron, cmd *exec.Cmd) (err error) {
	// Start command
	l.Debugf("Starting cmd %s", strings.Join(cmd.Args, " "))
	if err = cmd.Start(); err != nil {
		err = fmt.Errorf("starting cmd %s failed: %w", strings.Join(cmd.Args, " "), err)
		return
	}

	// Watch command
	go a.watchCmd(cmd)
	return
}
