// Package generators provides code generation from OpenAPI specifications.
package generators

import (
	"sort"

	"github.com/archesai/archesai/internal/spec"
)

// AppContext encapsulates the standalone vs composition decision
// and provides common data needed by all app generators.
type AppContext struct {
	// IsStandalone is true if the project has its own entities.
	IsStandalone bool

	// IsComposition is true if the project composes other packages.
	IsComposition bool

	// Entities are the entity schemas owned by this project.
	// Only populated for standalone apps.
	Entities []*spec.Schema

	// Operations are the operations owned by this project.
	// Only populated for standalone apps.
	Operations []spec.Operation

	// InternalPackages are the packages being composed.
	// Only populated for composition apps.
	InternalPackages []InternalPackage

	// InternalPackagesWithEntities includes repository info per package.
	// Only populated for composition apps (container generator).
	InternalPackagesWithEntities []InternalPackageWithEntities

	// Repositories are the unique repository names needed.
	// Only populated for standalone apps.
	Repositories []string

	// NeedsPublisher is true if any non-GET operation exists.
	// Only populated for standalone apps.
	NeedsPublisher bool
}

// BuildAppContext creates the shared context for all app generators.
// It determines whether this is a standalone app or composition app and
// computes all common data needed by app templates.
func BuildAppContext(ctx *GeneratorContext) *AppContext {
	actx := &AppContext{}

	// Get own operations and entities
	ownOperations := ctx.OwnOperations()
	entities := ctx.OwnEntitySchemas()

	// Standalone app: has own operations OR own entities
	if len(ownOperations) > 0 || len(entities) > 0 {
		actx.IsStandalone = true
		actx.Entities = entities
		actx.Operations = ownOperations

		// Sort operations by ID
		sort.Slice(actx.Operations, func(i, j int) bool {
			return actx.Operations[i].ID < actx.Operations[j].ID
		})

		// Compute repositories and publisher need
		actx.Repositories, actx.NeedsPublisher = ExtractRepositoryInfo(actx.Operations)
		return actx
	}

	// Check for composition app (composes other packages)
	composedPkgs := ctx.ComposedPackages()
	if len(composedPkgs) > 0 {
		actx.IsComposition = true
		actx.InternalPackages = BuildInternalPackages(composedPkgs)
		actx.InternalPackagesWithEntities = BuildInternalPackagesWithEntities(ctx, composedPkgs)
		return actx
	}

	// Neither standalone nor composition
	return actx
}

// BuildInternalPackages creates InternalPackage entries for the given package names.
func BuildInternalPackages(pkgNames []string) []InternalPackage {
	var packages []InternalPackage
	for _, name := range pkgNames {
		packages = append(packages, InternalPackage{
			Name:       name,
			Alias:      name,
			ImportPath: InternalPackageImportPath(name),
		})
	}
	sort.Slice(packages, func(i, j int) bool {
		return packages[i].Name < packages[j].Name
	})
	return packages
}

// BuildInternalPackagesWithEntities creates InternalPackageWithEntities entries
// with repository info extracted from operations.
func BuildInternalPackagesWithEntities(
	ctx *GeneratorContext,
	pkgNames []string,
) []InternalPackageWithEntities {
	// Group repositories by their x-internal package (from operations' tags)
	reposByPkg := make(map[string]map[string]bool)
	pkgNeedsPublisher := make(map[string]bool)

	for _, op := range ctx.Spec.Operations {
		if op.Internal == "" {
			continue
		}
		if op.CustomHandler {
			continue
		}

		// Initialize map for this package if needed
		if reposByPkg[op.Internal] == nil {
			reposByPkg[op.Internal] = make(map[string]bool)
		}

		// Add repo for this tag
		reposByPkg[op.Internal][op.Tag] = true

		// Check if this package needs a publisher
		if op.Method != "GET" {
			pkgNeedsPublisher[op.Internal] = true
		}
	}

	var packages []InternalPackageWithEntities
	for _, pkgName := range pkgNames {
		// Get sorted repositories for this package
		var repos []string
		for repo := range reposByPkg[pkgName] {
			repos = append(repos, repo)
		}
		sort.Strings(repos)

		packages = append(packages, InternalPackageWithEntities{
			InternalPackage: InternalPackage{
				Name:       pkgName,
				Alias:      pkgName,
				ImportPath: InternalPackageImportPath(pkgName),
			},
			Repositories:   repos,
			NeedsPublisher: pkgNeedsPublisher[pkgName],
		})
	}
	sort.Slice(packages, func(i, j int) bool {
		return packages[i].Name < packages[j].Name
	})

	return packages
}

// ExtractRepositoryInfo extracts unique repository names and determines
// if a publisher is needed from the given operations.
func ExtractRepositoryInfo(operations []spec.Operation) (repos []string, needsPublisher bool) {
	repoMap := make(map[string]bool)
	for _, op := range operations {
		if op.CustomHandler {
			continue
		}
		repoMap[op.Tag] = true

		// Check if we need a publisher (any non-GET, non-custom operations)
		if op.Method != "GET" {
			needsPublisher = true
		}
	}

	for repo := range repoMap {
		repos = append(repos, repo)
	}
	sort.Strings(repos)
	return repos, needsPublisher
}

// ExtractRepositories extracts unique repository names from operations.
// This is a simpler version that just returns repository names.
func ExtractRepositories(operations []spec.Operation) []string {
	repos, _ := ExtractRepositoryInfo(operations)
	return repos
}

// NeedsPublisher determines if any operation requires a publisher.
// An operation needs a publisher if it's not a GET and not a custom handler.
func NeedsPublisher(operations []spec.Operation) bool {
	for _, op := range operations {
		if !op.CustomHandler && op.Method != "GET" {
			return true
		}
	}
	return false
}

// ShouldSkip returns true if there's nothing to generate for app.
func (actx *AppContext) ShouldSkip() bool {
	return !actx.IsStandalone && !actx.IsComposition
}

// HasCustomHandlers returns true if any operation has a custom handler.
func HasCustomHandlers(operations []spec.Operation) bool {
	for _, op := range operations {
		if op.CustomHandler {
			return true
		}
	}
	return false
}
