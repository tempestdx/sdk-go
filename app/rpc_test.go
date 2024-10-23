package app

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appv1 "github.com/tempestdx/protobuf/gen/go/tempestdx/app/v1"
)

func TestDescribe(t *testing.T) {
	parsedSchemaStruct, err := MustParseJSONSchema(GenericEmptySchema).toStruct()
	require.NoError(t, err)

	testCases := []struct {
		desc   string
		addFns []string
		want   *connect.Response[appv1.DescribeResponse]
		err    error
	}{
		{
			desc:   "OK - Fully Loaded",
			addFns: []string{"create", "read", "update", "delete", "list", "healthcheck"},
			want: connect.NewResponse(&appv1.DescribeResponse{
				ResourceDefinitions: []*appv1.ResourceDefinition{
					{
						Type:             "example",
						DisplayName:      "Example",
						Description:      "An example resource",
						PropertiesSchema: parsedSchemaStruct,
						LifecycleStage:   appv1.LifecycleStage_LIFECYCLE_STAGE_OPERATE,
						Links: []*appv1.Link{
							{
								Url:   "http://example.com",
								Title: "Example",
								Type:  appv1.LinkType_LINK_TYPE_DOCUMENTATION,
							},
						},
						InstructionsMarkdown: "This is an example resource",
						ListSupported:        true,
						HealthcheckSupported: true,
						ReadSupported:        true,
						CreateSupported:      true,
						UpdateSupported:      true,
						DeleteSupported:      true,
						CreateInputSchema:    parsedSchemaStruct,
						UpdateInputSchema:    parsedSchemaStruct,
						Actions: []*appv1.ActionDefinition{
							{
								Name:         "do_something",
								DisplayName:  "Do Something",
								Description:  "Do something with the resource",
								InputSchema:  parsedSchemaStruct,
								OutputSchema: parsedSchemaStruct,
							},
						},
					},
				},
			}),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			app := &App{
				resourceDefinitions: []ResourceDefinition{generateRD(tc.addFns)},
			}

			res, err := app.Describe(context.Background(), connect.NewRequest(&appv1.DescribeRequest{}))
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, res)
		})
	}
}

func generateRD(fns []string) ResourceDefinition {
	parsedSchema := MustParseJSONSchema(GenericEmptySchema)

	rd := ResourceDefinition{
		Type:             "example",
		DisplayName:      "Example",
		Description:      "An example resource",
		PropertiesSchema: parsedSchema,
		LifecycleStage:   LifecycleStageOperate,
		Links: []Link{
			{
				URL:   "http://example.com",
				Title: "Example",
				Type:  LinkTypeDocumentation,
			},
		},
		InstructionsMarkdown: "This is an example resource",
		actions: []ActionDefinition{
			{
				Name:         "do_something",
				DisplayName:  "Do Something",
				Description:  "Do something with the resource",
				InputSchema:  parsedSchema,
				OutputSchema: parsedSchema,
			},
		},
	}

	for _, fn := range fns {
		switch fn {
		case "create":
			rd.CreateFn(simpleOpFn, parsedSchema)
		case "update":
			rd.UpdateFn(simpleOpFn, parsedSchema)
		case "read":
			rd.ReadFn(simpleOpFn)
		case "delete":
			rd.DeleteFn(simpleOpFn)
		case "list":
			rd.ListFn(simpleListFn)
		case "healthcheck":
			rd.HealthCheckFn(simpleHealthcheckFn)
		}
	}

	return rd
}
