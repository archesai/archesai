package ref

import (
	"fmt"

	"go.yaml.in/yaml/v4"
)

// State indicates whether a Ref holds a reference, resolved value, or both.
type State uint8

const (
	// StateUnresolved means only RefPath is set, Value is nil.
	StateUnresolved State = iota
	// StateResolved means both RefPath and Value are set.
	StateResolved
	// StateInline means no RefPath, only inline Value.
	StateInline
)

// Ref is a generic container for references or inline values.
// It supports lazy resolution and round-trip YAML serialization.
type Ref[T any] struct {
	// RefPath is the $ref string (e.g., "#/components/schemas/User" or "./User.yaml").
	// Empty for inline values.
	RefPath string

	// Value is the resolved or inline value. Nil until resolved for references.
	Value *T

	// state tracks resolution status.
	state State
}

// NewRef creates a reference that needs resolution.
func NewRef[T any](refPath string) *Ref[T] {
	return &Ref[T]{
		RefPath: refPath,
		state:   StateUnresolved,
	}
}

// NewInline creates an inline value (not a reference).
func NewInline[T any](value *T) *Ref[T] {
	return &Ref[T]{
		Value: value,
		state: StateInline,
	}
}

// NewResolved creates a resolved reference with both path and value.
func NewResolved[T any](refPath string, value *T) *Ref[T] {
	return &Ref[T]{
		RefPath: refPath,
		Value:   value,
		state:   StateResolved,
	}
}

// IsRef returns true if this is a reference (resolved or not).
func (r *Ref[T]) IsRef() bool {
	if r == nil {
		return false
	}
	return r.RefPath != ""
}

// IsInline returns true if this is an inline value (not a reference).
func (r *Ref[T]) IsInline() bool {
	if r == nil {
		return false
	}
	return r.state == StateInline
}

// IsResolved returns true if the value has been resolved.
func (r *Ref[T]) IsResolved() bool {
	if r == nil {
		return false
	}
	return r.Value != nil
}

// Get returns the value. Panics if not resolved.
func (r *Ref[T]) Get() *T {
	if r == nil {
		panic("Ref.Get() called on nil Ref")
	}
	if r.Value == nil {
		panic(fmt.Sprintf("Ref.Get() called on unresolved reference: %s", r.RefPath))
	}
	return r.Value
}

// GetOrNil returns the value or nil if not resolved.
func (r *Ref[T]) GetOrNil() *T {
	if r == nil {
		return nil
	}
	return r.Value
}

// Resolve sets the resolved value for a reference.
func (r *Ref[T]) Resolve(value *T) {
	if r == nil {
		return
	}
	r.Value = value
	if r.RefPath != "" {
		r.state = StateResolved
	} else {
		r.state = StateInline
	}
}

// State returns the current resolution state.
func (r *Ref[T]) State() State {
	if r == nil {
		return StateUnresolved
	}
	return r.state
}

// UnmarshalYAML handles parsing of either $ref or inline content.
func (r *Ref[T]) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		// Try to decode directly as the target type
		var value T
		if err := node.Decode(&value); err != nil {
			return fmt.Errorf("failed to decode inline value: %w", err)
		}
		r.Value = &value
		r.state = StateInline
		return nil
	}

	// Check for $ref first
	for i := 0; i < len(node.Content); i += 2 {
		if node.Content[i].Value == "$ref" {
			r.RefPath = node.Content[i+1].Value
			r.state = StateUnresolved
			return nil
		}
	}

	// No $ref, parse as inline value
	var value T
	if err := node.Decode(&value); err != nil {
		return fmt.Errorf("failed to decode inline value: %w", err)
	}
	r.Value = &value
	r.state = StateInline
	return nil
}

// MarshalYAML serializes as $ref if it's a reference, otherwise inline.
func (r *Ref[T]) MarshalYAML() (any, error) {
	if r == nil {
		return nil, nil
	}
	if r.RefPath != "" {
		// Serialize as $ref
		return map[string]string{"$ref": r.RefPath}, nil
	}
	// Serialize inline value
	return r.Value, nil
}
