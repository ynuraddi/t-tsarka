package config

import "github.com/ory/viper"

type Config struct {
	LogLevel int    `mapstructure:"LOG_LEVEL"`
	LogPath  string `mapstructure:"LOG_PATH"`

	HttpHost string `mapstructure:"HTTP_HOST"`
	HttpPort string `mapstructure:"HTTP_PORT"`

	RedisHost       string `mapstructure:"REDIS_HOST"`
	RedisPort       string `mapstructure:"REDIS_PORT"`
	RedisCounterKey string `mapstructure:"REDIS_COUNTER_KEY"`

	PostgresHost   string `mapstructure:"POSTGRES_HOST"`
	PostgresPort   string `mapstructure:"POSTGRES_PORT"`
	PostgresUser   string `mapstructure:"POSTGRES_USER"`
	PostgresPass   string `mapstructure:"POSTGRES_PASS"`
	PostgresDBName string `mapstructure:"POSTGRES_DBNAME"`
}

func LoadConfig(path string) (config *Config, err error) {
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
