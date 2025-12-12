package generators

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/archesai/archesai/internal/strutil"
)

// GenerateFrontendPackageJSON generates the package.json file.
func GenerateFrontendPackageJSON(ctx *GeneratorContext) error {
	path := filepath.Join("web", "package.json")

	data := &FrontendConfigTemplateData{
		ProjectName: ctx.ProjectName,
		PackageName: "@" + strings.ReplaceAll(ctx.ProjectName, "/", "-") + "/web",
		APIHost:     "http://localhost:3001",
		Port:        3000,
	}

	return ctx.RenderTSXToFile("package.json.tmpl", path, data)
}

// GenerateFrontendViteConfig generates the vite.config.ts file.
func GenerateFrontendViteConfig(ctx *GeneratorContext) error {
	path := filepath.Join("web", "vite.config.ts")

	data := &FrontendConfigTemplateData{
		ProjectName: ctx.ProjectName,
		APIHost:     "http://localhost:3001",
		Port:        3000,
	}

	return ctx.RenderTSXToFile("vite.config.ts.tmpl", path, data)
}

// GenerateFrontendTSConfig generates the tsconfig.json file.
func GenerateFrontendTSConfig(ctx *GeneratorContext) error {
	path := filepath.Join("web", "tsconfig.json")

	data := &FrontendConfigTemplateData{
		ProjectName: ctx.ProjectName,
	}

	return ctx.RenderTSXToFile("tsconfig.json.tmpl", path, data)
}

// GenerateFrontendTSConfigApp generates the tsconfig.app.json file.
func GenerateFrontendTSConfigApp(ctx *GeneratorContext) error {
	path := filepath.Join("web", "tsconfig.app.json")

	return ctx.RenderTSXToFile("tsconfig.app.json.tmpl", path, nil)
}

// GenerateFrontendTSConfigSpec generates the tsconfig.spec.json file.
func GenerateFrontendTSConfigSpec(ctx *GeneratorContext) error {
	path := filepath.Join("web", "tsconfig.spec.json")

	return ctx.RenderTSXToFile("tsconfig.spec.json.tmpl", path, nil)
}

// GenerateFrontendGlobalsCSS generates the globals.css file.
func GenerateFrontendGlobalsCSS(ctx *GeneratorContext) error {
	path := filepath.Join("web", "src", "styles", "globals.css")

	return ctx.RenderTSXToFile("globals.css.tmpl", path, nil)
}

// GenerateFrontendRootRoute generates the __root.tsx file.
func GenerateFrontendRootRoute(ctx *GeneratorContext) error {
	path := filepath.Join("web", "src", "routes", "__root.tsx")

	data := &FrontendConfigTemplateData{
		ProjectName:     ctx.ProjectName,
		SiteName:        extractSiteName(ctx.ProjectName),
		SiteDescription: "Generated with Arches",
	}

	return ctx.RenderTSXToFile("__root.tsx.tmpl", path, data)
}

// GenerateFrontendAppRoute generates the _app/route.tsx file.
func GenerateFrontendAppRoute(ctx *GeneratorContext) error {
	path := filepath.Join("web", "src", "routes", "_app", "route.tsx")

	data := &FrontendConfigTemplateData{
		ProjectName: ctx.ProjectName,
	}

	return ctx.RenderTSXToFile("app-route.tsx.tmpl", path, data)
}

// GenerateFrontendAppIndex generates the _app/index.tsx file.
func GenerateFrontendAppIndex(ctx *GeneratorContext) error {
	path := filepath.Join("web", "src", "routes", "_app", "index.tsx")

	data := &FrontendConfigTemplateData{
		ProjectName:     ctx.ProjectName,
		SiteName:        extractSiteName(ctx.ProjectName),
		SiteDescription: "Generated with Arches",
	}

	return ctx.RenderTSXToFile("app-index.tsx.tmpl", path, data)
}

// GenerateFrontendSiteConfig generates the site-config.ts file.
func GenerateFrontendSiteConfig(ctx *GeneratorContext) error {
	path := filepath.Join("web", "src", "lib", "site-config.ts")

	entities := make([]*EntityFrontendData, 0)
	for _, schema := range ctx.AllEntitySchemas() {
		entities = append(entities, BuildEntityFrontendData(ctx, schema))
	}

	data := &FrontendConfigTemplateData{
		ProjectName:     ctx.ProjectName,
		SiteName:        extractSiteName(ctx.ProjectName),
		SiteDescription: "Generated with Arches",
		SiteURL:         "http://localhost:3000",
		Entities:        entities,
	}

	return ctx.RenderTSXToFile("site-config.ts.tmpl", path, data)
}

