// Package capabilities provides system dependency detection and validation.
package capabilities

// Capability represents a system dependency that can be checked.
type Capability struct {
	Name        string // Human-readable name (e.g., "Go", "Node.js")
	Command     string // Binary/command to check (e.g., "go", "node")
	Version     string // Detected version (empty if not found)
	Required    bool   // Whether this capability is required
	Found       bool   // Whether the capability was detected
	Error       error  // Any error during detection
	Description string // Brief description of what this capability provides
}

// CheckResult contains the results of checking multiple capabilities.
type CheckResult struct {
	Capabilities []Capability
	AllRequired  bool // True if all required capabilities are satisfied
	AllFound     bool // True if all capabilities (required and optional) are found
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
