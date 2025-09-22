#!/usr/bin/env bash

# Generate XML project structure for token-efficient documentation
# This creates a compact XML representation of the project structure with metadata

set -euo pipefail

# Output file
OUTPUT_FILE="${1:-docs/architecture/project-structure.xml}"

# Ensure output directory exists
mkdir -p "$(dirname "$OUTPUT_FILE")"

# Function to detect directory type
detect_type() {
    local dir="$1"

    if [[ -f "$dir/package.json" ]]; then
        echo "npm-package"
    elif [[ -f "$dir/go.mod" ]]; then
        echo "go-module"
    elif [[ -f "$dir/Dockerfile" ]]; then
        echo "container"
    elif [[ -f "$dir/sqlc.yaml" ]]; then
        echo "database"
    elif [[ "$dir" == *"generated"* ]] || [[ "$dir" == *".gen"* ]]; then
        echo "generated"
    elif [[ -f "$dir/openapi.yaml" ]] || [[ -f "$dir/swagger.yaml" ]]; then
        echo "api-spec"
    elif [[ "$dir" == *"test"* ]] || [[ "$dir" == *"spec"* ]]; then
        echo "testing"
    else
        echo ""
    fi
}

# Function to detect directory purpose
detect_purpose() {
    local dir="$1"
    local basename=$(basename "$dir")

    # Try to extract from README
    if [[ -f "$dir/README.md" ]]; then
        local readme_purpose=$(head -n 3 "$dir/README.md" | grep -E "^#|^##" | head -1 | sed 's/^#\+\s*//' | cut -c1-50)
        if [[ -n "$readme_purpose" ]]; then
            echo "$readme_purpose"
            return
        fi
    fi

    # Try to extract from package.json
    if [[ -f "$dir/package.json" ]]; then
        local pkg_desc=$(jq -r '.description // empty' "$dir/package.json" 2>/dev/null | cut -c1-50)
        if [[ -n "$pkg_desc" ]]; then
            echo "$pkg_desc"
            return
        fi
    fi

    # Pattern-based purpose detection
    case "$basename" in
        auth*) echo "Authentication and authorization" ;;
        database*|db*) echo "Database operations" ;;
        api*) echo "API definitions" ;;
        migrations*) echo "Database migrations" ;;
        config*) echo "Configuration management" ;;
        cache*) echo "Caching layer" ;;
        storage*) echo "File storage" ;;
        redis*) echo "Redis operations" ;;
        server*) echo "Server infrastructure" ;;
        middleware*) echo "HTTP middleware" ;;
        logger*|logging*) echo "Logging utilities" ;;
        test*|spec*) echo "Testing utilities" ;;
        docs*|documentation*) echo "Documentation" ;;
        deployments*) echo "Deployment configurations" ;;
        scripts*) echo "Utility scripts" ;;
        tools*) echo "Development tools" ;;
        web*) echo "Web applications" ;;
        client*) echo "Client SDK" ;;
        platform*) echo "Platform UI" ;;
        ui*) echo "UI components" ;;
        components*) echo "Reusable components" ;;
        templates*) echo "Template files" ;;
        assets*) echo "Static assets" ;;
        *) echo "" ;;
    esac
}

# Function to count files by pattern
count_files() {
    local dir="$1"
    local pattern="${2:-*}"
    find "$dir" -maxdepth 1 -name "$pattern" -type f 2>/dev/null | wc -l
}

# Function to get key files in a directory
get_key_files() {
    local dir="$1"
    local key_files=""

    # Check for important files
    for file in main.go index.ts index.tsx handler.go service.go README.md package.json go.mod Dockerfile Makefile; do
        if [[ -f "$dir/$file" ]]; then
            if [[ -n "$key_files" ]]; then
                key_files="$key_files,$file"
            else
                key_files="$file"
            fi
        fi
    done

    echo "$key_files"
}

