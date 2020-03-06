package main

import (
	"context"
	"os"
	"strings"
	"sync"

	"github.com/google/go-github/v28/github"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

const org = "mattermost"

var verbose bool
var useDefault bool
var useRepo string
var useCoreLabels bool
var usePluginLabels bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	rootCmd.Flags().BoolVar(&useDefault, "default", false, "Updates default repos with the default labels, can not be combined with other flags.")
	rootCmd.Flags().StringVar(&useRepo, "repo", "", "Github Mattermost repository name.")
	rootCmd.Flags().BoolVar(&useCoreLabels, "core-labels", false, "Use a core set of Mattermost labels.")
	rootCmd.Flags().BoolVar(&usePluginLabels, "plugin-labels", false, "Use Mattermost plugin-specific labels.")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Fatal("command failed")
	}
}

var rootCmd = &cobra.Command{
	Use:   "labels",
	Short: "This tools allows syncing defined sets of labels across multiple repositories in a GitHub organization.",
	Run: func(cmd *cobra.Command, args []string) {
		token := os.Getenv("GITHUB_TOKEN")
		if token == "" {
			log.Fatal("You need to provide an access token")
		}

		if verbose {
			log.SetLevel(log.DebugLevel)
		}

		var mapping map[string][]Label
		switch {
		case useDefault:
			mapping = defaultMapping

		case useRepo != "":
			labels := defaultLabels
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

		log.Info("Starting syncing labels")

		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(context.Background(), ts)
		client := github.NewClient(tc)

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

func createOrUpdateLabels(client *github.Client, repo string, labels []Label) {
	remoteLabels, err := fetchLabels(client, repo)
	if err != nil {
		log.WithField("repo", repo).WithError(err).Error("Failed to fetch labels")
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
			if label.Equal(remoteLabel) {
				found = true
				break
			}
		}

		if !found {
			log.WithFields(log.Fields{"label": remoteLabel.GetName(), "repo": repo}).Warn("found untracked label")
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
