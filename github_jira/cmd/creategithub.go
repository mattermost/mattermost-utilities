package cmd

import (
	"errors"
	"fmt"

	"github.com/mattermost/mattermost-utilities/github_jira/github"
	"github.com/mattermost/mattermost-utilities/github_jira/jira"
	"github.com/spf13/cobra"
)

var createGithubCmd = &cobra.Command{
	Use:     "creategithub",
	Short:   "Create github issues from list of JIRA issue numbers",
	Example: `  creategithub -u <user> -j <jira token> -g <github token> -r mattermost/mattermost-server -l 'Tech/Go,Help Wanted,Up For Grabs' 19977 12345`,
	Args:    cobra.MinimumNArgs(1),
	RunE:    createGithubCmdF,
}

func init() {
	createGithubCmd.Flags().StringP("jira-token", "j", "", "The token used to authenticate the user against Jira.")
	createGithubCmd.MarkFlagRequired("jira-token")
	createGithubCmd.Flags().StringP("jira-username", "u", "", "Username of the user to get the ticket information.")
	createGithubCmd.MarkFlagRequired("jira-username")
	createGithubCmd.Flags().StringP("github-token", "g", "", "The token used to authenticate the user against Github.")
	createGithubCmd.MarkFlagRequired("github-token")
	createGithubCmd.Flags().StringP("repo", "r", "", "The repository which contains the issues. E.g. mattermost/mattermost-server")
	createGithubCmd.MarkFlagRequired("repo")
	createGithubCmd.Flags().StringSliceP("labels", "l", []string{}, "The labels to set to the issues")
	createGithubCmd.MarkFlagRequired("labels")
	createGithubCmd.Flags().Bool("dry-run", false, "Skip actually creating any tickets")
	createGithubCmd.Flags().Bool("debug", false, "Dump debugging information.")

	RootCmd.AddCommand(createGithubCmd)
}

func getStr(command *cobra.Command, name string) (string, error) {
	str, err := command.Flags().GetString(name)
	if err != nil {
		return "", errors.New(fmt.Sprintf("invalid %s parameter", name))
	}
	if str == "" {
		return "", errors.New(fmt.Sprintf("expected %s to not be empty", name))
	}
	return str, nil
}

func createGithubCmdF(command *cobra.Command, args []string) error {
	jiraUsername, err := getStr(command, "jira-username")
	if err != nil {
		return err
	}
	jiraToken, err := getStr(command, "jira-token")
	if err != nil {
		return err
	}
	ghToken, err := getStr(command, "github-token")
	if err != nil {
		return err
	}
	repo, err := getStr(command, "repo")
	if err != nil {
		return err
	}
	ghRepo, err := github.ParseRepo(repo)
	if err != nil {
		return err
	}

	// TODO: Should at least one label be required?
	labels, err := command.Flags().GetStringSlice("labels")
	if err != nil {
		return errors.New("invalid labels parameter")
	}
	dryRun, err := command.Flags().GetBool("dry-run")
	if err != nil {
		return errors.New("invalid dry-run parameter")
	}
	debug, err := command.Flags().GetBool("dry-run")
	if err != nil {
		return errors.New("invalid debug parameter")
	}

	jiraBasicAuth := jira.MakeBasicAuthStr(jiraUsername, jiraToken)
	jiraIssues, err := jira.SearchByNumber(jiraBasicAuth, debug, args)
	fmt.Printf("ISSUES %+v\n", jiraIssues)
	fmt.Printf("err %+v\n", err)
	if err != nil {
		return errors.New("Error searching jira issues by number: " + err.Error())
	}

	if debug {
		for _, issue := range jiraIssues {
			fmt.Println("DEBUG: Jira issues:")
			fmt.Printf("%+v\n", issue)
		}
	}
	fmt.Println(github.CreateIssues(jiraBasicAuth, ghToken, ghRepo, labels, jiraIssues, dryRun))

	return nil
}
