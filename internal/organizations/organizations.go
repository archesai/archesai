// Package organizations provides organization management functionality including
// organization CRUD operations, member management, and invitation handling.
package organizations

//go:generate go tool oapi-codegen --config=domain/models.cfg.yaml ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=adapters/http/server.cfg.yaml ../../api/openapi.bundled.yaml

// Domain constants
const (
	// DefaultPlan is the default organization plan for new organizations
	DefaultPlan = "free"

	// MaxMembersPerOrganization defines the maximum number of members per organization
	MaxMembersPerOrganization = 100

	// InvitationExpiryDays defines how long invitations remain valid
	InvitationExpiryDays = 7
)
