package app

import (
	"github.com/tempestdx/sdk-go/jsonschema"
)

// getResourceDefinition returns the ResourceDefinition with the given type, if it exists.
func (a *App) getResourceDefinition(t string) (*ResourceDefinition, bool) {
	for _, rd := range a.resourceDefinitions {
		if rd.Type == t {
			return &rd, true
		}
	}
	return nil, false
}

// ResourceDefinition represents a type of resource that can be managed by Tempest.
type ResourceDefinition struct {
	// The type by which this resource is identified. Must match `ResourceTypePattern`.
	Type string
	// The display name of the type for Tempest to show in the UI.
	DisplayName string
	// A description of the Resource type.
	Description string
	// PropertiesSchema is the parsed JSON schema for the Properties.
	PropertiesSchema *jsonschema.Schema
	// LifecycleStage represents how the Resource fits in the Developer Journey.
	LifecycleStage LifecycleStage
	// Links are links to documentation or other resources that can help users
	// understand how to use this Resource.
	Links []Link
	// Markdown formatted instructions for setting up or using the resource.
	// This field supports resource property variables in the format of {{ resource.<property name> }}.
	InstructionsMarkdown string

	// The CRUD operations that can be performed on this resource. These operations are optional.
	// These operations must be added by using the appropriate methods on the ResourceDefinition.
	create *operation
	read   *operation
	update *operation
	delete *operation
	// List is an optional operation that can be performed on this resource.
	// This operation must be added by using the appropriate method on the ResourceDefinition.
	list *listOperation

	// healthcheck is an optional function that returns the provisioning health of the resource.
	healthcheck HealthCheckFunc

	// Actions are additional actions related to the resource that can be performed.
	// A good example of an action might be "Trigger a Build" on a CI/CD resource.
	// Actions must be added by using the AddAction method on the ResourceDefinition.
	actions []ActionDefinition
}

// CreateFn adds a Create operation Handler to the ResourceDefinition.
//
// Tempest will use the inputSchema to validate the input data before calling the Handler.
// The Handler should create the resource in the external system and return the resource's ExternalID and properties.
// See the Create operation in the Printer example for an example implementation.
func (rd *ResourceDefinition) CreateFn(fn OperationFunc, inputSchema *jsonschema.Schema) {
	if inputSchema == nil {
		panic("input schema is nil")
	}

	if rd.PropertiesSchema == nil {
		panic("Properties must be set before adding a Create Operation")
	}

	if fn == nil {
		panic("OperationFunc must be set for a Create Operation")
	}

	rd.create = &operation{
		schema: schema{
			input:  inputSchema,
			output: rd.PropertiesSchema,
		},
		fn: fn,
	}
}

// UpdateFn adds an Update operation Handler to the ResourceDefinition.
//
// Tempest will use the inputSchema to validate the input data before calling the Handler.
// The Handler should update the resource in the external system and return the updated resource's properties.
// See the Update operation in the Printer example for an example implementation.
func (rd *ResourceDefinition) UpdateFn(fn OperationFunc, inputSchema *jsonschema.Schema) {
	if inputSchema == nil {
		panic("input schema is nil")
	}

	if rd.PropertiesSchema == nil {
		panic("properties schema must be set before adding an Update Operation")
	}

	if fn == nil {
		panic("OperationFunc must be set for an Update Operation")
	}

	rd.update = &operation{
		schema: schema{
			input:  inputSchema,
			output: rd.PropertiesSchema,
		},
		fn: fn,
	}
}

// DeleteFn adds a Delete operation Handler to the ResourceDefinition.
//
// The Handler should delete the resource in the external system.
// See the Delete operation in the Printer example for an example implementation.
func (rd *ResourceDefinition) DeleteFn(fn OperationFunc) {
	if rd.PropertiesSchema == nil {
		panic("Properties must be set before adding a Delete handler")
	}

	if fn == nil {
		panic("OperationFunc must be set for a Delete Operation")
	}

	rd.delete = &operation{
		schema: schema{
			input:  jsonschema.MustParseSchema(jsonschema.GenericEmptySchema),
			output: rd.PropertiesSchema,
		},
		fn: fn,
	}
}

// ReadFn adds a Read Operation to the ResourceDefinition.
//
// The Handler should query the external system for the resource's current state and return it.
// See the Read operation in the Printer example for an example implementation.
func (rd *ResourceDefinition) ReadFn(fn OperationFunc) {
	if rd.PropertiesSchema == nil {
		panic("Properties must be set before adding a Read handler")
	}

	if fn == nil {
		panic("OperationFunc must be set for a Read Operation")
	}

	rd.read = &operation{
		schema: schema{
			input:  jsonschema.MustParseSchema(jsonschema.GenericEmptySchema),
			output: rd.PropertiesSchema,
		},
		fn: fn,
	}
}

// ListFn adds a List Handler to the ResourceDefinition.
// The output will be validated against the ResourceDefinition's Properties schema.
//
// The Handler should query the external system for all resources of this type and return them.
// See the List operation in the Printer example for an example implementation.
func (rd *ResourceDefinition) ListFn(fn ListFunc) {
	if rd.PropertiesSchema == nil {
		panic("Properties must be set before adding a List handler")
	}

	if fn == nil {
		panic("ListFunc must be set for a List Operation")
	}

	rd.list = &listOperation{
		schema: schema{
			input:  jsonschema.MustParseSchema(jsonschema.GenericEmptySchema),
			output: rd.PropertiesSchema,
		},
		fn: fn,
	}
}

// HealthCheckFn adds a HealthCheck Handler to the ResourceDefinition.
//
// The Handler should return the provisioning health of the resource.
func (rd *ResourceDefinition) HealthCheckFn(fn HealthCheckFunc) {
	if fn == nil {
		panic("HealthCheckFunc must be set for a HealthCheck Operation")
	}

	rd.healthcheck = fn
}
