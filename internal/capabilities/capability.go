// Package capabilities provides system dependency detection and validation.
package capabilities

// Capability represents a system dependency that can be checked.
type Capability struct {
	// Name is the human-readable name (e.g., "Go", "Node.js")
	Name string `json:"name"`
	// Command is the binary/command to check (e.g., "go", "node")
	Command string `json:"command"`
	// Version is the detected version (empty if not found)
	Version string `json:"version,omitempty"`
	// Required indicates whether this capability is required
	Required bool `json:"required"`
	// Found indicates whether the capability was detected
	Found bool `json:"found"`
	// Error is any error during detection
	Error error `json:"error,omitempty"`
	// Description is a brief description of what this capability provides
	Description string `json:"description"`
}

// CheckResult contains the results of checking multiple capabilities.
type CheckResult struct {
	Capabilities []Capability `json:"capabilities"`
	// AllRequired is true if all required capabilities are satisfied
	AllRequired bool `json:"allRequired"`
	// AllFound is true if all capabilities (required and optional) are found
	AllFound bool `json:"allFound"`
}

// RequiredMissing returns capabilities that are required but not found.
func (r *CheckResult) RequiredMissing() []Capability {
	var missing []Capability
	for _, c := range r.Capabilities {
		if c.Required && !c.Found {
			missing = append(missing, c)
		}
	}
	return missing
}

// OptionalMissing returns capabilities that are optional but not found.
func (r *CheckResult) OptionalMissing() []Capability {
	var missing []Capability
	for _, c := range r.Capabilities {
		if !c.Required && !c.Found {
			missing = append(missing, c)
		}
	}
	return missing
}

// FoundCapabilities returns all capabilities that were found.
func (r *CheckResult) FoundCapabilities() []Capability {
	var found []Capability
	for _, c := range r.Capabilities {
		if c.Found {
			found = append(found, c)
		}
	}
	return found
}
