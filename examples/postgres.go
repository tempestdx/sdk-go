package example

import (
	"fmt"

	"github.com/tempestdx/sdk-go/jsonschema"
	"github.com/tempestdx/sdk-go/resource"
)

func NewPostgresResource() (*resource.Resource, error) {
	// Define properties schema
	properties := &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"instance_name": {
				Type:        "string",
				Pattern:     "^[a-z][a-z0-9-]+[a-z0-9]$",
				MinLength:   3,
				MaxLength:   63,
				Description: "Name of the PostgreSQL instance",
			},
			"version": {
				Type:        "string",
				Enum:        []interface{}{"11", "12", "13", "14", "15"},
				Description: "PostgreSQL version",
			},
			"storage_gb": {
				Type:        "integer",
				Minimum:     10,
				Maximum:     1000,
				Description: "Storage size in GB",
			},
			"region": {
				Type:        "string",
				Description: "Region where the database will be deployed",
			},
		},
		Required: []string{"instance_name", "version", "storage_gb", "region"},
	}

	// Create the resource
	r, err := resource.New(resource.Config{
		Name:           "postgres_database",
		LifecycleStage: resource.LifecycleStageGA,
		Properties:     properties,
	})
	if err != nil {
		return nil, err
	}

	// Add documentation
	r.WithLinks(
		resource.Link{Title: "Documentation", URL: "https://docs.example.com/postgres"},
		resource.Link{Title: "Support", URL: "https://support.example.com"},
	)

	r.WithInstructions(`
## PostgreSQL Database
This resource provisions and manages PostgreSQL databases.

### Connection Information
After installation, use the connection string from the output to connect to your database.
`)

	// Register install operation
	r.Install(handleInstall)

	// Register destroy operation with confirmation
	r.RegisterOperation("_destroy",
		handleDestroy,
		resource.On(resource.Destroy),
		resource.EnableAction(&resource.ActionConfig{
			Title:                "Delete Database",
			Description:          "Permanently deletes the database instance",
			RequiresConfirmation: true,
		}),
	)

	// Register backup operation as a custom action
	r.RegisterOperation("backup",
		handleBackup,
		resource.EnableAction(&resource.ActionConfig{
			Title:       "Create Backup",
			Description: "Creates a point-in-time backup",
		}),
		resource.WithPre(validateBackup),
		resource.WithPost(notifyBackupComplete),
	)

	return r, nil
}

// Operation handlers
func handleInstall(ctx *resource.Context) (interface{}, error) {
	// Parse config
	var config struct {
		InstanceName string `json:"instance_name"`
		Version      string `json:"version"`
		StorageGB    int    `json:"storage_gb"`
		Region       string `json:"region"`
	}

	// In a real implementation, you would:
	// 1. Call your cloud provider's API to create the database
	// 2. Wait for it to be provisioned
	// 3. Return connection details

	return map[string]interface{}{
		"status":            "creating",
		"connection_string": fmt.Sprintf("postgres://admin:password@%s.postgres.example.com:5432/postgres", config.InstanceName),
		"host":              fmt.Sprintf("%s.postgres.example.com", config.InstanceName),
		"port":              5432,
	}, nil
}

func handleDestroy(ctx *resource.Context) (interface{}, error) {
	// In a real implementation, you would:
	// 1. Call your cloud provider's API to delete the database
	// 2. Return status information

	return map[string]interface{}{
		"status": "deleting",
	}, nil
}

func handleBackup(ctx *resource.Context) (interface{}, error) {
	// Create a backup
	return map[string]interface{}{
		"backup_id": "bkp-123456",
		"status":    "creating",
	}, nil
}

func validateBackup(ctx *resource.Context) (interface{}, error) {
	// Validate the database is in a state where it can be backed up
	return nil, nil
}

func notifyBackupComplete(ctx *resource.Context) (interface{}, error) {
	// Send notification that backup is complete
	return nil, nil
}

// package example

// import (
// 	"github.com/tempestdx/sdk-go/app"
// 	"github.com/tempestdx/sdk-go/resource"
// )

// func NewPostgresResource() (*resource.V2, error) {
//     config := resource.Config{
//         Name:           "postgres_database",
//         LifecycleStage: resource.LifecycleStageGA,
//         Properties:     postgresPropertiesSchema(),
//     }

//     resource, err := app.NewResourceV2(config)
//     if err != nil {
//         return nil, err
//     }

//     // Define input/output schemas
//     deleteInput := &app.JSONSchema{
//         Type: "object",
//         Properties: map[string]*app.Schema{
//             "force": {
//                 Type:        "boolean",
//                 Description: "Force deletion even if database is in use",
//             },
//             "backup": {
//                 Type:        "boolean",
//                 Description: "Create final backup before deletion",
//             },
//         },
//     }

//     deleteJSONOutput := &app.JSONSchema{
//         Type: "object",
//         Properties: map[string]*app.JSONSchema{
//             "status": {
//                 Type: "string",
//                 Enum: []string{"deleting", "deleted"},
//             },
//             "backup_id": {
//                 Type:        "string",
//                 Description: "ID of final backup if requested",
//             },
//         },
//     }

//     // Register operation with input/output schemas
//     resource.RegisterOperation("deletePostgresSchema",
//         handleDelete,
//         app.On(app.CanonicalOperationDestroy, 100), // Using Destroy instead of Delete
// 			  app.WithAuth("apikey"),
//         app.EnableAction(&app.ActionConfig{
//             title:                "Delete Database",
//             description:          "Permanently deletes the database instance",
//             requiresConfirmation: true,
//         }),
//         app.WithInput(deleteInput),
//         app.WithOutput(deleteOutput),
// 		app.WithPre(),
// 		app.WithPost(),
// 		resource.WithInstructions(""),
//     )

// 	resource.Create()

//     return resource, nil
// }

// // Operation handler with typed input/output
// func handleDelete(ctx *app.Context) (interface{}, error) {
//     var input struct {
//         Force  bool `json:"force"`
//         Backup bool `json:"backup"`
//     }
//     if err := ctx.Config.Decode(&input); err != nil {
//         return nil, err
//     }

//     // Implementation...

//     return map[string]interface{}{
//         "status":    "deleting",
//         "backup_id": "bkp-final-123",
//     }, nil
// }
