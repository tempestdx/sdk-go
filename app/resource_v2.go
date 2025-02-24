package app

import "fmt"

type statefulMetadata struct {
	version        string
	published      bool
	publishedID    string
	organizationID string
	appID          string
}

type ResourceConfig struct {
	Name           string
	LifecycleStage LifecycleStage
	Properties     *JSONSchema
}

type ResourceOpt func(*resourceV2) error

type resourceV2 struct {
	operations     map[string]*OperationV2
	name           string
	LifecycleStage LifecycleStage
	Links          []Link
	Properties     *JSONSchema

	instructionsMarkdown string
	canonicalMap         map[CanonicalOperation]*OperationV2

	// If the resource is published, this will be set to the published state upon connecting with the API.
	publishedState *statefulMetadata
}

// NewResourceV2 creates a new ResourceV2. A constructor is used here as we validate the resource implemetation
// against your Tempest instance to see if this is an existing or new resource.
func NewResourceV2(config ResourceConfig, opts ...ResourceOpt) (*resourceV2, error) {
	r := &resourceV2{
		name:           config.Name,
		LifecycleStage: config.LifecycleStage,
		Properties:     config.Properties,
	}

	for _, opt := range opts {
		if err := opt(r); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func WithLinks(links ...Link) ResourceOpt {
	return func(r *resourceV2) error {
		r.Links = links
		return nil
	}
}

func WithInstructions(markdown string) ResourceOpt {
	return func(r *resourceV2) error {
		r.instructionsMarkdown = markdown
		return nil
	}
}

type OperationOpt func(*OperationV2) error

type OperationV2 struct {
	name string

	input  *JSONSchema
	output *JSONSchema

	pre  OperationFunc
	fn   OperationFunc
	post OperationFunc

	actionConfig *ActionConfig

	isCanonical bool
	canonical   CanonicalOperation
}

type ActionConfig struct {
	title                string
	description          string
	requiresConfirmation bool
}

// CanonicalOperation is an enum representing the canonical operations that a resource can perform
// that are part of Tempest's core management capabilities.
type CanonicalOperation int

const (
	CanonicalOperationInstall CanonicalOperation = iota
	CanonicalOperationUninstall
	CanonicalOperationUpgrade
	CanonicalOperationRollback
	CanonicalOperationDestroy
	CanonicalOperationConfigure
	CanonicalOperationRead
	CanonicalOperationWrite
	CanonicalOperationList
	CanonicalOperationGet
	CanonicalOperationUpdate
	CanonicalOperationDelete
	CanonicalOperationTest
	CanonicalOperationHealthz
	CanonicalOperationSync
)

func isCanonicalOperation(op CanonicalOperation) bool {
	return op >= CanonicalOperationInstall && op <= CanonicalOperationSync
}

// setCanonicalOperation sets the canonical operation for the resource. These can only be set once.
func (r *resourceV2) setCanonicalOperation(op CanonicalOperation, operation *OperationV2) error {
	if !isCanonicalOperation(op) {
		return fmt.Errorf("invalid canonical operation: %d", op)
	}

	if r.canonicalMap[op] != nil {
		return fmt.Errorf("canonical operation already set: %d", op)
	}

	r.canonicalMap[op] = operation
	return nil
}

func (r *resourceV2) RegisterOperation(name string, fn OperationFunc, opts ...OperationOpt) *resourceV2 {
	op := &OperationV2{
		name: name,
		fn:   fn,
	}
	for _, opt := range opts {
		if err := opt(op); err != nil {
			return r
		}
	}

	if isCanonicalOperation(op.canonical) {
		op.isCanonical = true
	}

	// If operation is neither canonical or an action, it won't be called
	// In this case, we don't register and return early
	if !op.isCanonical && op.actionConfig == nil {
		return r
	}

	r.operations[name] = op

	return r
}

func EnableAction(config *ActionConfig) OperationOpt {
	return func(op *OperationV2) error {
		op.actionConfig = config
		return nil
	}
}

func WithPre(fn OperationFunc) OperationOpt {
	return func(op *OperationV2) error {
		op.pre = fn
		return nil
	}
}

func WithPost(fn OperationFunc) OperationOpt {
	return func(op *OperationV2) error {
		op.post = fn
		return nil
	}
}

func On(canonical CanonicalOperation) OperationOpt {
	return func(op *OperationV2) error {
		op.canonical = canonical
		return nil
	}
}

func (r *resourceV2) Install(fn OperationFunc) *resourceV2 {
	return r.RegisterOperation("install", fn, On(CanonicalOperationInstall))
}

func (r *resourceV2) Uninstall(fn OperationFunc) *resourceV2 {
	return r.RegisterOperation("uninstall", fn, On(CanonicalOperationUninstall))
}
