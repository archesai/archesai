package codegen

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/archesai/archesai/internal/spec"
	"github.com/archesai/archesai/internal/strutil"
)

// GroupWeb is the generator group for frontend web files.
const GroupWeb = "web"

const (
	// Config generators
	genWebPackageJSON  = "web_package_json"
	genWebViteConfig   = "web_vite_config"
	genWebTSConfig     = "web_tsconfig"
	genWebTSConfigApp  = "web_tsconfig_app"
	genWebTSConfigSpec = "web_tsconfig_spec"
	genWebOrvalConfig  = "web_orval_config"

	// Auth generators
	genWebAuthRoute           = "web_auth_route"
	genWebAuthLogin           = "web_auth_login"
	genWebAuthForgotPassword  = "web_auth_forgot_password"
	genWebAuthMagicLinkVerify = "web_auth_magic_link_verify"
	genWebAuthOAuthCallback   = "web_auth_oauth_callback"

	// Lib generators
	genWebGlobalsCSS    = "web_globals_css"
	genWebSiteConfig    = "web_site_config"
	genWebFetcher       = "web_fetcher"
	genWebLibIndex      = "web_lib_index"
	genWebGetSessionSSR = "web_get_session_ssr"
	genWebConfig        = "web_config"
	genWebClientIndex   = "web_client_index"
	genWebSiteUtils     = "web_site_utils"
	genWebValidators    = "web_validators"

	// Route generators
	genWebRootRoute   = "web_root_route"
	genWebAppRoute    = "web_app_route"
	genWebAppIndex    = "web_app_index"
	genWebRouter      = "web_router"
	genWebListRoute   = "web_list_route"
	genWebDetailRoute = "web_detail_route"

	// Component generators
	genWebDatatable = "web_datatable"
	genWebForm      = "web_form"
)

// WebConfigTemplateData holds data for generating frontend config files.
type WebConfigTemplateData struct {
	ProjectName     string
	PackageName     string
	Entities        []*spec.EntityFrontendData
	Tags            []string
	APIHost         string
	PlatformURL     string
	Port            int
	SiteName        string
	SiteDescription string
	SiteURL         string
	SpecPath        string
}

// WebRouteTemplateData holds data for generating route files.
type WebRouteTemplateData struct {
	Entity      *spec.EntityFrontendData
	RouteType   string
	RoutePath   string
	ProjectName string
}

// ----------- Config generators -----------

// generateWebPackageJSON generates the package.json file.
func (c *Codegen) generateWebPackageJSON(s *spec.Spec) error {
	path := filepath.Join("web", "package.json")

	data := &WebConfigTemplateData{
		ProjectName: s.ProjectName,
		PackageName: "@" + strings.ReplaceAll(s.ProjectName, "/", "-") + "/web",
		APIHost:     "http://localhost:3001",
		Port:        3000,
	}

	return c.RenderTSXToFile("package.json.tmpl", path, data)
}

// generateWebViteConfig generates the vite.config.ts file.
func (c *Codegen) generateWebViteConfig(s *spec.Spec) error {
	path := filepath.Join("web", "vite.config.ts")

	data := &WebConfigTemplateData{
		ProjectName: s.ProjectName,
		APIHost:     "http://localhost:3001",
		Port:        3000,
	}

	return c.RenderTSXToFile("vite.config.ts.tmpl", path, data)
}

// generateWebTSConfig generates the tsconfig.json file.
func (c *Codegen) generateWebTSConfig(s *spec.Spec) error {
	path := filepath.Join("web", "tsconfig.json")

	data := &WebConfigTemplateData{
		ProjectName: s.ProjectName,
	}

	return c.RenderTSXToFile("tsconfig.json.tmpl", path, data)
}

// generateWebTSConfigApp generates the tsconfig.app.json file.
func (c *Codegen) generateWebTSConfigApp(_ *spec.Spec) error {
	path := filepath.Join("web", "tsconfig.app.json")

	return c.RenderTSXToFile("tsconfig.app.json.tmpl", path, nil)
}

