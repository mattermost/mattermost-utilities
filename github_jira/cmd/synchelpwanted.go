package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

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
		return errors.New("Error searching jira issues by number: " + err.Error())
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

	log := github.CreateIssues(jiraBasicAuth, ghToken, ghRepo, []string{"Help Wanted", "Up For Grabs"}, jiraIssues, dryRun)

	if log == "" {
		return nil
	}

	if webhookUrl == "" {
		fmt.Println(log)
		return nil
	}

	msg, err := json.Marshal(map[string]string{"text": log})
	if err != nil {
		fmt.Printf("Unable to send log to webhook: %+v\n", err)
		return nil
	}
	req, err := http.NewRequest("POST", webhookUrl, bytes.NewReader(msg))
	if err != nil {
		fmt.Printf("Unable to send log to webhook: %+v\n", err)
		return nil
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending log %s\n", err)
	}
	if resp.StatusCode >= 400 {
		fmt.Printf("Sending log failed with status code %d\n", resp.StatusCode)
	}

	return nil
}
