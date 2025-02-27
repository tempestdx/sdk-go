package jsonschema

import (
	"errors"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	invalidSchema = []byte(`{
    "$schema": "http://json-schema.org/draft-07/schema#",
    "$id":
}`)

	emptySchema = []byte(`{}`)

	schemaWithDefaults = []byte(`{
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
	"required": ["property1"]
}`)
)

func TestParseJSONSchema(t *testing.T) {
	testCases := []struct {
		desc     string
		schema   []byte
		expected *Schema
		err      error
	}{
		{
			desc:   "OK - Empty JSON Schema",
			schema: emptySchema,
			expected: &Schema{
				Schema: &jsonschema.Schema{},
				Raw:    []byte(`{}`),
			},
		},
		{
			desc:   "OK - Generic Empty Schema",
			schema: GenericEmptySchema,
			expected: &Schema{
				Schema: &jsonschema.Schema{
					Properties:           map[string]*jsonschema.Schema{},
					AdditionalProperties: true,
				},
				Raw: GenericEmptySchema,
			},
		},
		{
			desc:   "ERR - empty schema",
			schema: []byte{},
			err:    errors.New("schema is empty"),
		},
		{
			desc: "ERR - nil schema",
			err:  errors.New("schema is empty"),
		},
		{
			desc:   "ERR - invalid schema",
			schema: invalidSchema,
			err:    errors.New("un marshall schema: invalid character '}' looking for beginning of value"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			out, err := ParseSchema(tc.schema)
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expected.Schema.Properties, out.Properties)
			assert.Equal(t, tc.expected.AdditionalProperties, out.AdditionalProperties)
			assert.Equal(t, tc.expected.Required, out.Required)

			assert.Equal(t, tc.expected.Raw, out.Raw)
		})
	}
}

func TestMustParseJSONSchema(t *testing.T) {
	testCases := []struct {
		desc        string
		schema      []byte
		shouldPanic bool
	}{
		{
			desc:   "OK - No panic",
			schema: []byte(`{}`),
		},
		{
			desc:        "PANIC - empty schema",
			shouldPanic: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			if tc.shouldPanic {
				require.Panics(t, func() {
					MustParseSchema(tc.schema)
				})
				return
			}

			assert.NotNil(t, MustParseSchema(tc.schema))
		})
	}
}

func TestJSONSchemaToStruct(t *testing.T) {
	testCases := []struct {
		desc   string
		schema *Schema
		spb    *structpb.Struct
		err    error
	}{
		{
			desc:   "OK - Empty JSON Schema",
			schema: MustParseSchema(emptySchema),
			spb: &structpb.Struct{
				Fields: map[string]*structpb.Value{},
			},
		},
		{
			desc:   "OK - Generic Empty Schema",
			schema: MustParseSchema(GenericEmptySchema),
			spb: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"$comment": {Kind: &structpb.Value_StringValue{
						StringValue: "This is a generic empty schema that can be used as a placeholder for schemas that are still in development.",
					}},
					"$id": {Kind: &structpb.Value_StringValue{
						StringValue: "https://schema.tempestdx.io/sdk/generic_empty_schema.json",
					}},
					"$schema": {Kind: &structpb.Value_StringValue{
						StringValue: "http://json-schema.org/draft-07/schema#",
					}},
					"additionalProperties": {Kind: &structpb.Value_BoolValue{
						BoolValue: true,
					}},
					"properties": {Kind: &structpb.Value_StructValue{
						StructValue: &structpb.Struct{
							Fields: map[string]*structpb.Value{},
						},
					}},
					"required": {Kind: &structpb.Value_ListValue{
						ListValue: &structpb.ListValue{
							Values: []*structpb.Value{},
						},
					}},
					"type": {Kind: &structpb.Value_StringValue{
						StringValue: "object",
					}},
				},
			},
		},
		{
			desc: "OK - Empty Raw",
			schema: &Schema{
				Schema: &jsonschema.Schema{},
				Raw:    nil,
			},
			spb: &structpb.Struct{
				Fields: map[string]*structpb.Value{},
			},
		},
		{
			desc: "ERR - Unmarshal",
			schema: &Schema{
				Schema: &jsonschema.Schema{},
				Raw:    invalidSchema,
			},
			err: errors.New("unmarshal schema: invalid character '}' looking for beginning of value"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			spb, err := tc.schema.ToStruct()
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.spb, spb)
		})
	}
}

func TestInjectDefaults(t *testing.T) {
	testCases := []struct {
		desc   string
		input  map[string]any
		output map[string]any
	}{
		{
			desc: "Inject Defaults",
			input: map[string]any{
				"property2": "test",
			},
			output: map[string]any{
				"property1": "default1",
				"property2": "test",
			},
		},
		{
			desc: "Don't Inject Defaults - Already Set",
			input: map[string]any{
				"property1": "test",
				"property2": "test",
			},
			output: map[string]any{
				"property1": "test",
				"property2": "test",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			schema := MustParseSchema(schemaWithDefaults)
			schema.InjectDefaults(tc.input)

			assert.Equal(t, tc.output, tc.input)
		})
	}
}

func TestValidateJSONSchema(t *testing.T) {
	testCases := []struct {
		desc     string
		schema   []byte
		expected error
	}{
		{
			desc: "Valid schema",
			schema: []byte(`{
                "properties": {
                    "property1": {
                        "type": "string"
                    },
                    "property2": {
                        "type": "number"
                    }
                }
            }`),
			expected: nil,
		},
		{
			desc: "Invalid schema - property of type object",
			schema: []byte(`{
                "properties": {
                    "property1": {
                        "type": "object"
                    }
                }
            }`),
			expected: errPropertiesShouldNotBeObject,
		},
		{
			desc: "Invalid schema - property with $ref",
			schema: []byte(`{
                "properties": {
                    "property1": {
                        "$ref": "#/definitions/someDefinition"
                    }
                }
            }`),
			expected: errPropertiesShouldNotBeReferences,
		},
		{
			desc: "Invalid schema - properties is not an object",
			schema: []byte(`{
                "properties": "not an object"
            }`),
			expected: errPropertiesShouldBeObject,
		},
		{
			desc: "Invalid schema - property is an array of objects",
			schema: []byte(`{
                "properties": {
                    "property1": {
                        "type": "array",
                        "items": {
                            "type": "object"
                        }
                    }
                }
            }`),
			expected: errPropertiesShouldNotBeArrayOfObjects,
		},
		{
			desc: "Valid schema - no properties",
			schema: []byte(`{
                "title": "Example Schema"
            }`),
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			err := validateSchema(tc.schema)
			assert.Equal(t, tc.expected, err)
		})
	}
}
