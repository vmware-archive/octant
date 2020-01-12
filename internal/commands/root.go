/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// Execute executes octant.
func Execute(version string, gitCommit string, buildTime string) {
	// remove timestamp from log
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	rootCmd := newRoot(version, gitCommit, buildTime)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRoot(version string, gitCommit string, buildTime string) *cobra.Command {
	rootCmd := newOctantCmd(version)
	rootCmd.AddCommand(newVersionCmd(version, gitCommit, buildTime))

	return rootCmd
}
