#!/usr/bin/env bash
set -euo pipefail

# Function to log messages
log() {
    echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')]: $*"
}

# Determine the directory where the script resides
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Define paths for environment files and Redis CA certificate
PROD_ENV_FILE="$SCRIPT_DIR/../kubernetes/overlays/gcp/production/.env.production"
STAGE_ENV_FILE="$SCRIPT_DIR/../kubernetes/overlays/gcp/stage/.env.stage"
REDIS_CA_FILE="$SCRIPT_DIR/../kubernetes/overlays/gcp/base/redis-ca.pem"

# Define the path to the .env.template file
TEMPLATE_FILE="$SCRIPT_DIR/../.env.template"

# Clear existing content in the target files
>"$PROD_ENV_FILE"
>"$STAGE_ENV_FILE"
>"$REDIS_CA_FILE"

log "Processing template file: $TEMPLATE_FILE"

# Function to replace placeholders with secrets
replace_placeholders() {
    local input_file="$1"
    local output_file="$2"

    while IFS= read -r line || [[ -n "$line" ]]; do
        # Check if the line contains any placeholders
        if [[ "$line" =~ \{\{([^}:]+):([^}]+)\}\} ]]; then
            # Extract all placeholders in the line
            while [[ "$line" =~ \{\{([^}:]+):([^}]+)\}\} ]]; do
                full_placeholder="${BASH_REMATCH[0]}"
                secret_key="${BASH_REMATCH[1]}"
                secret_version="${BASH_REMATCH[2]}"

                log "Fetching secret: $secret_key, version: $secret_version"

                # Fetch the secret value
                secret_value=$(gcloud secrets versions access "$secret_version" --secret="$secret_key" 2>/dev/null) || {
                    log "Error: Failed to access secret '$secret_key' with version '$secret_version'"
                    exit 1
                }

                # If the line contains DATABASE_URL, escape ampersands
                if [[ "$line" == *"DATABASE_URL"* ]]; then
                    secret_value=$(printf '%s' "$secret_value" | sed 's/&/\\&/g')
                    log "Escaped ampersands in DATABASE_URL"
                fi

                # Replace the placeholder with the secret value
                line="${line//"$full_placeholder"/"$secret_value"}"
                log "Replaced placeholder '$full_placeholder' with secret value"
            done
        fi
        # Append the processed line to the output file
        echo "$line" >>"$output_file"
    done <"$input_file"
}

# Replace placeholders in the production environment file
replace_placeholders "$TEMPLATE_FILE" "$PROD_ENV_FILE"

log "Copied production environment file to staging"
cp -f "$PROD_ENV_FILE" "$STAGE_ENV_FILE"

log "Fetching stage-specific secrets"

# Fetch all secrets prefixed with STAGE_
stage_secrets=$(gcloud secrets list --project=archesai --filter="name:STAGE_*" --format="value(name)") || {
    log "Error: Failed to list stage secrets"
    exit 1
}

for secret_full_name in $stage_secrets; do
    # Extract the actual secret key by removing the STAGE_ prefix
    secret_key="${secret_full_name#STAGE_}"
    log "Fetching stage secret: $secret_full_name"

    # Fetch the latest version of the secret
    secret_value=$(gcloud secrets versions access latest --secret="$secret_full_name" --project=archesai 2>/dev/null) || {
        log "Error: Failed to access stage secret '$secret_full_name'"
        exit 1
    }

    # Remove any existing entry for the secret in the staging environment file
    sed -i "/^$secret_key=/d" "$STAGE_ENV_FILE" || {
        log "Warning: Failed to remove existing entry for '$secret_key' in staging environment file"
    }

    # Append the secret to the staging environment file
    echo "$secret_key=$secret_value" >>"$STAGE_ENV_FILE"
    log "Added stage secret '$secret_key' to staging environment file"
done

log "Fetching PROD_REDIS_CA secret"
# Fetch the PROD_REDIS_CA secret and write it to the Redis CA file
gcloud secrets versions access latest --secret=PROD_REDIS_CA --project=archesai >"$REDIS_CA_FILE" || {
    log "Error: Failed to access PROD_REDIS_CA secret"
    exit 1
}
log "Written Redis CA certificate to $REDIS_CA_FILE"

log "Environment files generated successfully."
