"""
Default execute function that raises an error.
This file should be replaced by mounting a custom execute.py
or by building a derived image with a custom implementation.
"""

from typing import Any, Dict


def execute_function(input_data: Dict[str, Any]) -> Dict[str, Any]:
    """
    Execute function - must be implemented by user.

    Args:
        input_data: The input data

    Returns:
        The output data

    Raises:
        NotImplementedError: Always, as this is a placeholder
    """
    raise NotImplementedError(
        "No execution function provided. "
        "Either mount a custom execute.py at /app/execute.py "
        "or use a pre-built generator image (e.g., archesai/generator-custom)."
    )
