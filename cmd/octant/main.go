/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"math/rand"
	"time"

	"github.com/vmware/octant/internal/commands"
)

// Default variables overridden by ldflags
var (
	version   = "(dev-version)"
	gitCommit = "(dev-commit)"
	buildTime = "(dev-buildtime)"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	commands.Execute(version, gitCommit, buildTime)
}
