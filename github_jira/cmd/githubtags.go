package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/mattermost/mattermost-utilities/github_jira/github"
	"github.com/spf13/cobra"
)

var GithubLabelsCmd = &cobra.Command{
	Use:     "labelgithub",
	Short:   "Add labels to github issues provided in list of issue-numbers",
	Example: `  labelgithub -t <github token> -r mattermost/mattermost-server -l 'Tech/Go,Help Wanted,Up For Grabs' 19977 12345`,
	Args:    cobra.MinimumNArgs(1),
	RunE:    createGithubLabels,
}

func init() {
	GithubLabelsCmd.Flags().StringP("token", "t", "", "The token used to authenticate the user against github.")
	GithubLabelsCmd.MarkFlagRequired("token")
	GithubLabelsCmd.Flags().StringP("repo", "r", "mattermost/mattermost-server", "The repository which contains the issues. E.g. mattermost/mattermost-server")
	GithubLabelsCmd.MarkFlagRequired("repo")
	GithubLabelsCmd.Flags().StringSliceP("labels", "l", []string{}, "The labels to set to the issues")
	GithubLabelsCmd.MarkFlagRequired("labels")
	GithubLabelsCmd.Flags().Bool("dry-run", false, "Skip actually creating any labels")
	GithubLabelsCmd.Flags().Bool("debug", false, "Dump debugging information.")

	RootCmd.AddCommand(GithubLabelsCmd)
}

func createGithubLabels(command *cobra.Command, args []string) error {

	ghToken, err := getNonEmptyString(command, "token")
	if err != nil {
		return err
	}
	repo, err := getNonEmptyString(command, "repo")
	if err != nil {
		return err
	}
	ghRepo, err := github.ParseRepo(repo)
	if err != nil {
		return err
	}

	labels, err := getNonEmptyStringSlice(command, "labels")
	if err != nil {
		return err
	}
	dryRun, err := command.Flags().GetBool("dry-run")
	if err != nil {
		return errors.New("invalid dry-run parameter")
	}
	debug, err := command.Flags().GetBool("debug")
	if err != nil {
		return errors.New("invalid debug parameter")
	}

	if debug {
		for _, issue := range args {
			fmt.Println("DEBUG: Github issues:")
			fmt.Printf("%+v\n", issue)
		}
	}
	var intArgs []int

	for _, arg := range args {
		i, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Println(err)
		}
		intArgs = append(intArgs, i)
	}

	client := github.GetClient(ghToken)

	// Get valid labels
	validLabels, errLabels := github.GetLabelsList(client, ghRepo, labels)
	if errLabels != nil {
		return errLabels
	}
	validIssues, errIssues := github.GetIssuesList(client, ghRepo, intArgs)
	if errLabels != nil {
		return errIssues
	}

	multiError := github.SetLabels(client, ghRepo, validLabels, validIssues)
	var newErr error
	if len(multiError) > 0 {
		newErr = errors.New("multiple errors found")
	}

	for mErr := range multiError {
		fmt.Println(mErr)
	}

	return newErr
}