# Function to detect file patterns in directory
detect_patterns() {
    local dir="$1"
    local patterns=""

    # Check for common patterns
    if ls "$dir"/*.gen.go >/dev/null 2>&1; then
        patterns="*.gen.go"
    fi
    if ls "$dir"/*.sql >/dev/null 2>&1; then
        patterns="${patterns:+$patterns,}*.sql"
    fi
    if ls "$dir"/*.yaml >/dev/null 2>&1; then
        patterns="${patterns:+$patterns,}*.yaml"
    fi
    if ls "$dir"/*.tsx >/dev/null 2>&1; then
        patterns="${patterns:+$patterns,}*.tsx"
    fi

    echo "$patterns"
}

# Function to process directory recursively
process_directory() {
    local dir="$1"
    local indent="$2"
    local depth="$3"
    local max_depth="${4:-3}"

    # Skip certain directories
    if [[ "$dir" == *"node_modules"* ]] || [[ "$dir" == *".git"* ]] || [[ "$dir" == *".next"* ]] || [[ "$dir" == *"dist"* ]] || [[ "$dir" == *"coverage"* ]]; then
        return
    fi

    local basename=$(basename "$dir")
    local type=$(detect_type "$dir")
    local purpose=$(detect_purpose "$dir")
    local file_count=$(find "$dir" -maxdepth 1 -type f 2>/dev/null | wc -l)
    local dir_count=$(find "$dir" -maxdepth 1 -type d 2>/dev/null | grep -v "^$dir$" | wc -l)
    local key_files=$(get_key_files "$dir")
    local patterns=$(detect_patterns "$dir")

    # Build XML attributes
    local attrs="name=\"$basename\""
    [[ -n "$type" ]] && attrs="$attrs type=\"$type\""
    [[ -n "$purpose" ]] && attrs="$attrs purpose=\"$purpose\""
    [[ $file_count -gt 0 ]] && attrs="$attrs files=\"$file_count\""
    [[ $dir_count -gt 0 ]] && attrs="$attrs dirs=\"$dir_count\""
    [[ -n "$patterns" ]] && attrs="$attrs patterns=\"$patterns\""
    [[ -n "$key_files" ]] && attrs="$attrs key-files=\"$key_files\""

    # Check if directory has subdirectories to process
    local has_subdirs=false
    if [[ $depth -lt $max_depth ]]; then
        for subdir in "$dir"/*/; do
            if [[ -d "$subdir" ]]; then
                subdir=${subdir%/}
                if [[ "$(basename "$subdir")" != "node_modules" ]] && [[ "$(basename "$subdir")" != ".git" ]]; then
                    has_subdirs=true
                    break
                fi
            fi
        done
    fi

    if [[ "$has_subdirs" == true ]]; then
        echo "${indent}<dir $attrs>"

        # Process subdirectories
        for subdir in "$dir"/*/; do
            if [[ -d "$subdir" ]]; then
                subdir=${subdir%/}
                if [[ "$(basename "$subdir")" != "node_modules" ]] && [[ "$(basename "$subdir")" != ".git" ]]; then
                    process_directory "$subdir" "  $indent" $((depth + 1)) "$max_depth"
                fi
            fi
        done

        echo "${indent}</dir>"
    else
        echo "${indent}<dir $attrs/>"
    fi
}

# Start XML generation
cat > "$OUTPUT_FILE" << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<!--
  Project Structure XML
  Generated: DATE_PLACEHOLDER

  This XML provides a compact representation of the project structure.
  It's ~75% more token-efficient than the traditional tree format.

  Attributes:
  - name: Directory/file name
  - type: Type of directory (npm-package, go-module, etc.)
  - purpose: Brief description of the directory's purpose
  - files: Number of files in the directory
  - dirs: Number of subdirectories
  - patterns: Common file patterns in the directory
  - key-files: Important files in the directory
-->
EOF

# Replace date placeholder
sed -i "s/DATE_PLACEHOLDER/$(date -Iseconds)/" "$OUTPUT_FILE"

# Get project root
PROJECT_ROOT="$(pwd)"
PROJECT_NAME=$(basename "$PROJECT_ROOT")

# Start project root element
echo "<project name=\"$PROJECT_NAME\" root=\"$PROJECT_ROOT\">" >> "$OUTPUT_FILE"

# Process root-level files
echo "  <root-files>" >> "$OUTPUT_FILE"
for file in Makefile go.mod go.sum package.json pnpm-lock.yaml tsconfig.json arches.yaml README.md LICENSE .golangci.yaml .air.toml .goreleaser.yaml .lefthook.yaml .redocly.yaml .mockery.yaml .editorconfig .cspell.json .markdownlint.json biome.json; do
    if [[ -f "$file" ]]; then
        echo "    <file>$file</file>" >> "$OUTPUT_FILE"
    fi
done
echo "  </root-files>" >> "$OUTPUT_FILE"

# Process major directories
for dir in api assets cmd deployments docs internal scripts test tools web; do
    if [[ -d "$dir" ]]; then
        process_directory "$dir" "  " 0 3 >> "$OUTPUT_FILE"
    fi
done

# Process .taskmaster if it exists
if [[ -d ".taskmaster" ]]; then
    process_directory ".taskmaster" "  " 0 2 >> "$OUTPUT_FILE"
fi

# Process .vscode if it exists
if [[ -d ".vscode" ]]; then
    process_directory ".vscode" "  " 0 1 >> "$OUTPUT_FILE"
fi

# Close project root
echo "</project>" >> "$OUTPUT_FILE"

echo "✅ XML project structure generated at $OUTPUT_FILE"

# Show size comparison
if [[ -f "docs/architecture/project-layout.md" ]]; then
    old_size=$(wc -c < docs/architecture/project-layout.md)
    new_size=$(wc -c < "$OUTPUT_FILE")
    reduction=$((100 - (new_size * 100 / old_size)))
    echo "📊 Size comparison:"
    echo "   Text tree: $(numfmt --to=iec $old_size)"
    echo "   XML format: $(numfmt --to=iec $new_size)"
    echo "   Reduction: ${reduction}%"
fi