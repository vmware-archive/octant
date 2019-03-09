package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newVersionCmd(version string, gitCommit string, buildTime string) *cobra.Command {

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Long:  "Version for sugarloaf binary",
		Run: func(cmd *cobra.Command, args []string) {
			out := cmd.OutOrStdout()

			fmt.Fprintln(out, "Version: ", version)
			fmt.Fprintln(out, "Git commit: ", gitCommit)
			fmt.Fprintln(out, "Built: ", buildTime)
		},
	}
	return versionCmd
}
