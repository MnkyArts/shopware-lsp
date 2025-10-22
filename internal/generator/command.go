package generator

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/shopware/shopware-lsp/internal/lsp"
	"github.com/shopware/shopware-lsp/internal/lsp/protocol"
)

// CommandProvider provides LSP commands for code generation
type CommandProvider struct {
	generator *Generator
}

// NewCommandProvider creates a new generator command provider
func NewCommandProvider(projectRoot string) *CommandProvider {
	return &CommandProvider{
		generator: NewGenerator(projectRoot),
	}
}

// GetCommands returns the available generator commands
func (p *CommandProvider) GetCommands(ctx context.Context) map[string]lsp.CommandFunc {
	return map[string]lsp.CommandFunc{
		"shopware/generate/adminComponent": p.generateAdminComponent,
		"shopware/generate/configXml":      p.generateConfigXml,
	}
}

// generateAdminComponent handles the admin component generation command
func (p *CommandProvider) generateAdminComponent(ctx context.Context, args *json.RawMessage) (interface{}, error) {
	var params struct {
		Name     string `json:"name"`
		Location string `json:"location"`
	}

	if err := json.Unmarshal(*args, &params); err != nil {
		return nil, fmt.Errorf("failed to unmarshal arguments: %w", err)
	}

	if params.Name == "" {
		return protocol.NewLspError("Component name is required", "name.required"), nil
	}

	if params.Location == "" {
		// Default location
		params.Location = fmt.Sprintf("src/Administration/Resources/app/administration/src/component/%s/index.js", params.Name)
	}

	path, err := p.generator.GenerateAdminComponent(params.Name, params.Location)
	if err != nil {
		return protocol.NewLspError(err.Error(), "generation.failed"), nil
	}

	return map[string]any{
		"uri":  "file://" + path,
		"line": 0,
	}, nil
}

// generateConfigXml handles the config.xml generation command
func (p *CommandProvider) generateConfigXml(ctx context.Context, args *json.RawMessage) (interface{}, error) {
	var params struct {
		PluginName string `json:"pluginName"`
		Location   string `json:"location"`
	}

	if err := json.Unmarshal(*args, &params); err != nil {
		return nil, fmt.Errorf("failed to unmarshal arguments: %w", err)
	}

	if params.PluginName == "" {
		return protocol.NewLspError("Plugin name is required", "pluginName.required"), nil
	}

	if params.Location == "" {
		// Default location
		params.Location = "src/Resources/config/config.xml"
	}

	path, err := p.generator.GenerateConfigXml(params.PluginName, params.Location)
	if err != nil {
		return protocol.NewLspError(err.Error(), "generation.failed"), nil
	}

	return map[string]any{
		"uri":  "file://" + path,
		"line": 0,
	}, nil
}
