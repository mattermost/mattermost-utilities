package commands

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var createGithubIssuesCommand = &cobra.Command{
	Use:   "labels",
	Short: "Add labels to github issues",
	Args:  cobra.MinimumNArgs(1),
	RunE:  createIssues,
}

func init() {
	createGithubIssuesCommand.Flags().StringP("repository", "r", "", "github repository in format owner/repo (required)")
	_ = createGithubIssuesCommand.MarkFlagRequired("repository")
	createGithubIssuesCommand.Flags().StringArrayP("label", "l", nil, "label name to add to issue")
	_ = createGithubIssuesCommand.MarkFlagRequired("label")
	createGithubIssuesCommand.Flags().String("github-token", "", "github token")
	_ = createGithubIssuesCommand.MarkFlagRequired("github-token")
	createGithubIssuesCommand.Flags().String("jira-username", "", "jira username")
	_ = createGithubIssuesCommand.MarkFlagRequired("jira-username")
	createGithubIssuesCommand.Flags().String("jira-token", "", "jira token")
	_ = createGithubIssuesCommand.MarkFlagRequired("jira-token")

	rootCmd.AddCommand(createGithubIssuesCommand)
}

func createIssues(cmd *cobra.Command, args []string) error {
	_, err := cmd.Flags().GetString("repository")
	if err != nil {
		return errors.Wrap(err, "could not read repository flag")
	}

	labels, err := cmd.Flags().GetStringArray("label")
	if err != nil {
		return errors.Wrap(err, "could not read label flag")
	}

	if len(labels) == 0 {
		return errors.New("at least one label should be applied")
	}
	return nil
}
