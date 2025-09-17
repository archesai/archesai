package sessions

import (
	"testing"
)

func TestGetRolePermissions(t *testing.T) {
	tests := []struct {
		name     string
		role     Role
		wantSome []Permission // Check that these permissions are included
	}{
		{
			name: "SuperAdmin has all permissions",
			role: RoleSuperAdmin,
			wantSome: []Permission{
				PermSystemManage,
				PermOrgDelete,
				PermWorkflowDelete,
				PermAIUsePremium,
			},
		},
		{
			name: "OrgOwner has organization management",
			role: RoleOrgOwner,
			wantSome: []Permission{
				PermOrgManageBilling,
				PermOrgDelete,
				PermWorkflowCreate,
				PermContentPublish,
			},
		},
		{
			name: "OrgMember has basic permissions",
			role: RoleOrgMember,
			wantSome: []Permission{
				PermOrgRead,
				PermWorkflowCreate,
				PermContentCreate,
				PermAIUseBasic,
			},
		},
		{
			name: "OrgViewer has read-only permissions",
			role: RoleOrgViewer,
			wantSome: []Permission{
				PermOrgRead,
				PermWorkflowRead,
				PermContentRead,
			},
		},
		{
			name: "OrgGuest has minimal permissions",
			role: RoleOrgGuest,
			wantSome: []Permission{
				PermOrgRead,
				PermContentRead,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetRolePermissions(tt.role)

			// Convert to map for easier checking
			gotMap := make(map[Permission]bool)
			for _, p := range got {
				gotMap[p] = true
			}

			// Check that expected permissions are present
			for _, want := range tt.wantSome {
				if !gotMap[want] {
					t.Errorf("GetRolePermissions(%v) missing permission %v", tt.role, want)
				}
			}

			// Ensure no duplicates
			if len(got) != len(gotMap) {
				t.Errorf("GetRolePermissions(%v) returned duplicate permissions", tt.role)
			}
		})
	}
}

func TestRoleHierarchy(t *testing.T) {
	t.Run("SuperAdmin inherits from SystemAdmin and OrgOwner", func(t *testing.T) {
		superPerms := GetRolePermissions(RoleSuperAdmin)
		systemPerms := GetRolePermissions(RoleSystemAdmin)
		ownerPerms := GetRolePermissions(RoleOrgOwner)

		// Convert to maps
		superMap := make(map[Permission]bool)
		for _, p := range superPerms {
			superMap[p] = true
		}

		// Check SystemAdmin permissions are included
		for _, p := range systemPerms {
			if !superMap[p] {
				t.Errorf("SuperAdmin missing SystemAdmin permission: %v", p)
			}
		}

		// Check OrgOwner permissions are included
		for _, p := range ownerPerms {
			if !superMap[p] {
				t.Errorf("SuperAdmin missing OrgOwner permission: %v", p)
			}
		}
	})

	t.Run("OrgAdmin inherits from OrgMember", func(t *testing.T) {
		adminPerms := GetRolePermissions(RoleOrgAdmin)
		memberPerms := GetRolePermissions(RoleOrgMember)

		adminMap := make(map[Permission]bool)
		for _, p := range adminPerms {
			adminMap[p] = true
		}

		for _, p := range memberPerms {
			if !adminMap[p] {
				t.Errorf("OrgAdmin missing OrgMember permission: %v", p)
			}
		}
	})

	t.Run("OrgMember inherits from OrgViewer", func(t *testing.T) {
		memberPerms := GetRolePermissions(RoleOrgMember)
		viewerPerms := GetRolePermissions(RoleOrgViewer)

		memberMap := make(map[Permission]bool)
		for _, p := range memberPerms {
			memberMap[p] = true
		}

		for _, p := range viewerPerms {
			if !memberMap[p] {
				t.Errorf("OrgMember missing OrgViewer permission: %v", p)
			}
		}
	})
}

