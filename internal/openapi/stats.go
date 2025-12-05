package openapi

import "fmt"

// Stats holds statistics about an OpenAPI document.
type Stats struct {
	TotalPaths           int
	TotalSchemas         int
	TotalParameters      int
	TotalResponses       int
	TotalSecuritySchemes int
}

// GetStats computes and returns statistics about the OpenAPI document.
// It should return the number of paths, operations, schemas, parameters, responses, and security schemes.
func (p *Parser) GetStats() error {
	index := p.doc.GetIndex()
	numSchemas := len(index.GetAllSchemas())
	numPaths := len(index.GetAllPaths())
	numParameters := len(index.GetAllParameters())
	numResponses := len(index.GetAllResponses())
	numSecuritySchemes := len(index.GetAllSecuritySchemes())

	stats := &Stats{
		TotalPaths:           numPaths,
		TotalSchemas:         numSchemas,
		TotalParameters:      numParameters,
		TotalResponses:       numResponses,
		TotalSecuritySchemes: numSecuritySchemes,
	}

	// Output stats (this could be extended to return structured data)
	fmt.Printf("OpenAPI Specification Statistics:\n")
	fmt.Printf("  Title: %s\n", p.doc.Info.Title)
	fmt.Printf("  Path: %s\n", p.doc.GetIndex().GetConfig().BasePath)
	fmt.Printf("  Total Paths: %d\n", stats.TotalPaths)
	fmt.Printf("  Total Schemas: %d\n", stats.TotalSchemas)
	fmt.Printf("  Total Parameters: %d\n", stats.TotalParameters)
	fmt.Printf("  Total Responses: %d\n", stats.TotalResponses)
	fmt.Printf("  Total Security Schemes: %d\n", stats.TotalSecuritySchemes)

	fmt.Printf("  OpenAPI Version: %s\n", p.doc.Version)

	return nil
}
