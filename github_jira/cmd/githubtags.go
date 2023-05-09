package cmd

import (
	"errors"
	"fmt"

	"github.com/mattermost/mattermost-utilities/github_jira/github"
	"github.com/mattermost/mattermost-utilities/github_jira/jira"
	"github.com/spf13/cobra"
)

var GithubLabelsCmd = &cobra.Command{
	Use:     "labelgithub",
	Short:   "Add labels to github issues provided in list of issue-numbers",
	Example: `  labelgithub -t <github token> -r mattermost/mattermost-server -l 'Tech/Go,Help Wanted,Up For Grabs' 19977 12345`,
	Args:    cobra.MinimumNArgs(1),
	RunE:    createGithubCmdF,
}

func init() {
	GithubLabelsCmd.Flags().StringP("token", "t", "", "The token used to authenticate the user against Jira.")
	GithubLabelsCmd.MarkFlagRequired("jira-token")
	GithubLabelsCmd.Flags().StringP("jira-username", "u", "", "Username of the user to get the ticket information.")
	GithubLabelsCmd.MarkFlagRequired("jira-username")
	GithubLabelsCmd.Flags().StringP("github-token", "g", "", "The token used to authenticate the user against Github.")
	GithubLabelsCmd.MarkFlagRequired("github-token")
	GithubLabelsCmd.Flags().StringP("repo", "r", "mattermost/mattermost-server", "The repository which contains the issues. E.g. mattermost/mattermost-server")
	GithubLabelsCmd.MarkFlagRequired("repo")
	GithubLabelsCmd.Flags().StringSliceP("labels", "l", []string{}, "The labels to set to the issues")
	GithubLabelsCmd.MarkFlagRequired("labels")
	GithubLabelsCmd.Flags().Bool("dry-run", false, "Skip actually creating any tickets")
	GithubLabelsCmd.Flags().Bool("debug", false, "Dump debugging information.")

	RootCmd.AddCommand(GithubLabelsCmd)
}

func createGithubCmdF(command *cobra.Command, args []string) error {
	jiraUsername, err := getNonEmptyString(command, "jira-username")
	if err != nil {
		return err
	}
	jiraToken, err := getNonEmptyString(command, "jira-token")
	if err != nil {
		return err
	}
	ghToken, err := getNonEmptyString(command, "github-token")
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

	jiraBasicAuth := jira.MakeBasicAuthStr(jiraUsername, jiraToken)
	jiraIssues, err := jira.SearchByNumber(jiraBasicAuth, debug, args)
	if err != nil {
		return fmt.Errorf("searching jira: %v", err)
	}

	if debug {
		for _, issue := range jiraIssues {
			fmt.Println("DEBUG: Jira issues:")
			fmt.Printf("%+v\n", issue)
		}
	}

	outcome, err := github.CreateIssues(jiraBasicAuth, ghToken, ghRepo, labels, jiraIssues, dryRun)

	if err != nil {
		fmt.Printf("Failed to create issues: %v\n", err)
	}
	fmt.Println(outcome.AsTables())

	return err
}
