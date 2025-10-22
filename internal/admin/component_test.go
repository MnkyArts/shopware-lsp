package admin

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseComponentFile(t *testing.T) {
	t.Run("Parse Component.register", func(t *testing.T) {
		content, err := os.ReadFile("testdata/sw-test-button.js")
		require.NoError(t, err)

		root, tree := ParseJavaScript(content)
		require.NotNil(t, root)
		require.NotNil(t, tree)
		defer tree.Close()

		components, err := parseComponentFile(root, content, "testdata/sw-test-button.js")
		require.NoError(t, err)
		require.Len(t, components, 1)

		component, ok := components["sw-test-button"]
		require.True(t, ok)

		assert.Equal(t, "sw-test-button", component.Name)
		assert.Equal(t, "testdata/sw-test-button.js", component.File)
		assert.Greater(t, component.Line, 0)

		// Check props
		assert.Contains(t, component.Props, "variant")
		assert.Contains(t, component.Props, "size")
		assert.Contains(t, component.Props, "disabled")

		// Check computed properties
		assert.Contains(t, component.Computed, "buttonClasses")
		assert.Contains(t, component.Computed, "isDisabled")
	})
}

func TestParseJavaScript(t *testing.T) {
	t.Run("Valid JavaScript", func(t *testing.T) {
		content := []byte(`
			Component.register('test-component', {
				props: {
					title: String
				}
			});
		`)

		root, tree := ParseJavaScript(content)
		require.NotNil(t, root)
		require.NotNil(t, tree)
		defer tree.Close()
		assert.Equal(t, "program", root.Kind())
	})

	t.Run("Empty content", func(t *testing.T) {
		content := []byte("")
		root, tree := ParseJavaScript(content)
		require.NotNil(t, root)
		require.NotNil(t, tree)
		defer tree.Close()
	})
}

func TestExtractComponentFromCallExpression(t *testing.T) {
	t.Run("Component.extend", func(t *testing.T) {
		content := []byte(`
			Component.extend('sw-custom-button', 'sw-button', {
				props: {
					customProp: String
				}
			});
		`)

		root, tree := ParseJavaScript(content)
		require.NotNil(t, root)
		require.NotNil(t, tree)
		defer tree.Close()

		components, err := parseComponentFile(root, content, "test.js")
		require.NoError(t, err)
		require.Len(t, components, 1)

		component, ok := components["sw-custom-button"]
		require.True(t, ok)
		assert.Equal(t, "sw-custom-button", component.Name)
		assert.Contains(t, component.Props, "customProp")
	})
}
