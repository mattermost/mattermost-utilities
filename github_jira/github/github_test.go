package github

import (
	"fmt"
	"testing"
)

func Test_ParseRepoHappyPath(t *testing.T) {
	expectedOwner := "mattermost"
	expectedRepo := "mattermost-utilities"
	r, err := ParseRepo(fmt.Sprintf("%s/%s", expectedOwner, expectedRepo))
	if err != nil {
		t.Errorf("Expected to parse repo, but got err: %s", err.Error())
	}
	if r.owner != expectedOwner {
		t.Errorf("Expected owner to be %s, but got %s", r.owner, expectedOwner)
	}
	if r.repo != expectedRepo {
		t.Errorf("Expected repo to be %s, but got %s", r.repo, expectedRepo)
	}
}

func Test_ParseRepoTooLong(t *testing.T) {
	repoStr := "https://github.com/mattermost/mattermost-utilities"
	r, err := ParseRepo(repoStr)
	if err == nil {
		t.Errorf("Expected to fail parsing repo, but got owner %s and repo %s", r.owner, r.repo)
	}
}

func Test_ParseRepoEmpty(t *testing.T) {
	repoStr := ""
	r, err := ParseRepo(repoStr)
	if err == nil {
		t.Errorf("Expected to fail parsing repo, but got owner %s and repo %s", r.owner, r.repo)
	}
}
