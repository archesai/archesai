package codegen

import (
	"os"
	"path/filepath"

	"github.com/archesai/archesai/internal/located"
	"github.com/archesai/archesai/internal/spec"
)

// GroupApp is the generator group for app bootstrap and wiring.
const GroupApp = "app"

const (
	genAppBootstrap  = "app_bootstrap"
	genAppContainer  = "app_container"
	genAppOperations = "app_operations"
	genAppPorts      = "app_ports"
	genAppHTTP       = "app_http"
)

// AppBootstrapData wraps spec with bundled OpenAPI spec for bootstrap template.
type AppBootstrapData struct {
	*spec.Spec
	OpenAPISpec string
	APITitle    string
}

// generateAppBootstrap generates the app bootstrap file.
func (c *Codegen) generateAppBootstrap(s *located.Located[spec.Spec]) error {
	if s.Value.ShouldSkipWire() {
		return nil
	}

	// Bundle the OpenAPI spec
	var openAPISpec string
	var apiTitle string
	if s.Path != "" {
		baseFS := os.DirFS(s.Dir())
		compositeFS := spec.BuildIncludeFS(baseFS, s.Value.EnabledIncludes)
		doc, err := spec.NewOpenAPIDocumentFromFS(compositeFS, filepath.Base(s.Path))
		if err == nil {
			bundler := spec.NewBundler(doc)
			specBytes, err := bundler.BundleToYAML()
			if err == nil {
				openAPISpec = string(specBytes)
			}
			apiTitle = doc.Raw().Info.Title
		}
	}

	data := &AppBootstrapData{
		Spec:        s.Value,
		OpenAPISpec: openAPISpec,
		APITitle:    apiTitle,
	}

	return c.RenderToFile(genAppBootstrap+".go.tmpl", "app/bootstrap.gen.go", data)
}

// generateAppContainer generates the app container file.
func (c *Codegen) generateAppContainer(s *spec.Spec) error {
	if s.ShouldSkipWire() {
		return nil
	}
	return c.RenderToFile(genAppContainer+".go.tmpl", "app/container.gen.go", s)
}

// generateAppOperations generates the app operations wiring file.
func (c *Codegen) generateAppOperations(s *spec.Spec) error {
	if s.ShouldSkipWire() {
		return nil
	}
	return c.RenderToFile(genAppOperations+".go.tmpl", "app/operations.gen.go", s)
}

// generateAppPorts generates the app ports (interfaces) file.
func (c *Codegen) generateAppPorts(s *spec.Spec) error {
	if s.ShouldSkipWire() {
		return nil
	}
	return c.RenderToFile(genAppPorts+".go.tmpl", "app/ports.gen.go", s)
}

// generateAppHTTP generates the app HTTP wiring file.
func (c *Codegen) generateAppHTTP(s *spec.Spec) error {
	if s.ShouldSkipWire() {
		return nil
	}
	return c.RenderToFile(genAppHTTP+".go.tmpl", "app/http.gen.go", s)
}
