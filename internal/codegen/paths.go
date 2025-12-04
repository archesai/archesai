package codegen

// InternalPackageBase is the base import path for internal packages.
const InternalPackageBase = "github.com/archesai/archesai/pkg/"

// InternalPackageImportPath returns the full import path for an internal package.
func InternalPackageImportPath(pkgName string) string {
	return InternalPackageBase + pkgName
}

// InternalPackageModelsPath returns the models import path for an internal package.
func InternalPackageModelsPath(pkgName string) string {
	return InternalPackageBase + pkgName + "/models"
}

// InternalPackageRepositoriesPath returns the repositories import path for an internal package.
func InternalPackageRepositoriesPath(pkgName string) string {
	return InternalPackageBase + pkgName + "/repositories"
}

// InternalPackageApplicationPath returns the application import path for an internal package.
func InternalPackageApplicationPath(pkgName string) string {
	return InternalPackageBase + pkgName + "/application"
}

// InternalPackageBootstrapPath returns the bootstrap import path for an internal package.
func InternalPackageBootstrapPath(pkgName string) string {
	return InternalPackageBase + pkgName + "/bootstrap"
}
