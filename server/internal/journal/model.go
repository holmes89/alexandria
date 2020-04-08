package journal

import "time"

type Entry struct {
	ID      string    `json:"id"`
	Content string    `json:"content"`
	Created time.Time `json:"created"`
}
