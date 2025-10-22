package definition

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

// EntityDefinitionProvider provides go-to-definition for entity names
type EntityDefinitionProvider struct {
	entityIndexer *entity.EntityIndexer
}

// NewEntityDefinitionProvider creates a new entity definition provider
func NewEntityDefinitionProvider(lspServer *lsp.Server) *EntityDefinitionProvider {
	entityIndexer, ok := lspServer.GetIndexer("entity.indexer")
	if !ok {
		return &EntityDefinitionProvider{}
	}

	return &EntityDefinitionProvider{
		entityIndexer: entityIndexer.(*entity.EntityIndexer),
	}
}

// GetDefinition returns the definition locations for entity names
func (p *EntityDefinitionProvider) GetDefinition(ctx context.Context, params *protocol.DefinitionParams) []protocol.Location {
	if p.entityIndexer == nil || params.Node == nil {
		return []protocol.Location{}
	}

	ext := strings.ToLower(filepath.Ext(params.TextDocument.URI))

	// Support both JavaScript and PHP files
	if ext != ".js" && ext != ".php" {
		return []protocol.Location{}
	}

	// Check if we're on a string node (potential entity name)
	if params.Node.Kind() != "string" && params.Node.Kind() != "string_content" {
		return []protocol.Location{}
	}

	// Get the entity name from the node
	entityName := treesitterhelper.GetNodeText(params.Node, params.DocumentContent)
	entityName = strings.Trim(entityName, `'"`)

	if entityName == "" {
		return []protocol.Location{}
	}

	// Look up the entity
	entities, err := p.entityIndexer.GetEntity(entityName)
	if err != nil || len(entities) == 0 {
		return []protocol.Location{}
	}

	// Return locations for all matching entities
	locations := make([]protocol.Location, 0, len(entities))
	for _, ent := range entities {
		locations = append(locations, protocol.Location{
			URI: fmt.Sprintf("file://%s", ent.File),
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      ent.Line - 1,
					Character: 0,
				},
				End: protocol.Position{
					Line:      ent.Line - 1,
					Character: 0,
				},
			},
		})
	}

	return locations
}
