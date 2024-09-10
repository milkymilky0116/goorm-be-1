package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/milkymilky0116/goorm-be-1/internal/customvalidator"
	"github.com/milkymilky0116/goorm-be-1/internal/db/repository"
	"github.com/rs/zerolog/log"
)

type Application struct {
	Addr       string
	DB         *pgxpool.Pool
	Repository *repository.Queries
	Validator  *validator.Validate
}

func Run(ctx context.Context, listener net.Listener, db *pgxpool.Pool) (*Application, error) {
	var app Application
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("password", customvalidator.PasswordValidator)
	addr := fmt.Sprintf("http://%s", listener.Addr())
	queries := repository.New(db)

	app.Addr = addr
	app.DB = db
	app.Repository = queries
	app.Validator = validate

	srv := &http.Server{
		Handler: app.routes(),
	}
	errChan := make(chan error, 1)

	go func() {
		log.Info().Msg(fmt.Sprintf("Starting Server on %s", app.Addr))
		if err := srv.Serve(listener); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()
	select {
	case <-ctx.Done():
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctxShutdown); err != nil {
			log.Fatal().Msg(fmt.Sprintf("Fail to shutting down server correctly: %v", err))
			return nil, err
		}
	case err := <-errChan:
		log.Fatal().Msg(fmt.Sprintf("Server encountered fatal error: %v", err))
		return nil, err
	}

	return &app, nil
}
