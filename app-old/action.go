package app

import (
	"context"
	"fmt"

	appv1 "github.com/tempestdx/protobuf/gen/go/tempestdx/app/v1"
)

func (a *App) getActionDefinition(resource, action string) (*ActionDefinition, bool) {
	rd, ok := a.getResourceDefinition(resource)
	if !ok {
		return nil, false
	}

	for _, a := range rd.actions {
		if a.Name == action {
			return &a, true
		}
	}
	return nil, false
}

func (rd *ResourceDefinition) AddActionDefinition(ad ActionDefinition) {
	for _, existing := range rd.actions {
		if existing.Name == ad.Name {
			panic(fmt.Sprintf("ActionDefinition with the same name '%s' already exists", ad.Name))
		}
	}

	if ad.InputSchema == nil {
		ad.InputSchema = MustParseJSONSchema(GenericEmptySchema)
	}

	if ad.OutputSchema == nil {
		ad.OutputSchema = MustParseJSONSchema(GenericEmptySchema)
	}

	rd.actions = append(rd.actions, ad)
}

type ActionDefinition struct {
	// Name is the unique identifier for the action.
	Name string
	// DisplayName is the name of the action as it should be displayed in the UI.
	DisplayName string
	// Description is a short description of the action.
	Description string
	// InputSchema is the parsed JSON schema for the input schema.
	InputSchema *JSONSchema
	// OutputSchema is the parsed JSON schema for the output schema.
	OutputSchema *JSONSchema
	// Handler is the function that will be called when the action is invoked.
	Handler func(context.Context, *ActionRequest) (*ActionResponse, error)
}

// ActionRequest contains the input data for an operation on a resource.
type ActionRequest struct {
	// Metadata contains information about the Project and User making the request.
	// This metadata does not contain information about the Resource being operated on.
	Metadata *Metadata
	// Resource is the resource being actioned, and contains the ExternalID of the resource,
	// as well as the properties at the time of the request.
	Resource *Resource
	// Action is the name of the action being performed.
	Action string
	// Input contains the input data for the request, after it has been validated against the schema.
	// Default values have already been applied to missing input properties.
	Input map[string]any
	// Environment contains the environment variables that are available to the operation.
	Environment map[string]EnvironmentVariable
}

func actionRequestFromProto(r *appv1.ExecuteResourceActionRequest) *ActionRequest {
	environment := make(map[string]EnvironmentVariable, len(r.EnvironmentVariables))
	for _, v := range r.EnvironmentVariables {
		var t EnvironmentVariableType
		switch v.Type {
		case appv1.EnvironmentVariableType_ENVIRONMENT_VARIABLE_TYPE_VAR:
			t = ENVIRONMENT_VARIABLE_TYPE_VAR
		case appv1.EnvironmentVariableType_ENVIRONMENT_VARIABLE_TYPE_SECRET:
			t = ENVIRONMENT_VARIABLE_TYPE_SECRET
		case appv1.EnvironmentVariableType_ENVIRONMENT_VARIABLE_TYPE_CERTIFICATE:
			t = ENVIRONMENT_VARIABLE_TYPE_CERTIFICATE
		case appv1.EnvironmentVariableType_ENVIRONMENT_VARIABLE_TYPE_PRIVATE_KEY:
			t = ENVIRONMENT_VARIABLE_TYPE_PRIVATE_KEY
		case appv1.EnvironmentVariableType_ENVIRONMENT_VARIABLE_TYPE_PUBLIC_KEY:
			t = ENVIRONMENT_VARIABLE_TYPE_PUBLIC_KEY
		}

		environment[v.Key] = EnvironmentVariable{
			Key:   v.Key,
			Value: v.Value,
			Type:  t,
		}
	}

	return &ActionRequest{
		Metadata:    metadataFromProto(r.Metadata),
		Resource:    resourceFromProto(r.Resource),
		Action:      r.Action,
		Input:       r.Input.AsMap(),
		Environment: environment,
	}
}

type ActionResponse struct {
	// Output contains the output data for the request. This data must validate against the output schema.
	Output map[string]any
}
