package app

import (
	"context"
	"fmt"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appv1 "github.com/tempestdx/protobuf/gen/go/tempestdx/app/v1"
	"github.com/tempestdx/sdk-go/jsonschema"
)

func TestDescribe(t *testing.T) {
	parsedSchemaStruct, err := jsonschema.MustParseSchema(jsonschema.GenericEmptySchema).ToStruct()
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

func TestExecuteResourceOperationErrors(t *testing.T) {
	testCases := []struct {
		desc string
		req  *appv1.ExecuteResourceOperationRequest
		err  error
	}{
		{
			desc: "ERR - nil resource",
			req: &appv1.ExecuteResourceOperationRequest{
				Resource: nil,
			},
			err: fmt.Errorf("invalid_argument: resource is required"),
		},
		{
			desc: "ERR - type not found",
			req: &appv1.ExecuteResourceOperationRequest{
				Resource: &appv1.Resource{
					Type: "not_found",
				},
			},
			err: fmt.Errorf("not_found: resource type not_found not found"),
		},
		{
			desc: "ERR - invalid operation",
			req: &appv1.ExecuteResourceOperationRequest{
				Resource: &appv1.Resource{
					Type: "example",
				},
				Operation: appv1.ResourceOperation_RESOURCE_OPERATION_UNSPECIFIED,
			},
			err: fmt.Errorf("invalid_argument: unsupported operation RESOURCE_OPERATION_UNSPECIFIED"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			app := &App{
				resourceDefinitions: []ResourceDefinition{generateRD(nil)},
			}

			res, err := app.ExecuteResourceOperation(context.Background(), connect.NewRequest(tC.req))
			assert.EqualError(t, err, tC.err.Error())
			assert.Nil(t, res)
		})
	}
}

var simpleSchema = []byte(`{
	"$schema": "http://json-schema.org/draft-07/schema#",
	"properties": {
		"property1": {
			"type": "string",
			"default": "default1"
		},
		"property2": {
			"type": "string"
		}
	},
	"required": ["property1", "property2"]
}`)

func TestExecuteResourceOperation_Create(t *testing.T) {
	parsedSchema := jsonschema.MustParseSchema(jsonschema.GenericEmptySchema)
	parsedSimpleSchema := jsonschema.MustParseSchema(simpleSchema)

	testCases := []struct {
		desc             string
		want             *connect.Response[appv1.ExecuteResourceOperationResponse]
		req              *appv1.ExecuteResourceOperationRequest
		inputSchema      *jsonschema.Schema
		propertiesSchema *jsonschema.Schema
		enableCreate     bool
		createErr        error
		err              error
	}{
		{
			desc:         "OK",
			enableCreate: true,
			inputSchema:  parsedSchema,
			req: &appv1.ExecuteResourceOperationRequest{
				Resource: &appv1.Resource{
					Type: "example",
				},
				Input: mustNewStruct(map[string]any{
					"key": "value",
				}),
				Operation: appv1.ResourceOperation_RESOURCE_OPERATION_CREATE,
			},
			want: connect.NewResponse(&appv1.ExecuteResourceOperationResponse{
				Resource: &appv1.Resource{
					Type:        "example",
					ExternalId:  "example-1",
					DisplayName: "Example",
					Properties: mustNewStruct(map[string]any{
						"key":       "value",
						"other_key": "other_value",
					}),
					Links: []*appv1.Link{
						{
							Url:   "http://example.com",
							Title: "Example",
							Type:  appv1.LinkType_LINK_TYPE_DOCUMENTATION,
						},
					},
				},
			}),
		},
		{
			desc:        "ERR - Create Disabled",
			inputSchema: parsedSchema,
			req: &appv1.ExecuteResourceOperationRequest{
				Resource: &appv1.Resource{
					Type: "example",
				},
				Input: mustNewStruct(map[string]any{
					"key": "value",
				}),
				Operation: appv1.ResourceOperation_RESOURCE_OPERATION_CREATE,
			},
			err: fmt.Errorf("invalid_argument: operation RESOURCE_OPERATION_CREATE not supported for resource type example"),
		},
		{
			desc:         "ERR - Invalid Input",
			enableCreate: true,
			inputSchema:  parsedSimpleSchema,
			req: &appv1.ExecuteResourceOperationRequest{
				Resource: &appv1.Resource{
					Type: "example",
				},
				Input: mustNewStruct(map[string]any{
					"key": "value",
				}),
				Operation: appv1.ResourceOperation_RESOURCE_OPERATION_CREATE,
			},
			err: fmt.Errorf("invalid_argument: validate create input: jsonschema validation failed"),
		},
		{
			desc:         "ERR - Create Error",
			enableCreate: true,
			inputSchema:  parsedSchema,
			req: &appv1.ExecuteResourceOperationRequest{
				Resource: &appv1.Resource{
					Type: "example",
				},
				Input: mustNewStruct(map[string]any{
					"key": "value",
				}),
				Operation: appv1.ResourceOperation_RESOURCE_OPERATION_CREATE,
			},
			createErr: fmt.Errorf("create error"),
			err:       fmt.Errorf("internal: create resource: create error"),
		},
		{
			desc:             "ERR - Invalid Output",
			enableCreate:     true,
			inputSchema:      parsedSchema,
			propertiesSchema: parsedSimpleSchema,
			req: &appv1.ExecuteResourceOperationRequest{
				Resource: &appv1.Resource{
					Type: "example",
				},
				Input: mustNewStruct(map[string]any{
					"key": "value",
				}),
				Operation: appv1.ResourceOperation_RESOURCE_OPERATION_CREATE,
			},
			err: fmt.Errorf("internal: validate create output: jsonschema validation failed"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			rd := generateRD(nil)
			if tc.propertiesSchema != nil {
				rd.PropertiesSchema = tc.propertiesSchema
			}

			if tc.enableCreate {
				rd.CreateFn(func(_ context.Context, req *OperationRequest) (*OperationResponse, error) {
					if tc.createErr != nil {
						return nil, tc.createErr
					}

					return &OperationResponse{
						Resource: &Resource{
							ExternalID:  "example-1",
							DisplayName: "Example",
							Type:        "example",
							Properties: map[string]any{
								"key":       "value",
								"other_key": "other_value",
							},
							Links: []*Link{
								{
									URL:   "http://example.com",
									Title: "Example",
									Type:  LinkTypeDocumentation,
								},
							},
						},
					}, nil
				}, tc.inputSchema)
			}

			app := &App{
				resourceDefinitions: []ResourceDefinition{rd},
			}

			res, err := app.ExecuteResourceOperation(context.Background(), connect.NewRequest(tc.req))
			if tc.err != nil {
				if tc.desc == "ERR - Invalid Output" || tc.desc == "ERR - Invalid Input" {
					assert.ErrorContains(t, err, tc.err.Error())
				} else {
					assert.EqualError(t, err, tc.err.Error())
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, res)
		})
	}
}

func TestExecuteResourceOperation_Update(t *testing.T) {
	parsedSchema := jsonschema.MustParseSchema(jsonschema.GenericEmptySchema)
	parsedSimpleSchema := jsonschema.MustParseSchema(simpleSchema)

	testCases := []struct {
		desc             string
		want             *connect.Response[appv1.ExecuteResourceOperationResponse]
		req              *appv1.ExecuteResourceOperationRequest
		inputSchema      *jsonschema.Schema
		propertiesSchema *jsonschema.Schema
		enableUpdate     bool
		updateErr        error
		err              error
	}{
		{
			desc:         "OK",
			enableUpdate: true,
			inputSchema:  parsedSchema,
			req: &appv1.ExecuteResourceOperationRequest{
				Resource: &appv1.Resource{
					Type:       "example",
					ExternalId: "example-1",
				},
				Input: mustNewStruct(map[string]any{
					"key": "value",
				}),
				Operation: appv1.ResourceOperation_RESOURCE_OPERATION_UPDATE,
			},
			want: connect.NewResponse(&appv1.ExecuteResourceOperationResponse{
				Resource: &appv1.Resource{
					Type:        "example",
					ExternalId:  "example-1",
					DisplayName: "Example",
					Properties: mustNewStruct(map[string]any{
						"key":       "value",
						"other_key": "other_value",
					}),
					Links: []*appv1.Link{
						{
							Url:   "http://example.com",
							Title: "Example",
							Type:  appv1.LinkType_LINK_TYPE_DOCUMENTATION,
						},
					},
				},
			}),
		},
		{
			desc:        "ERR - Update Disabled",
			inputSchema: parsedSchema,
			req: &appv1.ExecuteResourceOperationRequest{
				Resource: &appv1.Resource{
					Type:       "example",
					ExternalId: "example-1",
				},
				Input: mustNewStruct(map[string]any{
					"key": "value",
				}),
				Operation: appv1.ResourceOperation_RESOURCE_OPERATION_UPDATE,
			},
			err: fmt.Errorf("invalid_argument: operation RESOURCE_OPERATION_UPDATE not supported for resource type example"),
		},
		{
			desc:         "ERR - Invalid Input",
			enableUpdate: true,
			inputSchema:  parsedSimpleSchema,
			req: &appv1.ExecuteResourceOperationRequest{
				Resource: &appv1.Resource{
					Type:       "example",
					ExternalId: "example-1",
				},
				Input: mustNewStruct(map[string]any{
					"key": "value",
				}),
				Operation: appv1.ResourceOperation_RESOURCE_OPERATION_UPDATE,
			},
			err: fmt.Errorf("invalid_argument: validate update input: jsonschema validation failed"),
		},
		{
			desc:         "ERR - Update Error",
			enableUpdate: true,
			inputSchema:  parsedSchema,
			req: &appv1.ExecuteResourceOperationRequest{
				Resource: &appv1.Resource{
					Type:       "example",
					ExternalId: "example-1",
				},
				Input: mustNewStruct(map[string]any{
					"key": "value",
				}),
				Operation: appv1.ResourceOperation_RESOURCE_OPERATION_UPDATE,
			},
			updateErr: fmt.Errorf("update error"),
			err:       fmt.Errorf("internal: update resource: update error"),
		},
		{
			desc:             "ERR - Invalid Output",
			enableUpdate:     true,
			inputSchema:      parsedSchema,
			propertiesSchema: parsedSimpleSchema,
			req: &appv1.ExecuteResourceOperationRequest{
				Resource: &appv1.Resource{
					Type:       "example",
					ExternalId: "example-1",
				},
				Input: mustNewStruct(map[string]any{
					"key": "value",
				}),
				Operation: appv1.ResourceOperation_RESOURCE_OPERATION_UPDATE,
			},
			err: fmt.Errorf("internal: validate update output: jsonschema validation failed"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			rd := generateRD(nil)
			if tc.propertiesSchema != nil {
				rd.PropertiesSchema = tc.propertiesSchema
			}

			if tc.enableUpdate {
				rd.UpdateFn(func(_ context.Context, req *OperationRequest) (*OperationResponse, error) {
					if tc.updateErr != nil {
						return nil, tc.updateErr
					}

					return &OperationResponse{
						Resource: &Resource{
							ExternalID:  "example-1",
							DisplayName: "Example",
							Type:        "example",
							Properties: map[string]any{
								"key":       "value",
								"other_key": "other_value",
							},
							Links: []*Link{
								{
									URL:   "http://example.com",
									Title: "Example",
									Type:  LinkTypeDocumentation,
								},
							},
						},
					}, nil
				}, tc.inputSchema)
			}

			app := &App{
				resourceDefinitions: []ResourceDefinition{rd},
			}

			res, err := app.ExecuteResourceOperation(context.Background(), connect.NewRequest(tc.req))
			if tc.err != nil {
				if tc.desc == "ERR - Invalid Output" || tc.desc == "ERR - Invalid Input" {
					assert.ErrorContains(t, err, tc.err.Error())
				} else {
					assert.EqualError(t, err, tc.err.Error())
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, res)
		})
	}
}

func TestExecuteResourceOperation_Delete(t *testing.T) {
	testCases := []struct {
		desc         string
		want         *connect.Response[appv1.ExecuteResourceOperationResponse]
		req          *appv1.ExecuteResourceOperationRequest
		enableDelete bool
		deleteErr    error
		err          error
	}{
		{
			desc:         "OK",
			enableDelete: true,
			req: &appv1.ExecuteResourceOperationRequest{
				Resource: &appv1.Resource{
					Type:       "example",
					ExternalId: "example-1",
				},
				Operation: appv1.ResourceOperation_RESOURCE_OPERATION_DELETE,
			},
			want: connect.NewResponse(&appv1.ExecuteResourceOperationResponse{
				Resource: &appv1.Resource{
					Type:        "example",
					ExternalId:  "example-1",
					DisplayName: "Example",
					Properties: mustNewStruct(map[string]any{
						"key":       "value",
						"other_key": "other_value",
					}),
					Links: []*appv1.Link{
						{
							Url:   "http://example.com",
							Title: "Example",
							Type:  appv1.LinkType_LINK_TYPE_DOCUMENTATION,
						},
					},
				},
			}),
		},
		{
			desc: "ERR - Delete Disabled",
			req: &appv1.ExecuteResourceOperationRequest{
				Resource: &appv1.Resource{
					Type:       "example",
					ExternalId: "example-1",
				},
				Operation: appv1.ResourceOperation_RESOURCE_OPERATION_DELETE,
			},
			err: fmt.Errorf("invalid_argument: operation RESOURCE_OPERATION_DELETE not supported for resource type example"),
		},
		{
			desc:         "ERR - No External ID",
			enableDelete: true,
			req: &appv1.ExecuteResourceOperationRequest{
				Resource: &appv1.Resource{
					Type: "example",
				},
				Operation: appv1.ResourceOperation_RESOURCE_OPERATION_DELETE,
			},
			err: fmt.Errorf("invalid_argument: external ID is required for delete operation"),
		},
		{
			desc:         "ERR - Delete Error",
			enableDelete: true,
			req: &appv1.ExecuteResourceOperationRequest{
				Resource: &appv1.Resource{
					Type:       "example",
					ExternalId: "example-1",
				},
				Operation: appv1.ResourceOperation_RESOURCE_OPERATION_DELETE,
			},
			deleteErr: fmt.Errorf("delete error"),
			err:       fmt.Errorf("internal: delete resource: delete error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			rd := generateRD(nil)
			if tc.enableDelete {
				rd.DeleteFn(func(_ context.Context, req *OperationRequest) (*OperationResponse, error) {
					if tc.deleteErr != nil {
						return nil, tc.deleteErr
					}

					return &OperationResponse{
						Resource: &Resource{
							ExternalID:  "example-1",
							DisplayName: "Example",
							Type:        "example",
							Properties: map[string]any{
								"key":       "value",
								"other_key": "other_value",
							},
							Links: []*Link{
								{
									URL:   "http://example.com",
									Title: "Example",
									Type:  LinkTypeDocumentation,
								},
							},
						},
					}, nil
				})
			}

			app := &App{
				resourceDefinitions: []ResourceDefinition{rd},
			}

			res, err := app.ExecuteResourceOperation(context.Background(), connect.NewRequest(tc.req))
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, res)
		})
	}
}

func TestExecuteResourceOperation_Read(t *testing.T) {
	parsedSimpleSchema := jsonschema.MustParseSchema(simpleSchema)

	testCases := []struct {
		desc             string
		want             *connect.Response[appv1.ExecuteResourceOperationResponse]
		req              *appv1.ExecuteResourceOperationRequest
		propertiesSchema *jsonschema.Schema
		enableRead       bool
		readErr          error
		err              error
	}{
		{
			desc:       "OK",
			enableRead: true,
			req: &appv1.ExecuteResourceOperationRequest{
				Resource: &appv1.Resource{
					Type:       "example",
					ExternalId: "example-1",
				},
				Operation: appv1.ResourceOperation_RESOURCE_OPERATION_READ,
			},
			want: connect.NewResponse(&appv1.ExecuteResourceOperationResponse{
				Resource: &appv1.Resource{
					Type:        "example",
					ExternalId:  "example-1",
					DisplayName: "Example",
					Properties: mustNewStruct(map[string]any{
						"key":       "value",
						"other_key": "other_value",
					}),
					Links: []*appv1.Link{
						{
							Url:   "http://example.com",
							Title: "Example",
							Type:  appv1.LinkType_LINK_TYPE_DOCUMENTATION,
						},
					},
				},
			}),
		},
		{
			desc: "ERR - Read Disabled",
			req: &appv1.ExecuteResourceOperationRequest{
				Resource: &appv1.Resource{
					Type:       "example",
					ExternalId: "example-1",
				},
				Operation: appv1.ResourceOperation_RESOURCE_OPERATION_READ,
			},
			err: fmt.Errorf("invalid_argument: operation RESOURCE_OPERATION_READ not supported for resource type example"),
		},
		{
			desc:       "ERR - No External ID",
			enableRead: true,
			req: &appv1.ExecuteResourceOperationRequest{
				Resource: &appv1.Resource{
					Type: "example",
				},
				Operation: appv1.ResourceOperation_RESOURCE_OPERATION_READ,
			},
			err: fmt.Errorf("invalid_argument: external ID is required for read operation"),
		},
		{
			desc:       "ERR - Read Error",
			enableRead: true,
			req: &appv1.ExecuteResourceOperationRequest{
				Resource: &appv1.Resource{
					Type:       "example",
					ExternalId: "example-1",
				},
				Operation: appv1.ResourceOperation_RESOURCE_OPERATION_READ,
			},
			readErr: fmt.Errorf("read error"),
			err:     fmt.Errorf("internal: read resource: read error"),
		},
		{
			desc:             "ERR - Invalid Output",
			enableRead:       true,
			propertiesSchema: parsedSimpleSchema,
			req: &appv1.ExecuteResourceOperationRequest{
				Resource: &appv1.Resource{
					Type:       "example",
					ExternalId: "example-1",
				},
				Operation: appv1.ResourceOperation_RESOURCE_OPERATION_READ,
			},
			err: fmt.Errorf("internal: validate read output: jsonschema validation failed"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			rd := generateRD(nil)
			if tc.propertiesSchema != nil {
				rd.PropertiesSchema = tc.propertiesSchema
			}

			if tc.enableRead {
				rd.ReadFn(func(_ context.Context, req *OperationRequest) (*OperationResponse, error) {
					if tc.readErr != nil {
						return nil, tc.readErr
					}

					return &OperationResponse{
						Resource: &Resource{
							ExternalID:  "example-1",
							DisplayName: "Example",
							Type:        "example",
							Properties: map[string]any{
								"key":       "value",
								"other_key": "other_value",
							},
							Links: []*Link{
								{
									URL:   "http://example.com",
									Title: "Example",
									Type:  LinkTypeDocumentation,
								},
							},
						},
					}, nil
				})
			}

			app := &App{
				resourceDefinitions: []ResourceDefinition{rd},
			}

			res, err := app.ExecuteResourceOperation(context.Background(), connect.NewRequest(tc.req))
			if tc.err != nil {
				if tc.desc == "ERR - Invalid Output" {
					assert.ErrorContains(t, err, tc.err.Error())
				} else {
					assert.EqualError(t, err, tc.err.Error())
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, res)
		})
	}
}

func TestListResources(t *testing.T) {
	parsedSimpleSchema := jsonschema.MustParseSchema(simpleSchema)

	testCases := []struct {
		desc             string
		numResources     int
		want             *connect.Response[appv1.ListResourcesResponse]
		req              *appv1.ListResourcesRequest
		propertiesSchema *jsonschema.Schema
		enableList       bool
		listErr          error
		err              error
	}{
		{
			desc:         "OK - Single Resource",
			enableList:   true,
			numResources: 1,
			req: &appv1.ListResourcesRequest{
				Resource: &appv1.Resource{
					Type: "example",
				},
				Next: "1",
			},
			want: connect.NewResponse(&appv1.ListResourcesResponse{
				Resources: []*appv1.Resource{
					{
						Type:        "example",
						ExternalId:  "example-1",
						DisplayName: "Example",
						Properties: mustNewStruct(map[string]any{
							"key":       "value",
							"other_key": "other_value",
						}),
						Links: []*appv1.Link{
							{
								Url:   "http://example.com",
								Title: "Example",
								Type:  appv1.LinkType_LINK_TYPE_DOCUMENTATION,
							},
						},
					},
				},
			}),
		},
		{
			desc:         "OK - Multiple Resources",
			enableList:   true,
			numResources: 2,
			req: &appv1.ListResourcesRequest{
				Resource: &appv1.Resource{
					Type: "example",
				},
				Next: "1",
			},
			want: connect.NewResponse(&appv1.ListResourcesResponse{
				Resources: []*appv1.Resource{
					{
						Type:        "example",
						ExternalId:  "example-1",
						DisplayName: "Example",
						Properties: mustNewStruct(map[string]any{
							"key":       "value",
							"other_key": "other_value",
						}),
						Links: []*appv1.Link{
							{
								Url:   "http://example.com",
								Title: "Example",
								Type:  appv1.LinkType_LINK_TYPE_DOCUMENTATION,
							},
						},
					},
					{
						Type:        "example",
						ExternalId:  "example-2",
						DisplayName: "Example",
						Properties: mustNewStruct(map[string]any{
							"key":       "value",
							"other_key": "other_value",
						}),
						Links: []*appv1.Link{
							{
								Url:   "http://example.com",
								Title: "Example",
								Type:  appv1.LinkType_LINK_TYPE_DOCUMENTATION,
							},
						},
					},
				},
			}),
		},
		{
			desc: "ERR - List Disabled",
			req: &appv1.ListResourcesRequest{
				Resource: &appv1.Resource{
					Type: "example",
				},
				Next: "1",
			},
			err: fmt.Errorf("invalid_argument: list operation not supported for resource type example"),
		},
		{
			desc:       "ERR - List Error",
			enableList: true,
			req: &appv1.ListResourcesRequest{
				Resource: &appv1.Resource{
					Type: "example",
				},
				Next: "1",
			},
			listErr: fmt.Errorf("list error"),
			err:     fmt.Errorf("internal: list resources: list error"),
		},
		{
			desc:             "ERR - Invalid Output",
			enableList:       true,
			numResources:     2,
			propertiesSchema: parsedSimpleSchema,
			req: &appv1.ListResourcesRequest{
				Resource: &appv1.Resource{
					Type: "example",
				},
				Next: "1",
			},
			err: fmt.Errorf("internal: validate resource properties: jsonschema validation failed"),
		},
		{
			desc: "ERR - Resource Missing",
			req: &appv1.ListResourcesRequest{
				Resource: &appv1.Resource{
					Type: "not_found",
				},
				Next: "1",
			},
			err: fmt.Errorf("not_found: resource type not_found not found"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			rd := generateRD(nil)
			if tc.propertiesSchema != nil {
				rd.PropertiesSchema = tc.propertiesSchema
			}

			if tc.enableList {
				rd.ListFn(func(_ context.Context, req *ListRequest) (*ListResponse, error) {
					if tc.listErr != nil {
						return nil, tc.listErr
					}

					var resources []*Resource
					for i := 0; i < tc.numResources; i++ {
						resources = append(resources, &Resource{
							ExternalID:  fmt.Sprintf("example-%d", i+1),
							DisplayName: "Example",
							Type:        "example",
							Properties: map[string]any{
								"key":       "value",
								"other_key": "other_value",
							},
							Links: []*Link{
								{
									URL:   "http://example.com",
									Title: "Example",
									Type:  LinkTypeDocumentation,
								},
							},
						})
					}

					return &ListResponse{
						Resources: resources,
					}, nil
				})
			}

			app := &App{
				resourceDefinitions: []ResourceDefinition{rd},
			}

			res, err := app.ListResources(context.Background(), connect.NewRequest(tc.req))
			if tc.err != nil {
				if tc.desc == "ERR - Invalid Output" {
					assert.ErrorContains(t, err, tc.err.Error())
				} else {
					assert.EqualError(t, err, tc.err.Error())
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, res)
		})
	}
}

func TestHealthCheck(t *testing.T) {
	testCases := []struct {
		desc           string
		status         string
		healthcheckErr error
		req            *appv1.HealthCheckRequest
		want           *connect.Response[appv1.HealthCheckResponse]
		err            error
	}{
		{
			desc:   "OK - Healthy",
			status: "healthy",
			req:    &appv1.HealthCheckRequest{Type: "example"},
			want: connect.NewResponse(&appv1.HealthCheckResponse{
				Status: appv1.HealthCheckStatus_HEALTH_CHECK_STATUS_HEALTHY,
			}),
		},
		{
			desc:   "OK - Disrupted",
			status: "disrupted",
			req:    &appv1.HealthCheckRequest{Type: "example"},
			want: connect.NewResponse(&appv1.HealthCheckResponse{
				Status: appv1.HealthCheckStatus_HEALTH_CHECK_STATUS_DISRUPTED,
			}),
		},
		{
			desc:   "OK - Degraded",
			status: "degraded",
			req:    &appv1.HealthCheckRequest{Type: "example"},
			want: connect.NewResponse(&appv1.HealthCheckResponse{
				Status: appv1.HealthCheckStatus_HEALTH_CHECK_STATUS_DEGRADED,
			}),
		},
		{
			desc:           "OK - Healthcheck Error - Disrupted",
			healthcheckErr: fmt.Errorf("not ok"),
			req:            &appv1.HealthCheckRequest{Type: "example"},
			want: connect.NewResponse(&appv1.HealthCheckResponse{
				Status:  appv1.HealthCheckStatus_HEALTH_CHECK_STATUS_DISRUPTED,
				Message: "not ok",
			}),
		},
		{
			desc: "ERR - No Type",
			req:  &appv1.HealthCheckRequest{},
			err:  fmt.Errorf("invalid_argument: resource type is required"),
		},
		{
			desc: "ERR - Not Found",
			req:  &appv1.HealthCheckRequest{Type: "not_found"},
			err:  fmt.Errorf("not_found: resource type not_found not found"),
		},
		{
			desc:   "ERR - Unknown",
			status: "unknown",
			req:    &appv1.HealthCheckRequest{Type: "example"},
			err:    fmt.Errorf("internal: unknown health check status unknown"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			rd := generateRD(nil)
			rd.healthcheck = func(_ context.Context) (*HealthCheckResponse, error) {
				var status HealthCheckStatus
				switch tc.status {
				case "healthy":
					status = HealthCheckStatusHealthy
				case "disrupted":
					status = HealthCheckStatusDisrupted
				case "degraded":
					status = HealthCheckStatusDegraded
				case "unknown":
					status = HealthCheckStatusUnknown
				}
				return &HealthCheckResponse{
					Status: status,
				}, tc.healthcheckErr
			}

			app := &App{
				resourceDefinitions: []ResourceDefinition{rd},
			}

			res, err := app.HealthCheck(context.Background(), connect.NewRequest(tc.req))
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
	parsedSchema := jsonschema.MustParseSchema(jsonschema.GenericEmptySchema)

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
