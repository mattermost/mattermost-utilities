import requests
from requests.auth import HTTPBasicAuth
from github import Github
from jira_to_markdown import jira_to_markdown

def create_github_issues(jira_username, jira_token, github_token, repo, labels, issues, dry_run):
    footer = '''
----
If you're interested please comment here and come [join our "Contributors" community channel](https://community.mattermost.com/core/channels/tickets) on our daily build server, where you can discuss questions with community members and the Mattermost core team. For technical advice or questions, please  [join our "Developers" community channel](https://community.mattermost.com/core/channels/developers).

New contributors please see our [Developer's Guide](https://developers.mattermost.com/contribute/getting-started/).

JIRA: https://mattermost.atlassian.net/browse/{{TICKET}}
    '''

    g = Github(github_token)
    r = g.get_repo(repo)
    final_labels = []
    for label in r.get_labels():
        if label.name in labels:
            final_labels.append(label)

    if dry_run:
        print('We haven\'t created the github ticket because --dry-run flag was detected. Tickets information:')

    for issue in issues:
        title = issue['fields']['summary']
        key = issue['key']
        markdown_description = jira_to_markdown(issue['fields']['description'])
        description = markdown_description + "\n\n" + footer.replace("{{TICKET}}", key)

        if dry_run:
            print('------\n{}\n{}\n\n{}'.format(title, "="*len(title), description))
            continue

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

