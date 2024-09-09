package configuration

import "github.com/spf13/viper"

func InitConfiguration() (*AppConfig, error) {
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
	return &config, nil
}
