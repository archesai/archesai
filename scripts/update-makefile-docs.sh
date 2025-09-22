#!/usr/bin/env bash

# Update Makefile documentation with the latest make help output
# This script extracts the make help output and formats it into the documentation

set -euo pipefail

# Check if we're in the right directory (should have a Makefile)
if [ ! -f "Makefile" ]; then
    echo "Error: Makefile not found. Please run this script from the project root."
    exit 1
fi

# Ensure docs directory exists
mkdir -p docs/guides

# Create the documentation file
cat > docs/guides/makefile-commands.md << 'EOF'
# Makefile Commands

Run `make help` to see all available commands.

## Available Commands

```bash
EOF

# Append the make help output, removing ANSI color codes
make help | sed 's/\x1B\[[0-9;]\{1,\}[A-Za-z]//g' >> docs/guides/makefile-commands.md

# Close the code block
echo '```' >> docs/guides/makefile-commands.md

echo "✅ Makefile documentation updated at docs/guides/makefile-commands.md"