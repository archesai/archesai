package application

// NOTE: This file is user-editable. The generator will not overwrite it.

import (
	"context"
	"fmt"
)

// Ensure LinkAccountImpl implements LinkAccount
var _ LinkAccount = (*LinkAccountImpl)(nil)

// LinkAccountImpl implements the LinkAccount interface.
type LinkAccountImpl struct {
	// TODO: Add your dependencies here
}

// NewLinkAccount creates a new LinkAccount implementation.
func NewLinkAccount(
// TODO: Add your dependencies here
) LinkAccount {
	return &LinkAccountImpl{
		// TODO: Initialize dependencies
	}
}

// Execute performs the LinkAccount operation.
func (h *LinkAccountImpl) Execute(
	ctx context.Context,
	input *LinkAccountInput,
) (*LinkAccountOutput, error) {
	// TODO: Implement LinkAccount logic
	return nil, fmt.Errorf("not implemented")
}
