package spec

import (
	"os"
	"testing"
	"testing/fstest"

	"github.com/archesai/archesai/internal/schema"
)

// newTestParser creates a Parser from test files for use in tests.
func newTestParser(t *testing.T, files map[string]string) *Parser {
	t.Helper()
	fsys := fstest.MapFS{}
	for path, content := range files {
		fsys[path] = &fstest.MapFile{Data: []byte(content)}
	}
	doc, err := NewOpenAPIDocumentFromFS(fsys, "openapi.yaml")
	if err != nil {
		t.Fatalf("NewOpenAPIDocumentFromFS() error = %v", err)
	}
	return NewParser(doc)
}

// newTestParserWithIncludes creates a Parser from test files with includes.
func newTestParserWithIncludes(
	t *testing.T,
	files map[string]string,
	includeNames []string,
) *Parser {
	t.Helper()
	fsys := fstest.MapFS{}
	for path, content := range files {
		fsys[path] = &fstest.MapFile{Data: []byte(content)}
	}
	compositeFS := BuildIncludeFS(fsys, includeNames)
	doc, err := NewOpenAPIDocumentFromFS(compositeFS, "openapi.yaml")
	if err != nil {
		t.Fatalf("NewOpenAPIDocumentFromFS() error = %v", err)
	}
	return NewParser(doc).WithIncludes(includeNames)
}

func TestParser_Parse_BasicSpec(t *testing.T) {
	files := map[string]string{
		"openapi.yaml": `
openapi: 3.1.0
x-project-name: github.com/example/myapi
info:
  title: My API
  description: A test API
  version: v1.0.0
tags:
  - name: User
    description: User operations
components:
  schemas:
    User:
      $ref: components/schemas/User.yaml
paths:
  /users:
    $ref: paths/users.yaml
`,
		"components/schemas/User.yaml": `
title: User
type: object
x-codegen-schema-type: entity
properties:
  id:
    type: string
    format: uuid
  name:
    type: string
  email:
    type: string
    format: email
required:
  - id
  - name
  - email
`,
		"paths/users.yaml": `
x-path: /users
get:
  operationId: ListUsers
  summary: List users
  tags:
    - User
  responses:
    '200':
      description: Success
      content:
        application/json:
          schema:
            type: object
            properties:
              data:
                type: array
                items:
                  $ref: ../components/schemas/User.yaml
`,
	}

	p := newTestParser(t, files)
	s, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Verify basic metadata
	if s.ProjectName != "github.com/example/myapi" {
		t.Errorf("ProjectName = %q, want %q", s.ProjectName, "github.com/example/myapi")
	}
	if s.Title != "My API" {
		t.Errorf("Title = %q, want %q", s.Title, "My API")
	}
	if s.Version != "v1.0.0" {
		t.Errorf("Version = %q, want %q", s.Version, "v1.0.0")
	}

	// Verify tags
	if len(s.Tags) != 1 {
		t.Errorf("len(Tags) = %d, want 1", len(s.Tags))
	} else if s.Tags[0].Name != "User" {
		t.Errorf("Tags[0].Name = %q, want %q", s.Tags[0].Name, "User")
	}

	// Verify schemas
	found := false
	for _, sch := range s.Schemas {
		if sch.Title == "User" {
			found = true
			if sch.XCodegenSchemaType != schema.TypeEntity {
				t.Errorf(
					"User.XCodegenSchemaType = %q, want %q",
					sch.XCodegenSchemaType,
					schema.TypeEntity,
				)
			}
			// Entity schemas should have base fields added
			if _, ok := sch.Properties["ID"]; !ok {
				t.Error("User schema missing ID property")
			}
			if _, ok := sch.Properties["CreatedAt"]; !ok {
				t.Error("User schema missing CreatedAt property")
			}
			if _, ok := sch.Properties["UpdatedAt"]; !ok {
				t.Error("User schema missing UpdatedAt property")
			}
			break
		}
	}
	if !found {
		t.Error("User schema not found in parsed schemas")
	}

	// Verify operations
	if len(s.Operations) == 0 {
		t.Error("No operations parsed")
	} else {
		op := s.Operations[0]
		if op.ID != "ListUsers" {
			t.Errorf("Operation.ID = %q, want %q", op.ID, "ListUsers")
		}
		if op.Path != "/users" {
			t.Errorf("Operation.Path = %q, want %q", op.Path, "/users")
		}
		if op.Method != "GET" {
			t.Errorf("Operation.Method = %q, want %q", op.Method, "GET")
		}
	}
}

