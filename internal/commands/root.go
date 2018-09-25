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
func Execute() {
	rootCmd := newRoot()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func newRoot() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "hcli",
		Short: "hcli is the Heptio CLI",
	}

	rootCmd.AddCommand(newDashCmd())

	return rootCmd
}
