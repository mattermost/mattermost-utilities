package main

import (
	"context"
	"strings"
	"sync"

	"github.com/google/go-github/v32/github"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	useDefault      bool
	useRepo         string
	useCoreLabels   bool
	usePluginLabels bool
)

func init() {
	rootCmd.AddCommand(labelsCmd)

	labelsCmd.AddCommand(
		labelsSyncCmd,
		labelsMigrateCmd,
	)

	labelsMigrateCmd.Flags().BoolVar(&useDefault, "default", false, "Updates default repos with the default labels, can not be combined with other flags.")
	labelsMigrateCmd.Flags().StringVar(&useRepo, "repo", "", "Github Mattermost repository name.")

	labelsSyncCmd.Flags().BoolVar(&useDefault, "default", false, "Updates default repos with the default labels, can not be combined with other flags.")
	labelsSyncCmd.Flags().StringVar(&useRepo, "repo", "", "Github Mattermost repository name.")
	labelsSyncCmd.Flags().BoolVar(&useCoreLabels, "core-labels", false, "Use a core set of Mattermost labels.")
	labelsSyncCmd.Flags().BoolVar(&usePluginLabels, "plugin-labels", false, "Use Mattermost plugin-specific labels.")
}

var labelsCmd = &cobra.Command{
	Use:   "labels",
	Short: "Manage labels.",
}

var labelsSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "This tools allows syncing defined sets of labels across multiple repositories in a GitHub organization.",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getGitHubClient()
		if err != nil {
			log.Fatalf(err.Error())
		}

		var mapping map[string][]Label
		switch {
		case useDefault:
			mapping = defaultMapping

		case useRepo != "":
			var labels []Label
			switch {
			case usePluginLabels:
				labels = pluginLabels
			case useCoreLabels:
				labels = coreLabels
			default:
				log.Fatalf("You need to specify labels to set on %q, e.g. --core-labels or --plugin-labels", useRepo)
			}
			mapping = map[string][]Label{
				useRepo: labels,
			}
		default:
			log.Fatalf("You need to specify the repo to apply labels to, e.g. --repo=mattermost-test, or run the default settings with --default")
		}

		log.Info("Starting to sync labels")

		var wg sync.WaitGroup
		for repo, labels := range mapping {
			wg.Add(1)
			go func(repo string, labels []Label) {
				createOrUpdateLabels(client, repo, labels)
				wg.Done()
			}(repo, labels)
		}

		wg.Wait()
		log.WithFields(log.Fields{"number of repos": len(mapping)}).Info("Finished syncing labels")

	},
}

var labelsMigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate labels across multiple repositories in a GitHub organization to a new naming schema.",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getGitHubClient()
		if err != nil {
			log.Fatalf(err.Error())
		}

		var mapping map[string][]Label
		switch {
		case useDefault:
			mapping = defaultMapping

		case useRepo != "":
			mapping = map[string][]Label{
				useRepo: nil,
			}
		default:
			log.Fatalf("You need to specify the repo to apply labels to, e.g. --repo=mattermost-test, or run the default settings with --default")
		}

		log.Info("Start to migrate labels")

		var wg sync.WaitGroup
		for repo := range mapping {
			wg.Add(1)
			go func(repo string) {
				migrateLabels(client, repo)
				wg.Done()
			}(repo)
		}

		wg.Wait()
		log.WithFields(log.Fields{"number of repos": len(mapping)}).Info("Finished migrating labels")
	},
}

func migrateLabels(client *github.Client, repo string) {
	remoteLabels, err := fetchLabels(client, repo)
	if err != nil {
		log.WithField("repo", repo).WithError(err).Error("Failed to fetch labels")
		return
	}

	for old, new := range migrateMap {
		for _, remoteLabel := range remoteLabels {
			if strings.EqualFold(remoteLabel.GetName(), old) {
				if new == "" {
					logger := log.WithFields(log.Fields{"repo": repo, "old": old})

					_, err = client.Issues.DeleteLabel(context.Background(), org, repo, old)
					if err != nil {
						logger.WithError(err).Error("Failed to delete label")
						continue
					}

					logger.Info("Deleted unneeded label")
				} else {
					logger := log.WithFields(log.Fields{"repo": repo, "old": old, "new": new})

					newLabel := &github.Label{Name: &new}
					_, _, err = client.Issues.EditLabel(context.Background(), org, repo, old, newLabel)
					if err != nil {
						logger.WithError(err).Error("Failed to edit label")
						continue
					}

					logger.Info("Migrated label")
				}
			}
		}
	}
}

func removeAllLabels(client *github.Client, repo string) {
	remoteLabels, err := fetchLabels(client, repo)
	if err != nil {
		log.WithField("repo", repo).WithError(err).Error("Failed to fetch labels")
		return
	}

	for _, remoteLabel := range remoteLabels {
		logger := log.WithFields(log.Fields{"repo": repo, "label": remoteLabel.GetName()})

		_, err = client.Issues.DeleteLabel(context.Background(), org, repo, remoteLabel.GetName())
		if err != nil {
			logger.WithError(err).Error("Failed to delete label")
			continue
		}

		logger.Info("Deleted label")
	}
}

func createOrUpdateLabels(client *github.Client, repo string, labels []Label) {
	remoteLabels, err := fetchLabels(client, repo)
	if err != nil {
		log.WithField("repo", repo).WithError(err).Error("Failed to fetch labels")
		return
	}

	for _, label := range labels {
		logger := log.WithFields(log.Fields{"repo": repo, "label": label.Name})
		found := false

		for _, remoteLabel := range remoteLabels {
			if strings.EqualFold(remoteLabel.GetName(), label.Name) {
				if !label.Equal(remoteLabel) {
					_, _, err = client.Issues.EditLabel(context.Background(), org, repo, label.Name, label.ToGithubLabel())
					if err != nil {
						logger.WithError(err).Error("Failed to edit label")
						continue
					}

					logger.Info("Edited label")
				} else {
					logger.Debug("Label is in sync")
				}

				found = true
			}
		}

		if !found {
			_, _, err := client.Issues.CreateLabel(context.Background(), org, repo, label.ToGithubLabel())
			if err != nil {
				logger.WithError(err).Error("Failed to create label")
				continue
			}

			logger.Info("Created label")
		}
	}

	for _, remoteLabel := range remoteLabels {
		found := false

		for _, label := range labels {
			if strings.EqualFold(remoteLabel.GetName(), label.Name) {
				found = true
				break
			}
		}

		if !found {
			log.WithFields(log.Fields{"label": remoteLabel.GetName(), "repo": repo}).Warn("Found untracked label")
		}
	}
}

func fetchLabels(client *github.Client, repo string) ([]*github.Label, error) {
	opt := &github.ListOptions{
		PerPage: 50,
	}

	// Get all labels of results
	var allLabels []*github.Label

	for {
		repos, resp, err := client.Issues.ListLabels(context.Background(), org, repo, opt)
		if err != nil {
			return nil, err
		}

		allLabels = append(allLabels, repos...)

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return allLabels, nil
}
