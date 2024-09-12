package configuration

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func InitConfiguration() (*AppConfig, error) {
	log.Info().Msg("Starting to load configuration")
	var config AppConfig
	env := os.Getenv("APP_ENV")
	log.Info().Msg(env)
	if env == "" {
		env = "dev"
	}
	viper.AddConfigPath("configuration")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config." + env)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	log.Info().Msgf("Finish to load configuration in %v", env)
	return &config, nil
}
