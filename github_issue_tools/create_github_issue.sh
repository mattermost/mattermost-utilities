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
  while IFS="," read -r title owner repo description labels milestone
  do
    # Create API request body.
    body='{"title":"'"$title"'","body":"'"$description"'"'

    # Add label in the request body, if present.
    if [ "$labels" != "" ] 
    then
      # Extract a single label into multiple labels separated by /.
      labels=$(echo "$labels" | sed 's/\//","/g') 
      body=$body',"labels":["'"$labels"'"]'
    fi

    # Add milestone in the request body, if present.
    if [ "$milestone" != "" ] 
    then
      # Extract a single label into multiple labels separated by /.
      body=$body',"milestone":'$milestone''
    fi

    # Close the request body.
    body=$body"}"

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
