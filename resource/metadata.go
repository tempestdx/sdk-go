package resource

type Metadata struct {
	// ProjectID is the ID of the Tempest Project. This guaranteed to be unique.
	ProjectID string
	// ProjectName is the user defined name of the Tempest Project.
	ProjectName string
	// Owners are the user(s) who created or own the Project.
	Owners []Owner
	// Author is the user or team who created the Project.
	Author Owner
}

type OwnerType string

const (
	OwnerTypeUser OwnerType = "user"
	OwnerTypeTeam OwnerType = "team"
)

type Owner struct {
	Email string
	Name  string
	Type  OwnerType
}
