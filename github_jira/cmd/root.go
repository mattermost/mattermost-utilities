package cmd

import (
	"github.com/spf13/cobra"
)

type Command = cobra.Command

func Run(args []string) error {
	RootCmd.SetArgs(args)
	return RootCmd.Execute()
}

var RootCmd = &cobra.Command{
	Use:   "github_jira",
	Short: "Manage Mattermost Github & Jira",
	Long:  "Mattermost CLI for managing & synchronizing Mattermost github & jira",
}
