package movies

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/leak-streaming/leak-streaming/backend/internal/persistence/repository"
)

func TestCreateMovieSuccess(t *testing.T) {
	repo := repository.NewMovieRepository(nil)
	service := NewService(repo, NewInMemoryTokenSigner(), 5*time.Minute)

	start := time.Now().Add(2 * time.Hour).UTC().Format(time.RFC3339)
	end := time.Now().Add(26 * time.Hour).UTC().Format(time.RFC3339)

	input := CreateMovieInput{
		Title:             "Midnight Premiere",
		Synopsis:          "A suspense thriller that tests the streaming pipeline.",
		PosterURL:         "https://example.com/posters/midnight.jpg",
		AvailabilityStart: start,
		AvailabilityEnd:   end,
		IsVisible:         true,
		StreamURL:         "https://stream.example.com/movies/midnight/master.m3u8",
		AllowedHosts:      []string{"https://cdn.example.com", "stream.example.com"},
		Captions: []CaptionInput{
			{
				LanguageCode: "en",
				Label:        "English",
				CaptionURL:   "https://cdn.example.com/captions/midnight-en.vtt",
			},
		},
	}

	movie, err := service.CreateMovie(context.Background(), input)
	if err != nil {
		if valErr := (ValidationError{}); errors.As(err, &valErr) {
			t.Fatalf("expected success but got validation error: %+v", valErr.Fields)
		}
		t.Fatalf("CreateMovie returned error: %v", err)
	}

	if movie.Slug == "" {
		t.Fatalf("expected slug to be generated")
	}
	if movie.StreamURL != input.StreamURL {
		t.Fatalf("expected stream url %q, got %q", input.StreamURL, movie.StreamURL)
	}
	if !movie.IsVisible {
		t.Fatalf("expected movie to be visible")
	}
	if len(movie.Captions) != 1 {
		t.Fatalf("expected 1 caption, got %d", len(movie.Captions))
	}
	if movie.Captions[0].LanguageCode != "en" {
		t.Fatalf("expected caption language 'en', got %q", movie.Captions[0].LanguageCode)
	}
	if !containsHost(movie.AllowedStreamHosts, "stream.example.com") {
		t.Fatalf("expected allowed hosts to include stream host")
	}
	if !containsHost(movie.AllowedStreamHosts, "cdn.example.com") {
		t.Fatalf("expected allowed hosts to include provided host")
	}
	if movie.AvailabilityStart.IsZero() || movie.AvailabilityEnd.IsZero() {
		t.Fatalf("expected availability window to be set")
	}
}

func TestCreateMovieDuplicateTitle(t *testing.T) {
	repo := repository.NewMovieRepository(nil)
	service := NewService(repo, NewInMemoryTokenSigner(), 5*time.Minute)

	_, err := service.CreateMovie(context.Background(), CreateMovieInput{
		Title:     "ตัวอย่างภาพยนตร์",
		IsVisible: true,
		StreamURL: "https://stream.example.com/demo/master.m3u8",
	})
	if !errors.Is(err, ErrDuplicateMovieTitle) {
		t.Fatalf("expected ErrDuplicateMovieTitle, got %v", err)
	}
}

func TestCreateMovieValidationError(t *testing.T) {
	repo := repository.NewMovieRepository(nil)
	service := NewService(repo, NewInMemoryTokenSigner(), 5*time.Minute)

	_, err := service.CreateMovie(context.Background(), CreateMovieInput{})
	if err == nil {
		t.Fatalf("expected validation error")
	}
	var valErr ValidationError
	if !errors.As(err, &valErr) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
	if valErr.Fields["title"] == "" {
		t.Fatalf("expected validation error on title")
	}
	if valErr.Fields["streamUrl"] == "" {
		t.Fatalf("expected validation error on streamUrl")
	}
}

func containsHost(hosts []string, target string) bool {
	for _, host := range hosts {
		if host == target {
			return true
		}
	}
	return false
}
