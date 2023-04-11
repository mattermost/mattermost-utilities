package main

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/google/go-github/v32/github"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

func init() {
	rootCmd.AddCommand(prsCmd)

	prsCmd.AddCommand(
		prsListCmd,
		prsMergeCmd,
	)
}

var prsCmd = &cobra.Command{
	Use:   "prs",
	Short: "Manage PRs.",
}

var prsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List your plugin PRs.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getGitHubClient()
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), actionTimeout)
		defer cancel()

		prs, err := getPluginPRs(ctx, client)
		if err != nil {
			return err
		}

		slices.SortFunc(prs, func(a, b *prInfo) bool { return a.Repository < b.Repository })

		for _, pr := range prs {
			out := fmt.Sprintf("%s: %s", pr.Repository, pr.GetTitle())

			if pr.readyToBeMerged() {
				out += Green(checkMark)
			} else {
				out += Red(xMark)
			}
			fmt.Print(out + "\n")
		}

		return nil
	},
}

var prsMergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge plugin PRs that are ready.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getGitHubClient()
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), actionTimeout)
		defer cancel()

		prs, err := getPluginPRs(ctx, client)
		if err != nil {
			return err
		}

		if len(prs) == 0 {
			fmt.Println("No PRs to merge.")
			return nil
		}

		slices.SortFunc(prs, func(a, b *prInfo) bool { return a.Repository < b.Repository })

		for _, pr := range prs {
			if !pr.readyToBeMerged() {
				log.WithFields(log.Fields{"repo": pr.Repository, "title": pr.GetTitle()}).Debug("PR not ready to be merged")
				continue
			}

			ok, err := confirmPrompt(fmt.Sprintf("Merge %s/%s", pr.Repository, pr.GetTitle()))
			if err != nil {
				return err
			}
			if !ok {
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), actionTimeout)
			defer cancel()

			err = mergePR(ctx, client, pr)
			if err != nil {
				fmt.Printf("Failed to merged %s: %s\n\n", pr.GetHTMLURL(), err.Error())
				continue
			}

			fmt.Printf("Merged %s\n\n", pr.GetHTMLURL())
		}

		return nil
	},
}

type prInfo struct {
	Repository string
	*github.PullRequest
	reviews []*github.PullRequestReview
}

func getPluginPRs(ctx context.Context, client *github.Client) ([]*prInfo, error) {
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("is:pr is:open author:%v org:%v", user.GetLogin(), org)

	issues, _, err := client.Search.Issues(ctx, query, &github.SearchOptions{}) // TODO: handle pagination
	if err != nil {
		return nil, err
	}

	var prs []*prInfo

	var wg sync.WaitGroup
	for _, issue := range issues.Issues {
		if issue.IsPullRequest() {
			url := issue.GetHTMLURL()
			url = strings.TrimPrefix(url, "https://github.com/")
			url = strings.TrimPrefix(url, org+"/")

			repo := strings.Split(url, "/")[0]

			if !isPluginRepo(repo) {
				continue
			}

			wg.Add(1)
			go func(repo string, number int) {
				defer wg.Done()
				pr, _, err := client.PullRequests.Get(ctx, org, repo, number)
				if err != nil {
					log.WithError(err).Warn("failed to get PR")
				}

				reviews, _, err := client.PullRequests.ListReviews(ctx, org, repo, number, &github.ListOptions{}) // TODO: handle pagination
				if err != nil {
					log.WithError(err).Warn("failed to get PR")
				}
				prs = append(prs, &prInfo{
					Repository:  repo,
					PullRequest: pr,
					reviews:     reviews,
				})

			}(repo, issue.GetNumber())

		}
	}

	wg.Wait()

	return prs, nil
}

func (pr *prInfo) readyToBeMerged() bool {
	for _, l := range pr.Labels {
		if l.GetName() == "Do Not Merge" {
			return false
		}
		if l.GetName() == "Work In Progress" {
			return false
		}
		if l.GetName() == "3: QA Review" {
			return false
		}
	}

	if pr.GetMergeableState() != "clean" {
		return false
	}

	var hasOneApprovedReview bool
	for _, r := range pr.reviews {
		if r.GetState() == "APPROVED" {
			hasOneApprovedReview = true
		}
		if r.GetState() == "CHANGES_REQUESTED" {
			return false
		}

		if r.GetState() == "COMMENTED" {
			return false
		}

		// If a review from dylan is pending, the PR is not ready to be merged.
		if r.User.GetName() == "DHaussermann" && r.GetState() == "PENDING" {
			return false
		}
	}

	return hasOneApprovedReview
}

func isPluginRepo(repoName string) bool {
	if strings.HasPrefix(repoName, "mattermost-plugin-") {
		return true
	}

	for k := range defaultMapping {
		if k == repoName {
			return true
		}
	}

	return false
}

func mergePR(ctx context.Context, client *github.Client, pr *prInfo) error {
	for _, l := range pr.Labels {
		if l.GetName() == devReviewLabel.Name {
			_, err := client.Issues.RemoveLabelForIssue(ctx, org, pr.Repository, pr.GetNumber(), l.GetName())
			if err != nil {
				return err
			}
		}
	}

	_, _, err := client.Issues.AddLabelsToIssue(ctx, org, pr.Repository, pr.GetNumber(), []string{"4: Reviews Complete"})
	if err != nil {
		return err
	}

	_, _, err = client.PullRequests.Merge(ctx, org, pr.Repository, pr.GetNumber(), "", &github.PullRequestOptions{MergeMethod: "squash"})
	if err != nil {
		return err
	}

	return nil
}
