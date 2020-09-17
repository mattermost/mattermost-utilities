# Labels

This tools allows syncing defined sets of labels across multiple repositories in a GitHub organization.

## Getting Started
1. Generate a GitHub access token
2. Edit the labels in [`mapping.go`](https://github.com/mattermost/mattermost-utilities/blob/master/labels/mapping.go)
3. Run `cd labels && GITHUB_TOKEN=YOUR_GITHUB_TOKEN go run . --default`
4. Wait for the script to sync all labels

### Change labels for specific repo
For applying just the core labels run:
- `cd labels && GITHUB_TOKEN=YOUR_GITHUB_TOKEN go run . --repo $REPO_NAME --core-labels`

For applying the plugin repository labels run:
- `cd labels && GITHUB_TOKEN=YOUR_GITHUB_TOKEN go run . --repo $REPO_NAME --plugin-labels`