func TestParser_Parse_AutoDiscoveryPaths(t *testing.T) {
	// Test that paths are auto-discovered from paths/ directory
	// even without explicit references in openapi.yaml
	files := map[string]string{
		"openapi.yaml": `
openapi: 3.1.0
x-project-name: github.com/example/myapi
info:
  title: My API
  version: v1.0.0
`,
		"paths/users.yaml": `
x-path: /users
get:
  operationId: ListUsers
  summary: List users
  tags:
    - User
  responses:
    '200':
      description: Success
`,
		"paths/users_id.yaml": `
x-path: /users/{id}
get:
  operationId: GetUser
  summary: Get a user
  tags:
    - User
  responses:
    '200':
      description: Success
`,
		"paths/health.yaml": `
x-path: /health
get:
  operationId: GetHealth
  summary: Health check
  tags:
    - Health
  responses:
    '200':
      description: OK
`,
	}

	p := newTestParser(t, files)
	s, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Should have 3 operations from auto-discovered paths
	if len(s.Operations) != 3 {
		t.Errorf("len(Operations) = %d, want 3", len(s.Operations))
		for _, op := range s.Operations {
			t.Logf("  - %s %s (%s)", op.Method, op.Path, op.ID)
		}
	}

	// Verify all operations are present
	opIDs := make(map[string]bool)
	for _, op := range s.Operations {
		opIDs[op.ID] = true
	}

	expectedOps := []string{"ListUsers", "GetUser", "GetHealth"}
	for _, expected := range expectedOps {
		if !opIDs[expected] {
			t.Errorf("expected operation %q not found", expected)
		}
	}
}

func TestParser_Parse_AutoDiscoverySchemas(t *testing.T) {
	// Test that schemas are auto-discovered from components/schemas/
	files := map[string]string{
		"openapi.yaml": `
openapi: 3.1.0
x-project-name: github.com/example/myapi
info:
  title: My API
  version: v1.0.0
`,
		"components/schemas/User.yaml": `
title: User
type: object
properties:
  id:
    type: string
    format: uuid
  name:
    type: string
required:
  - id
  - name
`,
		"components/schemas/Organization.yaml": `
title: Organization
type: object
properties:
  id:
    type: string
    format: uuid
  name:
    type: string
required:
  - id
  - name
`,
	}

	p := newTestParser(t, files)
	s, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Should have auto-discovered schemas
	schemaNames := make(map[string]bool)
	for _, sch := range s.Schemas {
		schemaNames[sch.Title] = true
	}

	if !schemaNames["User"] {
		t.Error("User schema not auto-discovered")
	}
	if !schemaNames["Organization"] {
		t.Error("Organization schema not auto-discovered")
	}
}

func TestParser_Parse_InlineSchema(t *testing.T) {
	files := map[string]string{
		"openapi.yaml": `
openapi: 3.1.0
x-project-name: github.com/example/myapi
info:
  title: My API
  version: v1.0.0
components:
  schemas:
    User:
      title: User
      type: object
      properties:
        id:
          type: string
          format: uuid
        name:
          type: string
      required:
        - id
        - name
`,
	}

	p := newTestParser(t, files)
	s, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	found := false
	for _, sch := range s.Schemas {
		if sch.Title == "User" {
			found = true
			if _, ok := sch.Properties["ID"]; !ok {
				t.Error("User schema missing ID property")
			}
			if _, ok := sch.Properties["Name"]; !ok {
				t.Error("User schema missing Name property")
			}
			break
		}
	}
	if !found {
		t.Error("User schema not found")
	}
}

