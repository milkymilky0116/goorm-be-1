package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/milkymilky0116/goorm-be-1/internal/api"
	"github.com/milkymilky0116/goorm-be-1/internal/configuration"
)

func main() {
	config, err := configuration.InitConfiguration()
	if err != nil {
		log.Fatal(err)
	}
	conn, err := pgxpool.New(context.Background(), config.Database.GetDatabaseURL())
	if err != nil {
		log.Fatal(err)
	}

	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer conn.Close()
	listener, err := net.Listen("tcp", fmt.Sprintf("[::1]:%d", config.ApplicationPort))
	if err != nil {
		log.Fatal(err)
	}
	_, err = api.Run(ctx, listener, conn)
	if err != nil {
		log.Fatal(err)
	}
}
