package cmd

import (
	"fmt"

	"conhos-cli/pkg/connectors"
)

// Options represents configuration options for registry operations
type Options struct {
	List  bool
	Build bool
	Name  string
}

// Registry represents a registry management system
type Registry struct {
	options *Options
	ws      *connectors.WS
}

// NewRegistry creates a new Registry instance
func NewRegistry(options *Options) *Registry {
	return &Registry{
		options: options,
		ws:      connectors.NewWS(&connectors.Options{}),
	}
}

// Handle checks options and executes the appropriate registry command
func (r *Registry) Handle() error {
	// Validate options similar to JavaScript checkOptions
	if !r.options.List && !r.options.Build {
		return fmt.Errorf("one of options are required: --build | --list")
	}
	if r.options.List && r.options.Build {
		return fmt.Errorf("two options are not allowed to use together: --build & --list")
	}
	if r.options.Build && r.options.Name == "" {
		return fmt.Errorf("option 'name' is required for build: -n|--name [string]")
	}

	// Connect to WebSocket server
	if err := r.ws.Connect(); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer r.ws.Close()

	// Execute the requested operation
	if r.options.List {
		return r.handleList()
	}
	if r.options.Build {
		return r.handleBuild()
	}
	return nil
}

// handleList handles the registry listing operation
func (r *Registry) handleList() error {
	// Implement listing logic
	fmt.Println("Listing registries...")
	return nil
}

// handleBuild handles the registry build operation
func (r *Registry) handleBuild() error {
	// Implement build logic
	fmt.Printf("Building registry '%s'...\n", r.options.Name)
	return nil
}
