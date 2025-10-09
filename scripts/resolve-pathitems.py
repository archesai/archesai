#!/usr/bin/env python3
"""
Resolve pathItems references in OpenAPI bundled spec.

This script takes the bundled OpenAPI spec and inlines all pathItems content
directly into the paths section, removing the pathItems indirection.

Before:
    paths:
      /auth/login:
        $ref: '#/components/pathItems/auth_login'
    components:
      pathItems:
        auth_login:
          post:
            operationId: Login

After:
    paths:
      /auth/login:
        post:
          operationId: Login
    components:
      # pathItems removed
"""

import sys
from pathlib import Path
import yaml


def load_yaml(file_path):
    """Load YAML file preserving order"""
    with open(file_path, 'r') as f:
        return yaml.safe_load(f)


def save_yaml(file_path, data):
    """Save YAML file with consistent formatting"""
    with open(file_path, 'w') as f:
        yaml.dump(
            data,
            f,
            default_flow_style=False,
            sort_keys=False,
            allow_unicode=True,
            width=120
        )


def resolve_pathitems(spec):
    """
    Resolve pathItems references by inlining content into paths.

    Args:
        spec: The OpenAPI specification dictionary

    Returns:
        Modified specification with pathItems resolved
    """
    # Check if pathItems exist
    if 'components' not in spec or 'pathItems' not in spec['components']:
        print("No pathItems found in spec, skipping resolution")
        return spec

    path_items = spec['components']['pathItems']
    paths = spec.get('paths', {})

    resolved_count = 0

    # Iterate through all paths and resolve references
    for path_name, path_value in paths.items():
        # Check if this path is a reference to a pathItem
        if isinstance(path_value, dict) and '$ref' in path_value:
            ref = path_value['$ref']

            # Check if it's a pathItems reference
            if ref.startswith('#/components/pathItems/'):
                # Extract the pathItem name
                path_item_name = ref.split('/')[-1]

                # Get the actual pathItem content
                if path_item_name in path_items:
                    # Replace the reference with the actual content
                    paths[path_name] = path_items[path_item_name]
                    resolved_count += 1
                else:
                    print(f"Warning: pathItem '{path_item_name}' not found for path '{path_name}'")

    # Update the paths in the spec
    spec['paths'] = paths

    # Remove the pathItems section from components
    del spec['components']['pathItems']

    # If components is now empty, we could remove it, but let's keep it for other component types

    print(f"Resolved {resolved_count} pathItems references")

    return spec


def main():
    script_dir = Path(__file__).parent
    bundled_path = script_dir / "../api/openapi.bundled.yaml"

    # Resolve path
    bundled_path = bundled_path.resolve()

    if not bundled_path.exists():
        print(f"Error: Bundled OpenAPI file not found: {bundled_path}")
        sys.exit(1)

    print(f"Loading {bundled_path}")
    spec = load_yaml(bundled_path)

    print("Resolving pathItems references...")
    spec = resolve_pathitems(spec)

    print(f"Writing resolved spec to {bundled_path}")
    save_yaml(bundled_path, spec)

    print("✓ pathItems resolution complete!")


if __name__ == "__main__":
    main()