func TestHasPermission(t *testing.T) {
	tests := []struct {
		name       string
		role       Role
		permission Permission
		want       bool
	}{
		{
			name:       "SuperAdmin has system manage",
			role:       RoleSuperAdmin,
			permission: PermSystemManage,
			want:       true,
		},
		{
			name:       "OrgMember has workflow create",
			role:       RoleOrgMember,
			permission: PermWorkflowCreate,
			want:       true,
		},
		{
			name:       "OrgViewer cannot delete content",
			role:       RoleOrgViewer,
			permission: PermContentDelete,
			want:       false,
		},
		{
			name:       "OrgGuest cannot create workflows",
			role:       RoleOrgGuest,
			permission: PermWorkflowCreate,
			want:       false,
		},
		{
			name:       "OrgAdmin has inherited member permissions",
			role:       RoleOrgAdmin,
			permission: PermContentRead, // From OrgViewer via OrgMember
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasPermission(tt.role, tt.permission); got != tt.want {
				t.Errorf(
					"HasPermission(%v, %v) = %v, want %v",
					tt.role,
					tt.permission,
					got,
					tt.want,
				)
			}
		})
	}
}

// TODO: Implement GetScopePermissions function
/*
func TestGetScopePermissions(t *testing.T) {
	tests := []struct {
		name  string
		scope Scope
		want  []Permission
	}{
		{
			name:  "read profile scope",
			scope: ScopeReadProfile,
			want:  []Permission{PermUserUpdateProfile},
		},
		{
			name:  "read organizations scope",
			scope: ScopeReadOrganizations,
			want:  []Permission{PermOrgRead},
		},
		{
			name:  "write workflows scope",
			scope: ScopeWriteWorkflows,
			want:  []Permission{PermWorkflowCreate, PermWorkflowUpdate},
		},
		{
			name:  "execute workflows scope",
			scope: ScopeExecuteWorkflows,
			want:  []Permission{PermWorkflowExecute},
		},
		// TODO: Add ScopeAdminOrganizations when implemented
		// {
		// 	name:  "admin organizations scope",
		// 	scope: ScopeAdminOrganizations,
		// 	want:  []Permission{PermOrgManageMembers, PermOrgManageBilling, PermOrgDelete},
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetScopePermissions(tt.scope)

			// Sort for comparison
			sort.Slice(got, func(i, j int) bool {
				return string(got[i]) < string(got[j])
			})
			sort.Slice(tt.want, func(i, j int) bool {
				return string(tt.want[i]) < string(tt.want[j])
			})

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetScopePermissions(%v) = %v, want %v", tt.scope, got, tt.want)
			}
		})
	}
}
*/

// TODO: Implement ValidateScopes function
/*
func TestValidateScopes(t *testing.T) {
	tests := []struct {
		name           string
		permissions    []Permission
		requiredScopes []Scope
		want           bool
	}{
		{
			name: "has required scopes",
			permissions: []Permission{
				PermUserUpdateProfile,
				PermOrgRead,
				PermWorkflowRead,
			},
			requiredScopes: []Scope{
				ScopeReadProfile,
				ScopeReadOrganizations,
			},
			want: true,
		},
		{
			name: "missing required scope",
			permissions: []Permission{
				PermOrgRead,
			},
			requiredScopes: []Scope{
				ScopeReadProfile,
				ScopeReadOrganizations,
			},
			want: false,
		},
		{
			name: "has write scope permissions",
			permissions: []Permission{
				PermWorkflowCreate,
				PermWorkflowUpdate,
			},
			requiredScopes: []Scope{
				ScopeWriteWorkflows,
			},
			want: true,
		},
		{
			name: "partial write scope permissions",
			permissions: []Permission{
				PermWorkflowCreate,
				// Missing PermWorkflowUpdate
			},
			requiredScopes: []Scope{
				ScopeWriteWorkflows,
			},
			want: false,
		},
		{
			name:           "no required scopes",
			permissions:    []Permission{PermOrgRead},
			requiredScopes: []Scope{},
			want:           true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateScopes(tt.permissions, tt.requiredScopes); got != tt.want {
				t.Errorf("ValidateScopes() = %v, want %v", got, tt.want)
			}
		})
	}
}
*/

// TODO: Implement GetDefaultScopes and GetAPIKeyScopes functions

