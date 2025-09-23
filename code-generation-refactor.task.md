# Code Generation System Refactor

## Overview

Refactor the code generation system to have cleaner architecture with proper separation of concerns. Move template data structures to the templates package, consolidate generators into a single file, and ensure each layer has a single responsibility.

## Architecture Goals

### Three Clean Layers

1. **Parsing Layer** (`internal/parsers/`) - Only parses OpenAPI, no template knowledge
2. **Generation Layer** (`internal/codegen/`) - Transforms parsed data to template data
3. **Template Layer** (`internal/templates/`) - Defines contracts and renders output

### Data Flow

```
OpenAPI Spec
    ↓
[Parser] → Raw Schemas & Operations
    ↓
[Generator] → Standardized Template Data
    ↓
[Template] → Generated Code
```

## Tasks

### Phase 1: Create Template Data Structures ✅

#### Why This Matters

Currently, each generator creates anonymous structs inline, making it hard to know what data templates expect. We need a single source of truth for template contracts.

- [x] **Create `internal/templates/data.go` with base structures**

  **Base Structures:**
  - [x] `TemplateData` - The foundation all templates build on
    - `Package string` - The Go package name (e.g., "users", "accounts")
    - `Domain string` - The logical domain (might differ from package)
    - `Imports []string` - Required imports for generated code
    - Why: Every generated file needs these basics

  - [x] `OperationData` - Standardizes how we represent API operations
    - Should include everything from `parsers.Operation` but in template-friendly format
    - Fields like `Name`, `Method`, `Path`, `Description`
    - Parameter data transformed with proper Go types
    - Why: Templates shouldn't need to transform data, just use it

  - [x] `EntityData` - Represents a domain entity/model
    - Entity name in various forms (singular, plural, lower, upper)
    - Fields collection with proper type information
    - Relationships to other entities
    - Why: Multiple templates need entity information in consistent format

  - [x] `FieldData` - Represents a single field/property
    - Go-friendly field name (PascalCase)
    - Go type (string, int, uuid.UUID, etc.)
    - JSON/YAML tags
    - Validation rules
    - Database column info
    - Why: Fields need consistent representation across types, repositories, and services

- [x] **Add specialized template data structures**

  Each inherits from `TemplateData` and adds specific needs:
  - [x] `TypesData` - For type generation (extends TemplateData)
    - `Schemas []SchemaData` - All schemas to generate
    - `Constants []ConstantData` - Enums and constants
    - `TypeAliases []TypeAliasData` - Type definitions

  - [x] `HandlerData` - For HTTP handlers (extends TemplateData)
    - `Operations []OperationData` - HTTP operations to handle
    - `Middleware []string` - Required middleware
    - `HasFileUpload bool` - Whether multipart handling needed

  - [x] `ServiceData` - For service layer (extends TemplateData)
    - `Entities []EntityData` - Entities to create services for
    - `Methods []MethodData` - Service methods to generate

  - [x] `RepositoryData` - For repository interfaces (extends TemplateData)
    - `Entities []EntityData` - Entities needing repositories
    - `Operations []string` - CRUD operations to include
    - `CustomMethods []CustomMethodData` - Additional repository methods

  - [x] `EventsData` - For event publishers (extends TemplateData)
    - `Events []EventData` - Events to publish
    - `EventTypes []string` - Types of events (created, updated, deleted)

### Phase 2: Consolidate Generators ✅

#### Current Problem

We have 7+ generator files that all follow the same pattern. This is unnecessary duplication.

- [x] **Create `internal/codegen/generators.go`**

  A single file with a generic generator that can handle all types:

  ```go
  // GenericGenerator handles all template-based code generation
  type GenericGenerator struct {
      Name         string                    // "types", "service", "repository", etc.
      TemplateName string                    // "types.tmpl", "service.tmpl", etc.
      OutputFile   string                    // "types.gen.go", "service.gen.go", etc.
      PrepareFunc  DataPreparer              // Function to prepare template data
      Filter       SchemaFilter              // Which schemas to process
      Enabled      func(*Config) bool        // Check if generator is enabled
      UseOperations bool                     // Use operations instead of schemas
  }

  // DataPreparer transforms parsed data into template data
  type DataPreparer func(domain string, schemas []*parsers.Schema, operations []parsers.Operation) interface{}

  // SchemaFilter determines which schemas to include
  type SchemaFilter func(*parsers.Schema) bool
  ```

