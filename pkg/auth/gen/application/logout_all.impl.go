package application

// NOTE: This file is user-editable. The generator will not overwrite it.

import (
	"context"
	"fmt"
)

// Ensure LogoutAllImpl implements LogoutAll
var _ LogoutAll = (*LogoutAllImpl)(nil)

// LogoutAllImpl implements the LogoutAll interface.
type LogoutAllImpl struct {
	// TODO: Add your dependencies here
}

// NewLogoutAll creates a new LogoutAll implementation.
func NewLogoutAll(
// TODO: Add your dependencies here
) LogoutAll {
	return &LogoutAllImpl{
		// TODO: Initialize dependencies
	}
}

// Execute performs the LogoutAll operation.
func (h *LogoutAllImpl) Execute(_ context.Context, _ *LogoutAllInput) (*LogoutAllOutput, error) {
	// TODO: Implement LogoutAll logic
	return nil, fmt.Errorf("not implemented")
}
