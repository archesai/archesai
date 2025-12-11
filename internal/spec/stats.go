package spec

import "fmt"

// Stats holds statistics about an OpenAPI document.
type Stats struct {
	Title                string
	Version              string
	TotalPaths           int
	TotalOperations      int
	TotalSchemas         int
	TotalParameters      int
	TotalResponses       int
	TotalSecuritySchemes int
}

// GetStats computes and returns statistics about the parsed OpenAPI document.
func (p *Parser) GetStats() (*Stats, error) {
	if p.specResult == nil {
		return nil, fmt.Errorf("spec not parsed, call Parse() first")
	}

	// Count paths (unique paths from operations)
	pathSet := make(map[string]struct{})
	for _, op := range p.specResult.Operations {
		pathSet[op.Path] = struct{}{}
	}

	// Count components using discovery
	schemas, _ := DiscoverComponents(p.doc.fsys, ComponentSchemas)
	responses, _ := DiscoverComponents(p.doc.fsys, ComponentResponses)
	parameters, _ := DiscoverComponents(p.doc.fsys, ComponentParameters)

	stats := &Stats{
		Title:                p.specResult.Title,
		Version:              p.specResult.Version,
		TotalPaths:           len(pathSet),
		TotalOperations:      len(p.specResult.Operations),
		TotalSchemas:         len(schemas),
		TotalParameters:      len(parameters),
		TotalResponses:       len(responses),
		TotalSecuritySchemes: len(p.specResult.Security),
	}

	return stats, nil
}
