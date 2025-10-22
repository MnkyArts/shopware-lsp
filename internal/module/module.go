package module

// Module represents a Shopware admin module registration
type Module struct {
	Name       string            // Module name (e.g., "sw-product")
	File       string            // File path
	Line       int               // Line number
	Title      string            // Module title (snippet key)
	Labels     map[string]string // Label keys (e.g., "title", "description")
	Navigation []string          // Navigation entries
}
