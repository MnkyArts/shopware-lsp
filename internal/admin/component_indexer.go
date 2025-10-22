package admin

import (
	"path/filepath"
	"strings"

	"github.com/shopware/shopware-lsp/internal/indexer"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

// ComponentIndexer indexes Vue admin components
type ComponentIndexer struct {
	index *indexer.DataIndexer[Component]
}

// NewComponentIndexer creates a new component indexer
func NewComponentIndexer(configDir string) (*ComponentIndexer, error) {
	idx, err := indexer.NewDataIndexer[Component](filepath.Join(configDir, "admin_component.db"))
	if err != nil {
		return nil, err
	}

	return &ComponentIndexer{
		index: idx,
	}, nil
}

// ID returns the indexer identifier
func (ci *ComponentIndexer) ID() string {
	return "admin.component.indexer"
}

// Index processes a JavaScript file and indexes admin components
func (ci *ComponentIndexer) Index(path string, node *tree_sitter.Node, fileContent []byte) error {
	// Only index files in admin component directories
	// Typical paths: src/Administration/Resources/app/administration/src/
	if !strings.Contains(path, "/administration/src/") || 
	   !strings.HasSuffix(path, ".js") ||
	   strings.Contains(path, "/_fixtures/") ||
	   strings.Contains(path, "/node_modules/") {
		return nil
	}

	components, err := parseComponentFile(node, fileContent, path)
	if err != nil {
		return err
	}

	// Batch save all components found in this file
	batchSave := make(map[string]map[string]Component)
	for componentName, component := range components {
		if _, ok := batchSave[component.File]; !ok {
			batchSave[component.File] = make(map[string]Component)
		}
		batchSave[component.File][componentName] = component
	}

	return ci.index.BatchSaveItems(batchSave)
}

// RemovedFiles handles file deletions
func (ci *ComponentIndexer) RemovedFiles(paths []string) error {
	return ci.index.BatchDeleteByFilePaths(paths)
}

// Close closes the indexer
func (ci *ComponentIndexer) Close() error {
	return ci.index.Close()
}

// Clear clears the index
func (ci *ComponentIndexer) Clear() error {
	return ci.index.Clear()
}

// GetComponent retrieves a component by name
func (ci *ComponentIndexer) GetComponent(name string) ([]Component, error) {
	return ci.index.GetValues(name)
}

// GetAllComponents retrieves all indexed components
func (ci *ComponentIndexer) GetAllComponents() ([]Component, error) {
	return ci.index.GetAllValues()
}

// GetComponentNames retrieves all component names
func (ci *ComponentIndexer) GetComponentNames() ([]string, error) {
	return ci.index.GetAllKeys()
}