- [x] **Define all generators in one map**

  ```go
  var generators = map[string]*GenericGenerator{
      "types": {
          Name:         "types",
          TemplateName: "types.tmpl",
          OutputFile:   "types.gen.go",
          PrepareFunc:  prepareTypesData,
          Filter:       func(s *parsers.Schema) bool { return true }, // All schemas
          Enabled:      func(c *Config) bool { return c.Generators.Types != "" },
      },
      "service": {
          Name:         "service",
          TemplateName: "service.tmpl",
          OutputFile:   "service.gen.go",
          PrepareFunc:  prepareServiceData,
          Filter:       func(s *parsers.Schema) bool { return s.NeedsService() },
          Enabled:      func(c *Config) bool { return c.Generators.Service != nil },
      },
      // ... etc for all generators
  }
  ```

- [x] **Implement single Generate method**

  ```go
  func (g *GenericGenerator) Generate(ctx *GeneratorContext) error {
      if !g.Enabled(ctx.Config) {
          return nil
      }

      // Group schemas by domain with filter
      domainSchemas := groupByDomain(ctx.Schemas, g.Filter)

      // Get operations if needed
      operations := extractOperations(ctx) if g.UseOperations

      for domain, schemas := range domainSchemas {
          // Prepare data using the specific preparer
          data := g.PrepareFunc(domain, schemas, operations)

          // Write using template
          if err := writeTemplate(ctx, domain, g.TemplateName, g.OutputFile, data); err != nil {
              return fmt.Errorf("%s generation failed: %w", g.Name, err)
          }
      }

      return nil
  }
  ```

- [x] **Move all prepare functions to generators.go**
  - [x] `prepareTypesData()` - Transform schemas to TypesData
  - [x] `prepareServiceData()` - Transform to ServiceData
  - [x] `prepareRepositoryData()` - Transform to RepositoryData
  - [x] `prepareHandlerData()` - Transform operations to HandlerData
  - [x] `prepareEventsData()` - Transform to EventsData

- [x] **Handle special cases**

  Repository generator (multiple output files):

  ```go
  "repository": {
      Name:         "repository",
      PrepareFunc:  prepareRepositoryData,
      CustomGenerate: func(ctx *GeneratorContext, domain string, data interface{}) error {
          // Generate interface
          writeTemplate(ctx, domain, "repository.tmpl", "repository.gen.go", data)

          // Generate PostgreSQL if enabled
          if ctx.Config.Generators.Repository.Postgres != "" {
              writeTemplate(ctx, domain, "repository_postgres.tmpl",
                           ctx.Config.Generators.Repository.Postgres, data)
          }

          // Generate SQLite if enabled
          if ctx.Config.Generators.Repository.Sqlite != "" {
              writeTemplate(ctx, domain, "repository_sqlite.tmpl",
                           ctx.Config.Generators.Repository.Sqlite, data)
          }
      },
  }
  ```

- [x] **Delete individual generator files**
  - [x] Remove `generator_types.go`
  - [x] Remove `generator_service.go`
  - [x] Remove `generator_repository.go`
  - [x] Remove `generator_echo.go`
  - [x] Remove `generator_events.go`
  - [x] Remove `generator_cache.go`
  - [ ] Remove `generator_sql.go` (kept due to complex logic)

### Phase 3: Clean Up Parsers Package ✅

#### Current Problem

The parsers package is doing too much - it's parsing AND preparing template data. This violates single responsibility.

- [x] **Remove template preparation methods from `internal/parsers/openapi.go`**
  - [x] Remove `PrepareTypesData()` - Move logic to `prepareTypesData()` in generators.go
  - [x] Remove `PrepareEchoServerData()` - Move to `prepareHandlerData()` in generators.go
  - [x] Remove any other template-specific methods
  - Why: Parser should only extract raw data from OpenAPI, not format it for templates

- [x] **Fix `parsers.Operation` structure**
  - [x] Add `Type` field to `OperationParam`

    ```go
    // Current (missing type info)
    type OperationParam struct {
        Name     string
        In       string  // path, query, header
        Required bool
        // Missing: Type field!
    }

    // Should be:
    type OperationParam struct {
        Name     string
        In       string
        Type     string  // "string", "int", "uuid.UUID" - inferred from OpenAPI
        Format   string  // "uuid", "email", "date-time"
        Required bool
    }
    ```

  - [x] Ensure `ResponseType` is populated from response schemas
    - Extract from response → content → application/json → schema → $ref
    - Set on Operation during extraction in ExtractOperations()

- [x] **Keep parsers focused on parsing only**
  - [x] `parser.go` - Main orchestrator, loads OpenAPI
  - [x] `openapi.go` - Extracts operations from paths
  - [x] `jsonschema.go` - Parses schemas into raw data
  - [x] `operation.go` - Raw operation data structures
  - [x] `extension.go` - x-codegen extension definitions
  - [x] No template knowledge or data transformation

