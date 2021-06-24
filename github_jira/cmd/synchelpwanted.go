package cmd

import (
	"errors"
	"fmt"

	"github.com/mattermost/mattermost-utilities/github_jira/github"
	"github.com/mattermost/mattermost-utilities/github_jira/jira"
	"github.com/spf13/cobra"
)

var syncHelpWantedCmd = &cobra.Command{
	Use:     "synchelpwanted",
	Short:   "Sync help wanted issues from Jira to Github",
	Example: "  synchelpwanted -u <user> -j <jira token> -g <github token>",
	RunE:    syncHelpWantedCmdF,
}

func init() {
	syncHelpWantedCmd.Flags().StringP("jira-token", "j", "", "The token used to authenticate the user against Jira.")
	syncHelpWantedCmd.MarkFlagRequired("jira-token")
	syncHelpWantedCmd.Flags().StringP("jira-username", "u", "", "Username of the user to get the ticket information.")
	syncHelpWantedCmd.MarkFlagRequired("jira-username")
	syncHelpWantedCmd.Flags().StringP("github-token", "g", "", "The token used to authenticate the user against Github.")
	syncHelpWantedCmd.MarkFlagRequired("github-token")
	syncHelpWantedCmd.Flags().StringP("webhook-url", "w", "", "Webhook URL to send the list of created issues")
	syncHelpWantedCmd.Flags().Bool("dry-run", false, "Skip actually creating any tickets")
	syncHelpWantedCmd.Flags().Bool("debug", false, "Dump debugging information.")

	RootCmd.AddCommand(syncHelpWantedCmd)
}

func syncHelpWantedCmdF(command *cobra.Command, args []string) error {
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
	ghRepo, err := github.ParseRepo("mattermost/mattermost-server")
	if err != nil {
		return err
	}

	webhookUrl, err := command.Flags().GetString("webhook-url")
	if err != nil {
		return errors.New("invalid webhook-url parameter")
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
	jiraIssues, err := jira.SearchByStatus(jiraBasicAuth, debug)
	if err != nil {
		return fmt.Errorf("Error searching jira issues by number: %v", err)
	}

	if debug {
		for _, issue := range jiraIssues {
			fmt.Println("DEBUG: Jira issues:")
			fmt.Printf("%+v\n", issue)
		}
	}

	if len(jiraIssues) == 0 {
		return nil
	}

	outcome, err := github.CreateIssues(jiraBasicAuth, ghToken, ghRepo, []string{"Help Wanted", "Up For Grabs"}, jiraIssues, dryRun)

	outcomeToPrint := ""

	if err != nil {
		outcomeToPrint += fmt.Sprintf("Failed to create issues: %v\n", err)
	}
	outcomeToPrint += outcome.AsTables()

	if webhookUrl == "" {
		fmt.Println(outcomeToPrint)
	} else {
		sendWebhookMessage(webhookUrl, outcomeToPrint)
	}

	return nil
}
