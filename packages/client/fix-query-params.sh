#!/bin/bash

# Script to replace manual query parameter encoding with qs-based encoding
# in Orval-generated files

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GENERATED_DIR="$SCRIPT_DIR/src/generated"

echo "🔧 Fixing query parameter encoding in generated files..."

# First, find and fix files with manual normalizedParams.append patterns
find "$GENERATED_DIR" -name "*.ts" -exec grep -l "normalizedParams\.append" {} \; | while read -r file; do
    echo "📝 Processing manual encoding in: $file"
    
    # Create a temporary file for processing
    temp_file=$(mktemp)
    
    # Use awk to replace the entire block
    awk '
    /const normalizedParams = new URLSearchParams\(\)/ {
        print "  const stringifiedParams = qs.stringify(params || {}, { skipNulls: false, strictNullHandling: true })"
        # Skip the next lines until we find the const stringifiedParams line
        while ((getline line) > 0) {
            if (line ~ /const stringifiedParams = normalizedParams\.toString\(\)/) {
                break
            }
        }
        next
    }
    { print }
    ' "$file" > "$temp_file"
    
    # Replace original file with processed content
    mv "$temp_file" "$file"
    
    echo "✅ Updated manual encoding in: $file"
done

# Second, add missing qs imports to files that use qs.stringify but don't have the import
find "$GENERATED_DIR" -name "*.ts" -exec grep -l "qs\.stringify" {} \; | while read -r file; do
    if ! grep -q "import qs from 'qs'" "$file"; then
        echo "📝 Adding missing qs import to: $file"
        # Add qs import after useMutation import line
        sed -i '/^import { useMutation, useQuery, useSuspenseQuery } from/a\
\
import qs from '\''qs'\''' "$file"
        echo "✅ Added qs import to: $file"
    fi
done

echo "🎉 Query parameter encoding fix completed!"