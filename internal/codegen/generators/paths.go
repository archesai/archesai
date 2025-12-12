package generators

// InternalPackageBase is the base import path for internal packages.
const InternalPackageBase = "github.com/archesai/archesai/pkg"

// InternalPackageImportPath returns the full import path for an internal package.
func InternalPackageImportPath(pkgName string) string {
	return InternalPackageBase + "/" + pkgName
}

// InternalPackageModelsPath returns the models import path for an internal package.
// Models includes domain types (entities, value objects) and repository interfaces.
func InternalPackageModelsPath(pkgName string) string {
	return InternalPackageBase + "/" + pkgName + "/models"
}

// InternalPackageHandlersPath returns the handlers import path for an internal package.
// Handlers contain business logic (input/output DTOs and handler interfaces).
func InternalPackageHandlersPath(pkgName string) string {
	return InternalPackageBase + "/" + pkgName + "/handlers"
}

// InternalPackageRoutesPath returns the routes import path for an internal package.
// Routes contain HTTP layer code (request parsing, response writing).
func InternalPackageRoutesPath(pkgName string) string {
	return InternalPackageBase + "/" + pkgName + "/routes"
}

// InternalPackageAppPath returns the app import path for an internal package.
// App contains dependency injection and handler initialization code.
func InternalPackageAppPath(pkgName string) string {
	return InternalPackageBase + "/" + pkgName + "/app"
}
