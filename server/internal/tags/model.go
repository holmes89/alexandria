package tags

type Tag struct {
	ID string `json:"id"`
	DisplayName string `json:"display_name"`
}

type TaggedResource struct {
	ID string
	ResourceID string
	Type ResourceType
}

type ResourceType = string

const (
	BookResource = "book"
	PaperResource = "paper"
	LinksResource = "link"
)