package main

import (
	"encoding/json"
	"errors"
)

// executeFunction is the default execution function that returns an error.
// This file should be replaced by mounting a custom execute.go
// or by building a derived image with a custom implementation.
//
// To use:
//  1. Copy this file and implement your logic
//  2. Mount it to the container:
//     docker run -i --rm -v ./my-execute.go:/app/execute.go archesai/runner-go:latest
//
// Or build a custom image:
//  1. Create a Dockerfile:
//     FROM archesai/runner-go:latest
//     COPY execute.go ./execute.go
//     RUN go build -o /usr/local/bin/runner .
//  2. Build: docker build -t my-generator:latest .
func executeFunction(_ json.RawMessage) (json.RawMessage, error) {
	return nil, errors.New(
		"no execution function provided. " +
			"Either mount a custom execute.go at /app/execute.go " +
			"or use a pre-built generator image (e.g., archesai/generator-custom)",
	)
}
