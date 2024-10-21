package app

import (
	"fmt"

	appv1 "github.com/tempestdx/protobuf/gen/go/tempestdx/app/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

type Resource struct {
	// ExternalID is the unique identifier for the resource in the external system.
	ExternalID string
	// DisplayName represents the human-readable name of the resource, to be displayed in the Tempest UI.
	DisplayName string
	// Type is the name of the ResourceDefinition that this resource is an instance of.
	Type string
	// Links are resource-specific links that can help users understand how to use this Resource.
	// A good example of a link is a LinkTypeExternal with a URL to the external system's UI for this resource.
	Links []*Link
	// Properties contains the properties of the resource. These properties are validated against the resource Properties schema.
	Properties map[string]any
}

func (r *Resource) toProto() (*appv1.Resource, error) {
	if r == nil {
		return &appv1.Resource{}, nil
	}

	links := make([]*appv1.Link, 0, len(r.Links))
	for _, l := range r.Links {
		links = append(links, l.toProto())
	}

	properties, err := structpb.NewStruct(r.Properties)
	if err != nil {
		return nil, fmt.Errorf("convert properties to struct: %w", err)
	}

	return &appv1.Resource{
		ExternalId:  r.ExternalID,
		DisplayName: r.DisplayName,
		Type:        r.Type,
		Links:       links,
		Properties:  properties,
	}, nil
}

func resourceFromProto(r *appv1.Resource) *Resource {
	if r == nil {
		return &Resource{}
	}

	links := make([]*Link, 0, len(r.Links))
	for _, l := range r.Links {
		links = append(links, linkFromProto(l))
	}

	return &Resource{
		ExternalID:  r.GetExternalId(),
		DisplayName: r.GetDisplayName(),
		Type:        r.GetType(),
		Links:       links,
		Properties:  r.GetProperties().AsMap(),
	}
}
