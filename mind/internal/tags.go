package internal

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
)

type Tag struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	TagColor    string `json:"color"`
}

const baseTagPath = "/tags"

func (app *App) FindTags() ([]Tag, error) {
	endpoint := fmt.Sprintf("%s/%s/", app.Endpoint, baseTagPath)
	client := resty.New()
	results, err := client.R().Get(endpoint)
	if err != nil {
		return nil, err
	}

	var entities []Tag
	if err := json.Unmarshal(results.Body(), &entities); err != nil {
		return nil, err
	}

	return entities, nil
}
