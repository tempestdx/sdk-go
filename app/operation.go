package app

import (
	"context"

	appv1 "github.com/tempestdx/protobuf/gen/go/tempestdx/app/v1"
)

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
}

func operationRequestFromProto(r *appv1.ExecuteResourceOperationRequest) *OperationRequest {
	if r == nil {
		return &OperationRequest{
			Metadata: &Metadata{},
			Resource: &Resource{},
		}
	}

	return &OperationRequest{
		Metadata: metadataFromProto(r.Metadata),
		Resource: resourceFromProto(r.Resource),
		Input:    r.Input.AsMap(),
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
