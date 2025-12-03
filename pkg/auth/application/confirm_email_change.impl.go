// Package application provides handler implementations for auth operations.
package application

// NOTE: This file is user-editable. The generator will not overwrite it.

import (
	"context"
	"fmt"
)

// Ensure ConfirmEmailChangeHandlerImpl implements ConfirmEmailChangeHandler
var _ ConfirmEmailChangeHandler = (*ConfirmEmailChangeHandlerImpl)(nil)

// ConfirmEmailChangeHandlerImpl implements the ConfirmEmailChangeHandler interface.
type ConfirmEmailChangeHandlerImpl struct {
	// TODO: Add your dependencies here
}

// NewConfirmEmailChangeHandler creates a new ConfirmEmailChange handler.
func NewConfirmEmailChangeHandler(
// TODO: Add your dependencies here
) ConfirmEmailChangeHandler {
	return &ConfirmEmailChangeHandlerImpl{
		// TODO: Initialize dependencies
	}
}

// Execute performs the ConfirmEmailChange operation.
func (h *ConfirmEmailChangeHandlerImpl) Execute(
	_ context.Context,
	_ *ConfirmEmailChangeInput,
) error {
	// TODO: Implement ConfirmEmailChange logic
	return fmt.Errorf("not implemented")
}
