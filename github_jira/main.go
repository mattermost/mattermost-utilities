package main

import (
	"os"

	"github.com/mattermost/mattermost-utilities/github_jira/cmd"
)

func main() {
	if err := cmd.Run(os.Args[1:]); err != nil {
		os.Exit(1)
	}
}
