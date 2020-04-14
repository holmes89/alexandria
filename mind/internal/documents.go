package internal

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/url"
	"path"
	"time"
)

type Document struct {
	ID          string     `json:"id" yaml:"id"`
	DisplayName string     `json:"display_name" yaml:"display_name"`
	Name        string     `json:"name" yaml:"name"`
	Path        string     `json:"path" yaml:"-"`
	Type        string     `json:"type" yaml:"type"`
	Description string     `json:"description" yaml:"description"`
	Created     time.Time  `json:"created" yaml:"created"`
	Updated     *time.Time `json:"updated" yaml:"updated"`
}

const baseDocumentsPath = "/documents"

func (app *App) FindDocumentByID(id string) (*Document, error) {
	endpoint := fmt.Sprintf("%s/%s/%s", app.Endpoint, baseDocumentsPath, id)
	client := resty.New().SetAuthToken(app.Token)
	results, err := client.R().Get(endpoint)
	if err != nil {
		return nil, err
	}

	entity := &Document{}
	if err := json.Unmarshal(results.Body(), entity); err != nil {
		return nil, err
	}

	return entity, nil
}

func (app *App) FindDocuments() ([]*Document, error) {
	endpoint := fmt.Sprintf("%s/%s/", app.Endpoint, baseDocumentsPath)
	client := resty.New().SetAuthToken(app.Token)
	results, err := client.R().Get(endpoint)
	if err != nil {
		return nil, err
	}

	var entities []*Document
	if err := json.Unmarshal(results.Body(), &entities); err != nil {
		return nil, err
	}

	return entities, nil
}

func (app *App) DownloadDocument(id string) error {
	entity, err := app.FindDocumentByID(id)
	if err != nil {
		return err
	}
	p := entity.Path
	u, err := url.Parse(p)
	if err != nil {
		return err
	}
	fname := path.Base(u.Path)

	client := resty.New().SetAuthToken(app.Token)
	if _, err = client.R().SetOutput(fname).Get(p); err != nil {
		return err
	}

	return nil

}
