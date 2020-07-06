package rest

import (
	"encoding/json"
	log "github.com/dnozdrin/detask/internal/app/log"
	"github.com/dnozdrin/detask/internal/domain/services"
	"net/http"
	"strconv"
)

type responder struct {
	log log.Logger
}

// respond makes the response with the provided payload
func (r responder) respond(w http.ResponseWriter, status int, payload string) {
	w.WriteHeader(status)
	if _, err := w.Write([]byte(payload)); err != nil {
		r.log.Errorf("error while writing response: %v", err)
	}
}

// respondJSON makes the response with payload as JSON format
func (r responder) respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		r.log.Errorf("error while encoding JSON response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
	if _, err = w.Write(response); err != nil {
		r.log.Errorf("error while writing JSON response: %v", err)
	}
}

// respondError makes the error response with payload as JSON format
func (r responder) respondError(w http.ResponseWriter, code int, message string) {
	r.respondJSON(w, code, map[string]string{"error": message})
}

//parseFilter fetches filter parameter from the request query and parses
//it into services.Demand
func parseFilter(r *http.Request, demand services.Demand) error {
	for k, v := range r.URL.Query() {
		if k != "id" {
			intVal, err := strconv.Atoi(v[0])
			if err != nil {
				return err
			}
			if err = demand.Add(k, uint(intVal)); err != nil {
				return err
			}
		}
	}

	return nil
}
