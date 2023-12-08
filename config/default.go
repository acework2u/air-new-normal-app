package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	DBUrl  string `mapstructure:"DBUrl"`
	Port   string `mapstructure:"PORT"`
	Origin string `mapstructure:"CLIENT_ORIGIN"`
	DBName string `mapstructure:"DB_Name"`
}

func LoadConfig(path string) (config Config, err error) {

	if godotenv.Load() == nil {
		//fmt.Println("In working .env")

		return Config{
			DBUrl:  os.Getenv("DBUrl"),
			Port:   os.Getenv("PORT"),
			Origin: os.Getenv("CLIENT_ORIGIN"),
			DBName: os.Getenv("DB_Name"),
		}, err

	} else {

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

}
