package app

import (
	"context"

	appv1 "github.com/tempestdx/protobuf/gen/go/tempestdx/app/v1"
)

type EnvironmentVariableType string

const (
	ENVIRONMENT_VARIABLE_TYPE_VAR         EnvironmentVariableType = "variable"
	ENVIRONMENT_VARIABLE_TYPE_SECRET      EnvironmentVariableType = "secret"
	ENVIRONMENT_VARIABLE_TYPE_CERTIFICATE EnvironmentVariableType = "certificate"
	ENVIRONMENT_VARIABLE_TYPE_PRIVATE_KEY EnvironmentVariableType = "private_key"
	ENVIRONMENT_VARIABLE_TYPE_PUBLIC_KEY  EnvironmentVariableType = "public_key"
)

type EnvironmentVariable struct {
	Key   string
	Value string
	Type  EnvironmentVariableType
}

// OperationRequest contains the input data for an operation on a resource.
type OperationRequest struct {
	// Metadata contains information about the Project and User making the request.
	// This metadata does not contain information about the Resource being operated on.
	Metadata *Metadata
	// Resource is the resource being operated on, and contains the ExternalID of the resource,
	// as well as the properties at the time of the request.
	Resource *Resource
	// Input contains the input data for the request, after it has been validated against the schema.
	// Default values have already been applied to missing input properties.
	Input map[string]any
	// Environment contains the environment variables that are available to the operation.
	Environment map[string]EnvironmentVariable
}

func operationRequestFromProto(r *appv1.ExecuteResourceOperationRequest) *OperationRequest {
	if r == nil {
		return nil
	}

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

	return &OperationRequest{
		Metadata:    metadataFromProto(r.Metadata),
		Resource:    resourceFromProto(r.Resource),
		Input:       r.Input.AsMap(),
		Environment: environment,
	}
}

// OperationResponse contains the output data for an operation on a resource.
type OperationResponse struct {
	// Resource contains the properties of the resource after the operation has been performed.
	Resource *Resource
}

// operation is a struct that contains the schema and function for an operation.
// This must be constructed using the appropriate methods on the ResourceDefinition.
type operation struct {
	schema schema
	fn     OperationFunc
}

// schema contains the input and output JSON schemas for an operation.
type schema struct {
	input  *JSONSchema
	output *JSONSchema
}

type OperationFunc func(context.Context, *OperationRequest) (*OperationResponse, error)
