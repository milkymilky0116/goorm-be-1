package healthcheck

import (
	"context"
	"net/http"

	"github.com/milkymilky0116/goorm-be-1/internal/api/middleware"
	"github.com/milkymilky0116/goorm-be-1/internal/util"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var tracer = otel.Tracer("health_check")

type HealthCheckController struct{}

func (h *HealthCheckController) HealthCheck(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(context.Background(), "HealthCheck")
	requestID := util.GetRequestID(w, r)
	spanID := span.SpanContext().SpanID().String()
	span.SetAttributes(attribute.String(middleware.REQUEST_ID, *requestID))
	defer span.End()
	util.LogInfo(span, *requestID, spanID, "Health Check Complete")
	w.WriteHeader(http.StatusOK)
}

func InitHealthCheckController() *HealthCheckController {
	return &HealthCheckController{}
}
