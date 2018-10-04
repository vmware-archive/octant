package main

import "github.com/heptio/developer-dash/internal/commands"

// Default variables overridden by ldflags
var (
	version   = "(dev-version)"
	gitCommit = "(dev-commit)"
	buildTime = "(dev-buildtime)"
)

func main() {
	commands.Execute(version, gitCommit, buildTime)
}
