package entity

import (
	"os"
	"testing"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_php "github.com/tree-sitter/tree-sitter-php/bindings/go"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func parsePHP(content []byte) (*tree_sitter.Node, *tree_sitter.Tree) {
	parser := tree_sitter.NewParser()
	defer parser.Close()

	lang := tree_sitter.NewLanguage(tree_sitter_php.LanguagePHP())
	parser.SetLanguage(lang)

	tree := parser.Parse(content, nil)
	if tree == nil {
		return nil, nil
	}

	return tree.RootNode(), tree
}

func TestParseEntityDefinitions(t *testing.T) {
	t.Run("Parse EntityDefinition", func(t *testing.T) {
		content, err := os.ReadFile("testdata/CustomProductDefinition.php")
		require.NoError(t, err)

		root, tree := parsePHP(content)
		require.NotNil(t, root)
		require.NotNil(t, tree)
		defer tree.Close()

		entities, err := parseEntityDefinitions(root, content, "testdata/CustomProductDefinition.php")
		require.NoError(t, err)
		require.Len(t, entities, 1)

		entity, ok := entities["custom_product"]
		require.True(t, ok)

		assert.Equal(t, "custom_product", entity.EntityName)
		assert.Equal(t, "CustomProductDefinition", entity.ClassName)
		assert.Equal(t, "testdata/CustomProductDefinition.php", entity.File)
		assert.Greater(t, entity.Line, 0)
	})
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Product", "product"},
		{"CustomProduct", "custom_product"},
		{"ProductMedia", "product_media"},
		{"CustomerAddress", "customer_address"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := toSnakeCase(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
