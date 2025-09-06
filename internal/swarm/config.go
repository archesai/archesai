package swarm

import (
	"net/http"

	"github.com/archesai/archesai/internal/llm"
)

// ClientConfig represents the configuration for an LLM client
type ClientConfig struct {
	Provider           llm.Provider
	AuthToken          string
	BaseURL            string
	OrgID              string
	APIVersion         string
	AssistantVersion   string
	ModelMapperFunc    func(model string) string // replace model to provider-specific deployment name
	HTTPClient         *http.Client
	EmptyMessagesLimit uint
	Options            map[string]interface{} // Additional provider-specific options
}
