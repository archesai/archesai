package spec

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// RefState indicates whether a Ref holds a reference, resolved value, or both.
type RefState uint8

const (
	// RefStateUnresolved means only RefPath is set, Value is nil.
	RefStateUnresolved RefState = iota
	// RefStateResolved means both RefPath and Value are set.
	RefStateResolved
	// RefStateInline means no RefPath, only inline Value.
	RefStateInline
)

// Ref is a generic container for OpenAPI references or inline values.
// It supports lazy resolution and round-trip YAML serialization.
type Ref[T any] struct {
	// RefPath is the $ref string (e.g., "#/components/schemas/User" or "./User.yaml").
	// Empty for inline values.
	RefPath string

	// Value is the resolved or inline value. Nil until resolved for references.
	Value *T

	// state tracks resolution status.
	state RefState
}

// NewRef creates a reference that needs resolution.
func NewRef[T any](refPath string) *Ref[T] {
	return &Ref[T]{
		RefPath: refPath,
		state:   RefStateUnresolved,
	}
}

// NewInline creates an inline value (not a reference).
func NewInline[T any](value *T) *Ref[T] {
	return &Ref[T]{
		Value: value,
		state: RefStateInline,
	}
}

// NewResolved creates a resolved reference with both path and value.
func NewResolved[T any](refPath string, value *T) *Ref[T] {
	return &Ref[T]{
		RefPath: refPath,
		Value:   value,
		state:   RefStateResolved,
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
	return r.state == RefStateInline
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
		r.state = RefStateResolved
	} else {
		r.state = RefStateInline
	}
}

// State returns the current resolution state.
func (r *Ref[T]) State() RefState {
	if r == nil {
		return RefStateUnresolved
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
		r.state = RefStateInline
		return nil
	}

	// Check for $ref first
	for i := 0; i < len(node.Content); i += 2 {
		if node.Content[i].Value == "$ref" {
			r.RefPath = node.Content[i+1].Value
			r.state = RefStateUnresolved
			return nil
		}
	}

	// No $ref, parse as inline value
	var value T
	if err := node.Decode(&value); err != nil {
		return fmt.Errorf("failed to decode inline value: %w", err)
	}
	r.Value = &value
	r.state = RefStateInline
	return nil
}

// MarshalYAML serializes as $ref if it's a reference, otherwise inline.
func (r *Ref[T]) MarshalYAML() (interface{}, error) {
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
