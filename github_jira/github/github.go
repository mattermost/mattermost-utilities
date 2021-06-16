package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/go-github/v35/github"
	"github.com/mattermost/mattermost-utilities/github_jira/jira"
	"github.com/mattermost/mattermost-utilities/github_jira/verbiages"
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
		return repo{}, errors.New(fmt.Sprintf(`Expected repo to be of form "<owner>/<repo>", but got %s`, repoStr))
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

func CreateIssues(jiraBasicAuth string, ghToken string, repo repo, labels []string, jiraIssues []jira.Issue, dryRun bool) string {
	logs := ""

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	r, _, err := client.Repositories.Get(ctx, repo.owner, repo.repo)
	if err != nil {
		logs += fmt.Sprintf("could not get %s repo: %s\n", repo, err.Error())
		return logs
	}

	// ListTags
	finalLabels := []label{}
	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", *r.LabelsURL, nil)
	if err != nil {
		logs += fmt.Sprintf("could not get %s repo labels: %s\n", repo, err.Error())
		return logs
	}
	resp, err := httpClient.Do(req)

	if err != nil {
		logs += fmt.Sprintf("could not get %s repo labels: %s\n", repo, err.Error())
		return logs
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		logs += fmt.Sprintf("could not get %s repo labels: %s\n", repo, err.Error())
		return logs
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logs += fmt.Sprintf("could not read %s repo labels: %s\n", repo, err.Error())
		return logs
	}

	var candidateLabels []label
	err = json.Unmarshal(respBytes, &candidateLabels)

	if err != nil {
		logs += fmt.Sprintf("could not read %s repo labels: %s\n", repo, err.Error())
		return logs
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
		description := strings.Join(markdownDescription, "\n") + "\n\n" + strings.Replace(verbiages.TemplateContributing, "{{TICKET}}", key, 1)

		if dryRun {
			fmt.Printf("------\n%s\n%s\n\n%s\n", title, strings.Repeat("=", len(title)), description)
			continue
		}
		// Add one second sleep per https://docs.github.com/en/rest/guides/best-practices-for-integrators#dealing-with-abuse-rate-limits
		time.Sleep(1 * time.Second)
		issueRequest := github.IssueRequest{}
		newIssue, _, err := client.Issues.Create(ctx, repo.owner, repo.repo, &issueRequest)
		if err != nil {
			logs += fmt.Sprintf("Unable to create issue for jira issue %s. error: %s\n", key, err.Error())
			continue
		}
		err = jira.LinkToGithub(*newIssue.HTMLURL, key, jiraBasicAuth)
		if err != nil {
			logs += err.Error() + "\n"
			continue
		}
		logs += fmt.Sprintf("Created github issue for the jira issue %s here: %s\n", key, *newIssue.HTMLURL)
	}

	return logs
}
