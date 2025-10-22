package completion

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/shopware/shopware-lsp/internal/lsp"
	"github.com/shopware/shopware-lsp/internal/lsp/protocol"
	"github.com/shopware/shopware-lsp/internal/module"
	"github.com/shopware/shopware-lsp/internal/snippet"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

// ModuleCompletionProvider provides completions for module registrations
type ModuleCompletionProvider struct {
	moduleIndexer  *module.ModuleIndexer
	snippetIndexer *snippet.SnippetIndexer
}

// NewModuleCompletionProvider creates a new module completion provider
func NewModuleCompletionProvider(lspServer *lsp.Server) *ModuleCompletionProvider {
	moduleIndexer, _ := lspServer.GetIndexer("module.indexer")
	snippetIndexer, _ := lspServer.GetIndexer("snippet.indexer")

	var modIdx *module.ModuleIndexer
	var snipIdx *snippet.SnippetIndexer

	if moduleIndexer != nil {
		modIdx = moduleIndexer.(*module.ModuleIndexer)
	}
	if snippetIndexer != nil {
		snipIdx = snippetIndexer.(*snippet.SnippetIndexer)
	}

	return &ModuleCompletionProvider{
		moduleIndexer:  modIdx,
		snippetIndexer: snipIdx,
	}
}

// GetCompletions returns completion items for module contexts
func (p *ModuleCompletionProvider) GetCompletions(ctx context.Context, params *protocol.CompletionParams) []protocol.CompletionItem {
	if p.snippetIndexer == nil || params.Node == nil {
		return []protocol.CompletionItem{}
	}

	ext := strings.ToLower(filepath.Ext(params.TextDocument.URI))

	// Only provide completions in JavaScript files
	if ext != ".js" {
		return []protocol.CompletionItem{}
	}

	// Check if we're in a Module.register() context with snippet-like keys
	if !p.isInModuleContext(params.Node, params.DocumentContent) {
		return []protocol.CompletionItem{}
	}

	return p.getSnippetKeyCompletions()
}

// isInModuleContext checks if we're in a Module.register() context where snippet keys are expected
func (p *ModuleCompletionProvider) isInModuleContext(node *tree_sitter.Node, content []byte) bool {
	// Check if we're in a string node
	if node.Kind() != "string" && node.Kind() != "string_content" {
		return false
	}

	// Walk up to find if we're in Module.register()
	current := node.Parent()
	depth := 0
	for current != nil && depth < 10 {
		if current.Kind() == "call_expression" {
			// Check if this is Module.register
			if p.isModuleRegisterCall(current, content) {
				return true
			}
		}
		current = current.Parent()
		depth++
	}

	return false
}

// isModuleRegisterCall checks if a call expression is Module.register()
func (p *ModuleCompletionProvider) isModuleRegisterCall(callNode *tree_sitter.Node, content []byte) bool {
	if callNode.NamedChildCount() < 1 {
		return false
	}

	funcNode := callNode.NamedChild(0)
	if funcNode == nil || funcNode.Kind() != "member_expression" {
		return false
	}

	text := string(funcNode.Utf8Text(content))
	return strings.Contains(text, "Module") && strings.HasSuffix(text, ".register")
}

// getSnippetKeyCompletions returns admin snippet keys as completions
func (p *ModuleCompletionProvider) getSnippetKeyCompletions() []protocol.CompletionItem {
	snippets, err := p.snippetIndexer.GetFrontendSnippets()
	if err != nil {
		return []protocol.CompletionItem{}
	}

	completionItems := make([]protocol.CompletionItem, 0, len(snippets))
	for _, snippet := range snippets {
		completionItems = append(completionItems, protocol.CompletionItem{
			Label:  snippet,
			Kind:   int(protocol.TextCompletion),
			Detail: "Snippet Key",
		})
	}

	return completionItems
}

// GetTriggerCharacters returns the characters that trigger this completion provider
func (p *ModuleCompletionProvider) GetTriggerCharacters() []string {
	return []string{"'", "\"", "."}
}
