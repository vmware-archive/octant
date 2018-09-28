package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newVersionCmd(GitCommit string, BuildTime string) *cobra.Command {

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Long:  "Version for hcli binary",
		Run: func(cmd *cobra.Command, args []string) {
			out := cmd.OutOrStdout()

			fmt.Fprintln(out, "Version: ", "pre-alpha")
			fmt.Fprintln(out, "Git commit: ", GitCommit)
			fmt.Fprintln(out, "Built: ", BuildTime)
		},
	}
	return versionCmd
}
