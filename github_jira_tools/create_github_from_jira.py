import click
import requests
from requests.auth import HTTPBasicAuth
from github import Github

header = '''
If you're interested please comment here and come [join our "Contributors" community channel](https://community.mattermost.com/core/channels/tickets) on our daily build server, where you can discuss questions with community members and the Mattermost core team. For technical advice or questions, please  [join our "Developers" community channel](https://community.mattermost.com/core/channels/developers).

New contributors please see our [Developer's Guide](https://developers.mattermost.com/contribute/getting-started/).

----

**Notes**: [Jira ticket](https://mattermost.atlassian.net/browse/MM-{{TICKET}})
'''

@click.command()
@click.option('--jira-token', '-j', prompt='Your Jira access token', help='The token used to authenticate the user against Jira.')
@click.option('--jira-username', '-u', prompt='Your Jira username', help='Username of the user to get the ticket information.')
@click.option('--github-token', '-g', prompt='Your Github access token', help='The token used to authenticate the user against Github.')
@click.option('--repo', '-r', prompt='Repository', help='The repository which contains the issues. E.g. mattermost/mattermost-server')
@click.option('--labels', '-l', prompt='Labels', help='The labels to set to the issues', multiple=True)
@click.argument('issue-numbers', nargs=-1)
def cli(jira_token, jira_username, github_token, repo, labels, issue_numbers):
    if len(issue_numbers) < 1:
        print("You need to pass at least one issue number")
        return

    g = Github(github_token)
    r = g.get_repo(repo)
    final_labels = []
    for label in r.get_labels():
        if label.name in labels:
            final_labels.append(label)

    for issue_number in issue_numbers:
        resp = requests.get(
            "https://mattermost.atlassian.net/rest/api/3/issue/MM-"+issue_number,
            auth=HTTPBasicAuth(jira_username, jira_token)
        )
        data = resp.json()

        title = data['fields']['summary']
        description = header.replace("{{TICKET}}", issue_number) + "\n" + "\n\n".join(map(lambda c: c['content'][0]['text'], data['fields']['description']['content']))

        try:
            new_issue = r.create_issue(
                title=title,
                body=description,
                labels=final_labels,
            )
        except Exception as e:
            print("Unable to create issue for jira issue {}. error: {}".format(issue_number, e))
            return

        try:
            resp = requests.put(
                "https://mattermost.atlassian.net/rest/api/3/issue/MM-"+issue_number,
                json={
                    "fields": {
                        "customfield_11106": new_issue.html_url,
                        "fixVersions": [{"name": "Help Wanted"}]
                    },
                },
                auth=HTTPBasicAuth(jira_username, jira_token)
            )
        except Exception as e:
            print("Unable to update jira issue {}. error: {}".format(issue_number, e))
            return

if __name__ == "__main__":
    cli()
