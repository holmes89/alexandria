package internal

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/url"
	"path"
)

const basePapersPath = "/papers"

func (app *App) FindPaperByID(id string) (*Document, error) {
	endpoint := fmt.Sprintf("%s/%s/%s", app.Endpoint, basePapersPath, id)
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

func (app *App) FindPapers() ([]*Document, error) {
	endpoint := fmt.Sprintf("%s/%s/", app.Endpoint, basePapersPath)
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

func (app *App) DownloadPaper(id string) error {
	entity, err := app.FindPaperByID(id)
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

func (app *App) UploadPapers(path, name string) error {
	endpoint := fmt.Sprintf("%s/%s/", app.Endpoint, basePapersPath)
	client := resty.New().SetAuthToken(app.Token)
	_, err := client.R().SetFile("file", path).
		SetFormData(map[string]string{
			"name": name,
		}).Post(endpoint)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) TagPaper(id, tag string) error {
	endpoint := fmt.Sprintf("%s/%s/%s/tags/", app.Endpoint, basePapersPath, id)
	client := resty.New().SetAuthToken(app.Token)
	_, err := client.R().SetBody(tagRequest{Tag: tag}).Post(endpoint)
	if err != nil {
		return err
	}

	return nil
}
