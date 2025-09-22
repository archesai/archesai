#!/usr/bin/env bash

# Update project layout documentation with the current directory structure
# This script generates a tree structure and formats it into the documentation
#
# NOTE: For token-efficient project navigation, use the XML format instead:
#   Run: ./scripts/generate-project-structure-xml.sh
#   Output: docs/architecture/project-structure.xml (70% fewer tokens)

set -euo pipefail

# Check if tree command is available
if ! command -v tree &> /dev/null; then
    echo "Error: 'tree' command is not installed. Please install it first."
    echo "  Ubuntu/Debian: sudo apt-get install tree"
    echo "  macOS: brew install tree"
    echo "  RHEL/CentOS: sudo yum install tree"
    exit 1
fi

# Ensure docs directory exists
mkdir -p docs/architecture

# Generate tree structure (excluding hidden files and common build artifacts)
tree_output=$(tree -L 5 -a -I '.github|.claude|.git|node_modules|*.gen.go|dist|build|.DS_Store|*.log|coverage|.next|.turbo|.cache' --dirsfirst)

# Create a temporary file with the new content
cat > temp_project_layout.md << 'EOF'
# Project Layout

## Directory Structure

```text
EOF

echo "$tree_output" >> temp_project_layout.md
echo '```' >> temp_project_layout.md

# Check if the original file exists and has content after the directory structure
if [ -f "docs/architecture/project-layout.md" ]; then
    # Extract everything after the first closing ``` in the original file
    # This preserves any additional documentation that was added manually
    sed -n '/^```$/,$p' docs/architecture/project-layout.md | tail -n +2 >> temp_project_layout.md
fi

# Replace the original file
mv temp_project_layout.md docs/architecture/project-layout.md

echo "✅ Project layout documentation updated at docs/architecture/project-layout.md"