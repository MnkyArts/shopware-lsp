package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Generator provides functionality for code generation
type Generator struct {
	projectRoot string
}

// NewGenerator creates a new generator
func NewGenerator(projectRoot string) *Generator {
	return &Generator{
		projectRoot: projectRoot,
	}
}

// GenerateAdminComponent generates a new admin component
func (g *Generator) GenerateAdminComponent(name, location string) (string, error) {
	// Ensure the location directory exists
	fullPath := filepath.Join(g.projectRoot, location)
	dir := filepath.Dir(fullPath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Create index.js file
	indexPath := filepath.Join(dir, "index.js")
	content := g.generateAdminComponentTemplate(name)

	if err := os.WriteFile(indexPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return indexPath, nil
}

// generateAdminComponentTemplate generates the template for an admin component
func (g *Generator) generateAdminComponentTemplate(name string) string {
	template := `import template from './` + name + `.html.twig';

const { Component } = Shopware;

Component.register('` + name + `', {
    template,

    props: {
        // Add your props here
    },

    data() {
        return {
            // Add your data properties here
        };
    },

    computed: {
        // Add your computed properties here
    },

    methods: {
        // Add your methods here
    }
});
`
	return template
}

// GenerateConfigXml generates a config.xml file
func (g *Generator) GenerateConfigXml(pluginName, location string) (string, error) {
	fullPath := filepath.Join(g.projectRoot, location)
	dir := filepath.Dir(fullPath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	configPath := filepath.Join(dir, "config.xml")
	content := g.generateConfigXmlTemplate(pluginName)

	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return configPath, nil
}

// generateConfigXmlTemplate generates the template for config.xml
func (g *Generator) generateConfigXmlTemplate(pluginName string) string {
	template := `<?xml version="1.0" encoding="UTF-8"?>
<config xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
        xsi:noNamespaceSchemaLocation="https://raw.githubusercontent.com/shopware/platform/trunk/src/Core/System/SystemConfig/Schema/config.xsd">

    <card>
        <title>Basic configuration</title>
        <title lang="de-DE">Grundeinstellungen</title>

        <input-field>
            <name>exampleSetting</name>
            <label>Example Setting</label>
            <label lang="de-DE">Beispiel Einstellung</label>
            <helpText>This is an example setting</helpText>
            <helpText lang="de-DE">Dies ist eine Beispiel Einstellung</helpText>
        </input-field>
    </card>
</config>
`
	return template
}

// ToPascalCase converts a kebab-case string to PascalCase
func ToPascalCase(s string) string {
	parts := strings.Split(s, "-")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, "")
}

// ToKebabCase converts a string to kebab-case
func ToKebabCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('-')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
