package app

import (
	appv1 "github.com/tempestdx/protobuf/gen/go/tempestdx/app/v1"
)

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

func metadataFromProto(m *appv1.Metadata) *Metadata {
	if m == nil {
		return nil
	}

	owners := make([]Owner, 0, len(m.Owners))
	for _, owner := range m.GetOwners() {
		owners = append(owners, ownerFromProto(owner))
	}

	return &Metadata{
		ProjectID:   m.GetProjectId(),
		ProjectName: m.GetProjectName(),
		Owners:      owners,
		Author:      ownerFromProto(m.GetAuthor()),
	}
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

func ownerFromProto(o *appv1.Owner) Owner {
	if o == nil {
		return Owner{}
	}

	var t OwnerType
	switch o.GetType() {
	case appv1.OwnerType_OWNER_TYPE_USER:
		t = OwnerTypeUser
	case appv1.OwnerType_OWNER_TYPE_TEAM:
		t = OwnerTypeTeam
	}

	return Owner{
		Email: o.Email,
		Name:  o.Name,
		Type:  t,
	}
}
