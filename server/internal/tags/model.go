package tags

type Tag struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	TagColor    Color  `json:"color"`
}

type TaggedResource struct {
	ID          string       `json:"-"`
	ResourceID  string       `json:"id"`
	DisplayName string       `json:"display_name"`
	Type        ResourceType `json:"type"`
}

type ResourceType = string

const (
	BookResource  = "book"
	PaperResource = "paper"
	LinksResource = "link"
)
