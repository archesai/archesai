#!/bin/bash
set -e

# ArchesAI Kustomize + Helm Deployment Script
# Combines Helm templating with Kustomize component composition

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOYMENTS_DIR="$(dirname "$SCRIPT_DIR")"

# Default values
ENVIRONMENT=${1:-dev}
NAMESPACE=${2:-archesai-${ENVIRONMENT}}
DRY_RUN=${3:-false}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

echo_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

echo_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

echo_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Validate environment
if [[ ! -f "${DEPLOYMENTS_DIR}/helm-minimal/values-${ENVIRONMENT}.yaml" ]]; then
    echo_error "Environment '${ENVIRONMENT}' not found. Available: dev, prod"
    echo_info "Usage: $0 [environment] [namespace] [dry-run]"
    echo_info "Example: $0 dev archesai-dev false"
    exit 1
fi

# Check required tools
check_tools() {
    echo_info "Checking required tools..."

    if ! command -v helm &> /dev/null; then
        echo_error "helm is required but not installed."
        exit 1
    fi

    if ! command -v kustomize &> /dev/null; then
        echo_error "kustomize is required but not installed."
        exit 1
    fi

    if ! command -v kubectl &> /dev/null; then
        echo_error "kubectl is required but not installed."
        exit 1
    fi

    echo_success "All required tools are available"
}

# Generate kustomization.yaml from Helm template
generate_kustomization() {
    echo_info "Generating kustomization.yaml from Helm template..."

    local temp_dir=$(mktemp -d)
    local values_file="${DEPLOYMENTS_DIR}/helm-minimal/values-${ENVIRONMENT}.yaml"

    # Template the kustomization.yaml
    helm template archesai \
        "${DEPLOYMENTS_DIR}/helm-minimal" \
        -f "$values_file" \
        --set namespace="$NAMESPACE" \
        > "${temp_dir}/kustomization.yaml"

    echo_success "Kustomization file generated at ${temp_dir}/kustomization.yaml"
    echo "$temp_dir"
}

# Build and apply with Kustomize
deploy_with_kustomize() {
    local temp_dir=$1

    echo_info "Building manifests with Kustomize..."

    # Build the manifests
    local manifest_file="${temp_dir}/manifests.yaml"
    kustomize build "$temp_dir" > "$manifest_file"

    echo_success "Manifests built successfully"

    # Show what will be deployed
    echo_info "Resources to be deployed:"
    kubectl api-resources --verbs=list --namespaced -o name | xargs -n 1 kubectl get --show-kind --ignore-not-found -n "$NAMESPACE" --dry-run=client -o name 2>/dev/null || true

    if [[ "$DRY_RUN" == "true" ]]; then
        echo_warning "DRY RUN MODE - showing what would be deployed:"
        echo "=========================="
        cat "$manifest_file"
        echo "=========================="
        return 0
    fi

    # Apply the manifests
    echo_info "Applying manifests to cluster..."
    kubectl apply -f "$manifest_file"

    echo_success "Deployment completed successfully!"

    # Show deployed resources
    echo_info "Deployed resources in namespace '${NAMESPACE}':"
    kubectl get all -n "$NAMESPACE" 2>/dev/null || echo_warning "No resources found or namespace doesn't exist yet"
}

# Preview function
preview_deployment() {
    echo_info "Previewing deployment for environment: ${ENVIRONMENT}"
    echo_info "Namespace: ${NAMESPACE}"
    echo ""

    local temp_dir=$(generate_kustomization)

    echo_info "Generated kustomization.yaml:"
    echo "=================================="
    cat "${temp_dir}/kustomization.yaml"
    echo "=================================="
    echo ""

    echo_info "Resulting Kubernetes manifests:"
    echo "=================================="
    kustomize build "$temp_dir"
    echo "=================================="

    # Cleanup
    rm -rf "$temp_dir"
}

# Main deployment function
deploy() {
    echo_info "🚀 Deploying ArchesAI to ${ENVIRONMENT} environment"
    echo_info "Namespace: ${NAMESPACE}"
    echo_info "Dry run: ${DRY_RUN}"
    echo ""

    # Generate templated kustomization
    local temp_dir=$(generate_kustomization)

    # Deploy with kustomize
    deploy_with_kustomize "$temp_dir"

    # Cleanup temporary files
    rm -rf "$temp_dir"

    echo ""
    echo_success "🎉 Deployment completed!"
    echo_info "Access your application:"
    echo_info "  kubectl get services -n ${NAMESPACE}"
    echo_info "  kubectl port-forward -n ${NAMESPACE} svc/archesai-api 3001:3001"
}

# Handle script arguments
case "${1:-}" in
    "preview"|"--preview"|"-p")
        check_tools
        preview_deployment
        ;;
    "help"|"--help"|"-h")
        echo "ArchesAI Deployment Script"
        echo ""
        echo "Usage:"
        echo "  $0 [environment] [namespace] [dry-run]   Deploy to environment"
        echo "  $0 preview                               Preview deployment"
        echo "  $0 help                                  Show this help"
        echo ""
        echo "Examples:"
        echo "  $0 dev                                   Deploy to dev environment"
        echo "  $0 prod archesai true                    Dry run production deployment"
        echo "  $0 preview                               Preview what will be deployed"
        ;;
    *)
        check_tools
        deploy
        ;;
esac