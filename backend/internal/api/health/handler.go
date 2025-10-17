package health

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type response struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

func RegisterRoutes(r chi.Router) {
	r.Get("/healthz", handleHealthz)
}

func handleHealthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response{
		Status:    "ok",
		Timestamp: time.Now().UTC(),
	})
}