func TestParser_Parse_ServerInclude(t *testing.T) {
	files := map[string]string{
		"openapi.yaml": `
openapi: 3.1.0
x-project-name: github.com/example/myapi
info:
  title: My API
  version: v1.0.0
`,
	}

	p := newTestParserWithIncludes(t, files, []string{"server"})
	s, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Server include should be in enabled includes
	found := false
	for _, inc := range s.EnabledIncludes {
		if inc == "server" {
			found = true
			break
		}
	}
	if !found {
		t.Error("server include not found in EnabledIncludes")
	}

	// Server include schemas should be loaded
	schemaNames := make(map[string]bool)
	for _, sch := range s.Schemas {
		schemaNames[sch.Title] = true
	}

	// These are standard server include schemas
	serverSchemas := []string{
		"Base",
		"Page",
		"PaginationMeta",
		"FilterNode",
		"Problem",
		"UUID",
		"Health",
	}
	for _, name := range serverSchemas {
		if !schemaNames[name] {
			t.Errorf("server include schema %q not found", name)
		}
	}

	// Server include should add health operation
	opIDs := make(map[string]bool)
	for _, op := range s.Operations {
		opIDs[op.ID] = true
	}
	if !opIDs["GetHealth"] {
		t.Error("GetHealth operation from server include not found")
	}
}

func TestParser_Parse_SecuritySchemes(t *testing.T) {
	files := map[string]string{
		"openapi.yaml": `
openapi: 3.1.0
x-project-name: github.com/example/myapi
info:
  title: My API
  version: v1.0.0
security:
  - bearerAuth: []
components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      description: Bearer token authentication
    sessionCookie:
      type: apiKey
      name: session_token
      in: cookie
      description: Session cookie
`,
		"paths/users.yaml": `
x-path: /users
get:
  operationId: ListUsers
  summary: List users
  tags:
    - User
  responses:
    '200':
      description: Success
`,
	}

	p := newTestParser(t, files)
	s, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Check security schemes
	if len(s.Security) != 2 {
		t.Errorf("len(Security) = %d, want 2", len(s.Security))
	}

	if scheme, ok := s.Security["bearerAuth"]; !ok {
		t.Error("bearerAuth security scheme not found")
	} else {
		if scheme.Type != "http" {
			t.Errorf("bearerAuth.Type = %q, want %q", scheme.Type, "http")
		}
		if scheme.Scheme != "bearer" {
			t.Errorf("bearerAuth.Scheme = %q, want %q", scheme.Scheme, "bearer")
		}
	}

	// Check that operation inherits security
	if len(s.Operations) > 0 {
		op := s.Operations[0]
		if len(op.Security) == 0 {
			t.Error("Operation should inherit security from root")
		}
	}
}

func TestParser_Parse_RequestBody(t *testing.T) {
	files := map[string]string{
		"openapi.yaml": `
openapi: 3.1.0
x-project-name: github.com/example/myapi
info:
  title: My API
  version: v1.0.0
`,
		"paths/users.yaml": `
x-path: /users
post:
  operationId: CreateUser
  summary: Create a user
  tags:
    - User
  requestBody:
    description: User data
    required: true
    content:
      application/json:
        schema:
          type: object
          properties:
            name:
              type: string
              description: The user's name
            email:
              type: string
              format: email
              description: The user's email
          required:
            - name
            - email
  responses:
    '201':
      description: Created
`,
	}

	p := newTestParser(t, files)
	s, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if len(s.Operations) == 0 {
		t.Fatal("No operations parsed")
	}

	op := s.Operations[0]
	if op.RequestBody == nil {
		t.Fatal("RequestBody is nil")
	}

	if op.RequestBody.Schema == nil {
		t.Fatal("RequestBody.Schema is nil")
	}

	schema := op.RequestBody.Schema
	if _, ok := schema.Properties["Name"]; !ok {
		t.Error("RequestBody schema missing Name property")
	}
	if _, ok := schema.Properties["Email"]; !ok {
		t.Error("RequestBody schema missing Email property")
	}

	// Check required fields
	if len(schema.Required) != 2 {
		t.Errorf("len(Required) = %d, want 2", len(schema.Required))
	}
}

