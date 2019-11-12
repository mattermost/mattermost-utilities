package main

import (
	"os"

	"github.com/mattermost/mattermost-utilities/mm-github-jira/commands"
)

func main() {
	if err := commands.Run(); err != nil {
		// error was printed by cobra
		os.Exit(1)
	}
}
