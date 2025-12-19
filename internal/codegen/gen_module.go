package codegen

import "github.com/archesai/archesai/internal/spec"

// GroupModule is the generator group for Go module files.
const GroupModule = "module"

const (
	genGoMod = "module_go.mod"
	genMain  = "module_main"
)

// generateGoMod generates go.mod file.
func (c *Codegen) generateGoMod(s *spec.Spec) error {
	path := "go.mod"
	exists, err := c.storage.Exists(path)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	return c.RenderToFile(genGoMod+".tmpl", path, s)
}

// generateMain generates main.go for composition apps.
func (c *Codegen) generateMain(s *spec.Spec) error {
	// Only generate for composition apps
	if len(s.ComposedPackages()) == 0 {
		return nil
	}

	return c.RenderToFile(genMain+".go.tmpl", "main.gen.go", s)
}
