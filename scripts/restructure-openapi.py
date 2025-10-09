#!/usr/bin/env python3
"""
Script to restructure OpenAPI specification.

This script:
1. Reads the main openapi.yaml to get actual API paths
2. For each path, reads the referenced path file
3. Extracts each HTTP method from the path file
4. Creates individual operation files in api/operations/
5. Updates openapi.yaml to reference operations directly
6. Deletes the api/paths/ directory
"""

import os
import re
import yaml
import shutil
from pathlib import Path
from typing import Dict, List, Any

# Project root
ROOT = Path(__file__).parent.parent
OPENAPI_FILE = ROOT / "api" / "openapi.yaml"
PATHS_DIR = ROOT / "api" / "paths"
OPERATIONS_DIR = ROOT / "api" / "operations"

# HTTP methods we care about
HTTP_METHODS = ['get', 'post', 'put', 'patch', 'delete', 'head', 'options']


def sanitize_path_to_directory(api_path: str) -> str:
    """
    Convert API path to directory path.

    Examples:
        /auth/login -> auth/login
        /pipelines/{id} -> pipelines/id
        /users/me -> users/me
        /auth/magic-links/request -> auth/magic-links/request
    """
    # Remove leading slash
    path = api_path.lstrip('/')
    # Remove path parameter braces
    path = re.sub(r'\{([^}]+)\}', r'\1', path)
    return path


def load_yaml(file_path: Path) -> Dict[str, Any]:
    """Load a YAML file."""
    with open(file_path, 'r') as f:
        return yaml.safe_load(f)


def save_yaml(file_path: Path, data: Dict[str, Any]):
    """Save data to a YAML file."""
    with open(file_path, 'w') as f:
        yaml.dump(data, f, default_flow_style=False, sort_keys=False, allow_unicode=True)


def extract_operations(path_file: Path) -> Dict[str, Any]:
    """Extract all operations from a path file."""
    content = load_yaml(path_file)

    operations = {}
    for method in HTTP_METHODS:
        if method in content:
            operations[method] = content[method]

    return operations


def fix_refs(data: Any, depth: int) -> Any:
    """
    Recursively fix $ref paths in the operation data.

    Adjusts relative paths from api/paths/ to api/operations/<resource_path>/
    For example: ../components/... needs extra ../ for each level deep in operations/

    Args:
        data: YAML data structure (dict, list, or primitive)
        depth: How many levels deep we are in api/operations/ (e.g., auth/login = 2)

    Returns:
        Updated data with fixed refs
    """
    if isinstance(data, dict):
        result = {}
        for key, value in data.items():
            if key == '$ref' and isinstance(value, str):
                # Adjust the reference path
                # Original: ../components/... (from api/paths/)
                # New: ../../components/... (from api/operations/auth/)
                # New: ../../../components/... (from api/operations/auth/login/)
                if value.startswith('../components/'):
                    # Add extra ../ for each level of depth
                    extra_levels = '../' * depth
                    result[key] = extra_levels + value
                else:
                    result[key] = value
            else:
                result[key] = fix_refs(value, depth)
        return result
    elif isinstance(data, list):
        return [fix_refs(item, depth) for item in data]
    else:
        return data


def create_operation_file(resource_path: str, method: str, operation: Dict[str, Any]):
    """Create an operation file."""
    # Create directory structure
    operation_dir = OPERATIONS_DIR / resource_path
    operation_dir.mkdir(parents=True, exist_ok=True)

    # Calculate depth (number of extra ../ needed beyond the original ../)
    # This equals the number of directory segments in the resource path
    depth = len(resource_path.split('/'))

    # Fix $ref paths in the operation
    fixed_operation = fix_refs(operation, depth)

    # Write operation file
    operation_file = operation_dir / f"{method}.yaml"
    save_yaml(operation_file, fixed_operation)

    print(f"  Created: {operation_file.relative_to(ROOT)}")
    return f"operations/{resource_path}/{method}.yaml"


def main():
    """Main execution."""
    print(f"OpenAPI Restructuring Tool")
    print(f"=" * 60)
    print(f"OpenAPI file: {OPENAPI_FILE}")
    print(f"Paths directory: {PATHS_DIR}")
    print(f"Operations directory: {OPERATIONS_DIR}")
    print()

    # Ensure operations directory exists
    OPERATIONS_DIR.mkdir(parents=True, exist_ok=True)

    # Load the main OpenAPI spec to get actual paths
    openapi = load_yaml(OPENAPI_FILE)
    if 'paths' not in openapi:
        print("No paths found in openapi.yaml")
        return

    paths = openapi['paths']
    print(f"Found {len(paths)} paths in openapi.yaml")
    print()

    # Track the updated paths structure
    updated_paths = {}

    # Process each path
    for api_path, path_content in paths.items():
        # Skip if not a $ref (shouldn't happen in current structure)
        if not isinstance(path_content, dict) or '$ref' not in path_content:
            print(f"Skipping {api_path}: not a $ref")
            updated_paths[api_path] = path_content
            continue

        # Get the path file reference
        ref = path_content['$ref']
        path_file_name = ref.split('/')[-1]  # e.g., 'paths/auth_login.yaml' -> 'auth_login.yaml'
        path_file = PATHS_DIR / path_file_name

        if not path_file.exists():
            print(f"Warning: Path file not found: {path_file}")
            updated_paths[api_path] = path_content
            continue

        # Convert API path to directory path
        resource_path = sanitize_path_to_directory(api_path)
        print(f"Processing: {api_path} -> {resource_path}")
        print(f"  Path file: {path_file.name}")

        # Extract operations from the path file
        operations = extract_operations(path_file)

        if not operations:
            print(f"  No operations found, skipping")
            updated_paths[api_path] = path_content
            continue

        # Create operation files and build new path structure
        new_path = {}
        for method, operation in operations.items():
            operation_ref = create_operation_file(resource_path, method, operation)
            new_path[method] = {'$ref': operation_ref}

        updated_paths[api_path] = new_path
        print()

    # Update the OpenAPI spec with new paths structure
    openapi['paths'] = updated_paths

    # Save the updated OpenAPI spec
    print(f"Saving updated openapi.yaml...")
    save_yaml(OPENAPI_FILE, openapi)

    # Remove the paths directory
    if PATHS_DIR.exists():
        print(f"Removing {PATHS_DIR.relative_to(ROOT)}...")
        shutil.rmtree(PATHS_DIR)

    print(f"=" * 60)
    print(f"Restructuring complete!")
    print(f"Operations created in: {OPERATIONS_DIR.relative_to(ROOT)}")
    print(f"Paths directory removed")


if __name__ == "__main__":
    main()
