package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	appv1 "github.com/tempestdx/protobuf/gen/go/tempestdx/app/v1"
)

func TestListRequestFromProto(t *testing.T) {
	testCases := []struct {
		desc      string
		listReq   *ListRequest
		listReqpb *appv1.ListResourcesRequest
	}{
		{
			desc: "OK",
			listReq: &ListRequest{
				Resource: &Resource{
					ExternalID:  "external-id",
					DisplayName: "display-name",
					Type:        "type",
					Links:       []*Link{},
					Properties:  map[string]any{},
				},
				Next: "1",
			},
			listReqpb: &appv1.ListResourcesRequest{
				Resource: &appv1.Resource{
					ExternalId:  "external-id",
					DisplayName: "display-name",
					Type:        "type",
				},
				Metadata: nil,
				Next:     "1",
			},
		},
		{
			desc:      "OK - nil",
			listReq:   nil,
			listReqpb: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.listReq, listRequestFromProto(tc.listReqpb))
		})
	}
}
