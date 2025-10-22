package completion

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/shopware/shopware-lsp/internal/admin"
	"github.com/shopware/shopware-lsp/internal/lsp"
	"github.com/shopware/shopware-lsp/internal/lsp/protocol"
)

// AdminComponentCompletionProvider provides completions for admin components
type AdminComponentCompletionProvider struct {
	componentIndexer *admin.ComponentIndexer
}

// NewAdminComponentCompletionProvider creates a new admin component completion provider
func NewAdminComponentCompletionProvider(lspServer *lsp.Server) *AdminComponentCompletionProvider {
	componentIndexer, ok := lspServer.GetIndexer("admin.component.indexer")
	if !ok {
		return &AdminComponentCompletionProvider{}
	}

	return &AdminComponentCompletionProvider{
		componentIndexer: componentIndexer.(*admin.ComponentIndexer),
	}
}

// GetCompletions returns completion items for admin components
func (p *AdminComponentCompletionProvider) GetCompletions(ctx context.Context, params *protocol.CompletionParams) []protocol.CompletionItem {
	if p.componentIndexer == nil || params.Node == nil {
		return []protocol.CompletionItem{}
	}

	ext := strings.ToLower(filepath.Ext(params.TextDocument.URI))
	
	// Only provide completions in JavaScript files within administration directories
	if ext != ".js" || !strings.Contains(params.TextDocument.URI, "/administration/") {
		return []protocol.CompletionItem{}
	}

	// For now, provide basic component name completions
	// This could be enhanced to be context-aware (e.g., only in specific call sites)
	return p.getComponentCompletions()
}

// getComponentCompletions returns all indexed component names as completion items
func (p *AdminComponentCompletionProvider) getComponentCompletions() []protocol.CompletionItem {
	componentNames, err := p.componentIndexer.GetComponentNames()
	if err != nil {
		return []protocol.CompletionItem{}
	}

	completionItems := make([]protocol.CompletionItem, 0, len(componentNames))
	for _, name := range componentNames {
		components, err := p.componentIndexer.GetComponent(name)
		if err != nil || len(components) == 0 {
			continue
		}

		component := components[0] // Use first match

		// Build detail string with props info
		detail := "Admin Component"
		if len(component.Props) > 0 {
			detail += " (props: "
			propNames := make([]string, 0, len(component.Props))
			for propName := range component.Props {
				propNames = append(propNames, propName)
			}
			detail += strings.Join(propNames, ", ") + ")"
		}

		completionItems = append(completionItems, protocol.CompletionItem{
			Label:  name,
			Detail: detail,
			Kind:   int(protocol.ClassCompletion),
		})
	}

	return completionItems
}

// GetTriggerCharacters returns the characters that trigger this completion provider
func (p *AdminComponentCompletionProvider) GetTriggerCharacters() []string {
	return []string{"'", "\""}
}