func TestParser_Parse_PathParameters(t *testing.T) {
	files := map[string]string{
		"openapi.yaml": `
openapi: 3.1.0
x-project-name: github.com/example/myapi
info:
  title: My API
  version: v1.0.0
`,
		"paths/users_id.yaml": `
x-path: /users/{id}
get:
  operationId: GetUser
  summary: Get a user
  tags:
    - User
  responses:
    '200':
      description: Success
`,
		"paths/orgs_orgId_members_id.yaml": `
x-path: /organizations/{organizationId}/members/{id}
get:
  operationId: GetMember
  summary: Get a member
  tags:
    - Member
  responses:
    '200':
      description: Success
`,
	}

	p := newTestParser(t, files)
	s, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Find GetUser operation
	var getUserOp *Operation
	var getMemberOp *Operation
	for i := range s.Operations {
		if s.Operations[i].ID == "GetUser" {
			getUserOp = &s.Operations[i]
		}
		if s.Operations[i].ID == "GetMember" {
			getMemberOp = &s.Operations[i]
		}
	}

	if getUserOp == nil {
		t.Fatal("GetUser operation not found")
	}

	// Should have auto-extracted path parameter
	if len(getUserOp.Parameters) == 0 {
		t.Error("GetUser should have path parameters")
	} else {
		found := false
		for _, p := range getUserOp.Parameters {
			if p.In == "path" && p.Schema != nil && p.JSONTag == "id" {
				found = true
				if p.Format != "uuid" {
					t.Errorf("id parameter format = %q, want %q", p.Format, "uuid")
				}
				break
			}
		}
		if !found {
			t.Error("id path parameter not found")
		}
	}

	if getMemberOp == nil {
		t.Fatal("GetMember operation not found")
	}

	// Should have both organizationId and id parameters
	paramNames := make(map[string]bool)
	for _, p := range getMemberOp.Parameters {
		if p.In == "path" && p.Schema != nil {
			paramNames[p.JSONTag] = true
		}
	}
	if !paramNames["organizationId"] {
		t.Error("organizationId parameter not found")
	}
	if !paramNames["id"] {
		t.Error("id parameter not found")
	}
}

func TestParser_Parse_AllOfComposition(t *testing.T) {
	files := map[string]string{
		"openapi.yaml": `
openapi: 3.1.0
x-project-name: github.com/example/myapi
info:
  title: My API
  version: v1.0.0
`,
		"components/schemas/User.yaml": `
title: User
x-codegen-schema-type: entity
allOf:
  - $ref: Base.yaml
  - type: object
    properties:
      name:
        type: string
        description: User's name
      email:
        type: string
        format: email
    required:
      - name
      - email
`,
		"components/schemas/Base.yaml": `
title: Base
type: object
properties:
  id:
    type: string
    format: uuid
  createdAt:
    type: string
    format: date-time
  updatedAt:
    type: string
    format: date-time
required:
  - id
  - createdAt
  - updatedAt
`,
	}

	p := newTestParser(t, files)
	s, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Find User schema
	var userSchema *schema.Schema
	for _, sch := range s.Schemas {
		if sch.Title == "User" {
			userSchema = sch
			break
		}
	}

	if userSchema == nil {
		t.Fatal("User schema not found")
	}

	// Should have merged properties from allOf
	if _, ok := userSchema.Properties["Name"]; !ok {
		t.Error("User schema missing Name property from allOf")
	}
	if _, ok := userSchema.Properties["Email"]; !ok {
		t.Error("User schema missing Email property from allOf")
	}

	// Entity schemas get base fields added automatically
	if _, ok := userSchema.Properties["ID"]; !ok {
		t.Error("User schema missing ID property (base field)")
	}
}

