package ideas

import (
	"alexandria/internal/common"
	"time"
)

type Idea struct {
	ID          string    `json:"id"`
	DisplayName string    `json:"display_name"`
	Content     string    `json:"content"`
	Tags        []string  `json:"tag_ids"`
	Created     time.Time `json:"created"`
}

type IdeaResource struct {
	ID         string              `json:"id"`
	ResourceID string              `json:"resource_id"`
	Type       common.ResourceType `json:"type"`
}
