package tags

type Repository interface {
	FindAllTags() ([]Tag, error)
	CreateTag(Tag) (Tag, error)
	AddResourceTag(resourceID string, resourceType ResourceType, tagName string) error
	RemoveResourceTag(resourceID string, tagName string) error
	GetTaggedResources(id string) ([]TaggedResource, error)
}
