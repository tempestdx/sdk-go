package app

import (
	"fmt"
	"maps"

	"github.com/tempestdx/sdk-go/resource"
)

type App struct {
	name      string
	resources map[string]*resource.ResourceDefinition
}

type OptFunc func(*App) error

type Config struct {
	Name string
}

func New(config Config, opts ...OptFunc) (*App, error) {
	a := &App{
		name:      config.Name,
		resources: make(map[string]*resource.ResourceDefinition),
	}

	for _, opt := range opts {
		if err := opt(a); err != nil {
			return nil, err
		}
	}

	if a.name == "" {
		return nil, fmt.Errorf("app name is required")
	}

	return a, nil
}

func WithResource(r *resource.ResourceDefinition) OptFunc {
	return func(a *App) error {
		// Check if the resource name is already in use
		if _, ok := a.resources[r.UniqueID()]; ok {
			return fmt.Errorf("resource %s already exists", r.UniqueID())
		}
		a.resources[r.UniqueID()] = r

		// Validate that the resource's display name is unique
		for _, existingResource := range a.resources {
			if existingResource.DisplayName() == r.DisplayName() {
				return fmt.Errorf("resource display name %s already exists", r.DisplayName())
			}
		}

		return nil
	}
}

func (a *App) Resources() map[string]*resource.ResourceDefinition {
	// Create a copy of the resources map to prevent external modification
	resourcesCopy := make(map[string]*resource.ResourceDefinition, len(a.resources))
	maps.Copy(resourcesCopy, a.resources)
	return resourcesCopy
}
