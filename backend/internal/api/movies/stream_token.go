package movies

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	service "github.com/leak-streaming/leak-streaming/backend/internal/service/movies"
)

type StreamTokenHandler struct {
	service *service.Service
}

func NewStreamTokenHandler(service *service.Service) *StreamTokenHandler {
	return &StreamTokenHandler{service: service}
}

func (h *StreamTokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.service == nil {
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)
		return
	}

	slug := chi.URLParam(r, "slug")
	if slug == "" {
		http.Error(w, "missing slug", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	movie, err := h.service.GetMovie(ctx, slug)
	if err != nil {
		if errors.Is(err, service.ErrMovieNotFound) {
			http.Error(w, "movie not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to load movie", http.StatusInternalServerError)
		return
	}

	viewerID := viewerIDFromRequest(r)

	token, err := h.service.CreatePlaybackToken(ctx, movie, viewerID)
	if err != nil {
		if errors.Is(err, service.ErrMovieUnavailable) {
			http.Error(w, "movie unavailable", http.StatusConflict)
			return
		}
		http.Error(w, "failed to create token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(streamTokenResponse{
		Token: token,
	})
}

type streamTokenResponse struct {
	Token string `json:"token"`
}

func viewerIDFromRequest(r *http.Request) string {
	if id := r.Header.Get("X-Viewer-ID"); id != "" {
		return id
	}
	if cookie, err := r.Cookie("viewer_id"); err == nil {
		return cookie.Value
	}
	if ip := clientIP(r); ip != "" {
		return "ip:" + ip
	}
	return "anonymous"
}

func clientIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}
	host := r.RemoteAddr
	if host == "" {
		return ""
	}
	if idx := strings.LastIndex(host, ":"); idx >= 0 {
		return host[:idx]
	}
	return host
}