### Phase 4: Update Templates ✅

#### Current Problem

Templates don't document what they expect, leading to runtime errors when data structure doesn't match.

- [x] **Add documentation headers to each template**

  ```go
  {{- /*
  Template: echo_server.tmpl
  Expects: templates.HandlerData
  Generates: Echo framework HTTP handlers

  Expected structure:
  - Package: string (package name)
  - Domain: string (domain name)
  - Operations: []OperationData
    - Name: string (operation ID)
    - Method: string (GET, POST, etc.)
    - Path: string (URL path)
    - HasRequestBody: bool
    - RequestBodySchema: string (type name)
    - ResponseType: string (response type name)
    - PathParams: []ParamData
    - QueryParams: []ParamData
  */ -}}
  ```

- [x] **Update field references in templates**
  - [x] `echo_server.tmpl` - Use HandlerData fields consistently
  - [x] `types.tmpl` - Use TypesData fields
  - [x] `service.tmpl` - Use ServiceData fields
  - [x] `repository.tmpl` - Use RepositoryData fields
  - [x] `repository_postgres.tmpl` - Use RepositoryData fields
  - [x] `repository_sqlite.tmpl` - Use RepositoryData fields
  - [x] `events.tmpl` - Use EventsData fields

- [x] **Ensure consistent field naming**
  - [x] Use `HasRequestBody` not `RequestBody` for booleans
  - [x] Use `RequestBodySchema` not `RequestBodyType` for type names
  - [x] Document any template-specific fields

### Phase 5: Testing & Validation ✅

#### Why This Matters

We need to ensure our refactor doesn't break existing functionality.

- [x] **Clean slate test**

  ```bash
  make clean-generated
  make generate
  ```

  Should generate all files without errors

- [x] **Compilation test**

  ```bash
  go build ./...
  ```

  All generated code should compile

- [x] **Completeness test - verify each domain has:**
  - [x] `types.gen.go` - Type definitions present
  - [x] `repository.gen.go` - Repository interface present
  - [x] `repository_postgres.gen.go` - If PostgreSQL enabled
  - [x] `repository_sqlite.gen.go` - If SQLite enabled
  - [x] `service.gen.go` - Service implementation present
  - [x] `handler.gen.go` - HTTP handlers present

- [x] **Content validation**
  - [x] Operations have correct HTTP methods
  - [x] Types have proper JSON tags
  - [x] Repositories have CRUD methods
  - [x] Services have business logic methods
  - [x] Handlers map to correct paths

- [x] **Test with different configurations**
  - [x] Test with only types generation
  - [x] Test with PostgreSQL only
  - [x] Test with SQLite only
  - [x] Test with all generators enabled

### Phase 6: Update Documentation ⏳ (Partially Complete)

- [ ] **Update `docs/guides/code-generation.md`**
  - [ ] **Add Architecture Section**

    ```markdown
    ## Architecture

    The code generation system follows clean architecture principles:

    ### Layer 1: Parsing (internal/parsers/)

    Responsible for reading and parsing OpenAPI specifications:

    - `parser.go` - Orchestrates parsing, loads OpenAPI documents
    - `openapi.go` - Extracts operations from OpenAPI paths
    - `jsonschema.go` - Parses JSON schemas into Go types
    - `operation.go` - Raw operation data structures
    - `extension.go` - x-codegen extension definitions

    ### Layer 2: Generation (internal/codegen/)

    Transforms parsed data into template-ready data:

    - `generators.go` - All generators and data preparation
    - `codegen.go` - Main orchestrator
    - `types.go` - Configuration types

    ### Layer 3: Templates (internal/templates/)

    Defines contracts and renders output:

    - `data.go` - Standardized template data structures
    - `templates.go` - Template loading and management
    - `funcs.go` - Template helper functions
    - `filewriter.go` - File writing utilities
    - `tmpl/` - Go template files
    ```

  - [ ] **Add Data Flow Diagram**

    ```markdown
    ## Data Flow

    1. OpenAPI spec is parsed by `parsers` package
    2. Raw schemas and operations extracted
    3. Generator transforms raw data to template data
    4. Template renders final code

    Example:
    api/openapi.yaml
    → Parser extracts Operation{Method: "GET", Path: "/users/{id}"}
    → Generator creates HandlerData{Operations: []OperationData{...}}
    → Template generates handler.gen.go with GetUser() method
    ```

  - [ ] **Template Data Reference**

    ````markdown
    ## Template Data Structures

    All template data structures are defined in `internal/templates/data.go`.

    ### HandlerData

    Used by: `echo_server.tmpl`

    ```go
    type HandlerData struct {
        Package    string
        Domain     string
        Operations []OperationData
    }
    ```
    ````

    ### OperationData

    ```go
    type OperationData struct {
        Name              string // OperationID
        Method            string // HTTP method
        Path              string // URL path
        HasRequestBody    bool
        RequestBodySchema string
        ResponseType      string
        PathParams        []ParamData
        QueryParams       []ParamData
    }
    ```

    ```

    ```

  - [ ] **Generator Development Guide**

    ````markdown
    ## Adding a New Generator

    1. Define data structure in `internal/templates/data.go`
    2. Add generator entry to `generators` map in `internal/codegen/generators.go`
    3. Implement prepare function in same file
    4. Create template in `internal/templates/tmpl/`
    5. Add configuration to `Config` struct

    Example:

    ```go
    // In generators.go
    "my_generator": {
        Name:         "my_generator",
        TemplateName: "my_template.tmpl",
        OutputFile:   "my_output.gen.go",
        PrepareFunc:  prepareMyData,
        Filter:       func(s *parsers.Schema) bool { return true },
        Enabled:      func(c *Config) bool { return c.Generators.MyGenerator != nil },
    }
    ```
    ````

    ```

    ```

  - [ ] **Migration Guide**
    - What changed in the refactor
    - How to update custom generators
    - Breaking changes (if any)

