package main

import (
	"encoding/json"
	"errors"
	"math"
)

// Example custom execute function
//
// This file shows how to create a custom execute module for the Go runner.
//
// To use:
//  1. Copy this file and implement your logic
//  2. Rename it to execute.go
//  3. Mount it to the container:
//     docker run -i --rm -v ./execute.go:/app/execute.go archesai/runner-go:latest
//
// Or build a custom image:
//  1. Create a Dockerfile:
//     FROM archesai/runner-go:latest
//     COPY execute.go ./execute.go
//     RUN go build -o /usr/local/bin/runner .
//  2. Build: docker build -t my-generator:latest .

// ExampleInput represents the expected input structure
type ExampleInput struct {
	Values []float64 `json:"values"`
}

// ExampleOutput represents the output structure
type ExampleOutput struct {
	Count int     `json:"count"`
	Sum   float64 `json:"sum"`
	Mean  float64 `json:"mean"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
}

// Example_executeFunction is the main execution function.
// Implement your custom logic here.
func Example_executeFunction(input json.RawMessage) (json.RawMessage, error) {
	// Example: Simple data transformation
	var inputData ExampleInput
	if err := json.Unmarshal(input, &inputData); err != nil {
		return nil, err
	}

	if len(inputData.Values) == 0 {
		return nil, errors.New("expected 'values' array in input")
	}

	values := inputData.Values
	sum := 0.0
	min := math.Inf(1)
	max := math.Inf(-1)

	for _, v := range values {
		sum += v
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	mean := sum / float64(len(values))

	output := ExampleOutput{
		Count: len(values),
		Sum:   sum,
		Mean:  mean,
		Min:   min,
		Max:   max,
	}

	return json.Marshal(output)
}

/*
Example usage:

Input:
{
  "input": {
    "values": [1, 2, 3, 4, 5]
  }
}

Output:
{
  "ok": true,
  "output": {
    "count": 5,
    "sum": 15,
    "mean": 3,
    "min": 1,
    "max": 5
  }
}
*/
