package resource

type Resource struct {
	ExternalID  string
	DisplayName string
	Properties  map[string]any
	Links       []*Link
}
