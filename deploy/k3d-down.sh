#!/bin/bash

set -euo pipefail # Exit on error, undefined variable, and fail on pipe errors

# Constants and Configurations
CLUSTER_NAME="tower"

# Ensure required commands are available. run the shell file in the same directory
source ./deploy/commands.sh

echo "🔄 Starting teardown process..."

# Delete the k3d cluster if it exists
if k3d cluster list | grep -qw "$CLUSTER_NAME"; then
    echo "🔄 Deleting k3d cluster '$CLUSTER_NAME'..."
    k3d cluster delete "$CLUSTER_NAME"
    echo "✅ k3d cluster '$CLUSTER_NAME' deleted."
else
    echo "✅ k3d cluster '$CLUSTER_NAME' does not exist. Skipping deletion."
fi

echo "🎉 Teardown complete."
