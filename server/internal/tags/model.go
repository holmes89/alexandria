package tags

import "alexandria/internal/common"

type Tag struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	TagColor    Color  `json:"color"`
}

type TaggedResource struct {
	ID          string              `json:"-"`
	ResourceID  string              `json:"id"`
	DisplayName string              `json:"display_name"`
	Type        common.ResourceType `json:"type"`
}
