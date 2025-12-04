package application

// NOTE: This file is user-editable. The generator will not overwrite it.

import (
	"context"
	"fmt"
)

// Ensure LinkAccountHandlerImpl implements LinkAccountHandler
var _ LinkAccountHandler = (*LinkAccountHandlerImpl)(nil)

// LinkAccountHandlerImpl implements the LinkAccountHandler interface.
type LinkAccountHandlerImpl struct {
	// TODO: Add your dependencies here
}

// NewLinkAccountHandler creates a new LinkAccount handler.
func NewLinkAccountHandler(
// TODO: Add your dependencies here
) LinkAccountHandler {
	return &LinkAccountHandlerImpl{
		// TODO: Initialize dependencies
	}
}

// Execute performs the LinkAccount operation.
func (h *LinkAccountHandlerImpl) Execute(
	_ context.Context,
	_ *LinkAccountInput,
) (*LinkAccountOutput, error) {
	// TODO: Implement LinkAccount logic
	return nil, fmt.Errorf("not implemented")
}
