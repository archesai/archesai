package validation

// Validator is the interface that types implement to provide validation logic.
type Validator interface {
	// Validate performs validation and returns all errors found.
	Validate() Errors
}

// ValidateStruct validates a struct if it implements the Validator interface.
// Returns nil if the struct does not implement Validator or if validation passes.
func ValidateStruct[T any](v *T) Errors {
	if v == nil {
		return nil
	}
	if validator, ok := any(v).(Validator); ok {
		errs := validator.Validate()
		if errs.HasErrors() {
			return errs
		}
	}
	return nil
}
