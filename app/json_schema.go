package app

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/tidwall/gjson"
	"google.golang.org/protobuf/types/known/structpb"
)

//go:embed schema/generic_empty_schema.json
var GenericEmptySchema []byte

var (
	errPropertiesShouldNotBeObject         = errors.New("individual properties should not be of type 'object'")
	errPropertiesShouldNotBeArrayOfObjects = errors.New("individual properties should not be arrays of objects")
	errPropertiesShouldNotBeReferences     = errors.New("individual properties should not be references")
	errPropertiesShouldBeObject            = errors.New("properties should be of type 'object'")
)

type JSONSchema struct {
	*jsonschema.Schema
	// Raw holds the unparsed JSON schema.
	raw json.RawMessage
}

func (j *JSONSchema) toStruct() (*structpb.Struct, error) {
	if len(j.raw) == 0 {
		return structpb.NewStruct(nil)
	}

	var m map[string]any
	err := json.Unmarshal(j.raw, &m)
	if err != nil {
		return nil, fmt.Errorf("unmarshal schema: %w", err)
	}

	return structpb.NewStruct(m)
}

// injectDefaults takes the input map and injects default values from the schema.
func (j *JSONSchema) injectDefaults(input map[string]any) {
	for property, s := range j.Properties {
		// No default value, skip
		if s.Default == nil {
			continue
		}

		// Property already set, skip
		if _, ok := input[property]; ok {
			continue
		}

		// Set default value
		input[property] = s.Default
	}
}

// ParseJSONSchema parses a JSON schema and returns a JSONSchema object.
// The schema is compiled with annotations extraction enabled.
func ParseJSONSchema(schema []byte) (*JSONSchema, error) {
	if len(schema) == 0 {
		return nil, errors.New("schema is empty")
	}

	loader := jsonschema.SchemeURLLoader{
		"https": NewTempestSchemaLoader(),
	}
	compiler := jsonschema.NewCompiler()
	compiler.UseLoader(loader)

	if err := compiler.AddResource("", bytes.NewReader(schema)); err != nil {
		return nil, fmt.Errorf("load schema: %w", err)
	}

	s, err := compiler.Compile("")
	if err != nil {
		return nil, fmt.Errorf("compile schema: %w", err)
	}

	// Validate the schema to make sure it aligns with the Tempest product expectations.
	// This is done after compilation to avoid the need to re-implement the JSONSchema validation logic.
	if err := validateJSONSchema(schema); err != nil {
		return nil, fmt.Errorf("validate schema: %w", err)
	}

	return &JSONSchema{
		Schema: s,
		raw:    schema,
	}, nil
}

// MustParseJSONSchema parses a JSON schema and returns a JSONSchema object.
// It will panic if the schema cannot be parsed.
func MustParseJSONSchema(schema []byte) *JSONSchema {
	s, err := ParseJSONSchema(schema)
	if err != nil {
		panic(err)
	}

	return s
}

// validateJSONSchema validates the JSON schema against the Tempest product expectations.
// This is a client side check to assist users with an early feedback loop.
// The server will reject schemas that do not align with the product expectations.
func validateJSONSchema(schema []byte) error {
	properties := gjson.GetBytes(schema, "properties")
	if properties.Exists() {
		if properties.IsObject() {
			for _, p := range properties.Map() {
				// check that all properties are not "object" type
				if p.Get("type").String() == "object" {
					return errPropertiesShouldNotBeObject
				}

				// check that all properties are not arrays of objects
				if p.Get("type").String() == "array" {
					items := p.Get("items")
					if items.Exists() && items.Get("type").String() == "object" {
						return errPropertiesShouldNotBeArrayOfObjects
					}
				}

				// check that all properties are not references
				if p.Get("$ref").Exists() {
					return errPropertiesShouldNotBeReferences
				}
			}
		} else {
			return errPropertiesShouldBeObject
		}
	}

	return nil
}
