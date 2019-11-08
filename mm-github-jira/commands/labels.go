package commands

import (
	"github.com/mattermost/mattermost-utilities/mm-github-jira/gh"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var addLabelsCommand = &cobra.Command{
	Use:   "labels",
	Short: "Add labels to github issues",
	Args:  cobra.MinimumNArgs(1),
	RunE:  addLabels,
}

func init() {
	addLabelsCommand.Flags().StringP("repository", "r", "", "github repository in format owner/repo (required)")
	_ = addLabelsCommand.MarkFlagRequired("repository")
	addLabelsCommand.Flags().StringArrayP("label", "l", nil, "label name to add to issue")

	rootCmd.AddCommand(addLabelsCommand)
}

func addLabels(cmd *cobra.Command, args []string) error {
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

	client := gh.NewClient(GithubToken)
	err = client.AddLabelsToIssues(gh.AddLabelsRequest{
		Repository: repository,
		Labels:     labels,
	}, args...)
	if err != nil {
		return errors.Wrap(err, "an error ocurred")
	}

	return nil
}
