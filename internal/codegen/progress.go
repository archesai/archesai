package codegen

// ProgressCallback is called when a generator starts or completes.
type ProgressCallback func(event ProgressEvent)

// ProgressEvent represents a progress update from the codegen.
type ProgressEvent struct {
	Type          ProgressEventType
	GeneratorName string
	TotalCount    int
	CurrentIndex  int
	Error         error
}

// ProgressEventType indicates the type of progress event.
type ProgressEventType int

// Progress event types for codegen callbacks.
const (
	ProgressEventStart          ProgressEventType = iota // Generation started
	ProgressEventGeneratorStart                          // Individual generator started
	ProgressEventGeneratorDone                           // Individual generator completed
	ProgressEventDone                                    // All generation completed
	ProgressEventError                                   // Error occurred
)
