package resource

// Link represents a documentation or related resource link
type Link struct {
	Title    string       `json:"title"`
	URL      string       `json:"url"`
	Type     LinkType     `json:"type"`
	Category LinkCategory `json:"category"`
}

func (l *Link) isValid() bool {
	// Links must have a title and url
	if l.Title == "" || l.URL == "" {
		return false
	}

	// Links must have a valid type and category
	if l.Type == LinkTypeUnknown || l.Category == LinkCategoryUnknown {
		return false
	}

	return true
}

func (l *Link) setDefault() {
	if l.Type == LinkTypeUnknown {
		l.Type = LinkTypeWebsite
	}
	if l.Category == LinkCategoryUnknown {
		l.Category = LinkCategoryDocumentation
	}
}

type LinkType string

const (
	LinkTypeUnknown  LinkType = "unknown"
	LinkTypeWebsite  LinkType = "website"
	LinkTypeEndpoint LinkType = "endpoint"
	LinkTypeDSN      LinkType = "dsn"
	LinkTypeHTTPAPI  LinkType = "http_api"
	LinkTypeGraphQL  LinkType = "graphql"
)

type LinkCategory string

const (
	LinkCategoryUnknown       LinkCategory = "unknown"
	LinkCategoryDocumentation LinkCategory = "documentation"
	LinkCategoryWebsite       LinkCategory = "website"
	LinkCategoryCommunity     LinkCategory = "community"
	LinkCategorySourceCode    LinkCategory = "source_code"
	LinkCategoryIssueTracker  LinkCategory = "issue_tracker"
	LinkCategoryAPI           LinkCategory = "api"
	LinkCategoryInternal      LinkCategory = "internal"
)
