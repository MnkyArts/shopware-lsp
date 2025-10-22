package module

import (
	"strings"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

// parseModuleFile parses a JavaScript file for Module.register() calls
func parseModuleFile(root *tree_sitter.Node, document []byte, filePath string) (map[string]Module, error) {
	result := make(map[string]Module)

	// Find Module.register() calls
	var visitor func(*tree_sitter.Node)
	visitor = func(node *tree_sitter.Node) {
		if node.Kind() == "call_expression" {
			mod := extractModuleFromCallExpression(node, document, filePath)
			if mod != nil {
				result[mod.Name] = *mod
			}
		}

		// Recursively visit children
		for i := 0; i < int(node.NamedChildCount()); i++ {
			child := node.NamedChild(uint(i))
			if child != nil {
				visitor(child)
			}
		}
	}

	visitor(root)
	return result, nil
}

// extractModuleFromCallExpression extracts module data from Module.register() call
func extractModuleFromCallExpression(node *tree_sitter.Node, content []byte, filePath string) *Module {
	if node.NamedChildCount() < 2 {
		return nil
	}

	function := node.NamedChild(0)
	if function.Kind() != "member_expression" {
		return nil
	}

	// Check if it's Module.register
	if function.NamedChildCount() < 2 {
		return nil
	}

	objectNode := function.NamedChild(0)
	propertyNode := function.NamedChild(1)

	objectText := string(objectNode.Utf8Text(content))
	propertyText := string(propertyNode.Utf8Text(content))

	if objectText != "Module" || propertyText != "register" {
		return nil
	}

	// Get arguments
	arguments := node.NamedChild(1)
	if arguments == nil || arguments.Kind() != "arguments" || arguments.NamedChildCount() < 2 {
		return nil
	}

	// First argument is the module name
	nameArg := arguments.NamedChild(0)
	if nameArg == nil || nameArg.Kind() != "string" {
		return nil
	}

	moduleName := string(nameArg.Utf8Text(content))
	moduleName = strings.Trim(moduleName, `'"`)

	mod := &Module{
		Name:   moduleName,
		File:   filePath,
		Line:   int(node.Range().StartPoint.Row) + 1,
		Labels: make(map[string]string),
	}

	// Second argument is the config object
	configObject := arguments.NamedChild(1)
	if configObject != nil && configObject.Kind() == "object" {
		extractModuleConfig(configObject, content, mod)
	}

	return mod
}

// extractModuleConfig extracts labels and other config from the module config object
func extractModuleConfig(configNode *tree_sitter.Node, content []byte, mod *Module) {
	for i := 0; i < int(configNode.NamedChildCount()); i++ {
		pair := configNode.NamedChild(uint(i))
		if pair == nil || pair.Kind() != "pair" {
			continue
		}

		if pair.NamedChildCount() < 2 {
			continue
		}

		keyNode := pair.NamedChild(0)
		valueNode := pair.NamedChild(1)

		if keyNode == nil {
			continue
		}

		keyText := string(keyNode.Utf8Text(content))
		keyText = strings.Trim(keyText, `'"`)

		// Extract title
		if keyText == "title" && valueNode != nil && valueNode.Kind() == "string" {
			titleText := string(valueNode.Utf8Text(content))
			titleText = strings.Trim(titleText, `'"`)
			mod.Title = titleText
			mod.Labels["title"] = titleText
		}

		// Extract other snippet keys in the config
		if valueNode != nil && valueNode.Kind() == "string" {
			valueText := string(valueNode.Utf8Text(content))
			valueText = strings.Trim(valueText, `'"`)
			// Store keys that look like snippet keys (contain dots)
			if strings.Contains(valueText, ".") {
				mod.Labels[keyText] = valueText
			}
		}
	}
}
