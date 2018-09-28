package main

import "github.com/heptio/developer-dash/internal/commands"

// Default variables overridden by ldflags
var (
	gitCommit = "(dev-commit)"
	buildTime = "(dev-buildtime)"
)

func main() {
	commands.Execute(gitCommit, buildTime)
}
