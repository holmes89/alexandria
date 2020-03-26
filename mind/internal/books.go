package internal

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

const baseBooksPath = "/books"

func (app *App) UploadBook(path, name string) error {
	endpoint := fmt.Sprintf("%s/%s/", app.Endpoint, baseBooksPath)
	client := resty.New()
	_, err := client.R().SetFile("file", path).
		SetFormData(map[string]string{
			"name": name,
		}).Post(endpoint)
	if err != nil {
		return err
	}

	return nil
}
