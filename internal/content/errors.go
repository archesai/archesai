package content

import "errors"

// Domain errors
var (
	// ErrArtifactNotFound is returned when an artifact is not found
	ErrArtifactNotFound = errors.New("artifact not found")

	// ErrArtifactTooLarge is returned when an artifact exceeds the maximum size
	ErrArtifactTooLarge = errors.New("artifact too large")

	// ErrLabelNotFound is returned when a label is not found
	ErrLabelNotFound = errors.New("label not found")

	// ErrLabelExists is returned when a label already exists
	ErrLabelExists = errors.New("label already exists")
)
