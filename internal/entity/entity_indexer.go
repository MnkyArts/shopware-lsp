package entity

import (
	"path/filepath"
	"strings"

	"github.com/shopware/shopware-lsp/internal/indexer"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

// EntityIndexer indexes Shopware entity definitions
type EntityIndexer struct {
	index *indexer.DataIndexer[EntityDefinition]
}

// NewEntityIndexer creates a new entity indexer
func NewEntityIndexer(configDir string) (*EntityIndexer, error) {
	idx, err := indexer.NewDataIndexer[EntityDefinition](filepath.Join(configDir, "entity.db"))
	if err != nil {
		return nil, err
	}

	return &EntityIndexer{
		index: idx,
	}, nil
}

// ID returns the indexer identifier
func (ei *EntityIndexer) ID() string {
	return "entity.indexer"
}

// Index processes a PHP file and indexes entity definitions
func (ei *EntityIndexer) Index(path string, node *tree_sitter.Node, fileContent []byte) error {
	// Only index PHP files that end with "Definition.php"
	if !strings.HasSuffix(path, "Definition.php") ||
		strings.Contains(path, "/vendor/") ||
		strings.Contains(path, "/_fixtures/") {
		return nil
	}

	entities, err := parseEntityDefinitions(node, fileContent, path)
	if err != nil {
		return err
	}

	// Batch save all entities found in this file
	batchSave := make(map[string]map[string]EntityDefinition)
	for entityName, entity := range entities {
		if _, ok := batchSave[entity.File]; !ok {
			batchSave[entity.File] = make(map[string]EntityDefinition)
		}
		batchSave[entity.File][entityName] = entity
	}

	return ei.index.BatchSaveItems(batchSave)
}

// RemovedFiles handles file deletions
func (ei *EntityIndexer) RemovedFiles(paths []string) error {
	return ei.index.BatchDeleteByFilePaths(paths)
}

// Close closes the indexer
func (ei *EntityIndexer) Close() error {
	return ei.index.Close()
}

// Clear clears the index
func (ei *EntityIndexer) Clear() error {
	return ei.index.Clear()
}

// GetEntity retrieves an entity by name
func (ei *EntityIndexer) GetEntity(name string) ([]EntityDefinition, error) {
	return ei.index.GetValues(name)
}

// GetAllEntities retrieves all indexed entities
func (ei *EntityIndexer) GetAllEntities() ([]EntityDefinition, error) {
	return ei.index.GetAllValues()
}

// GetEntityNames retrieves all entity names
func (ei *EntityIndexer) GetEntityNames() ([]string, error) {
	return ei.index.GetAllKeys()
}
