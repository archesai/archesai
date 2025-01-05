#!/bin/bash

set -euo pipefail # Exit on error, undefined variable, and fail on pipe errors

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Ensure required commands are available
REQUIRED_COMMANDS=("k3d" "helm" "kubectl" "sudo")
for cmd in "${REQUIRED_COMMANDS[@]}"; do
    if ! command_exists "$cmd"; then
        echo "❌ Error: Required command '$cmd' is not installed."
        exit 1
    fi
done

# # Install dnsmasq if not installed
# if command_exists dnsmasq; then
#     echo "✅ dnsmasq is already installed."
# else
#     echo "🔄 Installing dnsmasq..."
#     sudo apt-get update
#     sudo apt-get install -y dnsmasq
#     echo "✅ dnsmasq installed successfully."
# fi
