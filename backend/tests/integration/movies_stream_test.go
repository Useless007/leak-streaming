package integration

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	apimiddleware "github.com/leak-streaming/leak-streaming/backend/internal/api/middleware"
	apimovies "github.com/leak-streaming/leak-streaming/backend/internal/api/movies"
	"github.com/leak-streaming/leak-streaming/backend/internal/domain/movies"
	"github.com/leak-streaming/leak-streaming/backend/internal/persistence/repository"
	service "github.com/leak-streaming/leak-streaming/backend/internal/service/movies"
)

func TestMovieStreamFlow(t *testing.T) {
	t.Parallel()

	upstreamSegmentBody := "FAKE-SEGMENT-DATA"
	var playbackToken string

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/movie.m3u8":
			w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
			io.WriteString(w, "#EXTM3U\n#EXT-X-VERSION:3\n#EXT-X-TARGETDURATION:4\n#EXTINF:4,\nsegment.ts\n")
		case "/segment.ts":
			w.Header().Set("Content-Type", "video/mp2t")
			io.WriteString(w, upstreamSegmentBody)
		default:
			http.NotFound(w, r)
		}
	}))
	defer upstream.Close()

	upstreamURL, err := url.Parse(upstream.URL + "/movie.m3u8")
	if err != nil {
		t.Fatalf("failed to parse upstream url: %v", err)
	}

	repo := repository.NewMovieRepository(nil)
	repo.UpsertSampleMovie(movies.Movie{
		ID:                 "movie-001",
		Slug:               "integration-movie",
		Title:              "Integration Test Movie",
		Synopsis:           "integration playback scenario",
		StreamURL:          upstreamURL.String(),
		AllowedStreamHosts: []string{},
		IsVisible:          true,
		AvailabilityStart:  time.Now().Add(-time.Hour),
		AvailabilityEnd:    time.Now().Add(time.Hour),
	})

	signer := service.NewInMemoryTokenSigner()
	movieService := service.NewService(repo, signer, time.Minute)

	r := chi.NewRouter()
	r.Use(apimiddleware.SecureHeaders())

	detailsHandler := apimovies.NewDetailsHandler(movieService)
	tokenHandler := apimovies.NewStreamTokenHandler(movieService)
	manifestHandler := apimovies.NewManifestHandler(movieService)
	segmentHandler := apimovies.NewSegmentHandler(movieService)

	r.Get("/movies/{slug}", detailsHandler.ServeHTTP)
	r.Post("/movies/{slug}/playback-token", tokenHandler.ServeHTTP)
	r.Get("/movies/{slug}/manifest.m3u8", manifestHandler.ServeHTTP)
	r.Get("/movies/{slug}/segment", segmentHandler.ServeHTTP)

	server := httptest.NewServer(r)
	defer server.Close()

	// Step 1: fetch movie details
	resp, err := http.Get(server.URL + "/movies/integration-movie")
	if err != nil {
		t.Fatalf("failed to get movie details: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var moviePayload struct {
		Slug  string `json:"slug"`
		Title string `json:"title"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&moviePayload); err != nil {
		t.Fatalf("failed to decode movie payload: %v", err)
	}

	if moviePayload.Slug != "integration-movie" {
		t.Fatalf("unexpected slug: %s", moviePayload.Slug)
	}

	// Step 2: request playback token
	resp, err = http.Post(server.URL+"/movies/integration-movie/playback-token", "application/json", http.NoBody)
	if err != nil {
		t.Fatalf("failed to request playback token: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 from token endpoint, got %d", resp.StatusCode)
	}

	var tokenPayload struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenPayload); err != nil {
		t.Fatalf("failed to decode token payload: %v", err)
	}
	if tokenPayload.Token == "" {
		t.Fatalf("expected non-empty token")
	}
	playbackToken = tokenPayload.Token

	// Step 3: fetch manifest via backend proxy
	manifestResp, err := http.Get(server.URL + "/movies/integration-movie/manifest.m3u8?token=" + playbackToken)
	if err != nil {
		t.Fatalf("failed to fetch manifest: %v", err)
	}
	defer manifestResp.Body.Close()

	if manifestResp.StatusCode != http.StatusOK {
		t.Fatalf("expected manifest 200, got %d", manifestResp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(manifestResp.Body)
	if err != nil {
		t.Fatalf("failed to read manifest body: %v", err)
	}
	body := string(bodyBytes)
	if !strings.Contains(body, "/movies/integration-movie/segment?") {
		t.Fatalf("expected manifest to contain rewritten segment URL, got: %s", body)
	}
	if !strings.Contains(body, "token="+playbackToken) {
		t.Fatalf("manifest should include playback token")
	}

	// Step 4: request proxied segment
	lines := strings.Split(body, "\n")
	var segmentPath string
	for _, line := range lines {
		if strings.HasPrefix(line, "/movies/") {
			segmentPath = line
			break
		}
	}
	if segmentPath == "" {
		t.Fatalf("could not find segment path in manifest")
	}

	segmentResp, err := http.Get(server.URL + segmentPath)
	if err != nil {
		t.Fatalf("failed to fetch segment: %v", err)
	}
	defer segmentResp.Body.Close()

	if segmentResp.StatusCode != http.StatusOK {
		t.Fatalf("expected segment 200, got %d", segmentResp.StatusCode)
	}

	data, err := io.ReadAll(segmentResp.Body)
	if err != nil {
		t.Fatalf("failed to read segment body: %v", err)
	}
	if string(data) != upstreamSegmentBody {
		t.Fatalf("unexpected segment body: %s", string(data))
	}

	// Step 5: invalid host should be blocked
	badURL := server.URL + "/movies/integration-movie/segment?token=" + playbackToken + "&target=" + url.QueryEscape("https://evil.example.com/segment.ts")
	badResp, err := http.Get(badURL)
	if err != nil {
		t.Fatalf("failed to fetch segment with bad host: %v", err)
	}
	defer badResp.Body.Close()
	if badResp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 for forbidden host, got %d", badResp.StatusCode)
	}

	// Step 6: invalid token should be rejected
	invalidResp, err := http.Get(server.URL + "/movies/integration-movie/manifest.m3u8?token=invalid")
	if err != nil {
		t.Fatalf("failed to call manifest with invalid token: %v", err)
	}
	defer invalidResp.Body.Close()
	if invalidResp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 for invalid token, got %d", invalidResp.StatusCode)
	}
}
