#!/usr/bin/env bash
set -e

# Get the directory of the script
SCRIPT_DIR=$(dirname $(realpath "$0"))

# Install Dependencies
apt-get -y update && apt-get install -y curl wget

export GITHUB_API_URL="https://api.github.com"

export API_SHA=$(curl -H "Authorization: token $GITHUB_TOKEN" "$GITHUB_API_URL/repos/archesai/api/commits" | jq -r '.[0].sha')
echo "API SHA: $API_SHA"

export NLP_SHA=$(curl -H "Authorization: token $GITHUB_TOKEN" "$GITHUB_API_URL/repos/archesai/nlp/commits" | jq -r '.[0].sha')
echo "NLP SHA: $NLP_SHA"

export UI_SHA=$(curl -H "Authorization: token $GITHUB_TOKEN" "$GITHUB_API_URL/repos/archesai/ui/commits" | jq -r '.[0].sha')
echo "UI SHA: $UI_SHA"

wget -O /workspace/yq https://github.com/mikefarah/yq/releases/download/v4.13.2/yq_linux_amd64
chmod +x /workspace/yq
/workspace/yq eval '.services.arches-api.image = "us-east4-docker.pkg.dev/archesai/images/arches-api:$API_SHA"' -i $SCRIPT_DIR/../docker-compose.yaml
/workspace/yq eval '.services.arches-nlp.image = "us-east4-docker.pkg.dev/archesai/images/arches-nlp:$NLP_SHA"' -i $SCRIPT_DIR/../docker-compose.yaml
/workspace/yq eval '.services.arches-ui.image = "us-east4-docker.pkg.dev/archesai/images/arches-ui:$UI_SHA"' -i $SCRIPT_DIR/../docker-compose.yaml
yes | cp -f $SCRIPT_DIR/../.env.minimal.example $SCRIPT_DIR/../.env.minimal
