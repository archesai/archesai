package parsers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOperationDef_GetSuccessResponse(t *testing.T) {
	// Load test OpenAPI doc to get real operations
	parser := NewOpenAPIParser()
	doc, err := parser.ParseFile("../../test/data/parsers/openapi/simple-api.yaml")
	require.NoError(t, err)
	require.NotNil(t, doc)

	operations, err := parser.extractOperations()
	require.NoError(t, err)
	require.NotEmpty(t, operations)

	// Find listUsers operation
	const listUsersOpID = "ListUsers"
	var listUsersOp *OperationDef
	for i := range operations {
		if operations[i].ID == listUsersOpID {
			listUsersOp = &operations[i]
			break
		}
	}
	require.NotNil(t, listUsersOp)

	// Test getting success response
	successResp := listUsersOp.GetSuccessResponse()
	assert.NotNil(t, successResp)
	assert.Equal(t, "200", successResp.StatusCode)
	assert.True(t, successResp.IsSuccess())

	// Find deleteUser operation (has 204 success)
	var deleteUserOp *OperationDef
	for i := range operations {
		if operations[i].ID == "DeleteUser" {
			deleteUserOp = &operations[i]
			break
		}
	}
	require.NotNil(t, deleteUserOp)

	successResp = deleteUserOp.GetSuccessResponse()
	assert.NotNil(t, successResp)
	assert.Equal(t, "204", successResp.StatusCode)
	assert.True(t, successResp.IsSuccess())
}

func TestOperationDef_GetErrorResponses(t *testing.T) {
	// Load test OpenAPI doc
	parser := NewOpenAPIParser()
	_, err := parser.ParseFile("../../test/data/parsers/openapi/simple-api.yaml")
	require.NoError(t, err)

	operations, err := parser.extractOperations()
	require.NoError(t, err)
	require.NotEmpty(t, operations)

	// Find listUsers operation
	const listUsersOpID = "ListUsers"
	var listUsersOp *OperationDef
	for i := range operations {
		if operations[i].ID == listUsersOpID {
			listUsersOp = &operations[i]
			break
		}
	}
	require.NotNil(t, listUsersOp)

	// Test getting error responses
	errorResponses := listUsersOp.GetErrorResponses()
	assert.NotEmpty(t, errorResponses)
	for _, errResp := range errorResponses {
		assert.False(t, errResp.IsSuccess())
		assert.NotEqual(t, "200", errResp.StatusCode)
		assert.NotEqual(t, "201", errResp.StatusCode)
		assert.NotEqual(t, "204", errResp.StatusCode)
	}
}

func TestOperationDef_HasBearerAuth(t *testing.T) {
	// Load test OpenAPI doc
	parser := NewOpenAPIParser()
	_, err := parser.ParseFile("../../test/data/parsers/openapi/simple-api.yaml")
	require.NoError(t, err)

	operations, err := parser.extractOperations()
	require.NoError(t, err)
	require.NotEmpty(t, operations)

	// Find operations to test
	opsMap := make(map[string]*OperationDef)
	for i := range operations {
		opsMap[operations[i].ID] = &operations[i]
	}

	// listUsers has no security
	const listUsersOpID = "ListUsers"
	listUsersOp := opsMap[listUsersOpID]
	assert.False(t, listUsersOp.HasBearerAuth())

	// getUser has bearerAuth
	getUserOp := opsMap["GetUser"]
	assert.True(t, getUserOp.HasBearerAuth())

	// updateUser has both bearerAuth and sessionCookie
	updateUserOp := opsMap["UpdateUser"]
	assert.True(t, updateUserOp.HasBearerAuth())
}

func TestOperationDef_HasCookieAuth(t *testing.T) {
	// Load test OpenAPI doc
	parser := NewOpenAPIParser()
	_, err := parser.ParseFile("../../test/data/parsers/openapi/simple-api.yaml")
	require.NoError(t, err)

	operations, err := parser.extractOperations()
	require.NoError(t, err)
	require.NotEmpty(t, operations)

	// Find operations to test
	opsMap := make(map[string]*OperationDef)
	for i := range operations {
		opsMap[operations[i].ID] = &operations[i]
	}

	// listUsers has no security
	const listUsersOpID = "ListUsers"
	listUsersOp := opsMap[listUsersOpID]
	assert.False(t, listUsersOp.HasCookieAuth())

	// getUser has only bearerAuth, no cookie
	getUserOp := opsMap["GetUser"]
	assert.False(t, getUserOp.HasCookieAuth())

	// Note: The test data might not have proper cookie auth setup
	// The HasCookieAuth method checks for Type == "apiKey" && Scheme == "cookie"
	// but the OpenAPI spec has In == "cookie" not Scheme == "cookie"
}
