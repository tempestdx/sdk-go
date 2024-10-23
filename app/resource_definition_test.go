package app

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetResourceDefinition(t *testing.T) {
	testCases := []struct {
		desc                       string
		app                        *App
		resourceType               string
		expectedResourceDefinition *ResourceDefinition
		found                      bool
	}{
		{
			desc: "OK - Found",
			app: &App{
				resourceDefinitions: []ResourceDefinition{
					{
						Type: "example",
					},
				},
			},
			resourceType: "example",
			expectedResourceDefinition: &ResourceDefinition{
				Type: "example",
			},
			found: true,
		},
		{
			desc: "OK - Not Found",
			app: &App{
				resourceDefinitions: []ResourceDefinition{
					{
						Type: "example",
					},
				},
			},
			resourceType: "example2",
			found:        false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			resourceDefinition, ok := tc.app.getResourceDefinition(tc.resourceType)
			if !tc.found {
				assert.Nil(t, resourceDefinition)
				assert.False(t, ok)
				return
			}

			assert.True(t, ok)
			assert.Equal(t, tc.expectedResourceDefinition, resourceDefinition)
		})
	}
}

var simpleOpFn = func(_ context.Context, _ *OperationRequest) (*OperationResponse, error) {
	return &OperationResponse{}, nil
}

func TestCreateFn(t *testing.T) {
	parsedEmptySchema := MustParseJSONSchema(emptySchema)

	testCases := []struct {
		desc        string
		fn          func(context.Context, *OperationRequest) (*OperationResponse, error)
		inputSchema *JSONSchema
		properties  *JSONSchema
		shouldPanic bool
	}{
		{
			desc:        "OK",
			fn:          simpleOpFn,
			inputSchema: parsedEmptySchema,
			properties:  parsedEmptySchema,
		},
		{
			desc:        "PANIC - No input schema",
			fn:          simpleOpFn,
			properties:  parsedEmptySchema,
			shouldPanic: true,
		},
		{
			desc:        "PANIC - No properties schema",
			fn:          simpleOpFn,
			inputSchema: parsedEmptySchema,
			shouldPanic: true,
		},
		{
			desc:        "PANIC - No fn",
			properties:  parsedEmptySchema,
			inputSchema: parsedEmptySchema,
			shouldPanic: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			expected := &ResourceDefinition{
				PropertiesSchema: tc.properties,
				create: &operation{
					fn: tc.fn,
					schema: schema{
						input:  tc.inputSchema,
						output: tc.properties,
					},
				},
			}

			rd := &ResourceDefinition{
				Type:             "example",
				PropertiesSchema: tc.properties,
			}

			if tc.shouldPanic {
				assert.Panics(t, func() {
					rd.CreateFn(tc.fn, tc.inputSchema)
				})
				return
			}

			rd.CreateFn(tc.fn, tc.inputSchema)
			assert.NotNil(t, rd.create)
			assert.Equal(t, expected.PropertiesSchema, rd.PropertiesSchema)
			assert.Equal(t, expected.create.schema, rd.create.schema)
		})
	}
}

func TestUpdateFn(t *testing.T) {
	parsedEmptySchema := MustParseJSONSchema(emptySchema)

	testCases := []struct {
		desc        string
		fn          func(context.Context, *OperationRequest) (*OperationResponse, error)
		inputSchema *JSONSchema
		properties  *JSONSchema
		shouldPanic bool
	}{
		{
			desc:        "OK",
			fn:          simpleOpFn,
			inputSchema: parsedEmptySchema,
			properties:  parsedEmptySchema,
		},
		{
			desc:        "PANIC - No input schema",
			fn:          simpleOpFn,
			properties:  parsedEmptySchema,
			shouldPanic: true,
		},
		{
			desc:        "PANIC - No properties schema",
			fn:          simpleOpFn,
			inputSchema: parsedEmptySchema,
			shouldPanic: true,
		},
		{
			desc:        "PANIC - No fn",
			properties:  parsedEmptySchema,
			inputSchema: parsedEmptySchema,
			shouldPanic: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			expected := &ResourceDefinition{
				PropertiesSchema: tc.properties,
				update: &operation{
					fn: tc.fn,
					schema: schema{
						input:  tc.inputSchema,
						output: tc.properties,
					},
				},
			}

			rd := &ResourceDefinition{
				Type:             "example",
				PropertiesSchema: tc.properties,
			}

			if tc.shouldPanic {
				assert.Panics(t, func() {
					rd.UpdateFn(tc.fn, tc.inputSchema)
				})
				return
			}

			rd.UpdateFn(tc.fn, tc.inputSchema)
			assert.NotNil(t, rd.update)
			assert.Equal(t, expected.PropertiesSchema, rd.PropertiesSchema)
			assert.Equal(t, expected.update.schema, rd.update.schema)
		})
	}
}

