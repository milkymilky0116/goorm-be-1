package main

import (
	"context"
	"fmt"
	"net"

	"github.com/milkymilky0116/goorm-be-1/internal/api"
	"github.com/milkymilky0116/goorm-be-1/internal/api/jwt"
	"github.com/milkymilky0116/goorm-be-1/internal/configuration"
	"github.com/milkymilky0116/goorm-be-1/internal/db"
	"github.com/milkymilky0116/goorm-be-1/internal/tracing"
	"github.com/rs/zerolog"

	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	config, err := configuration.InitConfiguration()
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to read configuration")
	}
	tracingProvider, err := tracing.InitTracing("goorm-class")
	defer func() {
		if err := tracingProvider.Shutdown(context.Background()); err != nil {
			log.Fatal().Err(err).Msg("Fail to shutdown tracer")
		}
	}()
	tracingProvider.Tracer("goorm-class")

	if err != nil {
		log.Fatal().Err(err).Msg("Fail to set tracer")
	}
	conn, err := db.InitDB(config.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to initialize database")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer conn.Close()
	listener, err := net.Listen("tcp", fmt.Sprintf("[::1]:%d", config.ApplicationPort))
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to bind address")
	}
	publicKey, privateKey, err := jwt.InitKey()
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to generate public/private key")
	}
	_, err = api.Run(ctx, listener, conn, publicKey, privateKey)
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to launch app")
	}
}
