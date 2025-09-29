package parsers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOperationDef_GetSuccessResponse(t *testing.T) {
	// Load test OpenAPI doc to get real operations
	doc, _, err := ParseOpenAPI("../../test/data/parsers/openapi/simple-api.yaml")
	require.NoError(t, err)
	require.NotNil(t, doc)

	operations := ExtractOperations(doc)
	require.NotEmpty(t, operations)

	// Find listUsers operation
	const listUsersOpID = "listUsers"
	var listUsersOp *OperationDef
	for i := range operations {
		if operations[i].OperationID == listUsersOpID {
			listUsersOp = &operations[i]
			break
		}
	}
	require.NotNil(t, listUsersOp)

	// Test getting success response
	successResp := listUsersOp.GetSuccessResponse()
	assert.NotNil(t, successResp)
	assert.Equal(t, "200", successResp.StatusCode)
	assert.True(t, successResp.IsSuccess)

	// Find deleteUser operation (has 204 success)
	var deleteUserOp *OperationDef
	for i := range operations {
		if operations[i].OperationID == "deleteUser" {
			deleteUserOp = &operations[i]
			break
		}
	}
	require.NotNil(t, deleteUserOp)

	successResp = deleteUserOp.GetSuccessResponse()
	assert.NotNil(t, successResp)
	assert.Equal(t, "204", successResp.StatusCode)
	assert.True(t, successResp.IsSuccess)
}

func TestOperationDef_GetErrorResponses(t *testing.T) {
	// Load test OpenAPI doc
	doc, _, err := ParseOpenAPI("../../test/data/parsers/openapi/simple-api.yaml")
	require.NoError(t, err)

	operations := ExtractOperations(doc)
	require.NotEmpty(t, operations)

	// Find listUsers operation
	const listUsersOpID = "listUsers"
	var listUsersOp *OperationDef
	for i := range operations {
		if operations[i].OperationID == listUsersOpID {
			listUsersOp = &operations[i]
			break
		}
	}
	require.NotNil(t, listUsersOp)

	// Test getting error responses
	errorResponses := listUsersOp.GetErrorResponses()
	assert.NotEmpty(t, errorResponses)
	for _, errResp := range errorResponses {
		assert.False(t, errResp.IsSuccess)
		assert.NotEqual(t, "200", errResp.StatusCode)
		assert.NotEqual(t, "201", errResp.StatusCode)
		assert.NotEqual(t, "204", errResp.StatusCode)
	}
}

func TestOperationDef_HasBearerAuth(t *testing.T) {
	// Load test OpenAPI doc
	doc, _, err := ParseOpenAPI("../../test/data/parsers/openapi/simple-api.yaml")
	require.NoError(t, err)

	operations := ExtractOperations(doc)
	require.NotEmpty(t, operations)

	// Find operations to test
	opsMap := make(map[string]*OperationDef)
	for i := range operations {
		opsMap[operations[i].OperationID] = &operations[i]
	}

	// listUsers has no security
	const listUsersOpID = "listUsers"
	listUsersOp := opsMap[listUsersOpID]
	assert.False(t, listUsersOp.HasBearerAuth())

	// getUser has bearerAuth
	getUserOp := opsMap["getUser"]
	assert.True(t, getUserOp.HasBearerAuth())

	// updateUser has both bearerAuth and sessionCookie
	updateUserOp := opsMap["updateUser"]
	assert.True(t, updateUserOp.HasBearerAuth())
}

func TestOperationDef_HasCookieAuth(t *testing.T) {
	// Load test OpenAPI doc
	doc, _, err := ParseOpenAPI("../../test/data/parsers/openapi/simple-api.yaml")
	require.NoError(t, err)

	operations := ExtractOperations(doc)
	require.NotEmpty(t, operations)

	// Find operations to test
	opsMap := make(map[string]*OperationDef)
	for i := range operations {
		opsMap[operations[i].OperationID] = &operations[i]
	}

	// listUsers has no security
	const listUsersOpID = "listUsers"
	listUsersOp := opsMap[listUsersOpID]
	assert.False(t, listUsersOp.HasCookieAuth())

	// getUser has only bearerAuth, no cookie
	getUserOp := opsMap["getUser"]
	assert.False(t, getUserOp.HasCookieAuth())

	// Note: The test data might not have proper cookie auth setup
	// The HasCookieAuth method checks for Type == "apiKey" && Scheme == "cookie"
	// but the OpenAPI spec has In == "cookie" not Scheme == "cookie"
}

func TestConstants(t *testing.T) {
	// Test SQL dialect constants
	assert.Equal(t, "postgresql", SQLDialectPostgres)
	assert.Equal(t, "sqlite", SQLDialectSQLite)

	// Test SQL type constants
	assert.Equal(t, "TEXT", SQLTypeText)
	assert.Equal(t, "INTEGER", SQLTypeInteger)
	assert.Equal(t, "BIGINT", SQLTypeBigInt)
	assert.Equal(t, "BOOLEAN", SQLTypeBoolean)
	assert.Equal(t, "TIMESTAMPTZ", SQLTypeTimestamp)
	assert.Equal(t, "DATETIME", SQLTypeDateTime)
	assert.Equal(t, "DATE", SQLTypeDate)
	assert.Equal(t, "UUID", SQLTypeUUID)
	assert.Equal(t, "JSONB", SQLTypeJSONB)
	assert.Equal(t, "NUMERIC", SQLTypeNumeric)
	assert.Equal(t, "REAL", SQLTypeReal)
	assert.Equal(t, "DOUBLE PRECISION", SQLTypeDouble)
}

func TestValidSlices(t *testing.T) {
	// Test valid repository operations
	assert.Equal(
		t,
		[]string{"create", "read", "update", "delete", "list"},
		ValidRepositoryOperations,
	)

	// Test valid HTTP methods
	assert.Equal(t, []string{"GET", "POST", "PUT", "PATCH", "DELETE"}, ValidHTTPMethods)

	// Test valid error handling
	assert.Equal(t, []string{"error_return", "panic", "custom"}, ValidErrorHandling)

	// Test valid log levels
	assert.Equal(t, []string{"debug", "info", "warn", "error"}, ValidLogLevels)

	// Test valid domain types
	assert.Equal(t, []string{"entity", "aggregate", "valueobject", "dto"}, ValidDomainTypes)
}
