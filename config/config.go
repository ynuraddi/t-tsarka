package config

import "github.com/ory/viper"

type Config struct {
	LogLevel int    `mapstructure:"LOG_LEVEL"`
	LogPath  string `mapstructure:"LOG_PATH"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
