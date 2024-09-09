package api

import (
	"net/http"

	"github.com/milkymilky0116/goorm-be-1/internal/api/auth"
	healthcheck "github.com/milkymilky0116/goorm-be-1/internal/api/healthCheck"
)

func (app *Application) routes() *http.ServeMux {
	healthCheckController := healthcheck.InitHealthCheckController()
	authController := auth.InitAuthController(app.Repository)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health_check", healthCheckController.HealthCheck)
	mux.HandleFunc("POST /auth/signup", authController.SignupController)
	return mux
}
