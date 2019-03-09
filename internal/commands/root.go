package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	// remove timestamp from log
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
}

// Execute executes sugarloaf.
func Execute(version string, gitCommit string, buildTime string) {
	rootCmd := newRoot(version, gitCommit, buildTime)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func newRoot(version string, gitCommit string, buildTime string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "sugarloaf",
		Short: "sugarloaf is a dashboard for high bandwidth cluster analysis operations",
	}

	rootCmd.AddCommand(newDashCmd())
	rootCmd.AddCommand(newVersionCmd(version, gitCommit, buildTime))

	return rootCmd
}
