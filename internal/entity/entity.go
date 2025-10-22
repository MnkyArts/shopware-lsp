package entity

// EntityDefinition represents a Shopware entity definition
type EntityDefinition struct {
	Name       string            // Entity name (e.g., "product", "customer")
	ClassName  string            // PHP class name (e.g., "ProductDefinition")
	File       string            // File path
	Line       int               // Line number
	Fields     map[string]string // Field name -> field type mapping
	EntityName string            // Value returned by getEntityName()
}
