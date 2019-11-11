package commands

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mm-github-jira",
	Short: "Mattermost Github/Jira Tools",
}

func Run() {
	// error ignored, cobra already prints the error
	_ = rootCmd.Execute()
}
