package application

// NOTE: This file is user-editable. The generator will not overwrite it.

import (
	"context"
	"fmt"
)

// Ensure LogoutAllHandlerImpl implements LogoutAllHandler
var _ LogoutAllHandler = (*LogoutAllHandlerImpl)(nil)

// LogoutAllHandlerImpl implements the LogoutAllHandler interface.
type LogoutAllHandlerImpl struct {
	// TODO: Add your dependencies here
}

// NewLogoutAllHandler creates a new LogoutAll handler.
func NewLogoutAllHandler(
// TODO: Add your dependencies here
) LogoutAllHandler {
	return &LogoutAllHandlerImpl{
		// TODO: Initialize dependencies
	}
}

// Execute performs the LogoutAll operation.
func (h *LogoutAllHandlerImpl) Execute(
	_ context.Context,
	_ *LogoutAllInput,
) (*LogoutAllOutput, error) {
	// TODO: Implement LogoutAll logic
	return nil, fmt.Errorf("not implemented")
}
