#!/bin/bash

# Script to add mapstructure tags to config structs in internal/config/
# This version preserves existing tags and only adds missing mapstructure tags
# It's idempotent - can be run multiple times safely

CONFIG_FILE="internal/config/types.gen.go"

# Process the file to add mapstructure tags where they're missing
# This will only add mapstructure tags if they don't already exist

# Case 1: Lines with json and yaml but NO mapstructure - with omitempty,omitzero
sed -i -E '/mapstructure:/!s/`json:"([^",]+),omitempty,omitzero" yaml:"([^"]+),omitempty"`$/`json:"\1,omitempty,omitzero" yaml:"\2,omitempty" mapstructure:"\1,omitempty"`/g' "$CONFIG_FILE"

# Case 2: Lines with json and yaml but NO mapstructure - with omitzero,omitempty
sed -i -E '/mapstructure:/!s/`json:"([^",]+),omitzero,omitempty" yaml:"([^"]+),omitempty"`$/`json:"\1,omitzero,omitempty" yaml:"\2,omitempty" mapstructure:"\1,omitempty"`/g' "$CONFIG_FILE"

# Case 3: Lines with json and yaml but NO mapstructure - with only omitempty
sed -i -E '/mapstructure:/!s/`json:"([^",]+),omitempty" yaml:"([^"]+),omitempty"`$/`json:"\1,omitempty" yaml:"\2,omitempty" mapstructure:"\1,omitempty"`/g' "$CONFIG_FILE"

# Case 4: Lines with json and yaml but NO mapstructure - with only omitzero
sed -i -E '/mapstructure:/!s/`json:"([^",]+),omitzero" yaml:"([^"]+),omitempty"`$/`json:"\1,omitzero" yaml:"\2,omitempty" mapstructure:"\1,omitempty"`/g' "$CONFIG_FILE"

# Case 5: Lines with json and yaml but NO mapstructure - simple fields without yaml modifiers
sed -i -E '/mapstructure:/!s/`json:"([^",]+)" yaml:"([^"]+)"`$/`json:"\1" yaml:"\2" mapstructure:"\1"`/g' "$CONFIG_FILE"

# Case 6: Lines with json (with modifiers) and yaml but NO mapstructure - yaml without modifiers
sed -i -E '/mapstructure:/!s/`json:"([^",]+),omitempty,omitzero" yaml:"([^"]+)"`$/`json:"\1,omitempty,omitzero" yaml:"\2" mapstructure:"\1,omitempty"`/g' "$CONFIG_FILE"
sed -i -E '/mapstructure:/!s/`json:"([^",]+),omitzero,omitempty" yaml:"([^"]+)"`$/`json:"\1,omitzero,omitempty" yaml:"\2" mapstructure:"\1,omitempty"`/g' "$CONFIG_FILE"
sed -i -E '/mapstructure:/!s/`json:"([^",]+),omitempty" yaml:"([^"]+)"`$/`json:"\1,omitempty" yaml:"\2" mapstructure:"\1,omitempty"`/g' "$CONFIG_FILE"
sed -i -E '/mapstructure:/!s/`json:"([^",]+),omitzero" yaml:"([^"]+)"`$/`json:"\1,omitzero" yaml:"\2" mapstructure:"\1,omitempty"`/g' "$CONFIG_FILE"

# Format the modified file silently
gofmt -w "$CONFIG_FILE" 2>/dev/null