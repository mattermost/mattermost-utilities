import click
import requests
import pprint
from requests.auth import HTTPBasicAuth
from github import Github

from utils import create_github_issues

footer = '''
----
If you're interested please comment here and come [join our "Contributors" community channel](https://community.mattermost.com/core/channels/tickets) on our daily build server, where you can discuss questions with community members and the Mattermost core team. For technical advice or questions, please  [join our "Developers" community channel](https://community.mattermost.com/core/channels/developers).

New contributors please see our [Developer's Guide](https://developers.mattermost.com/contribute/getting-started/).

JIRA: https://mattermost.atlassian.net/browse/MM-{{TICKET}}
'''

@click.command()
@click.option('--jira-token', '-j', prompt='Your Jira access token', help='The token used to authenticate the user against Jira.')
@click.option('--jira-username', '-u', prompt='Your Jira username', help='Username of the user to get the ticket information.')
@click.option('--github-token', '-g', prompt='Your Github access token', help='The token used to authenticate the user against Github.')
@click.option('--repo', '-r', prompt='Repository', help='The repository which contains the issues. E.g. mattermost/mattermost-server')
@click.option('--labels', '-l', prompt='Labels', help='The labels to set to the issues', multiple=True)
@click.option('--dry-run/--no-dry-run', help='Skip actually creating any tickets', default=False)
@click.option('--debug/--no-debug', help='Dump debugging information.', default=False)
@click.argument('issue-numbers', nargs=-1)
def cli(jira_token, jira_username, github_token, repo, labels, dry_run, debug, issue_numbers):
    if len(issue_numbers) < 1:
        print("You need to pass at least one issue number")
        return

    query = " OR ".join(map(lambda x: "key = MM-{}".format(x), issue_numbers))
    data = {
            "jql":"project = MM AND {}".format(query),
            "maxResults": len(issue_numbers),
            "fields": ["summary", "description"],
    }
    resp = requests.post(
        "https://mattermost.atlassian.net/rest/api/2/search",
        json=data,
        auth=HTTPBasicAuth(jira_username, jira_token)
    )
    if debug:
        pprint.pprint(resp.json())

    issues = resp.json()['issues']
    print(create_github_issues(jira_username, jira_token, github_token, repo, labels, issues, dry_run))

if __name__ == "__main__":
    cli()
