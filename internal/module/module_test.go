package module

import (
	"os"
	"testing"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_javascript "github.com/tree-sitter/tree-sitter-javascript/bindings/go"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func parseJavaScript(content []byte) (*tree_sitter.Node, *tree_sitter.Tree) {
	parser := tree_sitter.NewParser()
	defer parser.Close()

	lang := tree_sitter.NewLanguage(tree_sitter_javascript.Language())
	parser.SetLanguage(lang)

	tree := parser.Parse(content, nil)
	if tree == nil {
		return nil, nil
	}

	return tree.RootNode(), tree
}

func TestParseModuleFile(t *testing.T) {
	t.Run("Parse Module.register", func(t *testing.T) {
		content, err := os.ReadFile("testdata/sw-product-module.js")
		require.NoError(t, err)

		root, tree := parseJavaScript(content)
		require.NotNil(t, root)
		require.NotNil(t, tree)
		defer tree.Close()

		modules, err := parseModuleFile(root, content, "testdata/sw-product-module.js")
		require.NoError(t, err)
		require.Len(t, modules, 1)

		module, ok := modules["sw-product"]
		require.True(t, ok)

		assert.Equal(t, "sw-product", module.Name)
		assert.Equal(t, "testdata/sw-product-module.js", module.File)
		assert.Greater(t, module.Line, 0)

		// Check labels
		assert.Contains(t, module.Labels, "title")
		assert.Contains(t, module.Labels, "description")
		assert.Equal(t, "sw-product.general.mainMenuItemGeneral", module.Title)
		assert.Equal(t, "sw-product.general.mainMenuItemGeneral", module.Labels["title"])
	})
}
