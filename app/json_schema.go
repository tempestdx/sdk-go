package app

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"google.golang.org/protobuf/types/known/structpb"
)

//go:embed schema/generic_empty_schema.json
var GenericEmptySchema []byte

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

	// The package default compiler does not extract annotations.
	compiler := jsonschema.NewCompiler()
	compiler.ExtractAnnotations = true

	if err := compiler.AddResource("", bytes.NewReader(schema)); err != nil {
		return nil, fmt.Errorf("load schema: %w", err)
	}

	s, err := compiler.Compile("")
	if err != nil {
		return nil, fmt.Errorf("compile schema: %w", err)
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
