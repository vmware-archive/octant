package commands

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// Default variables overridden by main
var (
	GitCommit = "(unknown-commit)"
	BuildTime = "(unknown-buildtime)"
)

func newVersionCmd() *cobra.Command {

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Long:  "Version for hcli binary",
		Run: func(cmd *cobra.Command, args []string) {
			out := cmd.OutOrStdout()

			fmt.Fprintln(out, "Version: ", "pre-alpha")
			fmt.Fprintln(out, "Git commit: ", GitCommit)
			fmt.Fprintln(out, "Built: ", reformatDate(BuildTime))
		},
	}
	return versionCmd
}

func reformatDate(dateTime string) string {
	t, errTime := time.Parse(time.RFC3339Nano, dateTime)
	if errTime == nil {
		return t.Format(time.ANSIC)
	}
	return dateTime
}
