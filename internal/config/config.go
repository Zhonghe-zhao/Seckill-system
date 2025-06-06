package config

import "github.com/spf13/viper"

type Config struct {
	HTTPServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
	RedisAddress      string `mapstructure:"REDIS_ADDRESS"`
	DBSource          string `mapstructure:"DB_SOURCE"`
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