// generateWebTSConfigSpec generates the tsconfig.spec.json file.
func (c *Codegen) generateWebTSConfigSpec(_ *spec.Spec) error {
	path := filepath.Join("web", "tsconfig.spec.json")

	return c.RenderTSXToFile("tsconfig.spec.json.tmpl", path, nil)
}

// generateWebOrvalConfig generates the orval.config.ts file.
func (c *Codegen) generateWebOrvalConfig(_ *spec.Spec) error {
	path := filepath.Join("web", "orval.config.ts")

	// Compute the bundled spec path relative to the web directory
	specPath := "../spec/openapi.bundled.yaml" // default fallback
	if c.cfg.Spec != nil {
		// The spec path is relative to the working directory
		// The orval config is in web/, so we need to go up one level first
		specDir := filepath.Dir(*c.cfg.Spec)
		specPath = filepath.Join("..", specDir, "openapi.bundled.yaml")
	}

	data := &WebConfigTemplateData{
		SpecPath: specPath,
	}

	return c.RenderTSXToFile("orval.config.ts.tmpl", path, data)
}

// ----------- Auth generators -----------

// generateWebAuthRoute generates the auth/route.tsx file.
func (c *Codegen) generateWebAuthRoute(_ *spec.Spec) error {
	path := filepath.Join("web", "src", "routes", "auth", "route.tsx")

	return c.RenderTSXToFile("auth-route.tsx.tmpl", path, nil)
}

// generateWebAuthLoginPage generates the auth/login/index.tsx file.
func (c *Codegen) generateWebAuthLoginPage(_ *spec.Spec) error {
	path := filepath.Join("web", "src", "routes", "auth", "login", "index.tsx")

	return c.RenderTSXToFile("auth-login.tsx.tmpl", path, nil)
}

// generateWebAuthForgotPasswordPage generates the auth/forgot-password/index.tsx file.
func (c *Codegen) generateWebAuthForgotPasswordPage(_ *spec.Spec) error {
	path := filepath.Join("web", "src", "routes", "auth", "forgot-password", "index.tsx")

	return c.RenderTSXToFile("auth-forgot-password.tsx.tmpl", path, nil)
}

// generateWebAuthMagicLinkVerifyPage generates the auth/magic-link/verify.tsx file.
func (c *Codegen) generateWebAuthMagicLinkVerifyPage(_ *spec.Spec) error {
	path := filepath.Join("web", "src", "routes", "auth", "magic-link", "verify.tsx")

	return c.RenderTSXToFile("auth-magic-link-verify.tsx.tmpl", path, nil)
}

// generateWebAuthOAuthCallbackPage generates the auth/oauth/callback/index.tsx file.
func (c *Codegen) generateWebAuthOAuthCallbackPage(_ *spec.Spec) error {
	path := filepath.Join("web", "src", "routes", "auth", "oauth", "callback", "index.tsx")

	return c.RenderTSXToFile("auth-oauth-callback.tsx.tmpl", path, nil)
}

// ----------- Lib generators -----------

// generateWebGlobalsCSS generates the globals.css file.
func (c *Codegen) generateWebGlobalsCSS(_ *spec.Spec) error {
	path := filepath.Join("web", "src", "styles", "globals.css")

	return c.RenderTSXToFile("globals.css.tmpl", path, nil)
}

// generateWebSiteConfig generates the site-config.ts file.
func (c *Codegen) generateWebSiteConfig(s *spec.Spec) error {
	path := filepath.Join("web", "src", "lib", "site-config.ts")

	entities := make([]*spec.EntityFrontendData, 0)
	for _, sch := range s.AllEntitySchemas() {
		entities = append(entities, s.BuildEntityFrontendData(sch))
	}

	data := &WebConfigTemplateData{
		ProjectName:     s.ProjectName,
		SiteName:        extractSiteName(s.ProjectName),
		SiteDescription: "Generated with Arches",
		SiteURL:         "http://localhost:3000",
		Entities:        entities,
	}

	return c.RenderTSXToFile("site-config.ts.tmpl", path, data)
}