// GenerateFrontendDataTables generates datatable components for each entity.
func GenerateFrontendDataTables(ctx *GeneratorContext) error {
	for _, schema := range ctx.AllEntitySchemas() {
		data := BuildEntityFrontendData(ctx, schema)

		// Skip nested resources (entities that require parent IDs)
		if IsNestedResource(data.Operations) {
			continue
		}

		path := filepath.Join(
			"web",
			"src",
			"components",
			"datatables",
			strutil.KebabCase(schema.Name)+"-datatable.tsx",
		)

		if err := ctx.RenderTSXToFile("datatable.tsx.tmpl", path, data); err != nil {
			return err
		}
	}

	return nil
}

// GenerateFrontendForms generates form components for each entity.
func GenerateFrontendForms(ctx *GeneratorContext) error {
	for _, schema := range ctx.AllEntitySchemas() {
		data := BuildEntityFrontendData(ctx, schema)

		// Skip nested resources (entities that require parent IDs)
		if IsNestedResource(data.Operations) {
			continue
		}

		path := filepath.Join(
			"web",
			"src",
			"components",
			"forms",
			strutil.KebabCase(schema.Name)+"-form.tsx",
		)

		if err := ctx.RenderTSXToFile("form.tsx.tmpl", path, data); err != nil {
			return err
		}
	}

	return nil
}

// GenerateFrontendListRoutes generates list route pages for each entity.
func GenerateFrontendListRoutes(ctx *GeneratorContext) error {
	for _, schema := range ctx.AllEntitySchemas() {
		entityData := BuildEntityFrontendData(ctx, schema)

		// Skip nested resources (entities that require parent IDs)
		if IsNestedResource(entityData.Operations) {
			continue
		}

		path := filepath.Join(
			"web",
			"src",
			"routes",
			"_app",
			strutil.KebabCase(schema.Name)+"s",
			"index.tsx",
		)

		data := &FrontendRouteTemplateData{
			Entity:      entityData,
			RouteType:   "list",
			RoutePath:   "/" + strutil.KebabCase(schema.Name) + "s",
			ProjectName: ctx.ProjectName,
		}

		if err := ctx.RenderTSXToFile("list-route.tsx.tmpl", path, data); err != nil {
			return err
		}
	}

	return nil
}

// GenerateFrontendDetailRoutes generates detail route pages for each entity.
func GenerateFrontendDetailRoutes(ctx *GeneratorContext) error {
	for _, schema := range ctx.AllEntitySchemas() {
		entityData := BuildEntityFrontendData(ctx, schema)

		// Skip nested resources (entities that require parent IDs)
		if IsNestedResource(entityData.Operations) {
			continue
		}

		path := filepath.Join(
			"web",
			"src",
			"routes",
			"_app",
			strutil.KebabCase(schema.Name)+"s",
			"$"+strutil.CamelCase(schema.Name)+"ID",
			"index.tsx",
		)

		data := &FrontendRouteTemplateData{
			Entity:    entityData,
			RouteType: "detail",
			RoutePath: "/" + strutil.KebabCase(
				schema.Name,
			) + "s/$" + strutil.CamelCase(
				schema.Name,
			) + "ID",
			ProjectName: ctx.ProjectName,
		}

		if err := ctx.RenderTSXToFile("detail-route.tsx.tmpl", path, data); err != nil {
			return err
		}
	}

	return nil
}

// GenerateFrontendRouter generates the router.tsx file.
func GenerateFrontendRouter(ctx *GeneratorContext) error {
	path := filepath.Join("web", "src", "router.tsx")

	return ctx.RenderTSXToFile("router.tsx.tmpl", path, nil)
}

// GenerateFrontendFetcher generates the fetcher.ts file.
func GenerateFrontendFetcher(ctx *GeneratorContext) error {
	path := filepath.Join("web", "src", "lib", "fetcher.ts")

	return ctx.RenderTSXToFile("fetcher.ts.tmpl", path, nil)
}

// GenerateFrontendLibIndex generates the lib/index.ts file.
func GenerateFrontendLibIndex(ctx *GeneratorContext) error {
	path := filepath.Join("web", "src", "lib", "index.ts")

	entities := make([]*EntityFrontendData, 0)
	for _, schema := range ctx.AllEntitySchemas() {
		entities = append(entities, BuildEntityFrontendData(ctx, schema))
	}

	data := &FrontendConfigTemplateData{
		Entities: entities,
	}

	return ctx.RenderTSXToFile("index.ts.tmpl", path, data)
}

// GenerateFrontendGetSessionSSR generates the get-session-ssr.ts file.
func GenerateFrontendGetSessionSSR(ctx *GeneratorContext) error {
	path := filepath.Join("web", "src", "lib", "get-session-ssr.ts")

	return ctx.RenderTSXToFile("get-session-ssr.ts.tmpl", path, nil)
}

// GenerateFrontendAuthRoute generates the auth/route.tsx file.
func GenerateFrontendAuthRoute(ctx *GeneratorContext) error {
	path := filepath.Join("web", "src", "routes", "auth", "route.tsx")

	return ctx.RenderTSXToFile("auth-route.tsx.tmpl", path, nil)
}

