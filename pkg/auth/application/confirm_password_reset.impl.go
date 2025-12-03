package application

// NOTE: This file is user-editable. The generator will not overwrite it.

import (
	"context"
	"fmt"
)

// Ensure ConfirmPasswordResetHandlerImpl implements ConfirmPasswordResetHandler
var _ ConfirmPasswordResetHandler = (*ConfirmPasswordResetHandlerImpl)(nil)

// ConfirmPasswordResetHandlerImpl implements the ConfirmPasswordResetHandler interface.
type ConfirmPasswordResetHandlerImpl struct {
	// TODO: Add your dependencies here
}

// NewConfirmPasswordResetHandler creates a new ConfirmPasswordReset handler.
func NewConfirmPasswordResetHandler(
// TODO: Add your dependencies here
) ConfirmPasswordResetHandler {
	return &ConfirmPasswordResetHandlerImpl{
		// TODO: Initialize dependencies
	}
}

// Execute performs the ConfirmPasswordReset operation.
func (h *ConfirmPasswordResetHandlerImpl) Execute(
	_ context.Context,
	_ *ConfirmPasswordResetInput,
) error {
	// TODO: Implement ConfirmPasswordReset logic
	return fmt.Errorf("not implemented")
}
