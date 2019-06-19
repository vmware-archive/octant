package main

import (
	"math/rand"
	"time"

	"github.com/heptio/developer-dash/internal/commands"
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
