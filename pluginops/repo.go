package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/go-github/v32/github"
	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	mattermostBuild      = "mattermost-build"
	integrationsTeamSlug = "integrations"
	securityionsTeamSlug = "core-security"
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

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		repoSettings := &github.Repository{
			Private: github.Bool(true),

			AllowMergeCommit: github.Bool(true),
			AllowRebaseMerge: github.Bool(false),
			AllowSquashMerge: github.Bool(true),

			HasIssues:   github.Bool(false),
			HasWiki:     github.Bool(false),
			HasPages:    github.Bool(false),
			HasProjects: github.Bool(false),
		}
		_, _, err = client.Repositories.Edit(ctx, org, repo, repoSettings)
		if err != nil {
			log.WithField("repo", repo).WithError(err).Error("Failed to update repository")
			return
		}

		log.Info("Successfully updated repository settings")

		permssionPush := &github.RepositoryAddCollaboratorOptions{
			Permission: "push",
		}

		_, _, err = client.Repositories.AddCollaborator(ctx, org, repo, mattermostBuild, permssionPush)
		if err != nil {
			log.WithField("repo", repo).WithError(err).Errorf("Failed add %s to repo", mattermostBuild)
			return
		}

		permssionTriage := &github.TeamAddTeamRepoOptions{
			Permission: "triage",
		}

		_, err = client.Teams.AddTeamRepoBySlug(ctx, org, integrationsTeamSlug, org, repo, permssionTriage)
		if err != nil {
			log.WithField("repo", repo).WithError(err).Errorf("Failed add %s to repo", integrationsTeamSlug)
			return
		}

		_, err = client.Teams.AddTeamRepoBySlug(ctx, org, securityionsTeamSlug, org, repo, permssionTriage)
		if err != nil {
			log.WithField("repo", repo).WithError(err).Errorf("Failed add %s to repo", securityionsTeamSlug)
			return
		}

		log.Info("Successfully added needed permissions")

		log.Info("Cleaning up existing labels")
		removeAllLabels(client, repo)

		log.Info("Setting up new labels")
		createOrUpdateLabels(client, repo, communityPlugins)

		for {
			prompt := promptui.Prompt{
				Label:     "Do want to add a Collaborator to the repo",
				IsConfirm: true,
			}
			_, err := prompt.Run()
			if err != nil {
				if errors.Is(err, promptui.ErrAbort) {
					break
				}

				log.WithField("repo", repo).WithError(err).Error("Prompt failed")
				return
			}

			prompt = promptui.Prompt{
				Label: "Which Collaborator do you want to add?",
			}

			user, err := prompt.Run()
			if err != nil {
				log.WithField("repo", repo).WithError(err).Error("Prompt failed")
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			_, _, err = client.Users.Get(ctx, user)
			if err != nil {
				fmt.Printf("\nUser %s doesn't exist\n", user)
				continue
			}

			_, _, err = client.Repositories.AddCollaborator(ctx, org, repo, user, &github.RepositoryAddCollaboratorOptions{
				Permission: "pull",
			})
			if err != nil {
				log.WithField("repo", repo).WithError(err).Errorf("Failed add %s to repo", user)
				continue
			}

			fmt.Printf("Successfully added %s to the repo\n", user)
		}
	},
}
