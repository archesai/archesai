#!/bin/bash

set -euo pipefail # Exit on error, undefined variable, and fail on pipe errors

# -----------------------------
# Constants and Configurations
# -----------------------------
CLUSTER_NAME="tower"
K3D_CONFIG_FILE="deploy/k3d-config.yaml"

# -----------------------------
# Utility Functions
# -----------------------------

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to create k3d cluster if it doesn't exist
create_k3d_cluster() {
    if k3d cluster list | grep -qw "$CLUSTER_NAME"; then
        echo "✅ Cluster '$CLUSTER_NAME' already exists. Skipping creation."
    else
        echo "🔄 Creating cluster '$CLUSTER_NAME'..."
        k3d cluster create "$CLUSTER_NAME" --config "$K3D_CONFIG_FILE"
        echo "✅ Cluster '$CLUSTER_NAME' created successfully."
    fi
}

# -----------------------------
# Execute Setup
# -----------------------------

main() {
    for cmd in k3d helm kubectl sudo; do
        if ! command_exists "$cmd"; then
            echo "❌ Required command '$cmd' is not installed. Please install it before proceeding."
            exit 1
        fi
    done

    create_k3d_cluster

    source ./deploy/extract_tls.sh

    echo "🎉 Setup complete."
}

main
