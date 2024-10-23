package app

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appv1 "github.com/tempestdx/protobuf/gen/go/tempestdx/app/v1"
)

func TestResourceFromProto(t *testing.T) {
	testCases := []struct {
		desc       string
		resource   *Resource
		resourcepb *appv1.Resource
	}{
		{
			desc: "OK",
			resource: &Resource{
				ExternalID:  "external-id",
				DisplayName: "display-name",
				Type:        "type",
				Links: []*Link{
					{
						URL:   "http://example.com",
						Title: "Example",
						Type:  LinkTypeDocumentation,
					},
				},
				Properties: map[string]any{},
			},
			resourcepb: &appv1.Resource{
				ExternalId:  "external-id",
				DisplayName: "display-name",
				Type:        "type",
				Links: []*appv1.Link{
					{
						Url:   "http://example.com",
						Title: "Example",
						Type:  appv1.LinkType_LINK_TYPE_DOCUMENTATION,
					},
				},
				Properties: nil,
			},
		},
		{
			desc:       "OK - nil",
			resource:   nil,
			resourcepb: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.resource, resourceFromProto(tc.resourcepb))
		})
	}
}

func TestResourceToProto(t *testing.T) {
	testCases := []struct {
		desc       string
		resource   *Resource
		resourcepb *appv1.Resource
		err        error
	}{
		{
			desc: "OK",
			resource: &Resource{
				ExternalID:  "external-id",
				DisplayName: "display-name",
				Type:        "type",
				Links: []*Link{
					{
						URL:   "http://example.com",
						Title: "Example",
						Type:  LinkTypeDocumentation,
					},
				},
				Properties: map[string]any{
					"key": "value",
				},
			},
			resourcepb: &appv1.Resource{
				ExternalId:  "external-id",
				DisplayName: "display-name",
				Type:        "type",
				Links: []*appv1.Link{
					{
						Url:   "http://example.com",
						Title: "Example",
						Type:  appv1.LinkType_LINK_TYPE_DOCUMENTATION,
					},
				},
				Properties: mustNewStruct(map[string]any{"key": "value"}),
			},
		},
		{
			desc:       "ERR - nil",
			resource:   nil,
			resourcepb: nil,
			err:        fmt.Errorf("resource is nil"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			resourcepb, err := tc.resource.toProto()
			if tc.err != nil {
				assert.EqualError(t, tc.err, err.Error())
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.resourcepb, resourcepb)
		})
	}
}
