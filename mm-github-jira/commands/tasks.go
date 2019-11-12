package commands

import (
	"github.com/mattermost/mattermost-utilities/mm-github-jira/service"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var createGithubIssuesCommand = &cobra.Command{
	Use:   "tasks",
	Short: "Create Github issues from Jira tasks",
	Args:  cobra.MinimumNArgs(1),
	RunE:  createIssues,
}

func init() {
	createGithubIssuesCommand.Flags().StringP("repository", "r", "", "Github repository in format owner/repo (required)")
	_ = createGithubIssuesCommand.MarkFlagRequired("repository")

	createGithubIssuesCommand.Flags().StringArrayP("label", "l", nil, "label name to add to issue (required)")
	_ = createGithubIssuesCommand.MarkFlagRequired("label")

	createGithubIssuesCommand.Flags().String("github-token", "", "github token (required)")
	_ = createGithubIssuesCommand.MarkFlagRequired("github-token")

	createGithubIssuesCommand.Flags().String("jira-username", "", "jira username (required)")
	_ = createGithubIssuesCommand.MarkFlagRequired("jira-username")

	createGithubIssuesCommand.Flags().String("jira-token", "", "jira token (required)")
	_ = createGithubIssuesCommand.MarkFlagRequired("jira-token")

	rootCmd.AddCommand(createGithubIssuesCommand)
}

func createIssues(cmd *cobra.Command, args []string) error {
	jiraToken, err := cmd.Flags().GetString("jira-token")
	if err != nil {
		return errors.Wrap(err, "could not read jira-token flag")
	}

	jiraUsername, err := cmd.Flags().GetString("jira-username")
	if err != nil {
		return errors.Wrap(err, "could not read jira-username flag")
	}

	jiraClient, err := service.NewJiraClient(jiraUsername, jiraToken)
	if err != nil {
		return errors.Wrap(err, "could not create jira client")
	}

	issues, err := jiraClient.FindTasks(args...)
	if err != nil {
		return errors.Wrap(err, "could not retrieve tasks from jira")
	}

	repository, err := cmd.Flags().GetString("repository")
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
