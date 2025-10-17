package movies

import (
	"encoding/json"
	"net/http"
	"time"

	service "github.com/leak-streaming/leak-streaming/backend/internal/service/movies"
)

type ListHandler struct {
	service *service.Service
}

func NewListHandler(service *service.Service) *ListHandler {
	return &ListHandler{service: service}
}

func (h *ListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.service == nil {
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)
		return
	}

	movies, err := h.service.ListMovies(r.Context())
	if err != nil {
		http.Error(w, "failed to list movies", http.StatusInternalServerError)
		return
	}

	response := make([]movieListItem, 0, len(movies))
	for _, movie := range movies {
		item := movieListItem{
			ID:        movie.ID,
			Slug:      movie.Slug,
			Title:     movie.Title,
			Synopsis:  movie.Synopsis,
			PosterURL: movie.PosterURL,
			IsVisible: movie.IsVisible,
		}
		if !movie.AvailabilityStart.IsZero() {
			item.AvailabilityStart = movie.AvailabilityStart.Format(time.RFC3339)
		}
		if !movie.AvailabilityEnd.IsZero() {
			item.AvailabilityEnd = movie.AvailabilityEnd.Format(time.RFC3339)
		}
		response = append(response, item)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

type movieListItem struct {
	ID                string `json:"id"`
	Slug              string `json:"slug"`
	Title             string `json:"title"`
	Synopsis          string `json:"synopsis"`
	PosterURL         string `json:"posterUrl"`
	AvailabilityStart string `json:"availabilityStart"`
	AvailabilityEnd   string `json:"availabilityEnd"`
	IsVisible         bool   `json:"isVisible"`
}
