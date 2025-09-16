#!/bin/bash

# Generate Coverage Report Script
# This script generates a formatted markdown coverage report from Go test coverage data

set -e

# Colors for terminal output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Generating coverage report...${NC}"

# Ensure docs/guides directory exists
mkdir -p docs/guides

# Run tests with coverage
echo -e "${YELLOW}Running tests with coverage...${NC}"
go test -race -coverprofile=coverage.out -covermode=atomic ./... > test-output.txt 2>&1 || true
go tool cover -func=coverage.out > coverage.txt

# Extract total coverage
TOTAL_COVERAGE=$(tail -1 coverage.txt | awk '{print $NF}')
TOTAL_COVERAGE_NUM=${TOTAL_COVERAGE%.*}  # Remove decimal part
TOTAL_COVERAGE_NUM=${TOTAL_COVERAGE_NUM%\%}  # Remove % sign

# Determine status emoji for total coverage
if [[ $TOTAL_COVERAGE_NUM -ge 80 ]]; then
    TOTAL_STATUS="🟢"
elif [[ $TOTAL_COVERAGE_NUM -ge 40 ]]; then
    TOTAL_STATUS="🟡"
else
    TOTAL_STATUS="🔴"
fi

# Start generating the report
cat > docs/guides/test-coverage-report.md << EOF
# Test Coverage Report

Generated: $(date -u +"%Y-%m-%d %H:%M:%S UTC")

## 📊 Summary

**Total Coverage:** \`$TOTAL_COVERAGE\` $TOTAL_STATUS

EOF

# Add warning if coverage is low
if [[ $TOTAL_COVERAGE_NUM -lt 80 ]]; then
    echo "> ⚠️ **Warning:** Coverage is below recommended threshold of 80%" >> docs/guides/test-coverage-report.md
    echo "" >> docs/guides/test-coverage-report.md
else
    echo "> ✅ **Good:** Coverage meets recommended threshold of 80%" >> docs/guides/test-coverage-report.md
    echo "" >> docs/guides/test-coverage-report.md
fi

# Generate coverage by package table
echo "## Coverage by Package" >> docs/guides/test-coverage-report.md
echo "" >> docs/guides/test-coverage-report.md
echo "| Package | Coverage | Status |" >> docs/guides/test-coverage-report.md
echo "|---------|----------|--------|" >> docs/guides/test-coverage-report.md

# Process test output for package coverage
grep -E "^(ok|FAIL)" test-output.txt | while read status pkg time coverage_text; do
    if [[ "$coverage_text" == *"coverage:"* ]]; then
        # Extract just the percentage
        coverage=$(echo "$coverage_text" | sed 's/.*coverage: //' | sed 's/ of statements//')
        coverage_num=${coverage%.*}  # Remove decimal part
        coverage_num=${coverage_num%\%}  # Remove % sign
        
        # Determine status emoji
        if [[ "$coverage" == "[no test files]" ]] || [[ "$coverage" == "[no statements]" ]]; then
            status_emoji="⚫ No statements"
            coverage="-"
        elif [[ "$coverage_num" == "0" ]]; then
            status_emoji="⚫ None"
        elif [[ $coverage_num -ge 80 ]]; then
            status_emoji="🟢 Good"
        elif [[ $coverage_num -ge 40 ]]; then
            status_emoji="🟡 Medium"
        else
            status_emoji="🔴 Low"
        fi
        
        # Clean package name (remove github.com/archesai/archesai/)
        clean_pkg=$(echo $pkg | sed 's|github.com/archesai/archesai/||')
        
        echo "| \`$clean_pkg\` | $coverage | $status_emoji |" >> docs/guides/test-coverage-report.md
    fi
done

# Add legend
cat >> docs/guides/test-coverage-report.md << 'EOF'

## Coverage Trends

```
Legend: 🟢 >80% | 🟡 40-80% | 🔴 <40% | ⚫ 0%
```

## Top Uncovered Files

The following critical files have 0% coverage and should be prioritized:

EOF

# Find files with 0% coverage
echo "Analyzing uncovered files..." >&2
grep "0.0%" coverage.txt | head -5 | while read file line func coverage; do
    # Extract just the filename
    filename=$(echo $file | sed 's|github.com/archesai/archesai/||' | sed 's|:.*||')
    # Determine file type/purpose
    if [[ "$filename" == *"main.go"* ]]; then
        desc="Main application entry point"
    elif [[ "$filename" == *"app.go"* ]]; then
        desc="Core application setup"
    elif [[ "$filename" == *"routes.go"* ]]; then
        desc="API route registration"
    elif [[ "$filename" == *"cache"* ]]; then
        desc="Caching layer"
    elif [[ "$filename" == *"database"* ]]; then
        desc="Database layer"
    else
        desc=""
    fi
    echo "- **\`$filename\`** - $desc" >> docs/guides/test-coverage-report.md
done

# Add recommendations
cat >> docs/guides/test-coverage-report.md << 'EOF'

## Recommendations

### Immediate Actions Required
EOF

# Add specific recommendations based on coverage
if [[ $TOTAL_COVERAGE_NUM -lt 20 ]]; then
    cat >> docs/guides/test-coverage-report.md << 'EOF'
- [ ] Add basic unit tests for all packages
- [ ] Focus on critical business logic first
- [ ] Set up test infrastructure and mocks
- [ ] Aim for at least 40% coverage initially
EOF
elif [[ $TOTAL_COVERAGE_NUM -lt 40 ]]; then
    cat >> docs/guides/test-coverage-report.md << 'EOF'
- [ ] Increase coverage for core packages
- [ ] Add integration tests for API endpoints
- [ ] Test error handling paths
- [ ] Aim for 60% coverage next
EOF
elif [[ $TOTAL_COVERAGE_NUM -lt 80 ]]; then
    cat >> docs/guides/test-coverage-report.md << 'EOF'
- [ ] Add edge case testing
- [ ] Improve integration test coverage
- [ ] Add performance benchmarks
- [ ] Aim for 80% coverage target
EOF
else
    cat >> docs/guides/test-coverage-report.md << 'EOF'
- [ ] Maintain current coverage levels
- [ ] Add tests for new features before merging
- [ ] Consider adding mutation testing
- [ ] Keep coverage above 80%
EOF
fi

# Add test execution details
TOTAL_PACKAGES=$(grep -c "^ok\|^FAIL\|^\?" test-output.txt || echo "0")
TESTED_PACKAGES=$(grep -c "^ok" test-output.txt || echo "0")

cat >> docs/guides/test-coverage-report.md << EOF

### Next Steps
1. Focus on packages with lowest coverage
2. Add tests for critical user flows
3. Implement continuous coverage monitoring

## Test Execution Details

- **Total Packages:** $TOTAL_PACKAGES
- **Packages with Tests:** $TESTED_PACKAGES
- **Test Execution Time:** ~10s

---

*This report is automatically generated by GitHub Actions on every push to main.*
EOF

# Clean up temp files
rm -f test-output.txt

echo -e "${GREEN}✓ Coverage report generated at docs/guides/test-coverage-report.md${NC}"
echo -e "${GREEN}Total Coverage: $TOTAL_COVERAGE${NC}"