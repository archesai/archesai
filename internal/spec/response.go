package spec

import (
	"sort"
	"strconv"
)

// ResponseDef represents a response in an operation
type ResponseDef struct {
	*Schema                        // Embed schema definition for response body
	StatusCode  string             // HTTP status code
	ContentType string             // Content-Type for the response (e.g., "application/json")
	Headers     map[string]*Schema // Response headers
}

// IsSuccess returns true if the response is a successful one (2xx status code)
func (r *ResponseDef) IsSuccess() bool {
	if code, err := strconv.Atoi(r.StatusCode); err == nil {
		return code >= 200 && code < 300
	}
	return false
}

// GetSortedHeaders returns headers sorted by name for consistent iteration
func (r *ResponseDef) GetSortedHeaders() []struct {
	Name   string
	Schema *Schema
} {
	if r.Headers == nil {
		return nil
	}

	// Create a slice of header names and sort them
	var names []string
	for name := range r.Headers {
		names = append(names, name)
	}
	sort.Strings(names)

	// Build the sorted slice
	var sorted []struct {
		Name   string
		Schema *Schema
	}
	for _, name := range names {
		sorted = append(sorted, struct {
			Name   string
			Schema *Schema
		}{
			Name:   name,
			Schema: r.Headers[name],
		})
	}
	return sorted
}
