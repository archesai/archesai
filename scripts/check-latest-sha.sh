#!/usr/bin/env bash
set -e 

apt-get -y update && apt-get install -y jq

# Function to get the latest short SHA and check CI status
check_repo_ci_status() {
    local REPO_NAME=$1
    local GH_PAT=$2

    # Get the latest commit full SHA using GitHub API
    local full_sha=$(curl -s -H "Authorization: token $GH_PAT" "https://api.github.com/repos/archesai/${REPO_NAME}/commits" | jq -r '.[0].sha')
    # Extract the first 7 characters to get the short SHA
    local short_sha=${full_sha:0:7}
    echo "${REPO_NAME} Latest Short SHA: $short_sha"

    # Check CI status
    local response=$(curl -s -H "Authorization: token $GH_PAT" "https://api.github.com/repos/archesai/${REPO_NAME}/commits/$full_sha/check-runs")
    local all_checks_successful=$(echo "$response" | jq 'all(.check_runs[]; .status == "completed" and .conclusion == "success")')

    if [ "$all_checks_successful" = "false" ]; then
        echo "Some ${REPO_NAME} checks failed or are still running."
        exit 1
    else
        echo "All ${REPO_NAME} checks passed."
    fi

    # Export the short SHA to a file
    echo "export ${REPO_NAME}_sha=$short_sha" >> /workspace/tags.sh

}

> /workspace/tags.sh

# Check each repository
check_repo_ci_status "ui" $GH_PAT
check_repo_ci_status "api" $GH_PAT
check_repo_ci_status "nlp" $GH_PAT
