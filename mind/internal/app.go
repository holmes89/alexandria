package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"time"
)

type App struct {
	Endpoint string
	Token string
	Config   *viper.Viper
}

func (app *App) SetConfigurationValue(key, value string) error {
	switch key {
	case "endpoint":
		app.Config.Set("endpoint", value)
	default:
		return errors.New("invalid configuration setting")
	}
	return app.Config.WriteConfig()
}

func (app *App) GetAuthToken(username, password string) error {
	endpoint := fmt.Sprintf("%s/%s/", app.Endpoint, "auth")
	client := resty.New()
	resp, err := client.SetBasicAuth(username, password).R().Get(endpoint)

	if err != nil {
		return err
	}

	if resp.IsError() {
		fmt.Printf("error: %s\n", string(resp.Body()))
		return errors.New("failed login")
	}

	var t token
	if err := json.Unmarshal(resp.Body(), &t); err != nil {
		return err
	}

	app.Config.Set("token", t.AccessToken)
	return app.Config.WriteConfig()
}

type token struct {
	Type        string    `json:"type"`
	AccessToken string    `json:"token"`
	Expires     time.Time `json:"expires"`
}