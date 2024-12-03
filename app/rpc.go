package app

import (
	"context"
	"errors"
	"fmt"

	"connectrpc.com/connect"
	appv1 "github.com/tempestdx/protobuf/gen/go/tempestdx/app/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

func (a *App) Describe(_ context.Context, _ *connect.Request[appv1.DescribeRequest]) (*connect.Response[appv1.DescribeResponse], error) {
	resourceDefinitions := make([]*appv1.ResourceDefinition, 0, len(a.resourceDefinitions))

	for _, rd := range a.resourceDefinitions {
		r := &appv1.ResourceDefinition{
			Type:                 rd.Type,
			DisplayName:          rd.DisplayName,
			Description:          rd.Description,
			LifecycleStage:       appv1.LifecycleStage(rd.LifecycleStage),
			InstructionsMarkdown: rd.InstructionsMarkdown,
		}

		for _, link := range rd.Links {
			r.Links = append(r.Links, link.toProto())
		}

		if rd.healthcheck != nil {
			r.HealthcheckSupported = true
		}

		if rd.PropertiesSchema != nil {
			s, err := rd.PropertiesSchema.toStruct()
			if err != nil {
				return nil, fmt.Errorf("convert properties schema to struct: %w", err)
			}
			r.PropertiesSchema = s
		}

		if rd.create != nil {
			r.CreateSupported = true

			s, err := rd.create.schema.input.toStruct()
			if err != nil {
				return nil, fmt.Errorf("convert create input schema to struct: %w", err)
			}
			r.CreateInputSchema = s
		}

		if rd.update != nil {
			r.UpdateSupported = true

			s, err := rd.update.schema.input.toStruct()
			if err != nil {
				return nil, fmt.Errorf("convert update input schema to struct: %w", err)
			}
			r.UpdateInputSchema = s
		}

		if rd.delete != nil {
			r.DeleteSupported = true
		}

		if rd.read != nil {
			r.ReadSupported = true
		}

		if rd.list != nil {
			r.ListSupported = true
		}

		for _, a := range rd.actions {
			action := &appv1.ActionDefinition{
				Name:        a.Name,
				DisplayName: a.DisplayName,
				Description: a.Description,
			}

			if a.InputSchema != nil {
				s, err := a.InputSchema.toStruct()
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("convert action input schema to struct: %w", err))
				}
				action.InputSchema = s
			}

			if a.OutputSchema != nil {
				s, err := a.OutputSchema.toStruct()
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("convert action output schema to struct: %w", err))
				}
				action.OutputSchema = s
			}

			r.Actions = append(r.Actions, action)
		}

		resourceDefinitions = append(resourceDefinitions, r)
	}

	return connect.NewResponse(&appv1.DescribeResponse{
		ResourceDefinitions: resourceDefinitions,
	}), nil
}

func (a *App) ExecuteResourceOperation(ctx context.Context, req *connect.Request[appv1.ExecuteResourceOperationRequest]) (*connect.Response[appv1.ExecuteResourceOperationResponse], error) {
	if req.Msg.Resource == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("resource is required"))
	}

	rd, ok := a.getResourceDefinition(req.Msg.Resource.Type)
	if !ok {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("resource type %s not found", req.Msg.Resource.Type))
	}

	opReq := operationRequestFromProto(req.Msg)

	switch o := req.Msg.Operation; o {
	case appv1.ResourceOperation_RESOURCE_OPERATION_CREATE:
		op := operationForType(rd, o)
		if op == nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("operation %s not supported for resource type %s", o.String(), req.Msg.Resource.Type))
		}

		// Inject default values from the Schema into the input, then validate the input.
		op.schema.input.injectDefaults(opReq.Input)
		if err := op.schema.input.Validate(opReq.Input); err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("validate create input: %w", err))
		}

		res, err := op.fn(ctx, opReq)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("create resource: %w", err))
		}

		// Catch any validation errors before returning the resource.
		if err := op.schema.output.Validate(res.Resource.Properties); err != nil {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("validate create output: %w", err))
		}

		resource, err := res.Resource.toProto()
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("convert resource to proto: %w", err))
		}

		return connect.NewResponse(&appv1.ExecuteResourceOperationResponse{
			Resource: resource,
		}), nil
	case appv1.ResourceOperation_RESOURCE_OPERATION_UPDATE:
		op := operationForType(rd, o)
		if op == nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("operation %s not supported for resource type %s", o.String(), req.Msg.Resource.Type))
		}

		if req.Msg.Resource.ExternalId == "" {
			return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("external ID is required for update operation"))
		}

		// Inject default values from the Schema into the input, then validate the input.
		op.schema.input.injectDefaults(opReq.Input)
		if err := op.schema.input.Validate(opReq.Input); err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("validate update input: %w", err))
		}

		res, err := op.fn(ctx, opReq)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("update resource: %w", err))
		}

		// Catch any validation errors before returning the resource.
		if err := op.schema.output.Validate(res.Resource.Properties); err != nil {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("validate update output: %w", err))
		}

		resource, err := res.Resource.toProto()
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("convert resource to proto: %w", err))
		}

		return connect.NewResponse(&appv1.ExecuteResourceOperationResponse{
			Resource: resource,
		}), nil
	case appv1.ResourceOperation_RESOURCE_OPERATION_DELETE:
		op := operationForType(rd, o)
		if op == nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("operation %s not supported for resource type %s", o.String(), req.Msg.Resource.Type))
		}

		if req.Msg.Resource.ExternalId == "" {
			return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("external ID is required for delete operation"))
		}

		res, err := op.fn(ctx, opReq)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("delete resource: %w", err))
		}

		// We don't validate the output properties for a delete operation.
		resource, err := res.Resource.toProto()
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("convert resource to proto: %w", err))
		}

		return connect.NewResponse(&appv1.ExecuteResourceOperationResponse{
			Resource: resource,
		}), nil
	case appv1.ResourceOperation_RESOURCE_OPERATION_READ:
		op := operationForType(rd, o)
		if op == nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("operation %s not supported for resource type %s", o.String(), req.Msg.Resource.Type))
		}

		if req.Msg.Resource.ExternalId == "" {
			return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("external ID is required for read operation"))
		}

		res, err := op.fn(ctx, opReq)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("read resource: %w", err))
		}

		// Catch any validation errors before returning the resource.
		if err := op.schema.output.Validate(res.Resource.Properties); err != nil {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("validate read output: %w", err))
		}

		resource, err := res.Resource.toProto()
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("convert resource to proto: %w", err))
		}

		return connect.NewResponse(&appv1.ExecuteResourceOperationResponse{
			Resource: resource,
		}), nil
	default:
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("unsupported operation %s", o.String()))
	}
}

