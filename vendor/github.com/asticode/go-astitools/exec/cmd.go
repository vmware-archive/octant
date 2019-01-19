package astiexec

import (
	"context"
	"os/exec"
	"strings"
	"time"

	"github.com/asticode/go-astilog"
)

// Cmd represents a command
type Cmd struct {
	Args []string
	ctx  context.Context
}

// NewCmd creates a new command
func NewCmd(ctx context.Context, args ...string) (cmd *Cmd) {
	cmd = &Cmd{
		Args: args,
		ctx:  ctx,
	}
	return
}

// String allows Cmd to implements the stringify interface
func (c *Cmd) String() string {
	return strings.Join(c.Args, " ")
}

// Exec executes a command
var Exec = func(cmd *Cmd) (o []byte, d time.Duration, err error) {
	// Init
	defer func(t time.Time) {
		d = time.Since(t)
	}(time.Now())

	// Create exec command
	execCmd := exec.CommandContext(cmd.ctx, cmd.Args[0], cmd.Args[1:]...)

	// Execute command
	astilog.Debugf("Executing %s", cmd)
	o, err = execCmd.CombinedOutput()
	return
}
