# mattermost-utilities: Set of utilities to help in mattermost development

Mattermost is an open source, self-hosted Slack-alternative https://mattermost.org.

Currently this repo contains these utilities:

* **github_jira_tools**: is a CLI to help create Github issues from Mattermost Jira and a Docker container that runs in a cron-like fashion and keeps the help wanted tickets up to date.
* **github_jira**: is written in Go and is the successor to `github_jira_tools`.
* **mmgotool**: is a CLI to help with some task related to the mattermost-server development.
* **mmjstool**: is a CLI to help with some task related to the mattermost-webapp, mattermost-redux and mattermost-mobile development.
* **pluginops**: This tools allows syncing defined sets of labels across multiple repositories in a GitHub organization.

You can build the `github_jira` container with `docker build -f github_jira/Dockerfile.sync_helpwanted_tickets ./github_jira -t sync-helpwanted-tickets:latest` .
