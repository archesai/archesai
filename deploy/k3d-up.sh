#!/bin/bash

set -euo pipefail # Exit on error, undefined variable, and fail on pipe errors

# -----------------------------
# Constants and Configurations
# -----------------------------
CLUSTER_NAME="tower"
K3D_CONFIG_FILE="deploy/k3d-config.yaml"
NAMESPACE="default"
TLS_SECRET="archesai-tls"
INGRESS_NGINX_NAMESPACE="ingress-nginx"
INGRESS_NGINX_REPO="https://kubernetes.github.io/ingress-nginx"
INGRESS_NGINX_NAME="ingress-nginx"
JETSTACK_REPO="https://charts.jetstack.io"
CERT_MANAGER_NAME="cert-manager"
CERT_MANAGER_NAMESPACE="cert-manager"
CERT_MANAGER_VERSION="v1.16.2"

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

# Function to install Helm repository if not already added
add_helm_repo() {
    local repo_name="$1"
    local repo_url="$2"

    if helm repo list | grep -qw "$repo_name"; then
        echo "✅ Helm repository '$repo_name' already exists. Skipping."
    else
        echo "🔄 Adding Helm repository '$repo_name'..."
        helm repo add "$repo_name" "$repo_url"
        echo "✅ Helm repository '$repo_name' added successfully."
    fi
}

# Function to update Helm repositories
update_helm_repos() {
    echo "🔄 Updating Helm repositories..."
    helm repo update
    echo "✅ Helm repositories updated."
}

# Function to install or upgrade a Helm chart
helm_install_or_upgrade() {
    local release_name="$1"
    local chart="$2"
    local namespace="$3"
    local version="${4:-}"
    shift 4 || true
    local additional_args=("$@")

    if helm ls -n "$namespace" | grep -qw "$release_name"; then
        echo "✅ Helm release '$release_name' already exists in namespace '$namespace'. Upgrading..."
        helm upgrade "$release_name" "$chart" \
            --namespace "$namespace" \
            --version "$version" \
            "${additional_args[@]}"
        echo "✅ Helm release '$release_name' upgraded successfully."
    else
        echo "🔄 Installing Helm release '$release_name' in namespace '$namespace'..."
        helm install "$release_name" "$chart" \
            --namespace "$namespace" \
            --create-namespace \
            --version "$version" \
            "${additional_args[@]}"
        echo "✅ Helm release '$release_name' installed successfully."
    fi
}

# -----------------------------
# Main Setup Functions
# -----------------------------

install_ingress_nginx() {
    helm_install_or_upgrade \
        "$INGRESS_NGINX_NAME" \
        "ingress-nginx" \
        "$INGRESS_NGINX_NAMESPACE" \
        "" \
        --repo "$INGRESS_NGINX_REPO"
}

install_cert_manager() {
    add_helm_repo "jetstack" "$JETSTACK_REPO"
    update_helm_repos

    helm_install_or_upgrade \
        "$CERT_MANAGER_NAME" \
        "jetstack/cert-manager" \
        "$CERT_MANAGER_NAMESPACE" \
        "$CERT_MANAGER_VERSION" \
        --set crds.enabled=true
}

# -----------------------------
# Execute Setup
# -----------------------------

main() {
    # Ensure required commands are available
    for cmd in k3d helm kubectl sudo; do
        if ! command_exists "$cmd"; then
            echo "❌ Required command '$cmd' is not installed. Please install it before proceeding."
            exit 1
        fi
    done

    create_k3d_cluster

    # Extract TLS certificate and key
    source ./deploy/extract_tls.sh

    install_ingress_nginx
    install_cert_manager

    echo "🎉 Setup complete."
}

# Run the main function
main
