package internal

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

const basePapersPath = "/papers"

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
