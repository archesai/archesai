// Package handlers provides code generation for HTTP handler stubs
package handlers

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"sync"
)

var uppercaseAcronym = sync.Map{}

//"ID": "id",

// ConfigureAcronym allows you to add additional words which will be considered acronyms
func ConfigureAcronym(key, val string) {
	uppercaseAcronym.Store(key, val)
}

// Converts a string to CamelCase
func toCamelInitCase(s string, initCase bool) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	a, hasAcronym := uppercaseAcronym.Load(s)
	if hasAcronym {
		s = a.(string)
	}

	n := strings.Builder{}
	n.Grow(len(s))
	capNext := initCase
	prevIsCap := false
	for i, v := range []byte(s) {
		vIsCap := v >= 'A' && v <= 'Z'
		vIsLow := v >= 'a' && v <= 'z'
		if capNext {
			if vIsLow {
				v += 'A'
				v -= 'a'
			}
		} else if i == 0 {
			if vIsCap {
				v += 'a'
				v -= 'A'
			}
		} else if prevIsCap && vIsCap && !hasAcronym {
			v += 'a'
			v -= 'A'
		}
		prevIsCap = vIsCap

		if vIsCap || vIsLow {
			n.WriteByte(v)
			capNext = false
		} else if vIsNum := v >= '0' && v <= '9'; vIsNum {
			n.WriteByte(v)
			capNext = true
		} else {
			capNext = v == '_' || v == ' ' || v == '-' || v == '.'
		}
	}
	return n.String()
}

// ToCamel converts a string to CamelCase
func ToCamel(s string) string {
	return toCamelInitCase(s, true)
}

// ToLowerCamel converts a string to lowerCamelCase
func ToLowerCamel(s string) string {
	return toCamelInitCase(s, false)
}

// Method represents a handler method to generate
type Method struct {
	Name          string
	ServiceMethod string
	HTTPMethod    string
	Path          string
	HasRequest    bool
	HasParams     bool
	ReturnType    string
}

// Domain represents a domain package
type Domain struct {
	Name    string
	Package string
	Methods []Method
}

// Generator generates handler implementations
type Generator struct{}

// NewGenerator creates a new handler generator
func NewGenerator() *Generator {
	return &Generator{}
}

// Generate generates handlers for all domains
func (g *Generator) Generate() error {
	domains := []string{"auth", "organizations", "workflows", "content", "health", "config"}

	for _, domain := range domains {
		if err := g.generateDomain(domain); err != nil {
			return fmt.Errorf("failed to generate %s: %w", domain, err)
		}
	}

	return nil
}

// generateDomain generates handler for a specific domain
func (g *Generator) generateDomain(domain string) error {
	// Parse the handler.gen.go file to extract interface methods
	handlerGenPath := filepath.Join("internal", domain, "adapters", "http", "handler.gen.go")

	// Check if file exists
	if _, err := os.Stat(handlerGenPath); os.IsNotExist(err) {
		log.Printf("Skipping %s: handler.gen.go not found", domain)
		return nil
	}

	methods, err := g.parseHandlerInterface(handlerGenPath)
	if err != nil {
		return fmt.Errorf("failed to parse handler interface: %w", err)
	}

	if len(methods) == 0 {
		log.Printf("No methods found for %s", domain)
		return nil
	}

	// Generate handler implementation
	handlerPath := filepath.Join("internal", domain, "adapters", "http", "handler.go")

	// Check if handler already exists
	if _, err := os.Stat(handlerPath); err == nil {
		log.Printf("Handler already exists for %s, skipping", domain)
		return nil
	}

	d := Domain{
		Name:    ToCamel(domain),
		Package: domain,
		Methods: methods,
	}

	return g.writeHandler(handlerPath, d)
}

// parseHandlerInterface parses the StrictServerInterface to extract methods
func (g *Generator) parseHandlerInterface(path string) ([]Method, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var methods []Method

	// Find StrictServerInterface
	ast.Inspect(file, func(n ast.Node) bool {
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		// Check if it's StrictServerInterface
		if !strings.HasSuffix(typeSpec.Name.Name, "StrictServerInterface") {
			return true
		}

		interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
		if !ok {
			return true
		}

		// Extract methods
		for _, method := range interfaceType.Methods.List {
			funcType, ok := method.Type.(*ast.FuncType)
			if !ok {
				continue
			}

			methodName := method.Names[0].Name

			// Skip if it's a standard CRUD operation we'll handle
			m := Method{
				Name:          methodName,
				ServiceMethod: g.inferServiceMethod(methodName),
				HasRequest:    g.hasRequestBody(funcType),
				HasParams:     g.hasPathParams(funcType),
				ReturnType:    g.getReturnType(methodName),
			}

			methods = append(methods, m)
		}

		return false
	})

	return methods, nil
}

