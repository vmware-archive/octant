package main

import "github.com/heptio/developer-dash/internal/commands"

// Default variables overridden by ldflags
var (
	gitCommit = "(unknown-commit)"
	buildTime = "(unknown-buildtime)"
)

func main() {
	commands.GitCommit = gitCommit
	commands.BuildTime = buildTime

	commands.Execute()
}