// generateWebFetcher generates the fetcher.ts file.
func (c *Codegen) generateWebFetcher(_ *spec.Spec) error {
	path := filepath.Join("web", "src", "lib", "fetcher.ts")

	return c.RenderTSXToFile("fetcher.ts.tmpl", path, nil)
}

// generateWebLibIndex generates the lib/index.ts file.
func (c *Codegen) generateWebLibIndex(s *spec.Spec) error {
	path := filepath.Join("web", "src", "lib", "index.ts")

	entities := make([]*spec.EntityFrontendData, 0)
	for _, sch := range s.AllEntitySchemas() {
		entities = append(entities, s.BuildEntityFrontendData(sch))
	}

	data := &WebConfigTemplateData{
		Entities: entities,
	}

	return c.RenderTSXToFile("index.ts.tmpl", path, data)
}

// generateWebGetSessionSSR generates the get-session-ssr.ts file.
func (c *Codegen) generateWebGetSessionSSR(_ *spec.Spec) error {
	path := filepath.Join("web", "src", "lib", "get-session-ssr.ts")

	return c.RenderTSXToFile("get-session-ssr.ts.tmpl", path, nil)
}

// generateWebConfig generates the lib/config.ts file.
func (c *Codegen) generateWebConfig(_ *spec.Spec) error {
	path := filepath.Join("web", "src", "lib", "config.ts")

	data := &WebConfigTemplateData{
		APIHost:     "http://localhost:3001",
		PlatformURL: "http://localhost:3000",
	}

	return c.RenderTSXToFile("config.ts.tmpl", path, data)
}

// generateWebClientIndex generates the client/index.ts file that re-exports orval output.
func (c *Codegen) generateWebClientIndex(s *spec.Spec) error {
	path := filepath.Join("web", "src", "lib", "client", "index.ts")

	// Collect all unique tags from operations, converting to lowercase for folder names
	// Orval uses simple lowercase for folder names (e.g., "APIKey" -> "apikey")
	tagSet := make(map[string]bool)
	for _, op := range s.Operations {
		if op.Tag != "" {
			tagSet[strings.ToLower(op.Tag)] = true
		}
	}

	tags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		tags = append(tags, tag)
	}
	sort.Strings(tags)

	data := &WebConfigTemplateData{
		Tags: tags,
	}

	return c.RenderTSXToFile("client-index.ts.tmpl", path, data)
}

// generateWebSiteUtils generates the site-utils.ts file.
func (c *Codegen) generateWebSiteUtils(_ *spec.Spec) error {
	path := filepath.Join("web", "src", "lib", "site-utils.ts")

	return c.RenderTSXToFile("site-utils.ts.tmpl", path, nil)
}

// generateWebValidators generates the validators.ts file.
func (c *Codegen) generateWebValidators(_ *spec.Spec) error {
	path := filepath.Join("web", "src", "lib", "validators.ts")

	return c.RenderTSXToFile("validators.ts.tmpl", path, nil)
}

// ----------- Route generators -----------

// generateWebRootRoute generates the __root.tsx file.
func (c *Codegen) generateWebRootRoute(s *spec.Spec) error {
	path := filepath.Join("web", "src", "routes", "__root.tsx")

	data := &WebConfigTemplateData{
		ProjectName:     s.ProjectName,
		SiteName:        extractSiteName(s.ProjectName),
		SiteDescription: "Generated with Arches",
	}

	return c.RenderTSXToFile("__root.tsx.tmpl", path, data)
}

// generateWebAppRoute generates the _app/route.tsx file.
func (c *Codegen) generateWebAppRoute(s *spec.Spec) error {
	path := filepath.Join("web", "src", "routes", "_app", "route.tsx")

	data := &WebConfigTemplateData{
		ProjectName: s.ProjectName,
	}

	return c.RenderTSXToFile("app-route.tsx.tmpl", path, data)
}

// generateWebAppIndex generates the _app/index.tsx file.
func (c *Codegen) generateWebAppIndex(s *spec.Spec) error {
	path := filepath.Join("web", "src", "routes", "_app", "index.tsx")

	data := &WebConfigTemplateData{
		ProjectName:     s.ProjectName,
		SiteName:        extractSiteName(s.ProjectName),
		SiteDescription: "Generated with Arches",
	}

	return c.RenderTSXToFile("app-index.tsx.tmpl", path, data)
}

