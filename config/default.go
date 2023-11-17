package config

import (
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	DBUrl  string `mapstructure:"DBUrl"`
	Port   string `mapstructure:"PORT"`
	Origin string `mapstructure:"CLIENT_ORIGIN"`
}

func LoadConfig(path string) (config Config, err error) {

	appMode := os.Getenv("APP_MODE")
	fileEnvType := "env"
	fileEnvName := "app"

	envPath := "."
	if len(path) > 0 {
		envPath = path
	}

	if appMode == "dev" || appMode == "Dev" {
		fileEnvType = "dev"
		envPath = "."

	}

	viper.AddConfigPath(envPath)
	viper.SetConfigType(fileEnvType)
	viper.SetConfigName(fileEnvName)

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return

}
