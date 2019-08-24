import requests
import click
import json
from requests.auth import HTTPBasicAuth
from github import Github

def create_github_issues(jira_username, jira_token, github_token, issues):
    header = '''
If you're interested please comment here and come [join our "Contributors" community channel](https://community.mattermost.com/core/channels/tickets) on our daily build server, where you can discuss questions with community members and the Mattermost core team. For technical advice or questions, please  [join our "Developers" community channel](https://community.mattermost.com/core/channels/developers).

New contributors please see our [Developer's Guide](https://developers.mattermost.com/contribute/getting-started/).

----

**Notes**: [Jira ticket](https://mattermost.atlassian.net/browse/{{TICKET}})
    '''

    g = Github(github_token)
    r = g.get_repo('mattermost/mattermost-server')
    final_labels = []
    labels = ['Help Wanted', 'Up For Grabs']
    for label in r.get_labels():
        if label.name in labels:
            final_labels.append(label)

    for issue in issues:
        title = issue['fields']['summary']
        key = issue['key']
        description = header.replace("{{TICKET}}", key) + "\n" + issue['fields']['description']

        try:
            new_issue = r.create_issue(
                title=title,
                body=description,
                labels=final_labels,
            )
        except Exception as e:
            print("Unable to create issue for jira issue {}. error: {}".format(key, e))
            return

        try:
            resp = requests.put(
                "https://mattermost.atlassian.net/rest/api/3/issue/"+key,
                json={
                    "fields": {
                        "customfield_11106": new_issue.html_url,
                    },
                },
                auth=HTTPBasicAuth(jira_username, jira_token)
            )
        except Exception as e:
            print("Unable to update jira issue {}. error: {}".format(key, e))
            return

        print("Created github issue for the jira issue {} here: {}".format(key, new_issue.html_url))

@click.command()
@click.option('--jira-token', '-j', prompt='Your Jira access token', help='The token used to authenticate the user against Jira.')
@click.option('--jira-username', '-u', prompt='Your Jira username', help='Username of the user to get the ticket information.')
@click.option('--github-token', '-g', prompt='Your Github access token', help='The token used to authenticate the user against Github.')
def cli(jira_username, jira_token, github_token):
    data = {
            "jql":"project = MM AND status = Open AND fixversion = \"Help Wanted\" AND \"GITHUB ISSUE\" IS EMPTY",
            "maxResults": 100,
            "fields": ["summary", "description"],
    }
    resp = requests.post(
        "https://mattermost.atlassian.net/rest/api/2/search",
        json=data,
        auth=HTTPBasicAuth(jira_username, jira_token)
    )
    issues = resp.json()['issues']
    if len(issues) > 0:
        create_github_issues(jira_username, jira_token, github_token, issues)
    else:
        print("No new issues found")

if __name__ == "__main__":
    cli()
