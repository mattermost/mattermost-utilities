package service

import (
	"fmt"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/pkg/errors"
)

const (
	mattermostAtlassianUrl = "https://mattermost.atlassian.net"
)

type JiraClient struct {
	client *jira.Client
}

func NewJiraClient(username, token string) (*JiraClient, error) {
	tp := jira.BasicAuthTransport{
		Username: username,
		Password: token,
	}

	client, err := jira.NewClient(tp.Client(), mattermostAtlassianUrl)
	if err != nil {
		return nil, errors.Wrap(err, "could not create jira client")
	}

	return &JiraClient{client: client}, nil

}

func (jc *JiraClient) FindTasks(tasks ...string) ([]jira.Issue, error) {
	mmtasks := make([]string, len(tasks))
	for i := range tasks {
		mmtasks[i] = fmt.Sprintf("MM-%s", tasks[i])
	}

	jql := "project = MM AND " + strings.Join(mmtasks, "OR")
	issues, _, err := jc.client.Issue.Search(jql, &jira.SearchOptions{
		MaxResults: len(tasks),
		Fields:     []string{"summary", "description"},
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not get issues")
	}

	return issues, nil
}
