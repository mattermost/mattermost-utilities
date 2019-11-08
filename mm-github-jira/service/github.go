package service

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/andygrunwald/go-jira"

	"github.com/google/go-github/v28/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type GithubClient struct {
	client *github.Client
}

type AddLabelsRequest struct {
	Repository string
	Labels     []string
}

type CreateIssuesRequest struct {
	AddLabelsRequest
}

func NewGithubClient(token string) *GithubClient {
	tc := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))
	client := github.NewClient(tc)
	return &GithubClient{client: client}
}

func (gc *GithubClient) AddLabelsToIssues(req *AddLabelsRequest, issues ...string) error {
	repository := strings.TrimSpace(req.Repository)
	repoParts := strings.Split(repository, "/")

	if len(repoParts) != 2 {
		return errors.New("invalid repository name")
	}

	owner, repo := repoParts[0], repoParts[1]
	ctx := context.Background()
	for _, issueStr := range issues {
		issueNumber, err := strconv.Atoi(issueStr)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "invalid value for issue  %s\n", issueStr)
			continue
		}

		issue, _, err := gc.client.Issues.Get(ctx, owner, repo, issueNumber)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "could not get issue with number %d: %v\n", issueNumber, err)
			continue
		}

		_, _, err = gc.client.Issues.AddLabelsToIssue(ctx, owner, repo, *issue.Number, req.Labels)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "could not add labels to issue with number %d: %v\n", issueNumber, err)
			continue
		}
	}

	return nil
}

func (gc *GithubClient) CreateIssues(req *CreateIssuesRequest, issues []jira.Issue) {

}
