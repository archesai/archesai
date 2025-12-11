package generators

// InternalPackageBase is the base import path for internal packages.
const InternalPackageBase = "github.com/archesai/archesai/pkg"

// InternalPackageImportPath returns the full import path for an internal package.
func InternalPackageImportPath(pkgName string) string {
	return InternalPackageBase + "/" + pkgName
}

// InternalPackageModelsPath returns the models import path for an internal package.
func InternalPackageModelsPath(pkgName string) string {
	return InternalPackageBase + "/" + pkgName + "/models"
}

// InternalPackageRepositoriesPath returns the repositories import path for an internal package.
func InternalPackageRepositoriesPath(pkgName string) string {
	return InternalPackageBase + "/" + pkgName + "/repositories"
}

// InternalPackageHandlersPath returns the handlers import path for an internal package.
func InternalPackageHandlersPath(pkgName string) string {
	return InternalPackageBase + "/" + pkgName + "/handlers"
}

// InternalPackageBootstrapPath returns the bootstrap import path for an internal package.
func InternalPackageBootstrapPath(pkgName string) string {
	return InternalPackageBase + "/" + pkgName + "/bootstrap"
}
