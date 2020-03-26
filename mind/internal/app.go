package internal

import (
	"errors"
	"github.com/spf13/viper"
)

type App struct {
	Endpoint string
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
