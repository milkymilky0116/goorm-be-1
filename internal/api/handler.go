package api

import (
	"net/http"

	"github.com/milkymilky0116/goorm-be-1/internal/api/auth"
	healthcheck "github.com/milkymilky0116/goorm-be-1/internal/api/healthCheck"
	"github.com/milkymilky0116/goorm-be-1/internal/api/middleware"
)

func (app *Application) routes() *http.ServeMux {
	healthCheckController := healthcheck.InitHealthCheckController()
	authController := auth.InitAuthController(app.Repository, app.Validator)
	mux := http.NewServeMux()
	mux.Handle("GET /health_check", middleware.RequestIDMiddleware(http.HandlerFunc(healthCheckController.HealthCheck)))
	mux.Handle("POST /auth/signup", middleware.RequestIDMiddleware(http.HandlerFunc(authController.SignupController)))
	return mux
}
