package internal

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"time"
)

type Link struct {
	ID          string    `json:"id" yaml:"id"`
	Link        string    `json:"link" yaml:"link"`
	DisplayName string    `json:"display_name" yaml:"display_nam"`
	IconPath    string    `json:"icon_path" yaml:"-"`
	Tags        []string  `json:"tag_ids" yaml:"tags"`
	Created     time.Time `json:"created" yaml:"created"`
}

const baseLinkPath = "/links"

func (app *App) FindLinks() ([]Link, error) {
	endpoint := fmt.Sprintf("%s/%s/", app.Endpoint, baseLinkPath)
	client := resty.New().SetAuthToken(app.Token)
	results, err := client.R().Get(endpoint)
	if err != nil {
		return nil, err
	}

	var entities []Link
	if err := json.Unmarshal(results.Body(), &entities); err != nil {
		return nil, err
	}

	return entities, nil
}
func (app *App) FindLinkByID(id string) (entity Link, err error) {
	endpoint := fmt.Sprintf("%s/%s/%s", app.Endpoint, baseLinkPath, id)
	client := resty.New().SetAuthToken(app.Token)
	results, err := client.R().Get(endpoint)
	if err != nil {
		return entity, err
	}

	if err := json.Unmarshal(results.Body(), &entity); err != nil {
		return entity, err
	}

	return entity, nil
}

func (app *App) CreateLink(url string) error {
	endpoint := fmt.Sprintf("%s/%s/", app.Endpoint, baseLinkPath)
	client := resty.New().SetAuthToken(app.Token)
	_, err := client.R().SetBody(Link{Link: url}).Post(endpoint)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) TagLink(id, tag string) error {
	endpoint := fmt.Sprintf("%s/%s/%s/tags/", app.Endpoint, baseLinkPath, id)
	client := resty.New().SetAuthToken(app.Token)
	_, err := client.R().SetBody(tagRequest{Tag: tag}).Post(endpoint)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) UntagLink(id, tag string) error {
	endpoint := fmt.Sprintf("%s/%s/%s/tags/", app.Endpoint, baseLinkPath, id)
	client := resty.New().SetAuthToken(app.Token)
	_, err := client.R().SetBody(tagRequest{Tag: tag}).Delete(endpoint)
	if err != nil {
		return err
	}

	return nil
}

type tagRequest struct {
	Tag string `json"tag"`
}