func TestParser_Parse_NullableTypes(t *testing.T) {
	files := map[string]string{
		"openapi.yaml": `
openapi: 3.1.0
x-project-name: github.com/example/myapi
info:
  title: My API
  version: v1.0.0
components:
  schemas:
    User:
      title: User
      type: object
      properties:
        id:
          type: string
          format: uuid
        nickname:
          type:
            - string
            - 'null'
          description: Optional nickname
        age:
          type:
            - integer
            - 'null'
          description: Optional age
      required:
        - id
`,
	}

	p := newTestParser(t, files)
	s, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	var userSchema *schema.Schema
	for _, sch := range s.Schemas {
		if sch.Title == "User" {
			userSchema = sch
			break
		}
	}

	if userSchema == nil {
		t.Fatal("User schema not found")
	}

	// Check nullable property
	if nickname, ok := userSchema.Properties["Nickname"]; !ok {
		t.Error("Nickname property not found")
	} else if !nickname.GetOrNil().Nullable {
		t.Error("Nickname should be nullable")
	}

	if age, ok := userSchema.Properties["Age"]; !ok {
		t.Error("Age property not found")
	} else if !age.GetOrNil().Nullable {
		t.Error("Age should be nullable")
	}
}

func TestParser_Parse_XCodegenExtensions(t *testing.T) {
	files := map[string]string{
		"openapi.yaml": `
openapi: 3.1.0
x-project-name: github.com/example/myapi
info:
  title: My API
  version: v1.0.0
`,
	}

	p := newTestParser(t, files).
		WithCodegenOnly([]string{"models", "handlers"}).
		WithCodegenLint(true)
	s, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if len(s.CodegenOnly) != 2 {
		t.Errorf("len(CodegenOnly) = %d, want 2", len(s.CodegenOnly))
	}
	if !s.CodegenLint {
		t.Error("CodegenLint should be true")
	}
}

func TestParser_Parse_CustomHandler(t *testing.T) {
	files := map[string]string{
		"openapi.yaml": `
openapi: 3.1.0
x-project-name: github.com/example/myapi
info:
  title: My API
  version: v1.0.0
`,
		"paths/auth_login.yaml": `
x-path: /auth/login
post:
  operationId: Login
  summary: Login
  tags:
    - Auth
  x-codegen-custom-handler: true
  requestBody:
    content:
      application/json:
        schema:
          type: object
          properties:
            email:
              type: string
            password:
              type: string
  responses:
    '200':
      description: Success
`,
	}

	p := newTestParser(t, files)
	s, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	var loginOp *Operation
	for i := range s.Operations {
		if s.Operations[i].ID == "Login" {
			loginOp = &s.Operations[i]
			break
		}
	}

	if loginOp == nil {
		t.Fatal("Login operation not found")
	}

	if !loginOp.CustomHandler {
		t.Error("Login operation should have CustomHandler = true")
	}
}

func TestParser_Parse_DefaultValue(t *testing.T) {
	specPath := "../../pkg/auth/spec/openapi.yaml"
	baseFS := os.DirFS("../../pkg/auth/spec")
	doc, err := NewOpenAPIDocumentFromFS(baseFS, "openapi.yaml")
	if err != nil {
		t.Fatalf("NewOpenAPIDocumentFromFS(%s) error = %v", specPath, err)
	}
	parser := NewParser(doc)
	s, err := parser.Parse()
	if err != nil {
		t.Fatal(err)
	}

	// Find APIKey schema and verify scopes default
	for _, sch := range s.Schemas {
		if sch.Title == "APIKey" {
			for propName, propRef := range sch.Properties {
				if propName == "Scopes" {
					prop := propRef.GetOrNil()
					if prop.Default == nil {
						t.Error("Scopes.Default is nil, expected empty array")
					}

					// Check that Default is []any (same as []any)
					defaultSlice, ok := prop.Default.([]any)
					if !ok {
						t.Errorf("Scopes.Default type = %T, want []any", prop.Default)
					}
					if len(defaultSlice) != 0 {
						t.Errorf("Scopes.Default length = %d, want 0", len(defaultSlice))
					}
					return
				}
			}
			t.Error("Scopes property not found in APIKey schema")
			return
		}
	}
	t.Error("APIKey schema not found")
}

