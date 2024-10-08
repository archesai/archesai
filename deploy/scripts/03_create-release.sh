#!/usr/bin/env bash
set -e
set -o pipefail

# Load image shas from the file
source /workspace/values.sh
cat /workspace/values.sh
echo "Creating release ${RELEASE_NAME}"

# Set version in upgrade job FIXME - make this better somehow
SCRIPT_DIR=$(dirname $(realpath "$0"))
sed -i "s/nameSuffix: -v[0-9]*\.[0-9]*\.[0-9]*/nameSuffix: -${RELEASE_NAME}/" $SCRIPT_DIR/../kubernetes/base/upgrade/kustomization.yaml

# Create release in cloud dpeloy
RELEASE_NAME=$(echo "${RELEASE_NAME}" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9-]/-/g')
gcloud components install gke-gcloud-auth-plugin
gcloud container clusters get-credentials archesai-cluster --zone us-east1-c --project archesai
gcloud deploy releases create "$RELEASE_NAME" \
    --project=archesai \
    --region=us-east1 \
    --delivery-pipeline=arches-deployment \
    --images=us-east4-docker.pkg.dev/archesai/images/arches-api=us-east4-docker.pkg.dev/archesai/images/arches-api:${SHORT_SHA},us-east4-docker.pkg.dev/archesai/images/arches-ui=us-east4-docker.pkg.dev/archesai/images/arches-ui:${SHORT_SHA},us-east4-docker.pkg.dev/archesai/images/arches-nlp=us-east4-docker.pkg.dev/archesai/images/arches-nlp:${SHORT_SHA} \
    --skaffold-file=./deploy/skaffold.yaml
