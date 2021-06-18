#!/bin/bash
#
# Script to substitute variables in template files.

set -e

git_url="$(git config --get remote.origin.url)"

# Get project owner and project name from Git url.
if [[ "$git_url" == "http"* ]]; then
    project_owner="$(echo "$git_url" | awk -F'[:/.]' '{ print $6 }')"
    project_name="$(echo "$git_url" | awk -F'[:/.]' '{ print $7 }')"
    git_url="${git_url%.git}"
else
    project_owner="$(echo "$git_url" | awk -F'[:/.]' '{ print $3 }')"
    project_name="$(echo "$git_url" | awk -F'[:/.]' '{ print $4 }')"
    git_url="https://github.com/$project_owner/$project_name"
fi

# Prompt user for name of license.
read -rp "License Name: " license_name

export \
    PROJECT_OWNER="$project_owner" \
    PROJECT_NAME="$project_name" \
    GIT_URL="$git_url" \
    LICENSE_NAME="$license_name"

files=(".goreleaser.yml" "Dockerfile")
for file in "${files[@]}"; do
    contents="$(envsubst < "$file")"
    echo "$contents" > "$file"
done