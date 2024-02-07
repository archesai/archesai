#!/usr/bin/env bash
set -e
set -o pipefail

# Set variables
REPO_OWNER=archesai
REPO_NAME=archesai
GITHUB_API_URL="https://api.github.com"
RELEASE_NAME_PREFIX="Release"
GITHUB_TOKEN=$GITHUB_TOKEN

# Function to create a new tag
create_new_tag() {
    # Get the new tag name
    local previous_tag=$(
        curl -H "Authorization: token $GITHUB_TOKEN" "$GITHUB_API_URL/repos/$REPO_OWNER/$REPO_NAME/tags" | jq -r '.[0].name'
    )
    echo "Previous tag: $previous_tag"
    local base_version=${previous_tag%-rc.*}
    local rc_number=${previous_tag##*-rc.}
    rc_number=$((rc_number + 1))
    local new_tag=$(echo "${base_version}-rc.${rc_number}")
    echo "New tag: $new_tag"

    # Get the latest commit SHA
    local commit_sha=$(curl -H "Authorization: token $GITHUB_TOKEN" "$GITHUB_API_URL/repos/$REPO_OWNER/$REPO_NAME/commits" | jq -r '.[0].sha')
    echo "Commit SHA: $commit_sha"

    # Read the file content, escape it for JSON, and remove the surrounding quotes added by jq
    escaped_values=$(jq -Rs . /workspace/values.sh | sed 's/^"\(.*\)"$/\1/')
    echo "Escaped values: $escaped_values"

    local tag_object_response=$(curl -X POST -H "Authorization: token $GITHUB_TOKEN" -H "Content-Type: application/json" \
        -d "{\"tag\": \"$new_tag\", \"message\": \"$RELEASE_NAME_PREFIX $new_tag\n\", \"object\": \"$commit_sha\", \"type\": \"commit\"}" "$GITHUB_API_URL/repos/$REPO_OWNER/$REPO_NAME/git/tags")
    echo "Tag Object Response: $tag_object_response"
    local tag_sha=$(echo "$tag_object_response" | jq -r '.sha')
    echo "Tag SHA: $tag_sha"

    # Create the tag reference
    curl -X POST -H "Authorization: token $GITHUB_TOKEN" \
        -H "Content-Type: application/json" \
        -d "{\"ref\": \"refs/tags/$new_tag\", \"sha\": \"$tag_sha\"}" \
        "$GITHUB_API_URL/repos/$REPO_OWNER/$REPO_NAME/git/refs"
}

# Function to get the latest short SHA and check CI status
get_latest_short_sha() {
    set -e
    set -o pipefail
    local repo_name=$1

    # Get the latest commit full SHA using GitHub API
    local full_sha=$(curl -s -H "Authorization: token $GITHUB_TOKEN" "https://api.github.com/repos/archesai/${repo_name}/commits" | jq -r '.[0].sha')
    # Extract the first 7 characters to get the short SHA
    local short_sha=${full_sha:0:7}
    echo "${repo_name} Latest Short SHA: $short_sha"

    # Check CI status
    local response=$(curl -s -H "Authorization: token $GITHUB_TOKEN" "https://api.github.com/repos/archesai/${repo_name}/commits/$full_sha/check-runs")
    local all_checks_successful=$(echo "$response" | jq 'all(.check_runs[]; .status == "completed" and .conclusion == "success")')

    if [ "$all_checks_successful" = "false" ]; then
        if [ "$repo_name" = "archesai" ]; then
            echo "${repo_name} was unsuccessful, but we are continuing."
        else
            echo "Some ${repo_name} checks failed or are still running."
            exit 1
        fi
    else
        echo "All ${repo_name} checks passed."
    fi

    # Export the short SHA to a file
    echo "export ${repo_name^^}_SHA=$short_sha" >>/workspace/values.sh

}

# Install JQ
apt-get -y update && apt-get install -y jq

# Create file to hold tag values
>/workspace/values.sh

# Check each repository and get latest tag
get_latest_short_sha "ui"
get_latest_short_sha "api"
get_latest_short_sha "nlp"

# Load values into the environment
source /workspace/values.sh

# Create new tag
create_new_tag

echo "New tag $new_tag created."

echo "RELEASE_NAME=$new_tag" >>/workspace/values.sh
