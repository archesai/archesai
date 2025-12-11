package spec

// Tag represents a tag in the API specification.
type Tag struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// SecScheme represents a security scheme.
type SecScheme struct {
	Type        string `yaml:"type"`
	Scheme      string `yaml:"scheme"`
	Name        string `yaml:"name"`
	In          string `yaml:"in"`
	Description string `yaml:"description"`
}

// Security represents a security requirement.
type Security struct {
	Name   string   // Security scheme name
	Type   string   // Security type (http, apiKey, oauth2)
	Scheme string   // Security scheme (bearer for http, cookie for apiKey)
	Scopes []string // Required scopes
}
