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
	client := resty.New().SetAuthToken(app.Token)
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

func (app *App) TagMap() (map[string]string, error) {
	tags, err := app.FindTags()
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	for _, t := range tags {
		m[t.ID] = t.DisplayName
	}

	return m, nil
}

func (app *App) CreateTag(tag string) error {
	endpoint := fmt.Sprintf("%s/%s/", app.Endpoint, baseTagPath)
	client := resty.New().SetAuthToken(app.Token)
	_, err := client.R().SetBody(Tag{DisplayName: tag}).Post(endpoint)
	if err != nil {
		return err
	}

	return nil
}
