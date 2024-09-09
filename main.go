package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

type Application struct {
	Addr string
}

func run(ctx context.Context, listener net.Listener) (*Application, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health_check", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	addr := fmt.Sprintf("http://%s", listener.Addr())
	srv := &http.Server{
		Handler: mux,
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
	return &Application{Addr: addr}, nil
}
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	listener, err := net.Listen("tcp", "[::1]:8080")
	if err != nil {
		log.Fatal(err)
	}
	_, err = run(ctx, listener)
	if err != nil {
		log.Fatal(err)
	}
}
