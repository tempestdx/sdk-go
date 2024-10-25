package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	appv1 "github.com/tempestdx/protobuf/gen/go/tempestdx/app/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestOperationRequestFromProto(t *testing.T) {
	testCases := []struct {
		desc    string
		opReq   *OperationRequest
		opReqpb *appv1.ExecuteResourceOperationRequest
	}{
		{
			desc: "OK",
			opReq: &OperationRequest{
				Resource: &Resource{
					ExternalID:  "external-id",
					DisplayName: "display-name",
					Type:        "type",
					Links:       []*Link{},
					Properties:  map[string]any{},
				},
				Input: map[string]any{
					"key": "value",
				},
				Environment: map[string]EnvironmentVariable{
					"TELEMETRY_ENDPOINT": {
						Key:   "TELEMETRY_ENDPOINT",
						Value: "https://telemetry.example.com",
						Type:  ENVIRONMENT_VARIABLE_TYPE_VAR,
					},
					"DB_PASSWORD": {
						Key:   "DB_PASSWORD",
						Value: "so-secret",
						Type:  ENVIRONMENT_VARIABLE_TYPE_SECRET,
					},
				},
			},
			opReqpb: &appv1.ExecuteResourceOperationRequest{
				Resource: &appv1.Resource{
					ExternalId:  "external-id",
					DisplayName: "display-name",
					Type:        "type",
				},
				Metadata:  nil,
				Operation: appv1.ResourceOperation_RESOURCE_OPERATION_CREATE,
				Input:     mustNewStruct(map[string]any{"key": "value"}),
				EnvironmentVariables: []*appv1.EnvironmentVariable{
					{
						Key:   "TELEMETRY_ENDPOINT",
						Value: "https://telemetry.example.com",
						Type:  appv1.EnvironmentVariableType_ENVIRONMENT_VARIABLE_TYPE_VAR,
					},
					{
						Key:   "DB_PASSWORD",
						Value: "so-secret",
						Type:  appv1.EnvironmentVariableType_ENVIRONMENT_VARIABLE_TYPE_SECRET,
					},
				},
			},
		},
		{
			desc:    "OK - nil",
			opReq:   nil,
			opReqpb: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.opReq, operationRequestFromProto(tc.opReqpb))
		})
	}
}

func mustNewStruct(i map[string]any) *structpb.Struct {
	s, err := structpb.NewStruct(i)
	if err != nil {
		panic(err)
	}
	return s
}
