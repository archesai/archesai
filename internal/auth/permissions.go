package auth

// Role represents a user role in the system
type Role string

const (
	// RoleSuperAdmin has full system access
	RoleSuperAdmin Role = "super_admin"
	// RoleSystemAdmin can manage system settings
	RoleSystemAdmin Role = "system_admin"
	// RoleSupport provides support access
	RoleSupport Role = "support"

	// RoleOrgOwner has full control over an organization
	RoleOrgOwner Role = "org_owner"
	// RoleOrgAdmin can administer an organization
	RoleOrgAdmin Role = "org_admin"
	// RoleOrgMember is a regular organization member
	RoleOrgMember Role = "org_member"
	// RoleOrgViewer can only view organization resources
	RoleOrgViewer Role = "org_viewer"
	// RoleOrgGuest has limited guest access
	RoleOrgGuest Role = "org_guest"

	// RoleAPIUser represents API-only access
	RoleAPIUser Role = "api_user"
	// RoleServiceAccount is for service accounts
	RoleServiceAccount Role = "service_account"
)

// Permission represents a specific permission in the system
type Permission string

const (
	// PermSystemManage allows full system management
	PermSystemManage Permission = "system:manage"
	// PermSystemViewMetrics allows viewing system metrics
	PermSystemViewMetrics Permission = "system:view_metrics"
	// PermSystemViewLogs allows viewing system logs
	PermSystemViewLogs Permission = "system:view_logs"
	// PermSystemManageUsers allows managing system users
	PermSystemManageUsers Permission = "system:manage_users"

	// PermOrgCreate allows creating new organizations
	PermOrgCreate Permission = "org:create"
	// PermOrgRead allows reading organization data
	PermOrgRead Permission = "org:read"
	// PermOrgUpdate allows updating organization settings
	PermOrgUpdate Permission = "org:update"
	// PermOrgDelete allows deleting organizations
	PermOrgDelete Permission = "org:delete"
	// PermOrgManageMembers allows managing organization members
	PermOrgManageMembers Permission = "org:manage_members"
	// PermOrgManageBilling allows managing billing
	PermOrgManageBilling Permission = "org:manage_billing"
	// PermOrgManageSettings allows managing organization settings
	PermOrgManageSettings Permission = "org:manage_settings"
	// PermOrgInviteMembers allows inviting new members
	PermOrgInviteMembers Permission = "org:invite_members"
	// PermOrgViewAuditLog allows viewing audit logs
	PermOrgViewAuditLog Permission = "org:view_audit_log"

	// PermWorkflowCreate allows creating workflows
	PermWorkflowCreate Permission = "workflow:create"
	// PermWorkflowRead allows reading workflows
	PermWorkflowRead Permission = "workflow:read"
	// PermWorkflowUpdate allows updating workflows
	PermWorkflowUpdate Permission = "workflow:update"
	// PermWorkflowDelete allows deleting workflows
	PermWorkflowDelete Permission = "workflow:delete"
	// PermWorkflowExecute allows executing workflows
	PermWorkflowExecute Permission = "workflow:execute"
	// PermWorkflowViewRuns allows viewing workflow runs
	PermWorkflowViewRuns Permission = "workflow:view_runs"
	// PermWorkflowManageRuns allows managing workflow runs
	PermWorkflowManageRuns Permission = "workflow:manage_runs"

	// PermContentCreate allows creating content
	PermContentCreate Permission = "content:create"
	// PermContentRead allows reading content
	PermContentRead Permission = "content:read"
	// PermContentUpdate allows updating content
	PermContentUpdate Permission = "content:update"
	// PermContentDelete allows deleting content
	PermContentDelete Permission = "content:delete"
	// PermContentShare allows sharing content
	PermContentShare Permission = "content:share"
	// PermContentPublish allows publishing content
	PermContentPublish Permission = "content:publish"

	// PermAIUseBasic allows basic AI usage
	PermAIUseBasic Permission = "ai:use_basic"
	// PermAIUsePremium allows premium AI usage
	PermAIUsePremium Permission = "ai:use_premium"
	// PermAIManageKeys allows managing AI keys
	PermAIManageKeys Permission = "ai:manage_keys"
	// PermAIViewUsage allows viewing AI usage
	PermAIViewUsage Permission = "ai:view_usage"
	// PermAIManagePrompts allows managing AI prompts
	PermAIManagePrompts Permission = "ai:manage_prompts"

	// PermUserUpdateProfile allows updating user profile
	PermUserUpdateProfile Permission = "user:update_profile"
	// PermUserManageTokens allows managing user tokens
	PermUserManageTokens Permission = "user:manage_tokens"
	// PermUserManageSessions allows managing user sessions
	PermUserManageSessions Permission = "user:manage_sessions"
	// PermUserViewActivity allows viewing user activity
	PermUserViewActivity Permission = "user:view_activity"

	// PermAPICreateKeys allows creating API keys
	PermAPICreateKeys Permission = "api:create_keys"
	// PermAPIManageKeys allows managing API keys
	PermAPIManageKeys Permission = "api:manage_keys"
	// PermAPIViewUsage allows viewing API usage
	PermAPIViewUsage Permission = "api:view_usage"
)

// Scope represents API access scopes
type Scope string

