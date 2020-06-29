package rest

import (
	"github.com/dnozdrin/detask/internal/app/log"
	"net/http"
)

// HealthCheck provides a method to check the web service status
type HealthCheck struct {
	log  log.Logger
	resp *responder
}

// HealthCheck constructor
func NewHealthCheck(logger log.Logger) *HealthCheck {
	return &HealthCheck{
		log:  logger,
		resp: &responder{log: logger},
	}
}

// Status will response with status Ok
func (hc HealthCheck) Status(w http.ResponseWriter, _ *http.Request) {
	hc.log.Info("Health check is OK")
	hc.resp.respondJSON(w, http.StatusOK, struct{ Status string `json:"status"`}{"OK"})
}
