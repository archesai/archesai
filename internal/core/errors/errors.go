// Package errors defines domain-specific errors for the core business logic.
package errors

import (
	"errors"
	"fmt"
)

// Domain errors that can occur across all entities.
var (
	// ErrNotFound indicates that a requested entity was not found.
	ErrNotFound = errors.New("entity not found")

	// ErrAlreadyExists indicates that an entity already exists.
	ErrAlreadyExists = errors.New("entity already exists")

	// ErrInvalidInput indicates that the provided input is invalid.
	ErrInvalidInput = errors.New("invalid input")

	// ErrUnauthorized indicates that the operation is not authorized.
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden indicates that the operation is forbidden.
	ErrForbidden = errors.New("forbidden")

	// ErrOptimisticLock indicates an optimistic locking conflict.
	ErrOptimisticLock = errors.New("optimistic lock error: entity was modified")

	// ErrDomainViolation indicates a business rule violation.
	ErrDomainViolation = errors.New("domain rule violation")

	// ErrNoChange indicates that an update would result in no changes.
	ErrNoChange = errors.New("no changes detected")
)

// Label-specific errors.
var (
	// ErrLabelNotFound indicates that a label was not found.
	ErrLabelNotFound = errors.New("label not found")

	// ErrLabelAlreadyExists indicates that a label with the same name already exists.
	ErrLabelAlreadyExists = errors.New("label already exists")

	// ErrMaxLabelsExceeded indicates that the maximum number of labels has been exceeded.
	ErrMaxLabelsExceeded = errors.New("maximum number of labels exceeded")

	// ErrLabelInUse indicates that a label cannot be deleted because it's in use.
	ErrLabelInUse = errors.New("label is in use and cannot be deleted")

	// ErrInvalidLabelName indicates that the label name is invalid.
	ErrInvalidLabelName = errors.New("invalid label name")
)

// Organization-specific errors.
var (
	// ErrOrganizationNotFound indicates that an organization was not found.
	ErrOrganizationNotFound = errors.New("organization not found")

	// ErrOrganizationInactive indicates that the organization is inactive.
	ErrOrganizationInactive = errors.New("organization is inactive")

	// ErrInvalidOrganizationID indicates that the organization ID is invalid.
	ErrInvalidOrganizationID = errors.New("invalid organization ID")
)

// User-specific errors.
var (
	// ErrUserNotFound indicates that a user was not found.
	ErrUserNotFound = errors.New("user not found")

	// ErrUserAlreadyExists indicates that a user already exists.
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrUserInactive indicates that the user is inactive.
	ErrUserInactive = errors.New("user is inactive")

	// ErrInvalidUserID indicates that the user ID is invalid.
	ErrInvalidUserID = errors.New("invalid user ID")
)

// Tool-specific errors.
var (
	// ErrToolNotFound indicates that a tool was not found.
	ErrToolNotFound = errors.New("tool not found")

	// ErrToolInactive indicates that the tool is inactive.
	ErrToolInactive = errors.New("tool is inactive")

	// ErrUsageLimitExceeded indicates that the tool's usage limit has been exceeded.
	ErrUsageLimitExceeded = errors.New("usage limit exceeded")

	// ErrInvalidToolType indicates that the tool type is invalid.
	ErrInvalidToolType = errors.New("invalid tool type")

	// ErrToolExists is returned when a tool already exists.
	ErrToolExists = errors.New("tool already exists")
)

// Account specific errors.
var (
	// ErrAccountNotFound indicates that an account was not found.
	ErrAccountNotFound = errors.New("account not found")

	// ErrAccountInactive indicates that the account is inactive.
	ErrAccountInactive = errors.New("account is inactive")

	// ErrInvalidAccountIdentifier indicates that the account ID is invalid.
	ErrInvalidAccountIdentifier = errors.New("invalid account ID")
)

// ValidationError represents a validation failure
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface
func (e ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string, args ...interface{}) error {
	return ValidationError{
		Field:   field,
		Message: fmt.Sprintf(message, args...),
	}
}

// BusinessRuleError represents a business rule violation
type BusinessRuleError struct {
	Rule    string
	Message string
}

// Error implements the error interface
func (e BusinessRuleError) Error() string {
	return fmt.Sprintf("business rule violation: %s", e.Message)
}

// NewBusinessRuleError creates a new business rule error
func NewBusinessRuleError(rule, message string, args ...interface{}) error {
	return BusinessRuleError{
		Rule:    rule,
		Message: fmt.Sprintf(message, args...),
	}
}

// Missing entity-specific errors for generated code
var (
	// Artifact errors
	ErrArtifactNotFound = errors.New("artifact not found")
	ErrArtifactExists   = errors.New("artifact already exists")

	// Auth errors
	ErrAuthNotFound = errors.New("auth not found")
	ErrAuthExists   = errors.New("auth already exists")

	// Invitation errors
	ErrInvitationNotFound = errors.New("invitation not found")
	ErrInvitationExists   = errors.New("invitation already exists")

	// Member errors
	ErrMemberNotFound = errors.New("member not found")
	ErrMemberExists   = errors.New("member already exists")

	// Pipeline errors
	ErrPipelineNotFound = errors.New("pipeline not found")
	ErrPipelineExists   = errors.New("pipeline already exists")

	// PipelineStep errors
	ErrPipelineStepNotFound = errors.New("pipeline step not found")
	ErrPipelineStepExists   = errors.New("pipeline step already exists")

	// Run errors
	ErrRunNotFound = errors.New("run not found")
	ErrRunExists   = errors.New("run already exists")

	// Health errors
	ErrHealthNotFound = errors.New("health not found")
	ErrHealthExists   = errors.New("health already exists")

	// Config errors
	ErrConfigNotFound = errors.New("config not found")
	ErrConfigExists   = errors.New("config already exists")

	// Page errors
	ErrPageNotFound = errors.New("page not found")
	ErrPageExists   = errors.New("page already exists")

	// Base errors
	ErrBaseNotFound = errors.New("base not found")
	ErrBaseExists   = errors.New("base already exists")

	// FilterNode errors
	ErrFilterNodeNotFound = errors.New("filter node not found")
	ErrFilterNodeExists   = errors.New("filter node already exists")

	// Problem errors
	ErrProblemNotFound = errors.New("problem not found")
	ErrProblemExists   = errors.New("problem already exists")

	// Session errors
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExists   = errors.New("session already exists")

	// APIKey errors
	ErrAPIKeyNotFound = errors.New("api key not found")
	ErrAPIKeyExists   = errors.New("api key already exists")
)
