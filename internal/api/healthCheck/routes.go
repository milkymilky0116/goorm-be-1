package healthcheck

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("health_check")

type HealthCheckController struct{}

func (h *HealthCheckController) HealthCheck(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(context.Background(), "HealthCheck")
	defer span.End()
	w.WriteHeader(http.StatusOK)
}

func InitHealthCheckController() *HealthCheckController {
	return &HealthCheckController{}
}
