#!/usr/bin/env python3
"""
Python runner for container executor.
Reads JSON from stdin, validates against schemas, executes function, and writes result to stdout.
"""

import sys
import json
import traceback
from typing import Any, Dict
from jsonschema import validate, ValidationError

# Dynamically import the execute module
# This can be mounted at runtime or copied in derived images
from execute import execute_function


def main():
    """Main entry point for the runner."""
    try:
        # Read input from stdin
        raw_input = sys.stdin.read()
        if not raw_input:
            raise ValueError("No input provided")

        # Parse the request
        request = json.loads(raw_input)

        # Extract components
        schema_in = request.get("schema_in")
        schema_out = request.get("schema_out")
        input_data = request.get("input")

        if input_data is None:
            raise ValueError("Missing 'input' field in request")

        # Validate input against schema if provided
        if schema_in:
            try:
                validate(instance=input_data, schema=schema_in)
            except ValidationError as e:
                raise ValueError(f"Input validation failed: {e.message}")

        # Execute the function
        output_data = execute_function(input_data)

        # Validate output against schema if provided
        if schema_out:
            try:
                validate(instance=output_data, schema=schema_out)
            except ValidationError as e:
                raise ValueError(f"Output validation failed: {e.message}")

        # Return success response
        response = {
            "ok": True,
            "output": output_data
        }
        print(json.dumps(response))
        sys.exit(0)

    except Exception as e:
        # Return error response
        response = {
            "ok": False,
            "error": {
                "message": str(e),
                "details": traceback.format_exc() if sys.stderr.isatty() else None
            }
        }
        print(json.dumps(response))
        sys.exit(0)  # Exit with 0 even on error to ensure JSON response is returned


if __name__ == "__main__":
    main()