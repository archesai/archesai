package parsers

import (
	"fmt"
	"log/slog"
	"maps"
	"sort"

	"github.com/daveshanley/vacuum/model"
	"github.com/daveshanley/vacuum/motor"
	"github.com/daveshanley/vacuum/rulesets"
)

const (
	// ErrorSeverity represents the severity level for errors
	ErrorSeverity = "error"
)

// Lint performs strict linting on the OpenAPI specification using vacuum
// It uses ALL available rules for maximum strictness and returns an error
// if ANY violations are found
func (p *OpenAPIParser) Lint(specBytes []byte) error {
	return p.LintWithBasePath(specBytes, "")
}

// LintWithBasePath performs strict linting with a specified base path for resolving references
func (p *OpenAPIParser) LintWithBasePath(specBytes []byte, basePath string) error {
	logger := slog.Default()
	logger.Info("Starting OpenAPI specification linting...")
	defaultRuleSets := rulesets.BuildDefaultRuleSetsWithLogger(logger)
	selectedRS := defaultRuleSets.GenerateOpenAPIRecommendedRuleSet()
	owaspRules := rulesets.GetAllOWASPRules()
	maps.Copy(selectedRS.Rules, owaspRules)
	// Create rule execution with base path if provided
	execution := &motor.RuleSetExecution{
		RuleSet: selectedRS,
		Spec:    specBytes,
		Logger:  logger,
	}

	// Set base path if provided for resolving file references
	if basePath != "" {
		execution.Base = basePath
	}

	lintingResults := motor.ApplyRulesToRuleSet(execution)

	// Check if there are any results (violations)
	if len(lintingResults.Results) == 0 {
		// No violations found, spec is clean
		logger.Info("✅ OpenAPI specification passed linting - no violations found")
		return nil
	}

	// Create a rule result set for better organization
	resultSet := model.NewRuleResultSet(lintingResults.Results)

	// Sort results by line number for better readability
	resultSet.SortResultsByLineNumber()

	// Log violations using slog and return error
	logLintingErrors(resultSet, logger)

	// Return a simple error to block further processing
	return fmt.Errorf(
		"OpenAPI specification failed linting with %d violations",
		len(resultSet.Results),
	)
}

// logLintingErrors logs all linting violations using structured logging
func logLintingErrors(resultSet *model.RuleResultSet, logger *slog.Logger) {
	// Group violations by severity
	errorViolations := []*model.RuleFunctionResult{}
	warnViolations := []*model.RuleFunctionResult{}
	infoViolations := []*model.RuleFunctionResult{}

	for _, result := range resultSet.Results {
		switch result.RuleSeverity {
		case ErrorSeverity:
			errorViolations = append(errorViolations, result)
		case "warn", "warning":
			warnViolations = append(warnViolations, result)
		case "info", "hint":
			infoViolations = append(infoViolations, result)
		default:
			errorViolations = append(errorViolations, result) // Default to error
		}
	}

	// Sort each severity group by line number
	sortViolations := func(violations []*model.RuleFunctionResult) {
		sort.Slice(violations, func(i, j int) bool {
			if violations[i].StartNode == nil || violations[j].StartNode == nil {
				return false
			}
			if violations[i].StartNode.Line == violations[j].StartNode.Line {
				return violations[i].StartNode.Column < violations[j].StartNode.Column
			}
			return violations[i].StartNode.Line < violations[j].StartNode.Line
		})
	}

	sortViolations(errorViolations)
	sortViolations(warnViolations)
	sortViolations(infoViolations)

	// Helper function to log a violation
	logViolation := func(result *model.RuleFunctionResult) {
		line := 0
		column := 0
		if result.StartNode != nil {
			line = result.StartNode.Line
			column = result.StartNode.Column
		}

		severity := result.RuleSeverity
		if severity == "" {
			severity = ErrorSeverity
		}

		// Common attributes for all log entries
		attrs := []any{
			slog.String("rule", result.RuleId),
			slog.String("path", result.Path),
			slog.Int("line", line),
			slog.Int("column", column),
		}

		// Get file location if available
		if result.Origin != nil && result.Origin.AbsoluteLocation != "" {
			attrs = append(attrs, slog.String("file", result.Origin.AbsoluteLocation))
		}

		// Log with appropriate level based on severity
		msg := fmt.Sprintf("[%d:%d] %s", line, column, result.Message)

		switch severity {
		case ErrorSeverity:
			logger.Error(msg, attrs...)
		case "warn", "warning":
			logger.Warn(msg, attrs...)
		case "info", "hint":
			logger.Info(msg, attrs...)
		default:
			logger.Error(msg, attrs...)
		}
	}

	// Log all errors first
	for _, result := range errorViolations {
		logViolation(result)
	}

	// Then warnings
	for _, result := range warnViolations {
		logViolation(result)
	}

	// Then info
	for _, result := range infoViolations {
		logViolation(result)
	}

	// Count violations by rule
	ruleCounts := make(map[string]int)
	for _, result := range resultSet.Results {
		ruleCounts[result.RuleId]++
	}

	// Sort rule names for consistent output
	var ruleNames []string
	for rule := range ruleCounts {
		ruleNames = append(ruleNames, rule)
	}
	sort.Strings(ruleNames)

	// Log summary by rule
	logger.Info("Violations by rule:")
	for _, rule := range ruleNames {
		logger.Info(fmt.Sprintf("  %s: %d", rule, ruleCounts[rule]))
	}

	// Log total summary at the end
	totalViolations := len(resultSet.Results)
	logger.Error("❌ OpenAPI specification validation failed",
		slog.Int("total", totalViolations),
		slog.Int("errors", len(errorViolations)),
		slog.Int("warnings", len(warnViolations)),
		slog.Int("info", len(infoViolations)),
	)
}
