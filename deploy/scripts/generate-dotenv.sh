#!/usr/bin/env bash
set -e

SCRIPT_DIR=$(dirname $(realpath "$0"))

prod_env_file="$SCRIPT_DIR/../kubernetes/overlays/saas/production/.env.production"
stage_env_file="$SCRIPT_DIR/../kubernetes/overlays/saas/stage/.env.stage"
redis_ca_file="$SCRIPT_DIR/../kubernetes/overlays/saas/base/redis-ca.pem"

>"$prod_env_file"
>"$stage_env_file"
>"$redis_ca_file"

while IFS= read -r line || [[ -n "$line" ]]; do
    if [[ "$line" == *{{*}}* ]]; then
        placeholder=$(echo $line | grep -oP '{{\K[^}]+')
        secret_key=$(echo $placeholder | cut -d ':' -f 1)
        secret_version=$(echo $placeholder | cut -d ':' -f 2)
        secret_value=$(gcloud secrets versions access $secret_version --secret="$secret_key")
        if [[ "$line" == *"DATABASE_URL"* ]]; then
            secret_value=$(echo "$secret_value" | sed 's/&/\&/g')
        fi
        line=${line//\{\{$placeholder\}\}/$secret_value}
    fi
    echo "$line" >>"$prod_env_file"
done <$SCRIPT_DIR/../.env.template

yes | cp -f "$prod_env_file" "$stage_env_file"
secrets=$(gcloud secrets list --project=archesai --filter="name:STAGE_" --format="value(name)")
for secret_name in $secrets; do
    secret_value=$(gcloud secrets versions access latest --secret="$secret_name" --project=archesai)
    trimmed_key_name=${secret_name#STAGE_}
    sed -i "/^$trimmed_key_name=/d" "$stage_env_file"
    echo "$trimmed_key_name=$secret_value" >>"$stage_env_file"
done

echo "$(gcloud secrets versions access latest --secret=PROD_REDIS_CA --project=archesai)" >"$redis_ca_file"
