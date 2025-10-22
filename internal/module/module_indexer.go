package module

import (
	"path/filepath"
	"strings"

	"github.com/shopware/shopware-lsp/internal/indexer"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

// ModuleIndexer indexes Shopware admin module registrations
type ModuleIndexer struct {
	index *indexer.DataIndexer[Module]
}

// NewModuleIndexer creates a new module indexer
func NewModuleIndexer(configDir string) (*ModuleIndexer, error) {
	idx, err := indexer.NewDataIndexer[Module](filepath.Join(configDir, "module.db"))
	if err != nil {
		return nil, err
	}

	return &ModuleIndexer{
		index: idx,
	}, nil
}

// ID returns the indexer identifier
func (mi *ModuleIndexer) ID() string {
	return "module.indexer"
}

// Index processes a JavaScript file and indexes module registrations
func (mi *ModuleIndexer) Index(path string, node *tree_sitter.Node, fileContent []byte) error {
	// Only index JavaScript files in admin directories
	if !strings.HasSuffix(path, ".js") ||
		!strings.Contains(path, "/administration/") ||
		strings.Contains(path, "/node_modules/") ||
		strings.Contains(path, "/_fixtures/") {
		return nil
	}

	modules, err := parseModuleFile(node, fileContent, path)
	if err != nil {
		return err
	}

	// Batch save all modules found in this file
	batchSave := make(map[string]map[string]Module)
	for moduleName, module := range modules {
		if _, ok := batchSave[module.File]; !ok {
			batchSave[module.File] = make(map[string]Module)
		}
		batchSave[module.File][moduleName] = module
	}

	return mi.index.BatchSaveItems(batchSave)
}

// RemovedFiles handles file deletions
func (mi *ModuleIndexer) RemovedFiles(paths []string) error {
	return mi.index.BatchDeleteByFilePaths(paths)
}

// Close closes the indexer
func (mi *ModuleIndexer) Close() error {
	return mi.index.Close()
}

// Clear clears the index
func (mi *ModuleIndexer) Clear() error {
	return mi.index.Clear()
}

// GetModule retrieves a module by name
func (mi *ModuleIndexer) GetModule(name string) ([]Module, error) {
	return mi.index.GetValues(name)
}

// GetAllModules retrieves all indexed modules
func (mi *ModuleIndexer) GetAllModules() ([]Module, error) {
	return mi.index.GetAllValues()
}

// GetModuleNames retrieves all module names
func (mi *ModuleIndexer) GetModuleNames() ([]string, error) {
	return mi.index.GetAllKeys()
}
