package api

import (
	"net/http"

	"github.com/milkymilky0116/goorm-be-1/internal/api/auth"
	healthcheck "github.com/milkymilky0116/goorm-be-1/internal/api/healthCheck"
	"github.com/milkymilky0116/goorm-be-1/internal/api/middleware"
)

func (app *Application) routes() *http.ServeMux {
	healthCheckController := healthcheck.InitHealthCheckController()
	authController := auth.InitAuthController(app.Repository, app.Validator, app.DB, app.PrivateKey, app.PublicKey)
	mux := http.NewServeMux()
	mux.Handle("GET /health_check", middleware.RequestIDMiddleware(http.HandlerFunc(healthCheckController.HealthCheck)))
	mux.Handle("POST /auth/signup", middleware.RequestIDMiddleware(http.HandlerFunc(authController.SignupController)))
	mux.Handle("POST /auth/signin", middleware.RequestIDMiddleware(http.HandlerFunc(authController.SigninController)))
	return mux
}
