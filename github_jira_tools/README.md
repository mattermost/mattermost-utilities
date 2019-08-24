# Mattermost Github/Jira Tools

## github_tags.py

You can tag github issues, executing for example:

```sh
python github_tags.py --token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx -r mattermost/mattermost-server -l "Tech/Go" -l "Up For Grabs" -l "Difficulty/1:Easy" -l "Area/Technical Debt" -l "Help Wanted" 1234051
```

in this case the final number is the github issue number (you can pass multiple issue numbers if you want).

## create_github_from_jira.py

Copy tasks from Jira to Github and link them with the github issue field, executing for example:

```sh
python create_github_from_jira.py --jira-username jesus@mattermost.com --jira-token xxxxxxxxxxxxxxxxxxxxxxxx --github-token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx -r mattermost/mattermost-server -l "Tech/Go" -l "Up For Grabs" -l "Difficulty/1:Easy" -l "Area/Technical Debt" -l "Help Wanted" 1234051
```

in this case the final number is the jira issue number (you can pass multiple issue numbers if you want).

## sync_helpwanted_tickets.py

Copy all help wanted tickets without github issue linked yet to github and link them using the github issue field, executing for example:

```sh
python sync_helpwanted_tickets.py --jira-username jesus@mattermost.com --jira-token xxxxxxxxxxxxxxxxxxxxxxxx --github-token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

This script should be executed peridically somewhere.
