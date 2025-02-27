package resource

import (
	"context"

	"github.com/tempestdx/sdk-go/jsonschema"
)

type EnvironmentVariableType string

const (
	ENVIRONMENT_VARIABLE_TYPE_VAR         EnvironmentVariableType = "variable"
	ENVIRONMENT_VARIABLE_TYPE_SECRET      EnvironmentVariableType = "secret"
	ENVIRONMENT_VARIABLE_TYPE_CERTIFICATE EnvironmentVariableType = "certificate"
	ENVIRONMENT_VARIABLE_TYPE_PRIVATE_KEY EnvironmentVariableType = "private_key"
	ENVIRONMENT_VARIABLE_TYPE_PUBLIC_KEY  EnvironmentVariableType = "public_key"
)

type EnvironmentVariable struct {
	Key   string
	Value string
	Type  EnvironmentVariableType
}

// OperationRequest contains the input data for an operation on a resource.
type OperationRequest struct {
	ResourceDefinition *ResourceDefinition
	// Metadata contains information about the Project and User making the request.
	// This metadata does not contain information about the Resource being operated on.
	Metadata *Metadata
	// Resource is the resource being operated on, and contains the ExternalID of the resource,
	// as well as the properties at the time of the request.
	Resource *Resource
	// Input contains the input data for the request, after it has been validated against the schema.
	// Default values have already been applied to missing input properties.
	Args map[string]any
	// Environment contains the environment variables that are available to the operation.
	Environment map[string]EnvironmentVariable
}

// OperationResponse contains the output data for an operation on a resource.
type OperationResponse struct {
	// Resource contains the properties of the resource after the operation has been performed.
	Resource *Resource
	// Error is the error that occurred during the operation, if any.
	Error error
	// Message is the message of the operation.
	Message string
	// ResultProperties contains the properties of the resource after the operation has been performed.
	ResultProperties map[string]any
}

type OperationFunc func(ctx *context.Context, req *OperationRequest) (*OperationResponse, error)

type operation struct {
	name string
	args *jsonschema.Schema

	pre  OperationFunc
	fn   OperationFunc
	post OperationFunc

	actionConfig *ActionConfig
	canonicalOps []CanonicalOperation
}

func (o *operation) IsCanonical() bool {
	return len(o.canonicalOps) > 0
}

func (o *operation) IsAction() bool {
	return o.actionConfig != nil
}

type ActionConfig struct {
	Title                string
	Description          string
	RequiresConfirmation bool
}

func newOperation(name string, fn OperationFunc, opts ...operationOption) *operation {
	op := &operation{
		name:         name,
		fn:           fn,
		canonicalOps: []CanonicalOperation{},
	}

	for _, opt := range opts {
		if err := opt(op); err != nil {
			return op
		}
	}

	if !op.IsCanonical() && !op.IsAction() {
		return nil
	}

	return op
}

type operationOption func(*operation) error

var operationOptions = struct {
	EnableAction func(*ActionConfig) operationOption
	On           func(canonicals ...CanonicalOperation) operationOption
	WithPre      func(OperationFunc) operationOption
	WithPost     func(OperationFunc) operationOption
	WithArgs     func(*jsonschema.Schema) operationOption
}{
	EnableAction: enableAction,
	On:           on,
	WithPre:      withPre,
	WithPost:     withPost,
	WithArgs:     withArgs,
}

var (
	Op               = operationOptions
	OperationOptions = operationOptions
)

func enableAction(config *ActionConfig) operationOption {
	return func(op *operation) error {
		op.actionConfig = config
		return nil
	}
}

func withPre(fn OperationFunc) operationOption {
	return func(op *operation) error {
		op.pre = fn
		return nil
	}
}

func withPost(fn OperationFunc) operationOption {
	return func(op *operation) error {
		op.post = fn
		return nil
	}
}

func withArgs(args *jsonschema.Schema) operationOption {
	return func(op *operation) error {
		op.args = args
		return nil
	}
}

func on(canonicals ...CanonicalOperation) operationOption {
	return func(op *operation) error {
		op.canonicalOps = append(op.canonicalOps, canonicals...)
		return nil
	}
}
