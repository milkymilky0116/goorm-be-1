package testutil

import (
	"context"
	"time"

	"github.com/milkymilky0116/goorm-be-1/internal/configuration"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func LaunchPostgresContainer(dbConfig *configuration.DatabaseConfig) (*string, error) {
	postgresContainer, err := postgres.Run(context.Background(),
		"docker.io/postgres:latest",
		postgres.WithDatabase(dbConfig.DbName),
		postgres.WithUsername(dbConfig.Username),
		postgres.WithPassword(dbConfig.Password),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5*time.Second)))
	if err != nil {
		return nil, err
	}
	connectionString, err := postgresContainer.ConnectionString(context.Background())
	if err != nil {
		return nil, err
	}
	return &connectionString, nil
}
