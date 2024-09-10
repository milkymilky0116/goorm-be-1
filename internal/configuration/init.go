package configuration

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func InitConfiguration() (*AppConfig, error) {
	log.Info().Msg("Starting to load configuration")
	var config AppConfig
	viper.AddConfigPath("configuration")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	log.Info().Msg("Finish to load configuration")
	return &config, nil
}
