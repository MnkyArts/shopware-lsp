package generator

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAdminComponent(t *testing.T) {
	tmpDir := t.TempDir()
	gen := NewGenerator(tmpDir)

	componentName := "sw-test-component"
	location := "src/Administration/Resources/app/administration/src/component/sw-test-component/index.js"

	path, err := gen.GenerateAdminComponent(componentName, location)
	require.NoError(t, err)
	require.NotEmpty(t, path)

	// Check if file was created
	_, err = os.Stat(path)
	require.NoError(t, err)

	// Check file content
	content, err := os.ReadFile(path)
	require.NoError(t, err)

	contentStr := string(content)
	assert.Contains(t, contentStr, "Component.register('sw-test-component'")
	assert.Contains(t, contentStr, "template")
	assert.Contains(t, contentStr, "props:")
	assert.Contains(t, contentStr, "computed:")
	assert.Contains(t, contentStr, "methods:")
}

func TestGenerateConfigXml(t *testing.T) {
	tmpDir := t.TempDir()
	gen := NewGenerator(tmpDir)

	pluginName := "TestPlugin"
	location := "src/Resources/config/config.xml"

	path, err := gen.GenerateConfigXml(pluginName, location)
	require.NoError(t, err)
	require.NotEmpty(t, path)

	// Check if file was created
	_, err = os.Stat(path)
	require.NoError(t, err)

	// Check file content
	content, err := os.ReadFile(path)
	require.NoError(t, err)

	contentStr := string(content)
	assert.Contains(t, contentStr, "<?xml version")
	assert.Contains(t, contentStr, "<config")
	assert.Contains(t, contentStr, "<card>")
	assert.Contains(t, contentStr, "<input-field>")
}

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"sw-custom-button", "SwCustomButton"},
		{"test-component", "TestComponent"},
		{"single", "Single"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ToPascalCase(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToKebabCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"SwCustomButton", "sw-custom-button"},
		{"TestComponent", "test-component"},
		{"Single", "single"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ToKebabCase(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
