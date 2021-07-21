package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v35/github"
	"github.com/mattermost/mattermost-utilities/github_jira/jira"
	"golang.org/x/oauth2"
)

type label struct {
	Name string `json:"name"`
}

func has(needle string, haystack []string) bool {
	for _, item := range haystack {
		if needle == item {
			return true
		}
	}
	return false
}

func ParseRepo(repoStr string) (repo, error) {
	ownerAndRepo := strings.Split(repoStr, "/")
	if len(ownerAndRepo) != 2 {
		return repo{}, fmt.Errorf(`Expected repo to be of form "<owner>/<repo>", but got %s`, repoStr)
	}
	return repo{
		owner: ownerAndRepo[0],
		repo:  ownerAndRepo[1],
	}, nil
}

type repo struct {
	owner string
	repo  string
}

type LinkedIssue struct {
	JiraKey     string
	GithubIssue github.Issue
}

type FailedLink struct {
	JiraKey string
	Message string
}

type CreateOutcome struct {
	LinkedIssues []LinkedIssue
	FailedLinks  []FailedLink
}

func (o *CreateOutcome) AsTables() string {
	table := ""
	keyHeader := "Jira Key"
	keyHeaderLength := strconv.Itoa(len(keyHeader))

	if numCreated := len(o.LinkedIssues); numCreated > 0 {
		table += fmt.Sprintf(`Created %d github issues:
%s | Github URL
---------------------
`, numCreated, keyHeader)

		for _, linkedIssue := range o.LinkedIssues {
			table += fmt.Sprintf("%"+keyHeaderLength+"s | %s\n", linkedIssue.JiraKey, *linkedIssue.GithubIssue.HTMLURL)
		}
		table += "\n"
	}
	if numFailed := len(o.FailedLinks); numFailed > 0 {
		table += fmt.Sprintf(`Failed creating %d github issues:
%s | Error
%s
`, numFailed, keyHeader, strings.Repeat("-", len(keyHeader)+8))

		for _, failure := range o.FailedLinks {
			table += fmt.Sprintf("%"+keyHeaderLength+"s | %s\n", failure.JiraKey, failure.Message)
		}
		table += "\n"
	}
	return table
}

func CreateIssues(jiraBasicAuth string, ghToken string, repo repo, labels []string, jiraIssues []jira.Issue, dryRun bool) (CreateOutcome, error) {
	outcome := CreateOutcome{
		LinkedIssues: []LinkedIssue{},
		FailedLinks:  []FailedLink{},
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	r, _, err := client.Repositories.Get(ctx, repo.owner, repo.repo)
	if err != nil {
		return outcome, err
	}

	// ListTags
	finalLabels := []label{}
	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", *r.LabelsURL, nil)
	if err != nil {
		return outcome, err
	}
	resp, err := httpClient.Do(req)

	if err != nil {
		return outcome, err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return outcome, err
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return outcome, err
	}

	var candidateLabels []label
	err = json.Unmarshal(respBytes, &candidateLabels)

	if err != nil {
		return outcome, err
	}

	for _, label := range candidateLabels {
		if has(label.Name, labels) {
			finalLabels = append(finalLabels, label)
		}
	}

	if dryRun {
		fmt.Println("We haven't created the github ticket because --dry-run flag was detected. Tickets information:")
	}

	for _, issue := range jiraIssues {
		title := issue.Fields.Summary
		key := issue.Key
		markdownDescription := jira.ToMarkdown(strings.Split(issue.Fields.Description, "\n"))
		description := strings.Join(markdownDescription, "\n") + "\n\n" + strings.Replace(templateContributing, "{{TICKET}}", key, 1)

		if dryRun {
			fmt.Printf("------\n%s\n%s\n\n%s\n", title, strings.Repeat("=", len(title)), description)
			continue
		}
		// Add one second sleep per https://docs.github.com/en/rest/guides/best-practices-for-integrators#dealing-with-abuse-rate-limits
		time.Sleep(1 * time.Second)
		issueRequest := github.IssueRequest{}
		newIssue, _, err := client.Issues.Create(ctx, repo.owner, repo.repo, &issueRequest)
		if err != nil {
			outcome.FailedLinks = append(outcome.FailedLinks, FailedLink{
				JiraKey: key,
				Message: err.Error(),
			})
			continue
		}
		err = jira.LinkToGithub(*newIssue.HTMLURL, key, jiraBasicAuth)
		if err != nil {
			outcome.FailedLinks = append(outcome.FailedLinks, FailedLink{
				JiraKey: key,
				Message: err.Error(),
			})
			continue
		}
		outcome.LinkedIssues = append(outcome.LinkedIssues, LinkedIssue{
			JiraKey:     key,
			GithubIssue: *newIssue,
		})
	}

	if numFailures := len(outcome.FailedLinks); numFailures > 0 {
		return outcome, fmt.Errorf("Failed creating %d issues", numFailures)
	}
	return outcome, nil
}
