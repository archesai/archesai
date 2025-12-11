package codegen

import (
	"github.com/archesai/archesai/internal/located"
	"github.com/archesai/archesai/internal/spec"
)

// defaultGenerators returns all standard generators bound to the given spec and config.
func (c *Codegen) defaultGenerators(
	s *located.Located[spec.Spec],
) []generator {
	return []generator{
		// Module
		{
			name:     genGoMod,
			group:    GroupModule,
			priority: PriorityFirst,
			run:      func() error { return c.generateGoMod(s.Value) },
		},
		{
			name:     genMain,
			group:    GroupModule,
			priority: PriorityNormal,
			run:      func() error { return c.generateMain(s.Value) },
		},

		// Schemas
		{
			name:     genSchemas,
			group:    GroupSchemas,
			priority: PriorityNormal,
			run:      func() error { return c.generateSchemas(s.Value) },
		},

		// Operations
		{
			name:     genOperations,
			group:    GroupOperations,
			priority: PriorityNormal,
			run:      func() error { return c.generateOperations(s.Value) },
		},

		// HTTP
		{
			name:     genHTTP,
			group:    GroupHTTP,
			priority: PriorityNormal,
			run:      func() error { return c.generateHTTP(s.Value) },
		},

		// App
		{
			name:     genAppBootstrap,
			group:    GroupApp,
			priority: PriorityNormal,
			run:      func() error { return c.generateAppBootstrap(s) },
		},
		{
			name:     genAppContainer,
			group:    GroupApp,
			priority: PriorityNormal,
			run:      func() error { return c.generateAppContainer(s.Value) },
		},
		{
			name:     genAppOperations,
			group:    GroupApp,
			priority: PriorityNormal,
			run:      func() error { return c.generateAppOperations(s.Value) },
		},
		{
			name:     genAppPorts,
			group:    GroupApp,
			priority: PriorityNormal,
			run:      func() error { return c.generateAppPorts(s.Value) },
		},
		{
			name:     genAppHTTP,
			group:    GroupApp,
			priority: PriorityNormal,
			run:      func() error { return c.generateAppHTTP(s.Value) },
		},

		// Postgres
		{
			name:     genPostgres,
			group:    GroupPostgres,
			priority: PriorityNormal,
			run:      func() error { return c.generatePostgres(s.Value) },
		},
		{
			name:     genPostgresQueries,
			group:    GroupPostgres,
			priority: PriorityNormal,
			run:      func() error { return c.generatePostgresQueries(s.Value) },
		},
		{
			name:     genHCL,
			group:    GroupPostgres,
			priority: PriorityLast,
			run:      func() error { return c.generateHCL(s.Value) },
		},
		{
			name:     genSQLC,
			group:    GroupPostgres,
			priority: PriorityFinal,
			run:      func() error { return c.generateSQLC(s.Value) },
		},

		// SQLite
		{
			name:     genSQLite,
			group:    GroupSQLite,
			priority: PriorityNormal,
			run:      func() error { return c.generateSQLite(s.Value) },
		},
		{
			name:     genSQLiteDB,
			group:    GroupSQLite,
			priority: PriorityNormal,
			run:      func() error { return c.generateSQLiteDB(s.Value) },
		},
		{
			name:     genSQLiteQueries,
			group:    GroupSQLite,
			priority: PriorityNormal,
			run:      func() error { return c.generateSQLiteQueries(s.Value) },
		},

		// Web
		{
			name:     genWebPackageJSON,
			group:    GroupWeb,
			priority: PriorityFirst,
			run:      func() error { return c.generateWebPackageJSON(s.Value) },
		},
		{
			name:     genWebViteConfig,
			group:    GroupWeb,
			priority: PriorityFirst,
			run:      func() error { return c.generateWebViteConfig(s.Value) },
		},
		{
			name:     genWebTSConfig,
			group:    GroupWeb,
			priority: PriorityFirst,
			run:      func() error { return c.generateWebTSConfig(s.Value) },
		},
		{
			name:     genWebTSConfigApp,
			group:    GroupWeb,
			priority: PriorityFirst,
			run:      func() error { return c.generateWebTSConfigApp(s.Value) },
		},
		{
			name:     genWebTSConfigSpec,
			group:    GroupWeb,
			priority: PriorityFirst,
			run:      func() error { return c.generateWebTSConfigSpec(s.Value) },
		},
		{
			name:     genWebGlobalsCSS,
			group:    GroupWeb,
			priority: PriorityFirst,
			run:      func() error { return c.generateWebGlobalsCSS(s.Value) },
		},
		{
			name:     genWebRouter,
			group:    GroupWeb,
			priority: PriorityFirst,
			run:      func() error { return c.generateWebRouter(s.Value) },
		},
		{
			name:     genWebFetcher,
			group:    GroupWeb,
			priority: PriorityFirst,
			run:      func() error { return c.generateWebFetcher(s.Value) },
		},
		{
			name:     genWebLibIndex,
			group:    GroupWeb,
			priority: PriorityNormal,
			run:      func() error { return c.generateWebLibIndex(s.Value) },
		},
		{
			name:     genWebRootRoute,
			group:    GroupWeb,
			priority: PriorityFirst,
			run:      func() error { return c.generateWebRootRoute(s.Value) },
		},
		{
			name:     genWebAppRoute,
			group:    GroupWeb,
			priority: PriorityFirst,
			run:      func() error { return c.generateWebAppRoute(s.Value) },
		},
		{
			name:     genWebAppIndex,
			group:    GroupWeb,
			priority: PriorityFirst,
			run:      func() error { return c.generateWebAppIndex(s.Value) },
		},
		{
			name:     genWebSiteConfig,
			group:    GroupWeb,
			priority: PriorityNormal,
			run:      func() error { return c.generateWebSiteConfig(s.Value) },
		},
		{
			name:     genWebDatatable,
			group:    GroupWeb,
			priority: PriorityNormal,
			run:      func() error { return c.generateWebDataTables(s.Value) },
		},
		{
			name:     genWebForm,
			group:    GroupWeb,
			priority: PriorityNormal,
			run:      func() error { return c.generateWebForms(s.Value) },
		},
		{
			name:     genWebListRoute,
			group:    GroupWeb,
			priority: PriorityNormal,
			run:      func() error { return c.generateWebListRoutes(s.Value) },
		},
		{
			name:     genWebDetailRoute,
			group:    GroupWeb,
			priority: PriorityNormal,
			run:      func() error { return c.generateWebDetailRoutes(s.Value) },
		},
		{
			name:     genWebGetSessionSSR,
			group:    GroupWeb,
			priority: PriorityFirst,
			run:      func() error { return c.generateWebGetSessionSSR(s.Value) },
		},
		{
			name:     genWebAuthRoute,
			group:    GroupWeb,
			priority: PriorityFirst,
			run:      func() error { return c.generateWebAuthRoute(s.Value) },
		},
		{
			name:     genWebAuthLogin,
			group:    GroupWeb,
			priority: PriorityFirst,
			run:      func() error { return c.generateWebAuthLoginPage(s.Value) },
		},
		{
			name:     genWebAuthForgotPassword,
			group:    GroupWeb,
			priority: PriorityFirst,
			run:      func() error { return c.generateWebAuthForgotPasswordPage(s.Value) },
		},
		{
			name:     genWebAuthMagicLinkVerify,
			group:    GroupWeb,
			priority: PriorityFirst,
			run:      func() error { return c.generateWebAuthMagicLinkVerifyPage(s.Value) },
		},
		{
			name:     genWebAuthOAuthCallback,
			group:    GroupWeb,
			priority: PriorityFirst,
			run:      func() error { return c.generateWebAuthOAuthCallbackPage(s.Value) },
		},
		{
			name:     genWebConfig,
			group:    GroupWeb,
			priority: PriorityFirst,
			run:      func() error { return c.generateWebConfig(s.Value) },
		},
		{
			name:     genWebOrvalConfig,
			group:    GroupWeb,
			priority: PriorityFirst,
			run:      func() error { return c.generateWebOrvalConfig(s.Value) },
		},
		{
			name:     genWebClientIndex,
			group:    GroupWeb,
			priority: PriorityNormal,
			run:      func() error { return c.generateWebClientIndex(s.Value) },
		},
		{
			name:     genWebSiteUtils,
			group:    GroupWeb,
			priority: PriorityNormal,
			run:      func() error { return c.generateWebSiteUtils(s.Value) },
		},
		{
			name:     genWebValidators,
			group:    GroupWeb,
			priority: PriorityNormal,
			run:      func() error { return c.generateWebValidators(s.Value) },
		},
	}
}
