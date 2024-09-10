package configuration

import "fmt"

type AppConfig struct {
	ApplicationPort int            `yaml:"applicationPort"`
	Database        DatabaseConfig `yaml:"database"`
}

type DatabaseConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DbName   string `yaml:"dbName"`
}

func (d DatabaseConfig) GetDatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", d.Username, d.Password, d.Host, d.Port, d.DbName)
}
