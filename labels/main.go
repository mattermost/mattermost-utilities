package main

import (
	"context"
	"os"
	"sync"

	"github.com/google/go-github/v28/github"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

const org = "mattermost"

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("You need to provide an access token")
	}

	log.Info("Starting syncing labels")

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

	var wg sync.WaitGroup
	var i int
	for repo, labels := range mapping {
		for _, label := range labels {
			wg.Add(1)
			i++
			go createOrUpdateLabel(&wg, client, repo, label.ToGithubLabel())
		}
	}

	wg.Wait()
	log.WithFields(log.Fields{"number": i}).Info("Finished syncing labels")
}

func createOrUpdateLabel(wg *sync.WaitGroup, client *github.Client, repo string, label *github.Label) {
	defer wg.Done()
	logger := log.WithFields(log.Fields{"repo": repo, "label": label.GetName()})

	_, _, err := client.Issues.CreateLabel(context.Background(), org, repo, label)
	if err != nil {
		_, _, err = client.Issues.EditLabel(context.Background(), org, repo, label.GetName(), label)
		if err != nil {
			logger.WithError(err).Error("Failed to edit label")
			return
		}
		logger.Info("Edited label")
		return
	}
	logger.Info("Created label")
}
