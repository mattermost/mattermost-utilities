import click
from github import Github

@click.command()
@click.option('--token', '-t', prompt='Your access token', help='The token used to authenticate the user.')
@click.option('--repo', '-r', prompt='Repository', help='The repository which contains the issues. E.g. mattermost/mattermost-server')
@click.option('--labels', '-l', prompt='Labels', help='The labels to set to the issues', multiple=True)
@click.argument('issue-numbers', nargs=-1)
def cli(token, repo, labels, issue_numbers):
    if len(issue_numbers) < 1:
        print("You need to pass at least one issue number")
        return

    g = Github(token)
    r = g.get_repo(repo)
    final_labels = []
    for label in r.get_labels():
        if label.name in labels:
            final_labels.append(label)

    for issue_number in issue_numbers:
        try:
            issue = r.get_issue(int(issue_number))
            issue.set_labels(*final_labels)
        except Exception as e:
            print("Unable to update issue {}. error: {}".format(issue_number, e))

if __name__ == "__main__":
    cli()
