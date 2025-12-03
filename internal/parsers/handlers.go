package parsers

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// HandlerDef represents a parsed handler file
type HandlerDef struct {
	OperationID string          // e.g., "Login"
	FilePath    string          // e.g., "handlers/login.impl.go"
	Constructor *ConstructorDef // Parsed constructor
	HasHandler  bool            // Whether handler struct exists
}

// ConstructorDef represents a parsed constructor function
type ConstructorDef struct {
	Name       string     // e.g., "NewLoginHandler"
	Parameters []ParamDef // Constructor parameters (reusing ParamDef)
}

// DependencyDef represents a dependency parameter in a constructor
type DependencyDef struct {
	Name       string // e.g., "authService"
	Type       string // e.g., "*auth.Service"
	Package    string // e.g., "github.com/archesai/archesai/pkg/auth"
	IsPointer  bool   // Whether it's a pointer type
	Resolution string // e.g., "infra.AuthService" - how to resolve this dep
}

// HandlerParser parses Go handler files to extract constructor dependencies
type HandlerParser struct {
	handlersDir string
	depRegistry *DependencyRegistry
}

// DependencyRegistry maps type names to their resolution in the bootstrap
type DependencyRegistry struct {
	mappings map[string]DependencyResolution
}

// DependencyResolution describes how to resolve a dependency
type DependencyResolution struct {
	ImportPath string // e.g., "github.com/archesai/archesai/pkg/auth"
	Resolution string // e.g., "infra.AuthService"
}

// NewDependencyRegistry creates a new dependency registry with default mappings
func NewDependencyRegistry() *DependencyRegistry {
	return &DependencyRegistry{
		mappings: map[string]DependencyResolution{
			// Infrastructure services
			"*auth.Service": {
				ImportPath: "github.com/archesai/archesai/pkg/auth",
				Resolution: "infra.AuthService",
			},
			"events.Publisher": {
				ImportPath: "github.com/archesai/archesai/pkg/events",
				Resolution: "infra.EventPublisher",
			},
			"*events.Publisher": {
				ImportPath: "github.com/archesai/archesai/pkg/events",
				Resolution: "infra.EventPublisher",
			},
			"*executor.ExecutorService[map[string]any, map[string]any]": {
				ImportPath: "github.com/archesai/archesai/pkg/executor",
				Resolution: "infra.ExecutorService",
			},
			// Config
			"*config.Config": {
				ImportPath: "", // Same package
				Resolution: "cfg",
			},
		},
	}
}

// Resolve resolves a type to its bootstrap expression
func (r *DependencyRegistry) Resolve(typeName string) (DependencyResolution, bool) {
	// Direct match
	if res, ok := r.mappings[typeName]; ok {
		return res, true
	}

	// Check for repository pattern: repositories.XxxRepository -> repos.Xxxs
	if strings.HasPrefix(typeName, "repositories.") && strings.HasSuffix(typeName, "Repository") {
		entityName := strings.TrimSuffix(
			strings.TrimPrefix(typeName, "repositories."),
			"Repository",
		)
		return DependencyResolution{
			ImportPath: "", // Same package
			Resolution: "repos." + Pluralize(entityName),
		}, true
	}

	return DependencyResolution{}, false
}

// AddMapping adds a custom type mapping
func (r *DependencyRegistry) AddMapping(typeName string, resolution DependencyResolution) {
	r.mappings[typeName] = resolution
}

// NewHandlerParser creates a new handler parser
func NewHandlerParser(handlersDir string) *HandlerParser {
	return &HandlerParser{
		handlersDir: handlersDir,
		depRegistry: NewDependencyRegistry(),
	}
}

// ParseHandlers parses all handler files and extracts their dependencies
func (p *HandlerParser) ParseHandlers(operations []OperationDef) ([]HandlerDef, error) {
	handlers := make([]HandlerDef, 0, len(operations))

	for _, op := range operations {
		handler, err := p.parseHandler(op)
		if err != nil {
			return nil, fmt.Errorf("failed to parse handler for %s: %w", op.ID, err)
		}
		handlers = append(handlers, handler)
	}

	return handlers, nil
}

