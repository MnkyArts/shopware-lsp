package hover

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/shopware/shopware-lsp/internal/entity"
	"github.com/shopware/shopware-lsp/internal/lsp"
	"github.com/shopware/shopware-lsp/internal/lsp/protocol"
	treesitterhelper "github.com/shopware/shopware-lsp/internal/tree_sitter_helper"
)

// EntityHoverProvider provides hover information for entities
type EntityHoverProvider struct {
	entityIndexer *entity.EntityIndexer
}

// NewEntityHoverProvider creates a new entity hover provider
func NewEntityHoverProvider(lspServer *lsp.Server) *EntityHoverProvider {
	entityIndexer, ok := lspServer.GetIndexer("entity.indexer")
	if !ok {
		return &EntityHoverProvider{}
	}

	return &EntityHoverProvider{
		entityIndexer: entityIndexer.(*entity.EntityIndexer),
	}
}

// GetHover returns hover information for entities
func (p *EntityHoverProvider) GetHover(ctx context.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	if p.entityIndexer == nil || params.Node == nil {
		return nil, nil
	}

	ext := strings.ToLower(filepath.Ext(params.TextDocument.URI))

	// Support both JavaScript and PHP files
	if ext != ".js" && ext != ".php" {
		return nil, nil
	}

	// Check if we're on a string node (potential entity name)
	if params.Node.Kind() != "string" && params.Node.Kind() != "string_content" {
		return nil, nil
	}

	// Get the entity name from the node
	entityName := treesitterhelper.GetNodeText(params.Node, params.DocumentContent)
	entityName = strings.Trim(entityName, `'"`)

	if entityName == "" {
		return nil, nil
	}

	// Look up the entity
	entities, err := p.entityIndexer.GetEntity(entityName)
	if err != nil || len(entities) == 0 {
		return nil, nil
	}

	// Use the first matching entity
	ent := entities[0]

	// Build markdown content
	var markdownContent strings.Builder
	markdownContent.WriteString(fmt.Sprintf("**Entity**: `%s`\n\n", ent.EntityName))
	markdownContent.WriteString(fmt.Sprintf("**Class**: %s\n\n", ent.ClassName))
	markdownContent.WriteString(fmt.Sprintf("**File**: %s\n", filepath.Base(ent.File)))

	return &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  "markdown",
			Value: markdownContent.String(),
		},
	}, nil
}