func TestDeleteFn(t *testing.T) {
	parsedEmptySchema := MustParseJSONSchema(emptySchema)

	testCases := []struct {
		desc        string
		fn          func(context.Context, *OperationRequest) (*OperationResponse, error)
		properties  *JSONSchema
		shouldPanic bool
	}{
		{
			desc:       "OK",
			fn:         simpleOpFn,
			properties: parsedEmptySchema,
		},
		{
			desc:        "PANIC - No properties schema",
			fn:          simpleOpFn,
			shouldPanic: true,
		},
		{
			desc:        "PANIC - No fn",
			properties:  parsedEmptySchema,
			shouldPanic: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			expected := &ResourceDefinition{
				PropertiesSchema: tc.properties,
				delete: &operation{
					fn: tc.fn,
					schema: schema{
						input:  MustParseJSONSchema(GenericEmptySchema),
						output: tc.properties,
					},
				},
			}

			rd := &ResourceDefinition{
				Type:             "example",
				PropertiesSchema: tc.properties,
			}

			if tc.shouldPanic {
				assert.Panics(t, func() {
					rd.DeleteFn(tc.fn)
				})
				return
			}

			rd.DeleteFn(tc.fn)
			assert.NotNil(t, rd.delete)
			assert.Equal(t, expected.PropertiesSchema, rd.PropertiesSchema)
			assert.Equal(t, expected.delete.schema, rd.delete.schema)
		})
	}
}

func TestReadFn(t *testing.T) {
	parsedEmptySchema := MustParseJSONSchema(emptySchema)

	testCases := []struct {
		desc        string
		fn          func(context.Context, *OperationRequest) (*OperationResponse, error)
		properties  *JSONSchema
		shouldPanic bool
	}{
		{
			desc:       "OK",
			fn:         simpleOpFn,
			properties: parsedEmptySchema,
		},
		{
			desc:        "PANIC - No properties schema",
			fn:          simpleOpFn,
			shouldPanic: true,
		},
		{
			desc:        "PANIC - No fn",
			properties:  parsedEmptySchema,
			shouldPanic: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			expected := &ResourceDefinition{
				PropertiesSchema: tc.properties,
				read: &operation{
					fn: tc.fn,
					schema: schema{
						input:  MustParseJSONSchema(GenericEmptySchema),
						output: tc.properties,
					},
				},
			}

			rd := &ResourceDefinition{
				Type:             "example",
				PropertiesSchema: tc.properties,
			}

			if tc.shouldPanic {
				assert.Panics(t, func() {
					rd.ReadFn(tc.fn)
				})
				return
			}

			rd.ReadFn(tc.fn)
			assert.NotNil(t, rd.read)
			assert.Equal(t, expected.PropertiesSchema, rd.PropertiesSchema)
			assert.Equal(t, expected.read.schema, rd.read.schema)
		})
	}
}

func TestListFn(t *testing.T) {
	parsedEmptySchema := MustParseJSONSchema(emptySchema)

	listFn := func(_ context.Context, _ *ListRequest) (*ListResponse, error) {
		return &ListResponse{}, nil
	}

	testCases := []struct {
		desc        string
		fn          func(context.Context, *ListRequest) (*ListResponse, error)
		properties  *JSONSchema
		shouldPanic bool
	}{
		{
			desc:       "OK",
			fn:         listFn,
			properties: parsedEmptySchema,
		},
		{
			desc:        "PANIC - No properties schema",
			fn:          listFn,
			shouldPanic: true,
		},
		{
			desc:        "PANIC - No fn",
			properties:  parsedEmptySchema,
			shouldPanic: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			expected := &ResourceDefinition{
				PropertiesSchema: tc.properties,
				list: &listOperation{
					fn: tc.fn,
					schema: schema{
						input:  MustParseJSONSchema(GenericEmptySchema),
						output: tc.properties,
					},
				},
			}

			rd := &ResourceDefinition{
				Type:             "example",
				PropertiesSchema: tc.properties,
			}

			if tc.shouldPanic {
				assert.Panics(t, func() {
					rd.ListFn(tc.fn)
				})
				return
			}

			rd.ListFn(tc.fn)
			assert.NotNil(t, rd.list)
			assert.Equal(t, expected.PropertiesSchema, rd.PropertiesSchema)
			assert.Equal(t, expected.list.schema, rd.list.schema)
		})
	}
}

func TestHealthcheckFn(t *testing.T) {
	healthcheckFn := func(_ context.Context) (*HealthCheckResponse, error) {
		return &HealthCheckResponse{}, nil
	}

	testCases := []struct {
		desc        string
		fn          func(context.Context) (*HealthCheckResponse, error)
		shouldPanic bool
	}{
		{
			desc: "OK",
			fn:   healthcheckFn,
		},

		{
			desc:        "PANIC - No fn",
			shouldPanic: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			rd := &ResourceDefinition{
				Type: "example",
			}

			if tc.shouldPanic {
				assert.Panics(t, func() {
					rd.HealthCheckFn(tc.fn)
				})
				return
			}

			rd.HealthCheckFn(tc.fn)
			assert.NotNil(t, rd.healthcheck)
		})
	}
}