func (a *App) ListResources(ctx context.Context, req *connect.Request[appv1.ListResourcesRequest]) (*connect.Response[appv1.ListResourcesResponse], error) {
	if req.Msg.Resource == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("resource is required"))
	}

	rd, ok := a.getResourceDefinition(req.Msg.Resource.Type)
	if !ok {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("resource type %s not found", req.Msg.Resource.Type))
	}

	if rd.list == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("list operation not supported for resource type %s", req.Msg.Resource.Type))
	}

	res, err := rd.list.fn(ctx, listRequestFromProto(req.Msg))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("list resources: %w", err))
	}

	// Validate each resource before returning them.
	for _, r := range res.Resources {
		if err := rd.list.schema.output.Validate(r.Properties); err != nil {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("validate resource properties: %w", err))
		}
	}

	resources := make([]*appv1.Resource, 0, len(res.Resources))
	for _, r := range res.Resources {
		resource, err := r.toProto()
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("convert resource to proto: %w", err))
		}
		resources = append(resources, resource)
	}

	return connect.NewResponse(&appv1.ListResourcesResponse{
		Resources: resources,
		Next:      res.Next,
	}), nil
}

func (a *App) ExecuteResourceAction(ctx context.Context, req *connect.Request[appv1.ExecuteResourceActionRequest]) (*connect.Response[appv1.ExecuteResourceActionResponse], error) {
	if req.Msg.Resource == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("resource is required"))
	}

	action, ok := a.getActionDefinition(req.Msg.Resource.Type, req.Msg.Action)
	if !ok {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("action %s not found for resource type %s", req.Msg.Action, req.Msg.Resource.Type))
	}

	actionReq := actionRequestFromProto(req.Msg)

	action.InputSchema.injectDefaults(actionReq.Input)
	if err := action.InputSchema.Validate(actionReq.Input); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("validate action input: %w", err))
	}

	res, err := action.Handler(ctx, actionReq)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("execute action: %w", err))
	}

	if err := action.OutputSchema.Validate(res.Output); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("validate action output: %w", err))
	}

	o, err := structpb.NewStruct(res.Output)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("convert output to struct: %w", err))
	}

	return connect.NewResponse(&appv1.ExecuteResourceActionResponse{
		Output: o,
	}), nil
}

func (a *App) HealthCheck(ctx context.Context, req *connect.Request[appv1.HealthCheckRequest]) (*connect.Response[appv1.HealthCheckResponse], error) {
	if req.Msg.Type == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("resource type is required"))
	}

	rd, ok := a.getResourceDefinition(req.Msg.Type)
	if !ok {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("resource type %s not found", req.Msg.Type))
	}

	if rd.healthcheck == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("health check not supported for resource type %s", req.Msg.Type))
	}

	res, err := rd.healthcheck(ctx)
	if err != nil {
		return connect.NewResponse(&appv1.HealthCheckResponse{
			Status:  appv1.HealthCheckStatus_HEALTH_CHECK_STATUS_DISRUPTED,
			Message: err.Error(),
		}), nil
	}

	switch res.Status {
	case HealthCheckStatusHealthy:
		return connect.NewResponse(&appv1.HealthCheckResponse{
			Status:  appv1.HealthCheckStatus_HEALTH_CHECK_STATUS_HEALTHY,
			Message: res.Message,
		}), nil
	case HealthCheckStatusDegraded:
		return connect.NewResponse(&appv1.HealthCheckResponse{
			Status:  appv1.HealthCheckStatus_HEALTH_CHECK_STATUS_DEGRADED,
			Message: res.Message,
		}), nil
	case HealthCheckStatusDisrupted:
		return connect.NewResponse(&appv1.HealthCheckResponse{
			Status:  appv1.HealthCheckStatus_HEALTH_CHECK_STATUS_DISRUPTED,
			Message: res.Message,
		}), nil
	default:
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("unknown health check status %s", res.Status.String()))
	}
}

func operationForType(rd *ResourceDefinition, op appv1.ResourceOperation) *operation {
	switch op {
	case appv1.ResourceOperation_RESOURCE_OPERATION_CREATE:
		return rd.create
	case appv1.ResourceOperation_RESOURCE_OPERATION_UPDATE:
		return rd.update
	case appv1.ResourceOperation_RESOURCE_OPERATION_DELETE:
		return rd.delete
	case appv1.ResourceOperation_RESOURCE_OPERATION_READ:
		return rd.read
	default:
		return nil
	}
}
