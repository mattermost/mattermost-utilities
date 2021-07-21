package jira

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

const (
	mattermostAtlassianUrl = "https://mattermost.atlassian.net"
)

type customFields struct {
	Fields map[string]string `json:"fields"`
}

func MakeBasicAuthStr(user, pass string) string {
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(user+":"+pass)))
}

type issueSearch struct {
	Jql        string   `json:"jql"`
	MaxResults int      `json:"maxResults"`
	Fields     []string `json:"fields"`
}

type IssueFields struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`
}

type Issue struct {
	Key    string `json:"key"`
	Fields IssueFields
}

func search(basicAuth string, debug bool, jql string, maxResults int, fields []string) ([]Issue, error) {
	client := &http.Client{}
	searchBody := issueSearch{
		Jql:        jql,
		MaxResults: maxResults,
		Fields:     fields,
	}
	bodyReader, err := json.Marshal(searchBody)
	if err != nil {
		return nil, errors.Wrap(err, "searching jira")
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/rest/api/2/search", mattermostAtlassianUrl), bytes.NewReader(bodyReader))

	if err != nil {
		return nil, errors.Wrap(err, "searching jira")
	}
	req.Header.Set("Authorization", basicAuth)

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "searching jira")
	}

	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("searching jira, status code %d with text %s", resp.StatusCode, string(respBytes))
	}

	if err != nil {
		return nil, errors.Wrap(err, "parsing jira search results")
	}

	if debug {
		fmt.Println(string(respBytes))
	}

	var issues []Issue
	err = json.Unmarshal(respBytes, &issues)
	if err != nil {
		return nil, errors.Wrap(err, "parsing jira search results")
	}

	return issues, nil
}

func SearchByNumber(basicAuth string, debug bool, issueNumbers []string) ([]Issue, error) {
	issueNumbersQuery := []string{}
	for _, issueNumber := range issueNumbers {
		issueNumbersQuery = append(issueNumbersQuery, fmt.Sprintf("key = MM-%s", issueNumber))
	}
	jql := fmt.Sprintf("project = MM AND %s", strings.Join(issueNumbersQuery, " OR "))
	return search(basicAuth, debug, jql, len(issueNumbers), []string{"summary", "description"})
}

func SearchByStatus(basicAuth string, debug bool) ([]Issue, error) {
	statuses := strings.Join([]string{"Open", "Reopened"}, ", ")
	jql := fmt.Sprintf("project = MM AND status in (%s) AND fixversion = \"Help Wanted\" AND \"GITHUB ISSUE\" IS EMPTY AND type != EPIC", statuses)
	return search(basicAuth, debug, jql, 100, []string{"summary", "description"})
}

func LinkToGithub(ghUrl, jiraKey, basicAuth string) error {
	client := &http.Client{}
	jiraFields := customFields{Fields: map[string]string{"customfield_11106": ghUrl}}
	jiraFieldsReader, err := json.Marshal(jiraFields)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("creating request body for jira issue %s.", jiraKey))
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/rest/api/3/issue/%s", mattermostAtlassianUrl, jiraKey), bytes.NewReader(jiraFieldsReader))
	req.Header.Set("Authorization", basicAuth)

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("creating request for jira issue %s.", jiraKey))
	}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("unable to update jira issue %s", jiraKey))
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			respBytes = []byte("")
		}
		return fmt.Errorf("unable to update jira issue %s, response: %d, with text %s", jiraKey, resp.StatusCode, string(respBytes))
	}

	return nil
}
