package generators

import "github.com/archesai/archesai/internal/spec"

// SchemasTemplateData holds data for schema template rendering.
type SchemasTemplateData struct {
	Package string
	Schema  *spec.Schema
}

// ApplicationTemplateData holds data for handler template rendering.
type ApplicationTemplateData struct {
	Operation   *spec.Operation
	ProjectName string
}

// ApplicationStubTemplateData holds data for handler stub template rendering.
type ApplicationStubTemplateData struct {
	Operation   *spec.Operation
	ProjectName string
}

// RouteTemplateData holds data for route template rendering.
type RouteTemplateData struct {
	Operation   *spec.Operation
	ProjectName string
}

// AppBootstrapTemplateData holds data for app bootstrap template rendering.
type AppBootstrapTemplateData struct {
	ProjectName      string
	InternalPackages []InternalPackage
	OpenAPISpec      string
	APITitle         string
}

// AppContainerTemplateData holds data for app container template rendering.
type AppContainerTemplateData struct {
	Entities         []*spec.Schema
	ProjectName      string
	InternalPackages []InternalPackageWithEntities
}

// AppHandlersTemplateData holds data for app handlers template rendering.
type AppHandlersTemplateData struct {
	Operations        []spec.Operation
	Repositories      []string
	ProjectName       string
	NeedsPublisher    bool
	HasCustomHandlers bool
	InternalPackages  []InternalPackage
}

// AppRoutesTemplateData holds data for app routes template rendering.
type AppRoutesTemplateData struct {
	Operations       []spec.Operation
	ProjectName      string
	InternalPackages []InternalPackage
}

// InternalPackageWithEntities is an internal package with its entities and repository info.
type InternalPackageWithEntities struct {
	InternalPackage
	Entities       []*spec.Schema
	Repositories   []string
	NeedsPublisher bool
}