// parseHandler parses a single handler file
func (p *HandlerParser) parseHandler(op OperationDef) (HandlerDef, error) {
	fileName := SnakeCase(op.ID) + ".impl.go"
	filePath := filepath.Join(p.handlersDir, fileName)

	handler := HandlerDef{
		OperationID: op.ID,
		FilePath:    filePath,
		HasHandler:  false,
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return handler, nil
	}

	// Parse the Go file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return handler, fmt.Errorf("failed to parse %s: %w", filePath, err)
	}

	// Find the constructor function
	constructorName := "New" + op.ID + "Handler"
	constructor := p.findConstructor(node, constructorName)

	if constructor != nil {
		handler.HasHandler = true
		handler.Constructor = constructor
	}

	return handler, nil
}

// findConstructor finds and parses a constructor function in the AST
func (p *HandlerParser) findConstructor(node *ast.File, name string) *ConstructorDef {
	var constructor *ConstructorDef

	ast.Inspect(node, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		if fn.Name.Name != name {
			return true
		}

		// Found the constructor
		constructor = &ConstructorDef{
			Name:       name,
			Parameters: p.extractParameters(fn, node),
		}

		return false // Stop searching
	})

	return constructor
}

// extractParameters extracts parameters from a function declaration
func (p *HandlerParser) extractParameters(fn *ast.FuncDecl, file *ast.File) []ParamDef {
	if fn.Type.Params == nil {
		return nil
	}

	// Build import map for resolving package names
	imports := make(map[string]string) // alias/name -> full path
	for _, imp := range file.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		name := filepath.Base(path)
		if imp.Name != nil {
			name = imp.Name.Name
		}
		imports[name] = path
	}

	var params []ParamDef
	for _, field := range fn.Type.Params.List {
		typeStr := p.typeToString(field.Type)

		// For each name in the field (handles "a, b int" syntax)
		for _, name := range field.Names {
			param := ParamDef{
				SchemaDef: &SchemaDef{
					Name:   name.Name,
					GoType: typeStr,
				},
				In: "constructor", // Mark as constructor param
			}
			params = append(params, param)
		}
	}

	return params
}

// typeToString converts an AST type expression to a string
func (p *HandlerParser) typeToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + p.typeToString(t.X)
	case *ast.SelectorExpr:
		return p.typeToString(t.X) + "." + t.Sel.Name
	case *ast.ArrayType:
		return "[]" + p.typeToString(t.Elt)
	case *ast.MapType:
		return "map[" + p.typeToString(t.Key) + "]" + p.typeToString(t.Value)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.IndexExpr:
		// Generic type: Type[Param]
		return p.typeToString(t.X) + "[" + p.typeToString(t.Index) + "]"
	case *ast.IndexListExpr:
		// Generic type with multiple params: Type[Param1, Param2]
		var params []string
		for _, idx := range t.Indices {
			params = append(params, p.typeToString(idx))
		}
		return p.typeToString(t.X) + "[" + strings.Join(params, ", ") + "]"
	default:
		return fmt.Sprintf("%T", expr)
	}
}

// GetDependencies returns the resolved dependencies for a handler
func (p *HandlerParser) GetDependencies(handler HandlerDef) []DependencyDef {
	if handler.Constructor == nil {
		return nil
	}

	var deps []DependencyDef
	for _, param := range handler.Constructor.Parameters {
		resolution, found := p.depRegistry.Resolve(param.GoType)

		dep := DependencyDef{
			Name:      param.Name,
			Type:      param.GoType,
			IsPointer: strings.HasPrefix(param.GoType, "*"),
		}

		if found {
			dep.Package = resolution.ImportPath
			dep.Resolution = resolution.Resolution
		} else {
			// Unknown dependency - leave Resolution empty, will need manual wiring
			dep.Resolution = "/* TODO: resolve " + param.GoType + " */"
		}

		deps = append(deps, dep)
	}

	return deps
}
