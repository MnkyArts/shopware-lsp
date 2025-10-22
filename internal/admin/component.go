package admin

import (
	"strings"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_javascript "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
)

// Component represents a Vue admin component
type Component struct {
	Name     string            // Component name (e.g., "sw-button")
	File     string            // File path
	Line     int               // Line number
	Props    map[string]string // Component props
	Computed []string          // Computed properties
}

// parseComponentFile parses a JavaScript/Vue file for Component.register() calls
func parseComponentFile(root *tree_sitter.Node, document []byte, filePath string) (map[string]Component, error) {
	result := make(map[string]Component)

	// We need to find call expressions that match:
	// Component.register('component-name', { ... })
	// or Component.extend('component-name', 'parent', { ... })

	var visitor func(*tree_sitter.Node)
	visitor = func(node *tree_sitter.Node) {
		// Check if this is a call expression
		if node.Kind() == "call_expression" {
			component := extractComponentFromCallExpression(node, document, filePath)
			if component != nil {
				result[component.Name] = *component
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

// extractComponentFromCallExpression extracts component data from a Component.register() call
func extractComponentFromCallExpression(node *tree_sitter.Node, content []byte, filePath string) *Component {
	// Get the function being called (should be Component.register or Component.extend)
	if node.NamedChildCount() < 2 {
		return nil
	}

	function := node.NamedChild(0)
	if function.Kind() != "member_expression" {
		return nil
	}

	// Check if the object is "Component"
	if function.NamedChildCount() < 2 {
		return nil
	}

	objectNode := function.NamedChild(0)
	propertyNode := function.NamedChild(1)

	objectText := string(objectNode.Utf8Text(content))
	propertyText := string(propertyNode.Utf8Text(content))

	// We only care about Component.register or Component.extend
	if objectText != "Component" || (propertyText != "register" && propertyText != "extend") {
		return nil
	}

	// Get the arguments
	arguments := node.NamedChild(1)
	if arguments == nil || arguments.Kind() != "arguments" || arguments.NamedChildCount() < 1 {
		return nil
	}

	// First argument should be the component name (string)
	nameArg := arguments.NamedChild(0)
	if nameArg == nil || nameArg.Kind() != "string" {
		return nil
	}

	componentName := string(nameArg.Utf8Text(content))
	componentName = strings.Trim(componentName, `'"`)

	component := &Component{
		Name:     componentName,
		File:     filePath,
		Line:     int(node.Range().StartPoint.Row) + 1,
		Props:    make(map[string]string),
		Computed: make([]string, 0),
	}

	// Find the configuration object (last argument for both register and extend)
	var configObject *tree_sitter.Node
	lastArgIndex := arguments.NamedChildCount() - 1
	configObject = arguments.NamedChild(uint(lastArgIndex))

	if configObject != nil && configObject.Kind() == "object" {
		extractComponentConfig(configObject, content, component)
	}

	return component
}

// extractComponentConfig extracts props and computed properties from the config object
func extractComponentConfig(configNode *tree_sitter.Node, content []byte, component *Component) {
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

		// Extract props
		if keyText == "props" && valueNode != nil && valueNode.Kind() == "object" {
			extractProps(valueNode, content, component)
		}

		// Extract computed properties
		if keyText == "computed" && valueNode != nil && valueNode.Kind() == "object" {
			extractComputed(valueNode, content, component)
		}
	}
}

// extractProps extracts prop definitions from the props object
func extractProps(propsNode *tree_sitter.Node, content []byte, component *Component) {
	for i := 0; i < int(propsNode.NamedChildCount()); i++ {
		pair := propsNode.NamedChild(uint(i))
		if pair == nil || pair.Kind() != "pair" {
			continue
		}

		if pair.NamedChildCount() < 1 {
			continue
		}

		keyNode := pair.NamedChild(0)
		if keyNode == nil {
			continue
		}

		propName := string(keyNode.Utf8Text(content))
		propName = strings.Trim(propName, `'"`)

		// Store prop with empty type for now (could be enhanced to extract type info)
		component.Props[propName] = ""
	}
}

// extractComputed extracts computed property names
func extractComputed(computedNode *tree_sitter.Node, content []byte, component *Component) {
	for i := 0; i < int(computedNode.NamedChildCount()); i++ {
		pair := computedNode.NamedChild(uint(i))
		if pair == nil || (pair.Kind() != "pair" && pair.Kind() != "method_definition") {
			continue
		}

		if pair.NamedChildCount() < 1 {
			continue
		}

		keyNode := pair.NamedChild(0)
		if keyNode == nil {
			continue
		}

		computedName := string(keyNode.Utf8Text(content))
		computedName = strings.Trim(computedName, `'"`)

		component.Computed = append(component.Computed, computedName)
	}
}

// ParseJavaScript is a helper to parse JavaScript content with tree-sitter
// Note: The returned node is valid only as long as the tree exists
// Caller should handle tree lifecycle appropriately
func ParseJavaScript(content []byte) (*tree_sitter.Node, *tree_sitter.Tree) {
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
