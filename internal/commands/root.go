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

// Execute executes hcli.
func Execute(gitCommit string, buildTime string) {
	rootCmd := newRoot(gitCommit, buildTime)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func newRoot(gitCommit string, buildTime string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "hcli",
		Short: "hcli is the Heptio CLI",
	}

	rootCmd.AddCommand(newDashCmd())
	rootCmd.AddCommand(newVersionCmd(gitCommit, buildTime))

	return rootCmd
}
