# GitHub Issue Tools

GitHub issue tools provide a bash script to help automatically create GtiHub issues on any plugin.

## Setup

- Create an `issue.csv` file with the below format.  

| Title | Owner | Repo | Description (Optional) | Labels (Optional) | Milestone (Optional) |
| --- | --- | --- | --- | --- | --- |  
| issue_title | repo_owner | repo_name | issue_description | issue_labels (separated by `/`) | issue_milestone |
- Create a `.env` file.
    - Run command to copy the template `.env.template` file.
    ```
    cp .env.template .env
    ```
    - Configure the `.env` file created before. Enter your `github personal access token` in the `token` field and `issue.csv` file location in the `csv_file` field.

## How to use
- Provide the executable permissions to the script by the command:
```
chmod +x create_github_issue.sh
```
- Run the script to create GitHub issues in the repository:
```
./create_github_issue.sh
```
