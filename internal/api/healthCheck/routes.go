package healthcheck

import (
	"net/http"

	"github.com/milkymilky0116/goorm-be-1/internal/api/middleware"
	"github.com/milkymilky0116/goorm-be-1/internal/util"
	"github.com/rs/zerolog/log"
)

type HealthCheckController struct{}

func (h *HealthCheckController) HealthCheck(w http.ResponseWriter, r *http.Request) {
	requestID := util.GetRequestID(w, r)
	log.Info().Str(middleware.REQUEST_ID, *requestID).Msg("HealthCheck invoked")
	w.WriteHeader(http.StatusOK)
	log.Info().Str(middleware.REQUEST_ID, *requestID).Msg("HealthCheck finished")
}

func InitHealthCheckController() *HealthCheckController {
	return &HealthCheckController{}
}
