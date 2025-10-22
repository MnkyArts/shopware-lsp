package definition

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/shopware/shopware-lsp/internal/admin"
	"github.com/shopware/shopware-lsp/internal/lsp"
	"github.com/shopware/shopware-lsp/internal/lsp/protocol"
	treesitterhelper "github.com/shopware/shopware-lsp/internal/tree_sitter_helper"
)

// AdminComponentDefinitionProvider provides go-to-definition for admin components
type AdminComponentDefinitionProvider struct {
	componentIndexer *admin.ComponentIndexer
}

// NewAdminComponentDefinitionProvider creates a new admin component definition provider
func NewAdminComponentDefinitionProvider(lspServer *lsp.Server) *AdminComponentDefinitionProvider {
	componentIndexer, ok := lspServer.GetIndexer("admin.component.indexer")
	if !ok {
		return &AdminComponentDefinitionProvider{}
	}

	return &AdminComponentDefinitionProvider{
		componentIndexer: componentIndexer.(*admin.ComponentIndexer),
	}
}

// GetDefinition returns the definition locations for admin components
func (p *AdminComponentDefinitionProvider) GetDefinition(ctx context.Context, params *protocol.DefinitionParams) []protocol.Location {
	if p.componentIndexer == nil || params.Node == nil {
		return []protocol.Location{}
	}

	ext := strings.ToLower(filepath.Ext(params.TextDocument.URI))

	// Only provide definitions in JavaScript files within administration directories
	if ext != ".js" || !strings.Contains(params.TextDocument.URI, "/administration/") {
		return []protocol.Location{}
	}

	// Check if we're on a string node (potential component name)
	if params.Node.Kind() != "string" && params.Node.Kind() != "string_content" {
		return []protocol.Location{}
	}

	// Get the component name from the node
	componentName := treesitterhelper.GetNodeText(params.Node, params.DocumentContent)
	componentName = strings.Trim(componentName, `'"`)

	if componentName == "" {
		return []protocol.Location{}
	}

	// Look up the component
	components, err := p.componentIndexer.GetComponent(componentName)
	if err != nil || len(components) == 0 {
		return []protocol.Location{}
	}

	// Return locations for all matching components
	locations := make([]protocol.Location, 0, len(components))
	for _, component := range components {
		locations = append(locations, protocol.Location{
			URI: fmt.Sprintf("file://%s", component.File),
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      component.Line - 1,
					Character: 0,
				},
				End: protocol.Position{
					Line:      component.Line - 1,
					Character: 0,
				},
			},
		})
	}

	return locations
}