- [ ] **Create examples showing complete flow**
  - [ ] Example OpenAPI schema with x-codegen
  - [ ] Show what gets parsed
  - [ ] Show template data transformation
  - [ ] Show final generated code

## Success Criteria

1. **Clean Architecture**
   - [x] Each package has single responsibility
   - [x] No circular dependencies
   - [x] Clear data flow

2. **Code Reduction**
   - [x] ~85% less code in generators (7 files → 1 file)
   - [x] No duplicate generation logic
   - [x] Single pattern for all generators

3. **Type Safety**
   - [x] All template data strongly typed
   - [x] Compile-time checking of template inputs
   - [x] No runtime template errors

4. **Maintainability**
   - [x] Easy to add new generators (just add to map)
   - [x] Clear where to make changes
   - [x] Well-documented interfaces

5. **Backwards Compatibility**
   - [x] All existing generation still works
   - [x] Generated code remains the same
   - [x] No breaking changes to CLI

## Benefits of This Refactor

### Before (Current Issues)

- Template data preparation mixed with parsing logic
- No standardized template data structures
- 7+ generator files with duplicate logic
- Generators creating inline anonymous structs
- Templates expecting different field names for same concepts
- Hard to understand what data each template expects
- Runtime template errors

### After (Benefits)

- **Clean separation of concerns** - Each layer has one job
- **85% less generator code** - One file instead of seven
- **Type-safe template data** - Compiler catches mismatches
- **Single source of truth** - All generators in one place
- **Easier to test** - Each component isolated
- **Simpler to extend** - Just add to generators map
- **Better documentation** - Templates document expectations
- **No runtime surprises** - All checked at compile time

## File Structure After Refactor

```
internal/
├── parsers/              # Only parsing, no template knowledge
│   ├── parser.go         # Main orchestrator
│   ├── openapi.go        # Operation extraction (no PrepareXData methods)
│   ├── jsonschema.go     # Schema parsing
│   ├── operation.go      # Raw data structures (with Type field)
│   └── extension.go      # x-codegen definitions
├── codegen/              # Transforms parsed data to template data
│   ├── codegen.go        # Main orchestrator
│   ├── generators.go     # ALL generators + prepare functions
│   └── types.go          # Configuration types
└── templates/            # Template contracts and rendering
    ├── data.go           # Standardized data structures
    ├── templates.go      # Template management
    ├── funcs.go          # Helper functions
    ├── filewriter.go     # File output
    └── tmpl/             # Template files (with documentation headers)
        ├── types.tmpl
        ├── service.tmpl
        ├── repository.tmpl
        ├── repository_postgres.tmpl
        ├── repository_sqlite.tmpl
        └── echo_server.tmpl
```

## Commands for Testing

```bash
# Clean all generated files
make clean-generated

# Generate everything
make generate

# Test specific generator
make generate-codegen

# Verify compilation
go build ./...

# Run tests
go test ./...

# Check what was generated
find internal -name "*.gen.go" -type f | wc -l
```

## Implementation Order

1. **First**: Create template data structures (`templates/data.go`)
2. **Second**: Create consolidated generators (`codegen/generators.go`)
3. **Third**: Update templates to use new data structures
4. **Fourth**: Clean up parsers (remove template methods)
5. **Fifth**: Delete old generator files
6. **Sixth**: Test everything works
7. **Last**: Update documentation

This order ensures we can test at each step and roll back if needed.
