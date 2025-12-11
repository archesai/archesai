// Package ptr provides generic pointer utilities.
package ptr

// To returns a pointer to the given value.
// Useful for creating pointers to literals or inline values.
//
// Example:
//
//	name := ptr.To("hello")  // *string pointing to "hello"
//	count := ptr.To(42)      // *int pointing to 42
func To[T any](v T) *T {
	return &v
}

// From returns the value pointed to by p, or the zero value if p is nil.
//
// Example:
//
//	var name *string
//	s := ptr.From(name)  // returns ""
//
//	name = ptr.To("hello")
//	s = ptr.From(name)   // returns "hello"
func From[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

// FromOr returns the value pointed to by p, or the default value if p is nil.
//
// Example:
//
//	var name *string
//	s := ptr.FromOr(name, "default")  // returns "default"
//
//	name = ptr.To("hello")
//	s = ptr.FromOr(name, "default")   // returns "hello"
func FromOr[T any](p *T, defaultValue T) T {
	if p == nil {
		return defaultValue
	}
	return *p
}
