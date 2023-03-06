package main

import (
	"github.com/cqroot/doter/cmd"
)

func main() {
	// You can customize DotManager before executing the command.
	// cmd.DotManager = dotmanager.New()
	cmd.Execute()
}
