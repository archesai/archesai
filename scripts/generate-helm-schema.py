#!/usr/bin/env python3

"""
Generate values.schema.json for Helm from Config.yaml
This script properly resolves all $ref references recursively
"""

import json
import os
import sys
from pathlib import Path
import yaml


def load_yaml(file_path):
    """Load YAML file"""
    with open(file_path, 'r') as f:
        return yaml.safe_load(f)


def resolve_refs(obj, schemas_dir, visited=None):
    """Recursively resolve $ref references"""
    if visited is None:
        visited = set()

    if isinstance(obj, dict):
        if '$ref' in obj:
            ref_path = obj['$ref']
            if ref_path.startswith('./'):
                filename = ref_path[2:]  # Remove ./

                if filename in visited:
                    return obj

                visited.add(filename)
                ref_file_path = schemas_dir / filename

                if ref_file_path.exists():
                    referenced_schema = load_yaml(ref_file_path)
                    resolved = resolve_refs(referenced_schema, schemas_dir, visited.copy())
                    visited.remove(filename)
                    return resolved
                else:
                    print(f"Error: Referenced file not found: {filename}")
                    sys.exit(1)
            else:
                return obj
        else:
            # Process all values in the dictionary
            result = {}
            for key, value in obj.items():
                result[key] = resolve_refs(value, schemas_dir, visited)
            return result

    elif isinstance(obj, list):
        return [resolve_refs(item, schemas_dir, visited) for item in obj]

    else:
        return obj


def main():
    script_dir = Path(__file__).parent
    schemas_dir = script_dir / "../api/components/schemas"
    output_path = script_dir / "../deployments/helm-minimal/values.schema.json"

    # Resolve paths
    schemas_dir = schemas_dir.resolve()
    output_path = output_path.resolve()

    if not schemas_dir.exists():
        print(f"Error: Schemas directory not found: {schemas_dir}")
        sys.exit(1)

    arches_config_path = schemas_dir / "Config.yaml"
    if not arches_config_path.exists():
        print(f"Error: Config.yaml not found: {arches_config_path}")
        sys.exit(1)

    arches_config = load_yaml(arches_config_path)
    resolved_config = resolve_refs(arches_config, schemas_dir)

    helm_schema = {
        "$schema": "https://json-schema.org/draft/2020-12/schema",
        "title": "Arches Helm Values Schema",
        "description": "JSON Schema for Arches Helm chart values",
        **resolved_config
    }

    output_path.parent.mkdir(parents=True, exist_ok=True)

    with open(output_path, 'w') as f:
        json.dump(helm_schema, f, indent=2)


if __name__ == "__main__":
    main()