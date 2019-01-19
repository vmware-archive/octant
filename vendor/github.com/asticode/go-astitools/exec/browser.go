package astiexec

import (
	"context"
	"os/exec"
	"runtime"

	"github.com/pkg/errors"
)

// Cheers to https://gist.github.com/threeaccents/607f3bc3a57a2ddd9d57
func OpenBrowser(ctx context.Context, url string) error {
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.CommandContext(ctx, args[0], append(args[1:], url)...)
	if b, err := cmd.CombinedOutput(); err != nil {
		return errors.Wrapf(err, "astiexec: opening browser failed with body %s", b)
	}
	return nil
}
