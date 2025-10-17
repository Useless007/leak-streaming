package movies

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	domain "github.com/leak-streaming/leak-streaming/backend/internal/domain/movies"
	service "github.com/leak-streaming/leak-streaming/backend/internal/service/movies"
)

type DetailsHandler struct {
	service *service.Service
}

func NewDetailsHandler(service *service.Service) *DetailsHandler {
	return &DetailsHandler{service: service}
}

func (h *DetailsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.service == nil {
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)
		return
	}

	slug := chi.URLParam(r, "slug")
	if slug == "" {
		http.Error(w, "missing slug", http.StatusBadRequest)
		return
	}

	movie, err := h.service.GetMovie(r.Context(), slug)
	if err != nil {
		if errors.Is(err, service.ErrMovieNotFound) {
			http.Error(w, "movie not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to load movie", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movieResponseFromDomain(movie))
}

type movieResponse struct {
	ID                string         `json:"id"`
	Slug              string         `json:"slug"`
	Title             string         `json:"title"`
	Synopsis          string         `json:"synopsis"`
	PosterURL         string         `json:"posterUrl"`
	AvailabilityStart string         `json:"availabilityStart"`
	AvailabilityEnd   string         `json:"availabilityEnd"`
	IsVisible         bool           `json:"isVisible"`
	Captions          []captionModel `json:"captions"`
}

type captionModel struct {
	LanguageCode string `json:"languageCode"`
	Label        string `json:"label"`
	CaptionURL   string `json:"captionUrl"`
}

func movieResponseFromDomain(movie domain.Movie) movieResponse {
	captions := make([]captionModel, 0, len(movie.Captions))
	for _, c := range movie.Captions {
		captions = append(captions, captionModel{
			LanguageCode: c.LanguageCode,
			Label:        c.Label,
			CaptionURL:   c.CaptionURL,
		})
	}

	response := movieResponse{
		ID:        movie.ID,
		Slug:      movie.Slug,
		Title:     movie.Title,
		Synopsis:  movie.Synopsis,
		PosterURL: movie.PosterURL,
		IsVisible: movie.IsVisible,
		Captions:  captions,
	}
	if !movie.AvailabilityStart.IsZero() {
		response.AvailabilityStart = movie.AvailabilityStart.Format(time.RFC3339)
	}
	if !movie.AvailabilityEnd.IsZero() {
		response.AvailabilityEnd = movie.AvailabilityEnd.Format(time.RFC3339)
	}

	return response
}