// GenerateFrontendAuthLoginPage generates the auth/login/index.tsx file.
func GenerateFrontendAuthLoginPage(ctx *GeneratorContext) error {
	path := filepath.Join("web", "src", "routes", "auth", "login", "index.tsx")

	return ctx.RenderTSXToFile("auth-login.tsx.tmpl", path, nil)
}

// GenerateFrontendAuthForgotPasswordPage generates the auth/forgot-password/index.tsx file.
func GenerateFrontendAuthForgotPasswordPage(ctx *GeneratorContext) error {
	path := filepath.Join("web", "src", "routes", "auth", "forgot-password", "index.tsx")

	return ctx.RenderTSXToFile("auth-forgot-password.tsx.tmpl", path, nil)
}

// GenerateFrontendAuthMagicLinkVerifyPage generates the auth/magic-link/verify.tsx file.
func GenerateFrontendAuthMagicLinkVerifyPage(ctx *GeneratorContext) error {
	path := filepath.Join("web", "src", "routes", "auth", "magic-link", "verify.tsx")

	return ctx.RenderTSXToFile("auth-magic-link-verify.tsx.tmpl", path, nil)
}

// GenerateFrontendAuthOAuthCallbackPage generates the auth/oauth/callback/index.tsx file.
func GenerateFrontendAuthOAuthCallbackPage(ctx *GeneratorContext) error {
	path := filepath.Join("web", "src", "routes", "auth", "oauth", "callback", "index.tsx")

	return ctx.RenderTSXToFile("auth-oauth-callback.tsx.tmpl", path, nil)
}

// GenerateFrontendConfig generates the lib/config.ts file.
func GenerateFrontendConfig(ctx *GeneratorContext) error {
	path := filepath.Join("web", "src", "lib", "config.ts")

	data := &FrontendConfigTemplateData{
		APIHost:     "http://localhost:3001",
		PlatformURL: "http://localhost:3000",
	}

	return ctx.RenderTSXToFile("config.ts.tmpl", path, data)
}

// GenerateFrontendOrvalConfig generates the orval.config.ts file.
func GenerateFrontendOrvalConfig(ctx *GeneratorContext) error {
	path := filepath.Join("web", "orval.config.ts")

	// Calculate relative path from web/ (inside output dir) to the bundled spec.
	// The bundled spec is at {specDir}/openapi.bundled.yaml where specDir
	// is the directory containing the source spec.

	// Get absolute path to the bundled spec
	specDir := filepath.Dir(ctx.SpecPath)
	bundledSpecPath := filepath.Join(specDir, "openapi.bundled.yaml")
	absBundledSpec, err := filepath.Abs(bundledSpecPath)
	if err != nil {
		absBundledSpec = bundledSpecPath
	}

	// Get absolute path to the web directory inside the output
	webDir := filepath.Join(ctx.Storage.BaseDir(), "web")
	absWebDir, err := filepath.Abs(webDir)
	if err != nil {
		absWebDir = webDir
	}

	// Compute relative path from web dir to bundled spec
	relPath, err := filepath.Rel(absWebDir, absBundledSpec)
	if err != nil {
		// Fallback to default if we can't compute relative path
		relPath = "../api/openapi.bundled.yaml"
	}

	data := &FrontendConfigTemplateData{
		SpecPath: relPath,
	}

	return ctx.RenderTSXToFile("orval.config.ts.tmpl", path, data)
}

// GenerateFrontendClientIndex generates the client/index.ts file that re-exports orval output.
func GenerateFrontendClientIndex(ctx *GeneratorContext) error {
	path := filepath.Join("web", "src", "lib", "client", "index.ts")

	// Collect all unique tags from operations, converting to lowercase for folder names
	// Orval uses simple lowercase for folder names (e.g., "APIKey" -> "apikey")
	tagSet := make(map[string]bool)
	for _, op := range ctx.Spec.Operations {
		if op.Tag != "" {
			tagSet[strings.ToLower(op.Tag)] = true
		}
	}

	tags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		tags = append(tags, tag)
	}
	sort.Strings(tags)

	data := &FrontendConfigTemplateData{
		Tags: tags,
	}

	return ctx.RenderTSXToFile("client-index.ts.tmpl", path, data)
}

// GenerateFrontendSiteUtils generates the site-utils.ts file.
func GenerateFrontendSiteUtils(ctx *GeneratorContext) error {
	path := filepath.Join("web", "src", "lib", "site-utils.ts")

	return ctx.RenderTSXToFile("site-utils.ts.tmpl", path, nil)
}

// GenerateFrontendValidators generates the validators.ts file.
func GenerateFrontendValidators(ctx *GeneratorContext) error {
	path := filepath.Join("web", "src", "lib", "validators.ts")

	return ctx.RenderTSXToFile("validators.ts.tmpl", path, nil)
}

// extractSiteName extracts a display name from the project path.
func extractSiteName(projectName string) string {
	parts := strings.Split(projectName, "/")
	if len(parts) > 0 {
		name := parts[len(parts)-1]
		return strutil.PascalCase(name)
	}
	return "App"
}
