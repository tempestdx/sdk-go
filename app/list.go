package app

import (
	"context"

	appv1 "github.com/tempestdx/protobuf/gen/go/tempestdx/app/v1"
)

// ListRequest contains the input data for listing resources.
type ListRequest struct {
	// Metadata contains information about the Project and User making the request.
	// This metadata does not contain information about the Resource being operated on.
	Metadata *Metadata
	// Resource is the resource being listed, and contains at least the Resource Type.
	Resource *Resource
	// Next is a token that can be used to fetch the next page of results.
	Next string
}

func listRequestFromProto(r *appv1.ListResourcesRequest) *ListRequest {
	if r == nil {
		return nil
	}

	return &ListRequest{
		Metadata: metadataFromProto(r.Metadata),
		Resource: resourceFromProto(r.Resource),
		Next:     r.Next,
	}
}

// ListResponse contains the output data for listing resources.
type ListResponse struct {
	Resources []*Resource
	// Next is a token that can be used to fetch the next page of results.
	Next string
}

type listOperation struct {
	schema schema
	fn     ListFunc
}

type ListFunc func(context.Context, *ListRequest) (*ListResponse, error)
