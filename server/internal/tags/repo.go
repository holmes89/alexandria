package tags

import "alexandria/internal/common"

type Repository interface {
	FindAllTags() ([]Tag, error)
	CreateTag(Tag) (Tag, error)
	AddResourceTag(resourceID string, resourceType common.ResourceType, tagName string) error
	RemoveResourceTag(resourceID string, tagName string) error
	GetTaggedResources(id string) ([]TaggedResource, error)
}
