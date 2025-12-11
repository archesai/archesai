// Package generators provides code generation from OpenAPI specifications.
package generators

// Priority levels for generators.
const (
	PriorityFirst  = 0
	PriorityNormal = 100
	PriorityLast   = 200
	PriorityFinal  = 300
)

// DefaultGenerators returns all standard generators.
func DefaultGenerators() []Generator {
	return []Generator{
		// Package-level
		{Name: "package_gomod", Priority: PriorityFirst, Generate: GenerateGoMod},
		{Name: "package_main", Priority: PriorityNormal, Generate: GenerateMain},

		// Models
		{Name: "models", Priority: PriorityNormal, Generate: GenerateSchemas},

		// Handlers
		{Name: "handlers", Priority: PriorityNormal, Generate: GenerateHandlers},
		{Name: "handler_stubs", Priority: PriorityNormal, Generate: GenerateHandlerStubs},

		// Routes
		{Name: "routes", Priority: PriorityNormal, Generate: GenerateRoutes},

		// Database
		{Name: "postgres", Priority: PriorityNormal, Generate: GeneratePostgres},
		{Name: "postgres_queries", Priority: PriorityNormal, Generate: GeneratePostgresQueries},
		{Name: "sqlite", Priority: PriorityNormal, Generate: GenerateSQLite},
		{Name: "sqlite_db", Priority: PriorityNormal, Generate: GenerateSQLiteDB},
		{Name: "sqlite_queries", Priority: PriorityNormal, Generate: GenerateSQLiteQueries},

		// App
		{Name: "app_bootstrap", Priority: PriorityNormal, Generate: GenerateAppBootstrap},
		{Name: "app_container", Priority: PriorityNormal, Generate: GenerateAppContainer},
		{Name: "app_handlers", Priority: PriorityNormal, Generate: GenerateAppHandlers},
		{Name: "app_infrastructure", Priority: PriorityNormal, Generate: GenerateAppInfrastructure},
		{Name: "app_routes", Priority: PriorityNormal, Generate: GenerateAppRoutes},

		// Database config
		{Name: "hcl", Priority: PriorityLast, Generate: GenerateHCL},
		{Name: "sqlc", Priority: PriorityFinal, Generate: GenerateSQLC},

		// Frontend
		{
			Name:     "frontend_package_json",
			Priority: PriorityFirst,
			Generate: GenerateFrontendPackageJSON,
		},
		{
			Name:     "frontend_vite_config",
			Priority: PriorityFirst,
			Generate: GenerateFrontendViteConfig,
		},
		{Name: "frontend_tsconfig", Priority: PriorityFirst, Generate: GenerateFrontendTSConfig},
		{
			Name:     "frontend_tsconfig_app",
			Priority: PriorityFirst,
			Generate: GenerateFrontendTSConfigApp,
		},
		{
			Name:     "frontend_tsconfig_spec",
			Priority: PriorityFirst,
			Generate: GenerateFrontendTSConfigSpec,
		},
		{
			Name:     "frontend_globals_css",
			Priority: PriorityFirst,
			Generate: GenerateFrontendGlobalsCSS,
		},
		{Name: "frontend_router", Priority: PriorityFirst, Generate: GenerateFrontendRouter},
		{Name: "frontend_fetcher", Priority: PriorityFirst, Generate: GenerateFrontendFetcher},
		{Name: "frontend_lib_index", Priority: PriorityNormal, Generate: GenerateFrontendLibIndex},
		{Name: "frontend_root_route", Priority: PriorityFirst, Generate: GenerateFrontendRootRoute},
		{Name: "frontend_app_route", Priority: PriorityFirst, Generate: GenerateFrontendAppRoute},
		{Name: "frontend_app_index", Priority: PriorityFirst, Generate: GenerateFrontendAppIndex},
		{
			Name:     "frontend_site_config",
			Priority: PriorityNormal,
			Generate: GenerateFrontendSiteConfig,
		},
		{
			Name:     "frontend_datatable",
			Priority: PriorityNormal,
			Generate: GenerateFrontendDataTables,
		},
		{Name: "frontend_form", Priority: PriorityNormal, Generate: GenerateFrontendForms},
		{
			Name:     "frontend_list_route",
			Priority: PriorityNormal,
			Generate: GenerateFrontendListRoutes,
		},
		{
			Name:     "frontend_detail_route",
			Priority: PriorityNormal,
			Generate: GenerateFrontendDetailRoutes,
		},
		{
			Name:     "frontend_get_session_ssr",
			Priority: PriorityFirst,
			Generate: GenerateFrontendGetSessionSSR,
		},
		{Name: "frontend_auth_route", Priority: PriorityFirst, Generate: GenerateFrontendAuthRoute},
		{
			Name:     "frontend_auth_login",
			Priority: PriorityFirst,
			Generate: GenerateFrontendAuthLoginPage,
		},
		{
			Name:     "frontend_auth_forgot_password",
			Priority: PriorityFirst,
			Generate: GenerateFrontendAuthForgotPasswordPage,
		},
		{
			Name:     "frontend_auth_magic_link_verify",
			Priority: PriorityFirst,
			Generate: GenerateFrontendAuthMagicLinkVerifyPage,
		},
		{
			Name:     "frontend_auth_oauth_callback",
			Priority: PriorityFirst,
			Generate: GenerateFrontendAuthOAuthCallbackPage,
		},
		{Name: "frontend_config", Priority: PriorityFirst, Generate: GenerateFrontendConfig},
		{
			Name:     "frontend_orval_config",
			Priority: PriorityFirst,
			Generate: GenerateFrontendOrvalConfig,
		},
		{
			Name:     "frontend_client_index",
			Priority: PriorityNormal,
			Generate: GenerateFrontendClientIndex,
		},
		{
			Name:     "frontend_site_utils",
			Priority: PriorityNormal,
			Generate: GenerateFrontendSiteUtils,
		},
		{
			Name:     "frontend_validators",
			Priority: PriorityNormal,
			Generate: GenerateFrontendValidators,
		},
	}
}
