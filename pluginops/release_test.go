package main_test

import (
	"testing"

	pluginops "github.com/mattermost/mattermost-utilities/pluginops"
	"github.com/stretchr/testify/assert"
)

func TestGetRepoURL(t *testing.T) {
	testCases := map[string]struct {
		url           string
		org           string
		repo          string
		expectedError bool
	}{
		"http URL": {
			url:           "https://github.com/mattermost/mattermost-plugin-starter-template.git",
			org:           "mattermost",
			repo:          "mattermost-plugin-starter-template",
			expectedError: false,
		},
		"http URL with trailing slash": {
			url:           "https://github.com/mattermost/mattermost-plugin-starter-template.git/",
			org:           "",
			repo:          "",
			expectedError: true,
		},
		"ssh URL": {
			url:           "git@github.com:mattermost/mattermost-plugin-starter-template.git",
			org:           "mattermost",
			repo:          "mattermost-plugin-starter-template",
			expectedError: false,
		},
		"ssh URL with trailing slash": {
			url:           "git@github.com:mattermost/mattermost-plugin-starter-template.git/",
			org:           "",
			repo:          "",
			expectedError: true,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			org, repo, err := pluginops.GetRepoURL(tc.url)

			if tc.expectedError {
				assert.Error(t, err)
				assert.Equal(t, "", org)
				assert.Equal(t, "", repo)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.org, org)
				assert.Equal(t, tc.repo, repo)
			}
		})
	}

}
