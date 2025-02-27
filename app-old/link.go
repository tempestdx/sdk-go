package app

import appv1 "github.com/tempestdx/protobuf/gen/go/tempestdx/app/v1"

//go:generate go run golang.org/x/tools/cmd/stringer -type=LinkType -linecomment
type LinkType int

const (
	LinkTypeUnspecified    LinkType = iota // unspecified
	LinkTypeDocumentation                  // documentation
	LinkTypeAdministration                 // administration
	LinkTypeSupport                        // support
	LinkTypeEndpoint                       // endpoint
	LinkTypeExternal                       // external
)

type Link struct {
	URL   string
	Title string
	Type  LinkType
}

func (l *Link) toProto() *appv1.Link {
	if l == nil {
		return nil
	}

	return &appv1.Link{
		Url:   l.URL,
		Title: l.Title,
		Type:  appv1.LinkType(l.Type),
	}
}

func linkFromProto(l *appv1.Link) *Link {
	if l == nil {
		return nil
	}

	return &Link{
		URL:   l.GetUrl(),
		Title: l.GetTitle(),
		Type:  LinkType(l.GetType()),
	}
}
