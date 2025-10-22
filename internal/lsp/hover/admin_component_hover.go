package hover

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

// AdminComponentHoverProvider provides hover information for admin components
type AdminComponentHoverProvider struct {
	componentIndexer *admin.ComponentIndexer
}

// NewAdminComponentHoverProvider creates a new admin component hover provider
func NewAdminComponentHoverProvider(lspServer *lsp.Server) *AdminComponentHoverProvider {
	componentIndexer, ok := lspServer.GetIndexer("admin.component.indexer")
	if !ok {
		return &AdminComponentHoverProvider{}
	}

	return &AdminComponentHoverProvider{
		componentIndexer: componentIndexer.(*admin.ComponentIndexer),
	}
}

// GetHover returns hover information for admin components
func (p *AdminComponentHoverProvider) GetHover(ctx context.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	if p.componentIndexer == nil || params.Node == nil {
		return nil, nil
	}

	ext := strings.ToLower(filepath.Ext(params.TextDocument.URI))

	// Only provide hover in JavaScript files within administration directories
	if ext != ".js" || !strings.Contains(params.TextDocument.URI, "/administration/") {
		return nil, nil
	}

	// Check if we're on a string node (potential component name)
	if params.Node.Kind() != "string" && params.Node.Kind() != "string_content" {
		return nil, nil
	}

	// Get the component name from the node
	componentName := treesitterhelper.GetNodeText(params.Node, params.DocumentContent)
	componentName = strings.Trim(componentName, `'"`)

	if componentName == "" {
		return nil, nil
	}

	// Look up the component
	components, err := p.componentIndexer.GetComponent(componentName)
	if err != nil || len(components) == 0 {
		return nil, nil
	}

	// Use the first matching component
	component := components[0]

	// Build markdown content
	var markdownContent strings.Builder
	markdownContent.WriteString(fmt.Sprintf("**Admin Component**: `%s`\n\n", component.Name))
	markdownContent.WriteString(fmt.Sprintf("**File**: %s\n\n", filepath.Base(component.File)))

	// Add props information
	if len(component.Props) > 0 {
		markdownContent.WriteString("**Props**:\n")
		for propName := range component.Props {
			markdownContent.WriteString(fmt.Sprintf("- `%s`\n", propName))
		}
		markdownContent.WriteString("\n")
	}

	// Add computed properties information
	if len(component.Computed) > 0 {
		markdownContent.WriteString("**Computed Properties**:\n")
		for _, computed := range component.Computed {
			markdownContent.WriteString(fmt.Sprintf("- `%s`\n", computed))
		}
	}

	return &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  "markdown",
			Value: markdownContent.String(),
		},
	}, nil
}
