package links

import (
	"time"
)

type Link struct {
	ID          string    `json:"id"`
	Link        string    `json:"link"`
	DisplayName string    `json:"display_name"`
	IconPath    string    `json:"icon_path"`
	Tags        []string  `json:"tag_ids"`
	Created     time.Time `json:"created"`
}
