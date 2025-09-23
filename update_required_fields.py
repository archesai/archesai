#!/usr/bin/env python3
import yaml
import re
from pathlib import Path

def update_schema_file(filepath):
    with open(filepath, 'r') as f:
        content = f.read()

    # Parse YAML
    data = yaml.safe_load(content)

    if not isinstance(data, dict) or 'properties' not in data:
        return False

    # Find properties with defaults
    props_with_defaults = []
    for prop_name, prop_data in data.get('properties', {}).items():
        if isinstance(prop_data, dict) and 'default' in prop_data:
            props_with_defaults.append(prop_name)

    if not props_with_defaults:
        print(f"  No properties with defaults in {filepath}")
        return False

    print(f"  Found properties with defaults: {props_with_defaults}")

    # Get existing required fields
    existing_required = data.get('required', [])

    # Combine and deduplicate
    all_required = sorted(list(set(existing_required + props_with_defaults)))

    # Update the data
    data['required'] = all_required

    # Write back
    with open(filepath, 'w') as f:
        yaml.dump(data, f, default_flow_style=False, sort_keys=False, allow_unicode=True)

    print(f"  Updated {filepath} - required fields: {all_required}")
    return True

# Process all Config*.yaml files
config_files = Path('api/components/schemas').glob('Config*.yaml')
for config_file in sorted(config_files):
    print(f"Processing {config_file}")
    update_schema_file(config_file)