const (
	// ScopeOpenID is the OpenID Connect scope
	ScopeOpenID Scope = "openid"
	// ScopeEmail allows reading email address
	ScopeEmail Scope = "email"
	// ScopeProfile allows reading basic profile info
	ScopeProfile Scope = "profile"

	// ScopeReadProfile allows reading user profile
	ScopeReadProfile Scope = "read:profile"
	// ScopeReadOrganizations allows reading organization data
	ScopeReadOrganizations Scope = "read:organizations"
	// ScopeReadWorkflows allows reading workflows
	ScopeReadWorkflows Scope = "read:workflows"
	// ScopeReadContent allows reading content
	ScopeReadContent Scope = "read:content"
	// ScopeReadMetrics allows reading metrics
	ScopeReadMetrics Scope = "read:metrics"

	// ScopeWriteProfile allows writing to user profile
	ScopeWriteProfile Scope = "write:profile"
	// ScopeWriteOrganizations allows writing to organizations
	ScopeWriteOrganizations Scope = "write:organizations"
	// ScopeWriteWorkflows allows writing to workflows
	ScopeWriteWorkflows Scope = "write:workflows"
	// ScopeWriteContent allows writing content
	ScopeWriteContent Scope = "write:content"

	// ScopeExecuteWorkflows allows executing workflows
	ScopeExecuteWorkflows Scope = "execute:workflows"
	// ScopeExecuteCommands allows executing commands
	ScopeExecuteCommands Scope = "execute:commands"
	// ScopeExecuteAI allows executing AI operations
	ScopeExecuteAI Scope = "execute:ai"
)

// RolePermissions maps roles to their associated permissions
var RolePermissions = map[Role][]Permission{
	RoleSuperAdmin: {
		PermSystemManage,
		PermSystemViewMetrics,
		PermSystemViewLogs,
		PermSystemManageUsers,
		PermOrgCreate,
		PermOrgRead,
		PermOrgUpdate,
		PermOrgDelete,
		PermOrgManageMembers,
		PermOrgManageBilling,
		PermOrgManageSettings,
		PermOrgInviteMembers,
		PermOrgViewAuditLog,
		PermWorkflowCreate,
		PermWorkflowRead,
		PermWorkflowUpdate,
		PermWorkflowDelete,
		PermWorkflowExecute,
		PermContentCreate,
		PermContentRead,
		PermContentUpdate,
		PermContentDelete,
		PermContentShare,
		PermContentPublish,
		PermAIUsePremium,
		PermAIManageKeys,
		PermAPICreateKeys,
		PermAPIManageKeys,
	},
	RoleOrgOwner: {
		PermOrgRead,
		PermOrgUpdate,
		PermOrgDelete,
		PermOrgManageMembers,
		PermOrgManageBilling,
		PermOrgManageSettings,
		PermOrgInviteMembers,
		PermOrgViewAuditLog,
		PermWorkflowCreate,
		PermWorkflowRead,
		PermWorkflowUpdate,
		PermWorkflowDelete,
		PermWorkflowExecute,
		PermContentCreate,
		PermContentRead,
		PermContentUpdate,
		PermContentDelete,
		PermContentShare,
		PermContentPublish,
		PermAIUsePremium,
		PermAPICreateKeys,
	},
	RoleOrgAdmin: {
		PermOrgRead,
		PermOrgUpdate,
		PermOrgManageMembers,
		PermOrgInviteMembers,
		PermWorkflowCreate,
		PermWorkflowRead,
		PermWorkflowUpdate,
		PermWorkflowExecute,
		PermContentCreate,
		PermContentRead,
		PermContentUpdate,
		PermContentShare,
		PermAIUseBasic,
	},
	RoleOrgMember: {
		PermOrgRead,
		PermWorkflowCreate,
		PermWorkflowRead,
		PermWorkflowExecute,
		PermContentCreate,
		PermContentRead,
		PermContentUpdate,
		PermAIUseBasic,
	},
	RoleOrgViewer: {
		PermOrgRead,
		PermWorkflowRead,
		PermContentRead,
	},
	RoleOrgGuest: {
		PermOrgRead,
		PermContentRead,
	},
}

// HasPermission checks if a role has a specific permission
func HasPermission(role Role, permission Permission) bool {
	permissions, exists := RolePermissions[role]
	if !exists {
		return false
	}

	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// GetRolePermissions returns all permissions for a given role
func GetRolePermissions(role Role) []Permission {
	return RolePermissions[role]
}

// GetDefaultScopes returns the default OAuth scopes
func GetDefaultScopes() []Scope {
	return []Scope{
		ScopeOpenID,
		ScopeEmail,
		ScopeProfile,
		ScopeReadProfile,
		ScopeReadOrganizations,
	}
}

// GetAPIKeyScopes returns scopes available for API keys
func GetAPIKeyScopes() []Scope {
	return []Scope{
		ScopeReadWorkflows,
		ScopeWriteWorkflows,
		ScopeExecuteWorkflows,
		ScopeReadContent,
		ScopeWriteContent,
		ScopeExecuteAI,
	}
}

// ValidateScopes checks if permissions satisfy required scopes
func ValidateScopes(_ []Permission, requiredScopes []Scope) bool {
	// TODO: Implement proper scope validation logic
	// For now, return true if no scopes required
	return len(requiredScopes) == 0
}