// generateWebRouter generates the router.tsx file.
func (c *Codegen) generateWebRouter(_ *spec.Spec) error {
	path := filepath.Join("web", "src", "router.tsx")

	return c.RenderTSXToFile("router.tsx.tmpl", path, nil)
}

// generateWebListRoutes generates list route pages for each entity.
func (c *Codegen) generateWebListRoutes(s *spec.Spec) error {
	for _, sch := range s.AllEntitySchemas() {
		entityData := s.BuildEntityFrontendData(sch)

		// Skip nested resources (entities that require parent IDs)
		if entityData.Operations.IsNested() {
			continue
		}

		path := filepath.Join(
			"web",
			"src",
			"routes",
			"_app",
			strutil.KebabCase(sch.Title)+"s",
			"index.tsx",
		)

		data := &WebRouteTemplateData{
			Entity:      entityData,
			RouteType:   "list",
			RoutePath:   "/" + strutil.KebabCase(sch.Title) + "s",
			ProjectName: s.ProjectName,
		}

		if err := c.RenderTSXToFile("list-route.tsx.tmpl", path, data); err != nil {
			return err
		}
	}

	return nil
}

// generateWebDetailRoutes generates detail route pages for each entity.
func (c *Codegen) generateWebDetailRoutes(s *spec.Spec) error {
	for _, sch := range s.AllEntitySchemas() {
		entityData := s.BuildEntityFrontendData(sch)

		// Skip nested resources (entities that require parent IDs)
		if entityData.Operations.IsNested() {
			continue
		}

		path := filepath.Join(
			"web",
			"src",
			"routes",
			"_app",
			strutil.KebabCase(sch.Title)+"s",
			"$"+strutil.CamelCase(sch.Title)+"ID",
			"index.tsx",
		)

		data := &WebRouteTemplateData{
			Entity:    entityData,
			RouteType: "detail",
			RoutePath: "/" + strutil.KebabCase(
				sch.Title,
			) + "s/$" + strutil.CamelCase(
				sch.Title,
			) + "ID",
			ProjectName: s.ProjectName,
		}

		if err := c.RenderTSXToFile("detail-route.tsx.tmpl", path, data); err != nil {
			return err
		}
	}

	return nil
}

// ----------- Component generators -----------

// generateWebDataTables generates datatable components for each entity.
func (c *Codegen) generateWebDataTables(s *spec.Spec) error {
	for _, sch := range s.AllEntitySchemas() {
		data := s.BuildEntityFrontendData(sch)

		// Skip nested resources (entities that require parent IDs)
		if data.Operations.IsNested() {
			continue
		}

		path := filepath.Join(
			"web",
			"src",
			"components",
			"datatables",
			strutil.KebabCase(sch.Title)+"-datatable.tsx",
		)

		if err := c.RenderTSXToFile("datatable.tsx.tmpl", path, data); err != nil {
			return err
		}
	}

	return nil
}

// generateWebForms generates form components for each entity.
func (c *Codegen) generateWebForms(s *spec.Spec) error {
	for _, sch := range s.AllEntitySchemas() {
		data := s.BuildEntityFrontendData(sch)

		// Skip nested resources (entities that require parent IDs)
		if data.Operations.IsNested() {
			continue
		}

		path := filepath.Join(
			"web",
			"src",
			"components",
			"forms",
			strutil.KebabCase(sch.Title)+"-form.tsx",
		)

		if err := c.RenderTSXToFile("form.tsx.tmpl", path, data); err != nil {
			return err
		}
	}

	return nil
}

// ----------- Helper functions -----------

// extractSiteName extracts a display name from the project path.
func extractSiteName(projectName string) string {
	parts := strings.Split(projectName, "/")
	if len(parts) > 0 {
		name := parts[len(parts)-1]
		return strutil.PascalCase(name)
	}
	return "App"
}
