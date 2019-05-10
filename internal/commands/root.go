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

// Execute executes clustereye.
func Execute(version string, gitCommit string, buildTime string) {
	rootCmd := newRoot(version, gitCommit, buildTime)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRoot(version string, gitCommit string, buildTime string) *cobra.Command {
	rootCmd := newClusterEyeCmd()
	rootCmd.AddCommand(newVersionCmd(version, gitCommit, buildTime))

	return rootCmd
}
