package movies

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"

	service "github.com/leak-streaming/leak-streaming/backend/internal/service/movies"
)

type ManifestHandler struct {
	service *service.Service
}

func NewManifestHandler(service *service.Service) *ManifestHandler {
	return &ManifestHandler{service: service}
}

func (h *ManifestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.service == nil {
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)
		return
	}

	slug := chi.URLParam(r, "slug")
	token := r.URL.Query().Get("token")
	if slug == "" || token == "" {
		http.Error(w, "missing parameters", http.StatusBadRequest)
		return
	}

	streamAccess, err := h.service.ResolveStream(r.Context(), slug, token)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	resp, err := http.Get(streamAccess.URL) //nolint:gosec
	if err != nil {
		http.Error(w, "failed to fetch stream", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		http.Error(w, "upstream unavailable", http.StatusBadGateway)
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "failed to read manifest", http.StatusBadGateway)
		return
	}

	baseURL, err := url.Parse(streamAccess.URL)
	if err != nil {
		http.Error(w, "invalid upstream url", http.StatusBadGateway)
		return
	}

	rewritten := rewriteManifest(string(data), baseURL, slug, token)

	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	w.Header().Set("Cache-Control", "no-store")
	io.WriteString(w, rewritten)
}

func rewriteManifest(manifest string, base *url.URL, slug, token string) string {
	var builder strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(manifest))
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "#") {
			rel, err := url.Parse(trimmed)
			if err == nil {
				resolved := base.ResolveReference(rel)
				backendURL := url.URL{
					Path: fmt.Sprintf("/movies/%s/segment", slug),
				}
				q := backendURL.Query()
				q.Set("token", token)
				q.Set("target", resolved.String())
				backendURL.RawQuery = q.Encode()
				line = backendURL.String()
			}
		}
		builder.WriteString(line)
		builder.WriteByte('\n')
	}
	return builder.String()
}
