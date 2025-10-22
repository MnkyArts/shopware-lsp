package completion

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/shopware/shopware-lsp/internal/entity"
	"github.com/shopware/shopware-lsp/internal/lsp"
	"github.com/shopware/shopware-lsp/internal/lsp/protocol"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

// RepositoryCompletionProvider provides completions for repositoryFactory.create() calls
type RepositoryCompletionProvider struct {
	entityIndexer *entity.EntityIndexer
}

// NewRepositoryCompletionProvider creates a new repository completion provider
func NewRepositoryCompletionProvider(lspServer *lsp.Server) *RepositoryCompletionProvider {
	entityIndexer, ok := lspServer.GetIndexer("entity.indexer")
	if !ok {
		return &RepositoryCompletionProvider{}
	}

	return &RepositoryCompletionProvider{
		entityIndexer: entityIndexer.(*entity.EntityIndexer),
	}
}

// GetCompletions returns completion items for entity names in repositoryFactory.create() calls
func (p *RepositoryCompletionProvider) GetCompletions(ctx context.Context, params *protocol.CompletionParams) []protocol.CompletionItem {
	if p.entityIndexer == nil || params.Node == nil {
		return []protocol.CompletionItem{}
	}

	ext := strings.ToLower(filepath.Ext(params.TextDocument.URI))

	// Only provide completions in JavaScript files
	if ext != ".js" {
		return []protocol.CompletionItem{}
	}

	// Check if we're in a repositoryFactory.create() call
	if !p.isInRepositoryFactoryCreate(params.Node, params.DocumentContent) {
		return []protocol.CompletionItem{}
	}

	return p.getEntityCompletions()
}

// isInRepositoryFactoryCreate checks if the cursor is in a repositoryFactory.create() call
func (p *RepositoryCompletionProvider) isInRepositoryFactoryCreate(node *tree_sitter.Node, content []byte) bool {
	// Check if we're in a string node
	if node.Kind() != "string" && node.Kind() != "string_content" {
		return false
	}

	// Walk up the tree to find a call_expression
	current := node.Parent()
	for current != nil && current.Kind() != "program" {
		if current.Kind() == "call_expression" {
			// Check if this is a repositoryFactory.create() call
			if p.isRepositoryFactoryCreateCall(current, content) {
				return true
			}
		}
		current = current.Parent()
	}

	return false
}

// isRepositoryFactoryCreateCall checks if a call_expression is repositoryFactory.create()
func (p *RepositoryCompletionProvider) isRepositoryFactoryCreateCall(callNode *tree_sitter.Node, content []byte) bool {
	if callNode.NamedChildCount() < 1 {
		return false
	}

	// Get the function being called
	funcNode := callNode.NamedChild(0)
	if funcNode == nil || funcNode.Kind() != "member_expression" {
		return false
	}

	// Check if it's *.repositoryFactory.create
	text := string(funcNode.Utf8Text(content))
	return strings.Contains(text, "repositoryFactory") && strings.HasSuffix(text, ".create")
}

// getEntityCompletions returns all indexed entity names as completion items
func (p *RepositoryCompletionProvider) getEntityCompletions() []protocol.CompletionItem {
	entityNames, err := p.entityIndexer.GetEntityNames()
	if err != nil {
		return []protocol.CompletionItem{}
	}

	completionItems := make([]protocol.CompletionItem, 0, len(entityNames))
	for _, name := range entityNames {
		entities, err := p.entityIndexer.GetEntity(name)
		if err != nil || len(entities) == 0 {
			continue
		}

		entity := entities[0] // Use first match

		completionItems = append(completionItems, protocol.CompletionItem{
			Label:  name,
			Detail: "Entity: " + entity.ClassName,
			Kind:   int(protocol.ClassCompletion),
		})
	}

	return completionItems
}

// GetTriggerCharacters returns the characters that trigger this completion provider
func (p *RepositoryCompletionProvider) GetTriggerCharacters() []string {
	return []string{"'", "\""}
}
