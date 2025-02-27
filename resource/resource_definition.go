package resource

import (
	"errors"

	"github.com/tempestdx/sdk-go/jsonschema"
)

// resourceStatefulMetadata contains the published state of a resource
type resourceStatefulMetadata struct {
	version        string
	published      bool
	publishedID    string
	organizationID string
	recordID       string
	appID          string
}

type ResourceDefinition struct {
	// The display name of your resouce
	displayName string
	// The unique ID of your resource, if not specified, your display name will be slugified
	uniqueID string

	lifecycleStage       LifecycleStage
	categories           []Category
	links                []Link
	properties           *jsonschema.Schema
	instructionsMarkdown string

	operations  map[string]*operation
	healthCheck HealthCheckFunc

	publishedState *resourceStatefulMetadata
}

func (r *ResourceDefinition) UniqueID() string {
	return r.uniqueID
}

func (r *ResourceDefinition) DisplayName() string {
	return r.displayName
}

// Config represents the configuration for creating a new resource
type ResourceDefinitonConfig struct {
	Name           string
	DisplayName    string
	UniqueID       string
	LifecycleStage LifecycleStage
	Properties     *jsonschema.Schema
}

// New creates a new V2 resource
func New(config ResourceDefinitonConfig, opts ...resourceDefinitionOption) (*ResourceDefinition, error) {
	r := &ResourceDefinition{
		displayName:    config.DisplayName,
		uniqueID:       config.UniqueID,
		lifecycleStage: config.LifecycleStage,
		properties:     config.Properties,
		operations:     make(map[string]*operation),
	}

	if config.DisplayName == "" {
		return nil, errors.New("display name is required")
	}

	if config.UniqueID == "" {
		return nil, errors.New("unique id is required")
	}

	// validate links
	for _, link := range r.links {
		link.setDefault()
		if !link.isValid() {
			return nil, errors.New("invalid link")
		}
	}

	// validate categories
	for _, category := range r.categories {
		if !isValidCategory(category) {
			return nil, errors.New("invalid category")
		}
	}

	// validate properties
	if r.properties == nil {
		return nil, errors.New("properties are required")
	}

	// parse properties
	if _, err := jsonschema.ParseSchema(r.properties.Raw); err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt(r)
	}

	return r, nil
}

type resourceDefinitionOption func(*ResourceDefinition)

func WithDefaultLinks(links ...Link) resourceDefinitionOption {
	// Validate links
	return func(r *ResourceDefinition) {
		r.links = links
	}
}

func WithInstructions(markdown string) resourceDefinitionOption {
	return func(r *ResourceDefinition) {
		r.instructionsMarkdown = markdown
	}
}

func WithHealthCheck(fn HealthCheckFunc) resourceDefinitionOption {
	return func(r *ResourceDefinition) {
		r.healthCheck = fn
	}
}

func (r *ResourceDefinition) RegisterOperation(name string, fn OperationFunc, opts ...operationOption) *ResourceDefinition {
	op := newOperation(name, fn, opts...)
	r.operations[name] = op
	return r
}
