package commands

import (
	"github.com/spf13/cobra"
)

var (
	GithubToken string
)

var rootCmd = &cobra.Command{
	Use:   "mm-github-jira",
	Short: "Mattermost Github/Jira Tools",
}

func Run() {
	rootCmd.PersistentFlags().StringVar(&GithubToken, "token", "", "Github token")
	_ = rootCmd.MarkPersistentFlagRequired("token")

	// error ignored, cobra already prints the error
	_ = rootCmd.Execute()
}