// inferServiceMethod infers the service method name from handler method
func (g *Generator) inferServiceMethod(handlerMethod string) string {
	// Simple mapping - can be enhanced
	switch {
	case strings.HasPrefix(handlerMethod, "Create"):
		return handlerMethod
	case strings.HasPrefix(handlerMethod, "Get"):
		return handlerMethod
	case strings.HasPrefix(handlerMethod, "List"):
		return handlerMethod
	case strings.HasPrefix(handlerMethod, "Update"):
		return handlerMethod
	case strings.HasPrefix(handlerMethod, "Delete"):
		return handlerMethod
	default:
		return handlerMethod
	}
}

// hasRequestBody checks if method has a request body
func (g *Generator) hasRequestBody(funcType *ast.FuncType) bool {
	for _, param := range funcType.Params.List {
		for _, name := range param.Names {
			if strings.Contains(name.Name, "request") || strings.Contains(name.Name, "body") {
				return true
			}
		}
	}
	return false
}

// hasPathParams checks if method has path parameters
func (g *Generator) hasPathParams(funcType *ast.FuncType) bool {
	for _, param := range funcType.Params.List {
		for _, name := range param.Names {
			if strings.Contains(name.Name, "params") || strings.Contains(name.Name, "id") {
				return true
			}
		}
	}
	return false
}

// getReturnType gets the response type for a method
func (g *Generator) getReturnType(methodName string) string {
	return methodName + "ResponseObject"
}

// writeHandler writes the handler implementation
func (g *Generator) writeHandler(path string, domain Domain) error {
	tmpl := `// Package http provides HTTP handlers for {{ .Package }} operations
package http

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/archesai/archesai/internal/{{ .Package }}"
	"github.com/google/uuid"
)

// {{ .Name }}Handler handles {{ .Package }} operations
type {{ .Name }}Handler struct {
	service *{{ .Package }}.{{ .Name }}Service
	logger  *slog.Logger
}

// New{{ .Name }}Handler creates a new {{ .Package }} handler
func New{{ .Name }}Handler(service *{{ .Package }}.{{ .Name }}Service, logger *slog.Logger) *{{ .Name }}Handler {
	return &{{ .Name }}Handler{
		service: service,
		logger:  logger,
	}
}

{{ range .Methods }}
// {{ .Name }} implements the {{ .Name }} endpoint
func (h *{{ $.Name }}Handler) {{ .Name }}(ctx context.Context{{ if .HasParams }}, params {{ .Name }}Params{{ end }}{{ if .HasRequest }}, request {{ .Name }}RequestObject{{ end }}) ({{ .ReturnType }}, error) {
	h.logger.Debug("{{ .Name }} called")

	// TODO: Implement {{ .Name }}
	// Example implementation:
	{{ if eq .Name "GetHealth" }}
	status := h.service.CheckHealth(ctx)
	response := &{{ $.Package }}.HealthResponse{
		// Map status fields
	}
	return {{ .Name }}200JSONResponse(*response), nil
	{{ else if eq .Name "GetConfig" }}
	// Return current configuration
	return {{ .Name }}200JSONResponse(*h.service.GetConfig()), nil
	{{ else }}
	// Call service method
	// result, err := h.service.{{ .ServiceMethod }}(ctx, ...)
	// if err != nil {
	//     return nil, err
	// }
	// return {{ .Name }}200JSONResponse(result), nil
	return nil, fmt.Errorf("not implemented")
	{{ end }}
}
{{ end }}
`

	t, err := template.New("handler").Parse(tmpl)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, domain); err != nil {
		return err
	}

	// Write file
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	return os.WriteFile(path, buf.Bytes(), 0644)
}

// Run executes the generator
func Run() error {
	g := NewGenerator()
	return g.Generate()
}
