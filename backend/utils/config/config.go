package config

import (
	log2 "backend/utils/log"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	AppName     string `mapstructure:"APP_NAME"`
	Environment string `mapstructure:"ENVIRONMENT"`
	ServerPort  string `mapstructure:"SERVER_PORT"`
	Timeout     int    `mapstructure:"TIMEOUT"`
	DbItem      string `mapstructure:"DB_ITEM"`
	DbAuth      string `mapstructure:"DB_AUTH"`
	DbCashFlow  string `mapstructure:"DB_CASH_FLOW"`
	//AuthUrl     string `mapstructure:"AUTH_URL"`
	//UserUrl    string `mapstructure:"ADMIN_URL"`
	Host      string `mapstructure:"HOST"`
	SecretKey string `mapstructure:"SECRET_KEY"`
}

var MyConfig *Config

func init() {
	log2.Log.Formatter = &logrus.TextFormatter{
		ForceColors:   true,
		ForceQuote:    true,
		FullTimestamp: true,
	}

	log2.Log.Println("Start Load Config")

	var err error
	if MyConfig, err = LoadConfig("."); err != nil {
		log2.Log.Errorf("Load Config Failed : " + err.Error())
		panic(err)
	}
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&config)
	return
}
