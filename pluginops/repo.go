package main

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(repoCmd)

	repoCmd.AddCommand(
		setupCommunityCmd,
	)
}

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Manage repositories.",
}

var setupCommunityCmd = &cobra.Command{
	Use:     "setup-fork",
	Short:   "Setup a fork of a community repository to be used for the Marketplace.",
	Example: "   pluginops repo setup-fork mattermost-plugin-giphy-moussetc",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getGitHubClient()
		if err != nil {
			log.Fatalf(err.Error())
		}

		repo := args[0]

		log.Info("Cleaning up existing labels")

		ctx, cancel := context.WithTimeout(context.Background(), actionTimeout)
		defer cancel()

		removeAllLabels(ctx, client, repo)

		log.Info("Setting up new labels")
		createOrUpdateLabels(ctx, client, repo, communityPlugins)
	},
}
