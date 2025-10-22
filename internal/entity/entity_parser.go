package entity

import (
	"strings"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

// parseEntityDefinitions parses a PHP file for EntityDefinition classes
func parseEntityDefinitions(root *tree_sitter.Node, document []byte, filePath string) (map[string]EntityDefinition, error) {
	result := make(map[string]EntityDefinition)

	// Find all class declarations that extend EntityDefinition
	var visitor func(*tree_sitter.Node)
	visitor = func(node *tree_sitter.Node) {
		if node.Kind() == "class_declaration" {
			entity := extractEntityDefinition(node, document, filePath)
			if entity != nil {
				result[entity.EntityName] = *entity
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

// extractEntityDefinition extracts entity definition from a class declaration
func extractEntityDefinition(classNode *tree_sitter.Node, content []byte, filePath string) *EntityDefinition {
	// Get class name
	className := ""
	for i := 0; i < int(classNode.NamedChildCount()); i++ {
		child := classNode.NamedChild(uint(i))
		if child != nil && child.Kind() == "name" {
			className = string(child.Utf8Text(content))
			break
		}
	}

	if className == "" || !strings.HasSuffix(className, "Definition") {
		return nil
	}

	// Check if class extends EntityDefinition
	extendsEntityDefinition := false
	for i := 0; i < int(classNode.NamedChildCount()); i++ {
		child := classNode.NamedChild(uint(i))
		if child != nil && child.Kind() == "base_clause" {
			// Check if it extends EntityDefinition
			baseClass := extractBaseClassName(child, content)
			if baseClass == "EntityDefinition" {
				extendsEntityDefinition = true
				break
			}
		}
	}

	if !extendsEntityDefinition {
		return nil
	}

	entity := &EntityDefinition{
		ClassName: className,
		File:      filePath,
		Line:      int(classNode.Range().StartPoint.Row) + 1,
		Fields:    make(map[string]string),
	}

	// Find the getEntityName() method to get the entity name
	for i := 0; i < int(classNode.NamedChildCount()); i++ {
		child := classNode.NamedChild(uint(i))
		if child != nil && child.Kind() == "declaration_list" {
			// Look for getEntityName method
			entityName := extractEntityName(child, content)
			if entityName != "" {
				entity.EntityName = entityName
				entity.Name = entityName
			}
		}
	}

	// If no entity name found, derive from class name
	if entity.EntityName == "" {
		// Convert ProductDefinition -> product
		baseName := strings.TrimSuffix(className, "Definition")
		entity.EntityName = toSnakeCase(baseName)
		entity.Name = entity.EntityName
	}

	return entity
}

// extractBaseClassName extracts the base class name from a base_clause node
func extractBaseClassName(baseClause *tree_sitter.Node, content []byte) string {
	for i := 0; i < int(baseClause.NamedChildCount()); i++ {
		child := baseClause.NamedChild(uint(i))
		if child != nil && (child.Kind() == "name" || child.Kind() == "qualified_name") {
			text := string(child.Utf8Text(content))
			// Extract just the class name (e.g., "EntityDefinition" from "Shopware\Core\Framework\DataAbstractionLayer\EntityDefinition")
			parts := strings.Split(text, "\\")
			return parts[len(parts)-1]
		}
	}
	return ""
}

// extractEntityName extracts the entity name from the getEntityName() method
func extractEntityName(declarationList *tree_sitter.Node, content []byte) string {
	for i := 0; i < int(declarationList.NamedChildCount()); i++ {
		child := declarationList.NamedChild(uint(i))
		if child == nil || child.Kind() != "method_declaration" {
			continue
		}

		// Check if this is the getEntityName method
		methodName := ""
		for j := 0; j < int(child.NamedChildCount()); j++ {
			nameNode := child.NamedChild(uint(j))
			if nameNode != nil && nameNode.Kind() == "name" {
				methodName = string(nameNode.Utf8Text(content))
				break
			}
		}

		if methodName != "getEntityName" {
			continue
		}

		// Find the return statement
		for j := 0; j < int(child.NamedChildCount()); j++ {
			bodyNode := child.NamedChild(uint(j))
			if bodyNode == nil || bodyNode.Kind() != "compound_statement" {
				continue
			}

			return extractReturnValue(bodyNode, content)
		}
	}

	return ""
}

// extractReturnValue extracts the string value from a return statement
func extractReturnValue(bodyNode *tree_sitter.Node, content []byte) string {
	for i := 0; i < int(bodyNode.NamedChildCount()); i++ {
		child := bodyNode.NamedChild(uint(i))
		if child == nil || child.Kind() != "return_statement" {
			continue
		}

		// Look for string or class constant
		for j := 0; j < int(child.NamedChildCount()); j++ {
			returnExpr := child.NamedChild(uint(j))
			if returnExpr == nil {
				continue
			}

			// Handle direct string return
			if returnExpr.Kind() == "string" {
				text := string(returnExpr.Utf8Text(content))
				text = strings.Trim(text, `'"`)
				return text
			}

			// Handle self::ENTITY_NAME or static::ENTITY_NAME
			if returnExpr.Kind() == "class_constant_access_expression" {
				// For now, we'll skip this and rely on deriving from class name
				continue
			}
		}
	}

	return ""
}

// toSnakeCase converts PascalCase to snake_case
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
