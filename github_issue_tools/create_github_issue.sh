#!/bin/bash

# Load the environment variables from the .env file.
if [[ -f ".env" ]]; then
    source .env
else
    echo -e "Error: The '.env' file does not exist.\nYou can create one using the '.env.template' file."
    exit 1
fi

echo "Starting script"

# Read the CSV file and extract the column values.
{
  # Read the header of the CSV file.
  read
  while IFS="," read -r title owner repo description_location labels milestone
  do
    # Extract issue description from the description file location.
    description=$(cat "$description_location")

    # Extract a single label into multiple labels separated by |.
    [[ "$labels" != "" ]] && labels=['"'$(echo "$labels" | sed 's/|/","/g')'"'] || labels=[]

    # Extract the milestone value.
    [[ "$milestone" != "" ]] && milestone=$milestone || milestone=null

    # Create API request body.
    body=$( jq -n \
               --arg title "$title" \
               --arg body "$description" \
               --argjson labels "$labels" \
               --argjson milestone "$milestone" \
               '{title: $title, body: $body, labels: $labels, milestone: $milestone}' )

    # Make an API request to create GitHub issues.
    response=$(curl -L \
    -X POST \
    -H "Accept: application/vnd.github+json" \
    -H "Authorization: Bearer $token" \
    -H "X-GitHub-Api-Version: 2022-11-28" \
    https://api.github.com/repos/$owner/$repo/issues \
    -d "$body")
    
    # Use jq to extract the error message from the response.
    error_message=$(echo $response | jq -r '.message')

    # Handle responses from the API.
    if [ "$error_message" != "null" ] 
    then
      echo "Unable to create issue for owner: $owner, repo: $repo, title: $title, error: $error_message"
    else
      echo "Issue created for owner: $owner, repo: $repo, title: $title"
    fi

  done
} < $csv_file

echo "Finishing script"
