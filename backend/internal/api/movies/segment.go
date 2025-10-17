package movies

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"

	service "github.com/leak-streaming/leak-streaming/backend/internal/service/movies"
)

type SegmentHandler struct {
	service *service.Service
}

func NewSegmentHandler(service *service.Service) *SegmentHandler {
	return &SegmentHandler{service: service}
}

func (h *SegmentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.service == nil {
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)
		return
	}

	slug := chi.URLParam(r, "slug")
	token := r.URL.Query().Get("token")
	target := r.URL.Query().Get("target")
	if slug == "" || token == "" || target == "" {
		http.Error(w, "missing parameters", http.StatusBadRequest)
		return
	}

	streamAccess, err := h.service.ResolveStream(r.Context(), slug, token)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	baseURL, err := url.Parse(streamAccess.URL)
	if err != nil {
		http.Error(w, "invalid base url", http.StatusBadRequest)
		return
	}

	targetURL, err := url.Parse(target)
	if err != nil {
		http.Error(w, "invalid target", http.StatusBadRequest)
		return
	}

	if !targetURL.IsAbs() {
		targetURL = baseURL.ResolveReference(targetURL)
	}

	targetHost := targetURL.Hostname()
	if !isAllowedHost(targetHost, streamAccess.AllowedHosts) {
		http.Error(w, "forbidden host", http.StatusForbidden)
		return
	}

	if targetURL.Scheme != baseURL.Scheme {
		http.Error(w, "forbidden host", http.StatusForbidden)
		return
	}

	resp, err := http.Get(targetURL.String()) //nolint:gosec
	if err != nil {
		http.Error(w, "failed to fetch segment", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		http.Error(w, "upstream unavailable", http.StatusBadGateway)
		return
	}

	for key, values := range resp.Header {
		if key == "Content-Length" || key == "Content-Type" {
			for _, v := range values {
				w.Header().Add(key, v)
			}
		}
	}
	io.Copy(w, resp.Body)
}

func isAllowedHost(host string, allowed []string) bool {
	if host == "" {
		return false
	}
	for _, entry := range allowed {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		if strings.EqualFold(host, entry) {
			return true
		}
		if strings.HasPrefix(entry, ".") && strings.HasSuffix(host, entry) {
			return true
		}
		if !strings.HasPrefix(entry, ".") && strings.HasSuffix(host, "."+entry) {
			return true
		}
	}
	return false
}
