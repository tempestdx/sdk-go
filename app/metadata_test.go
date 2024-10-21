package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	appv1 "github.com/tempestdx/protobuf/gen/go/tempestdx/app/v1"
)

func TestMetadataFromProto(t *testing.T) {
	testCases := []struct {
		desc string
		md   *Metadata
		mdpb *appv1.Metadata
	}{
		{
			desc: "OK",
			md: &Metadata{
				ProjectID:   "project-id",
				ProjectName: "project-name",
				Owners: []Owner{
					{
						Name: "test",
						Type: OwnerTypeUser,
					},
					{
						Name: "test2",
						Type: OwnerTypeTeam,
					},
				},
				Author: Owner{
					Name: "test",
					Type: OwnerTypeUser,
				},
			},
			mdpb: &appv1.Metadata{
				ProjectId:   "project-id",
				ProjectName: "project-name",
				Owners: []*appv1.Owner{
					{
						Name: "test",
						Type: appv1.OwnerType_OWNER_TYPE_USER,
					},
					{
						Name: "test2",
						Type: appv1.OwnerType_OWNER_TYPE_TEAM,
					},
				},
				Author: &appv1.Owner{
					Name: "test",
					Type: appv1.OwnerType_OWNER_TYPE_USER,
				},
			},
		},
		{
			desc: "OK - nil Owners",
			md: &Metadata{
				ProjectID:   "project-id",
				ProjectName: "project-name",
				Owners:      []Owner{},
				Author: Owner{
					Name: "test",
					Type: OwnerTypeUser,
				},
			},
			mdpb: &appv1.Metadata{
				ProjectId:   "project-id",
				ProjectName: "project-name",
				Owners:      nil,
				Author: &appv1.Owner{
					Name: "test",
					Type: appv1.OwnerType_OWNER_TYPE_USER,
				},
			},
		},
		{
			desc: "OK - nil Author",
			md: &Metadata{
				ProjectID:   "project-id",
				ProjectName: "project-name",
				Owners: []Owner{
					{
						Name: "test",
						Type: OwnerTypeUser,
					},
					{
						Name: "test2",
						Type: OwnerTypeTeam,
					},
				},
				Author: Owner{},
			},
			mdpb: &appv1.Metadata{
				ProjectId:   "project-id",
				ProjectName: "project-name",
				Owners: []*appv1.Owner{
					{
						Name: "test",
						Type: appv1.OwnerType_OWNER_TYPE_USER,
					},
					{
						Name: "test2",
						Type: appv1.OwnerType_OWNER_TYPE_TEAM,
					},
				},
				Author: nil,
			},
		},
		{
			desc: "OK - nil",
			md:   &Metadata{},
			mdpb: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.md, metadataFromProto(tc.mdpb))
		})
	}
}
