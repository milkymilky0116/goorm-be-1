package healthcheck

import "net/http"

type HealthCheckController struct{}

func (h *HealthCheckController) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func InitHealthCheckController() *HealthCheckController {
	return &HealthCheckController{}
}
