"""
Example custom execute function

This file shows how to create a custom execute module for the Python runner.

To use:
1. Copy this file and implement your logic
2. Mount it to the container:
   docker run -i --rm -v ./my_execute.py:/app/execute.py archesai/runner-python:latest

Or build a custom image:
1. Create a Dockerfile:
   FROM archesai/runner-python:latest
   RUN pip install --no-cache-dir your-dependencies
   COPY execute.py ./execute.py
2. Build: docker build -t my-generator:latest .
"""

from typing import Any, Dict


def execute_function(input_data: Dict[str, Any]) -> Dict[str, Any]:
    """
    Execute function - implement your custom logic here.

    Args:
        input_data: The input data from the request

    Returns:
        The output data

    Raises:
        ValueError: If execution fails
    """
    # Example: Simple data transformation
    if "values" not in input_data or not isinstance(input_data["values"], list):
        raise ValueError("Expected 'values' array in input")

    values = input_data["values"]
    total = sum(values)
    mean = total / len(values)

    return {
        "count": len(values),
        "sum": total,
        "mean": mean,
        "min": min(values),
        "max": max(values),
    }


"""
Example usage:

Input:
{
  "input": {
    "values": [1, 2, 3, 4, 5]
  }
}

Output:
{
  "ok": true,
  "output": {
    "count": 5,
    "sum": 15,
    "mean": 3.0,
    "min": 1,
    "max": 5
  }
}
"""