func TestGetDefaultScopes(t *testing.T) {
	scopes := GetDefaultScopes()

	expectedScopes := []Scope{
		ScopeOpenID,
		ScopeEmail,
		ScopeProfile,
		ScopeReadProfile,
		ScopeReadOrganizations,
	}

	if len(scopes) != len(expectedScopes) {
		t.Errorf("GetDefaultScopes() returned %d scopes, want %d", len(scopes), len(expectedScopes))
	}

	// Check each expected scope is present
	scopeMap := make(map[Scope]bool)
	for _, s := range scopes {
		scopeMap[s] = true
	}

	for _, expected := range expectedScopes {
		if !scopeMap[expected] {
			t.Errorf("GetDefaultScopes() missing scope: %v", expected)
		}
	}
}

func TestGetAPIKeyScopes(t *testing.T) {
	scopes := GetAPIKeyScopes()

	expectedScopes := []Scope{
		ScopeReadWorkflows,
		ScopeWriteWorkflows,
		ScopeExecuteWorkflows,
		ScopeReadContent,
		ScopeWriteContent,
		ScopeExecuteAI,
	}

	if len(scopes) != len(expectedScopes) {
		t.Errorf("GetAPIKeyScopes() returned %d scopes, want %d", len(scopes), len(expectedScopes))
	}

	// Check each expected scope is present
	scopeMap := make(map[Scope]bool)
	for _, s := range scopes {
		scopeMap[s] = true
	}

	for _, expected := range expectedScopes {
		if !scopeMap[expected] {
			t.Errorf("GetAPIKeyScopes() missing scope: %v", expected)
		}
	}
}

func TestPermissionConsistency(t *testing.T) {
	t.Run("No role has duplicate permissions", func(t *testing.T) {
		for role := range RolePermissions {
			perms := GetRolePermissions(role)
			seen := make(map[Permission]bool)
			for _, p := range perms {
				if seen[p] {
					t.Errorf("Role %v has duplicate permission: %v", role, p)
				}
				seen[p] = true
			}
		}
	})

	t.Run("Higher roles have more permissions", func(t *testing.T) {
		// Check that higher roles in hierarchy have more or equal permissions
		ownerPerms := len(GetRolePermissions(RoleOrgOwner))
		adminPerms := len(GetRolePermissions(RoleOrgAdmin))
		memberPerms := len(GetRolePermissions(RoleOrgMember))
		viewerPerms := len(GetRolePermissions(RoleOrgViewer))
		guestPerms := len(GetRolePermissions(RoleOrgGuest))

		if ownerPerms < adminPerms {
			t.Errorf(
				"OrgOwner has fewer permissions (%d) than OrgAdmin (%d)",
				ownerPerms,
				adminPerms,
			)
		}
		if adminPerms < memberPerms {
			t.Errorf(
				"OrgAdmin has fewer permissions (%d) than OrgMember (%d)",
				adminPerms,
				memberPerms,
			)
		}
		if memberPerms < viewerPerms {
			t.Errorf(
				"OrgMember has fewer permissions (%d) than OrgViewer (%d)",
				memberPerms,
				viewerPerms,
			)
		}
		if viewerPerms < guestPerms {
			t.Errorf(
				"OrgViewer has fewer permissions (%d) than OrgGuest (%d)",
				viewerPerms,
				guestPerms,
			)
		}
	})
}

func BenchmarkGetRolePermissions(b *testing.B) {
	roles := []Role{
		RoleSuperAdmin,
		RoleOrgOwner,
		RoleOrgAdmin,
		RoleOrgMember,
		RoleOrgViewer,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, role := range roles {
			_ = GetRolePermissions(role)
		}
	}
}

func BenchmarkHasPermission(b *testing.B) {
	permissions := []Permission{
		PermOrgRead,
		PermWorkflowCreate,
		PermContentDelete,
		PermAIUsePremium,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, perm := range permissions {
			_ = HasPermission(RoleOrgAdmin, perm)
		}
	}
}

func BenchmarkValidateScopes(b *testing.B) {
	permissions := []Permission{
		PermUserUpdateProfile,
		PermOrgRead,
		PermWorkflowCreate,
		PermWorkflowUpdate,
		PermContentRead,
		PermContentCreate,
	}

	scopes := []Scope{
		ScopeReadProfile,
		ScopeWriteWorkflows,
		ScopeReadContent,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ValidateScopes(permissions, scopes)
	}
}
