package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/milkymilky0116/goorm-be-1/internal/configuration"
	"github.com/milkymilky0116/goorm-be-1/internal/util"
	"github.com/rs/zerolog/log"
)

func InitDB(config configuration.DatabaseConfig) (*pgxpool.Pool, error) {
	log.Info().Msg("Starting to connect to database")
	conn, err := pgxpool.New(context.Background(), config.GetDatabaseURL())
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to establish database connection")
		return nil, err
	}

	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to connect database")
		return nil, err
	}
	log.Info().Msg("Finish to connect to database")

	err = util.MigrateDB(config.GetDatabaseURL())
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to migrate database")
		return nil, err

	}
	return conn, nil
}