func TestParser_Parse_ResponseContentType(t *testing.T) {
	// Test that responses with application/problem+json are correctly parsed
	files := map[string]string{
		"openapi.yaml": `
openapi: 3.1.0
x-project-name: github.com/example/myapi
info:
  title: My API
  version: v1.0.0
`,
		"components/schemas/Problem.yaml": `
title: Problem
type: object
description: RFC 7807 Problem Details
properties:
  type:
    type: string
  title:
    type: string
  status:
    type: integer
  detail:
    type: string
  instance:
    type: string
`,
		"components/responses/BadRequest.yaml": `
description: Bad request
content:
  application/problem+json:
    schema:
      $ref: ../schemas/Problem.yaml
`,
		"components/responses/NotFound.yaml": `
description: Not found
content:
  application/problem+json:
    schema:
      $ref: ../schemas/Problem.yaml
`,
		"components/responses/UserResponse.yaml": `
description: User response
content:
  application/json:
    schema:
      type: object
      properties:
        data:
          $ref: ../schemas/User.yaml
`,
		"components/schemas/User.yaml": `
title: User
type: object
properties:
  id:
    type: string
    format: uuid
  name:
    type: string
required:
  - id
  - name
`,
		"paths/users_id.yaml": `
x-path: /users/{id}
get:
  operationId: GetUser
  summary: Get a user
  tags:
    - User
  responses:
    '200':
      $ref: ../components/responses/UserResponse.yaml
    '400':
      $ref: ../components/responses/BadRequest.yaml
    '404':
      $ref: ../components/responses/NotFound.yaml
`,
	}

	p := newTestParser(t, files)
	s, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Find GetUser operation
	var getUserOp *Operation
	for i := range s.Operations {
		if s.Operations[i].ID == "GetUser" {
			getUserOp = &s.Operations[i]
			break
		}
	}

	if getUserOp == nil {
		t.Fatal("GetUser operation not found")
	}

	// Verify responses have correct content types
	responseTests := []struct {
		statusCode  string
		contentType string
	}{
		{"200", "application/json"},
		{"400", "application/problem+json"},
		{"404", "application/problem+json"},
	}

	for _, tt := range responseTests {
		found := false
		for _, resp := range getUserOp.Responses {
			if resp.StatusCode == tt.statusCode {
				found = true
				if resp.ContentType != tt.contentType {
					t.Errorf(
						"Response %s ContentType = %q, want %q",
						tt.statusCode,
						resp.ContentType,
						tt.contentType,
					)
				}
				break
			}
		}
		if !found {
			t.Errorf("Response with status code %s not found", tt.statusCode)
		}
	}
}

func TestParser_Parse_InlineResponseContentType(t *testing.T) {
	// Test that inline responses with application/problem+json are correctly parsed
	files := map[string]string{
		"openapi.yaml": `
openapi: 3.1.0
x-project-name: github.com/example/myapi
info:
  title: My API
  version: v1.0.0
`,
		"paths/users.yaml": `
x-path: /users
post:
  operationId: CreateUser
  summary: Create a user
  tags:
    - User
  requestBody:
    content:
      application/json:
        schema:
          type: object
          properties:
            name:
              type: string
  responses:
    '201':
      description: Created
      content:
        application/json:
          schema:
            type: object
            properties:
              data:
                type: object
    '422':
      description: Validation error
      content:
        application/problem+json:
          schema:
            type: object
            properties:
              type:
                type: string
              title:
                type: string
              status:
                type: integer
              detail:
                type: string
`,
	}

	p := newTestParser(t, files)
	s, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Find CreateUser operation
	var createUserOp *Operation
	for i := range s.Operations {
		if s.Operations[i].ID == "CreateUser" {
			createUserOp = &s.Operations[i]
			break
		}
	}

	if createUserOp == nil {
		t.Fatal("CreateUser operation not found")
	}

	// Verify responses have correct content types
	responseTests := []struct {
		statusCode  string
		contentType string
	}{
		{"201", "application/json"},
		{"422", "application/problem+json"},
	}

	for _, tt := range responseTests {
		found := false
		for _, resp := range createUserOp.Responses {
			if resp.StatusCode == tt.statusCode {
				found = true
				if resp.ContentType != tt.contentType {
					t.Errorf(
						"Response %s ContentType = %q, want %q",
						tt.statusCode,
						resp.ContentType,
						tt.contentType,
					)
				}
				break
			}
		}
		if !found {
			t.Errorf("Response with status code %s not found", tt.statusCode)
		}
	}
}
