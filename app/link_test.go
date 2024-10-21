package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	appv1 "github.com/tempestdx/protobuf/gen/go/tempestdx/app/v1"
)

func TestLinkToProto(t *testing.T) {
	testCases := []struct {
		desc   string
		link   *Link
		linkpb *appv1.Link
	}{
		{
			desc: "OK",
			link: &Link{
				URL:   "https://example.com",
				Title: "Example",
				Type:  LinkTypeDocumentation,
			},
			linkpb: &appv1.Link{
				Url:   "https://example.com",
				Title: "Example",
				Type:  appv1.LinkType_LINK_TYPE_DOCUMENTATION,
			},
		},
		{
			desc:   "OK - nil",
			link:   nil,
			linkpb: &appv1.Link{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.linkpb, tc.link.toProto())
		})
	}
}

func TestLinkFromProto(t *testing.T) {
	testCases := []struct {
		desc   string
		linkpb *appv1.Link
		link   *Link
	}{
		{
			desc: "OK",
			linkpb: &appv1.Link{
				Url:   "https://example.com",
				Title: "Example",
				Type:  appv1.LinkType_LINK_TYPE_DOCUMENTATION,
			},
			link: &Link{
				URL:   "https://example.com",
				Title: "Example",
				Type:  LinkTypeDocumentation,
			},
		},
		{
			desc:   "OK - nil",
			linkpb: nil,
			link:   &Link{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.link, linkFromProto(tc.linkpb))
		})
	}
}
