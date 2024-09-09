package api

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	Addr string
	DB   *pgxpool.Pool
}

func Run(ctx context.Context, listener net.Listener, db *pgxpool.Pool) (*Application, error) {
	var app Application
	addr := fmt.Sprintf("http://%s", listener.Addr())
	srv := &http.Server{
		Handler: app.routes(),
	}
	go func() {
		if err := srv.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	<-ctx.Done()
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctxShutdown); err != nil {
		return nil, err
	}
	return &Application{Addr: addr, DB: db}, nil
}
