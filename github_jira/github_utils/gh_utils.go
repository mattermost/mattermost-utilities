package github_utils

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	gh "github.com/google/go-github/v52/github"
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

func validateLabel(foundLabel string, inputLabels []*gh.Label) bool {
	for _, item := range inputLabels {
		if foundLabel == *item.Name {
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
	GithubIssue gh.Issue
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
	client := GetClient(ghToken)

	// ListTags
	validLabels, errLabels := GetLabelsList(client, repo, labels)
	if errLabels != nil {
		return outcome, errLabels
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
		issueRequest := gh.IssueRequest{Title: &title,
			Body:   &description,
			Labels: &validLabels}
		fmt.Printf("creating github issue for jira key %s\n", issue.Key)
		newIssue, _, err := client.Issues.Create(ctx, repo.owner, repo.repo, &issueRequest)
		if err != nil {
			outcome.FailedLinks = append(outcome.FailedLinks, FailedLink{
				JiraKey: key,
				Message: err.Error(),
			})
			continue
		}
		fmt.Printf("Updating %s to jira \n", *newIssue.URL)
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
		return outcome, fmt.Errorf("failed creating %d issues", numFailures)
	}
	return outcome, nil
}

func GetClient(ghToken string) *gh.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	return gh.NewClient(tc)
}

func GetLabelsList(client *gh.Client, repo repo, labels []string) ([]string, error) {
	var foundLabels []string
	ctx := context.Background()
	existingGhLabels, _, err := client.Issues.ListLabels(ctx, repo.owner, repo.repo, nil)

	if err != nil {
		return foundLabels, err
	}

	var flag bool
	for _, ll := range labels {
		flag = validateLabel(ll, existingGhLabels)
		if flag {
			foundLabels = append(foundLabels, ll)
		} else {
			fmt.Printf("Unknown label %s \n", ll)
		}
	}
	return foundLabels, nil
}

func checkExistingIssue(issueId int, issues []*gh.Issue) bool {
	for _, issue := range issues {
		if *issue.Number == issueId {
			return true
		}
	}
	return false
}
func GetIssuesList(client *gh.Client, repo repo, issues []int) ([]int, error) {

	var foundIssues []int
	ctx := context.Background()

	issuesList, _, err := client.Issues.ListByRepo(ctx, repo.owner, repo.repo, nil)
	if err != nil {
		return foundIssues, err
	}
	var foundIssue bool
	for _, issue := range issues {
		foundIssue = checkExistingIssue(issue, issuesList)
		if foundIssue {
			foundIssues = append(foundIssues, issue)
		} else {
			fmt.Printf("unknown issue %d \n", issue)
		}
	}
	return foundIssues, nil
}

func SetLabels(client *gh.Client, repo repo, labels []string, issues []int) (multiError []error) {

	ctx := context.Background()
	for _, issue := range issues {
		err := setLabel(ctx, client, repo, issue, labels)
		if err != nil {
			multiError = append(multiError, fmt.Errorf("error %v setting label to issue %d ", err, issue))
		}
	}
	return multiError
}

func setLabel(ctx context.Context, client *gh.Client, repo repo, issue int, labels []string) error {
	_, _, err := client.Issues.AddLabelsToIssue(ctx, repo.owner, repo.repo, issue, labels)
	return err
}
