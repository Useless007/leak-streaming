package movies

import (
	"context"
	"database/sql"
	"errors"
	"net/url"
	"time"

	"github.com/leak-streaming/leak-streaming/backend/internal/domain/movies"
	"github.com/leak-streaming/leak-streaming/backend/internal/persistence/repository"
)

var (
	ErrMovieUnavailable = errors.New("movie unavailable")
	ErrMovieNotFound    = sql.ErrNoRows
)

type TokenSigner interface {
	SignToken(movieID, viewerID string, ttl time.Duration) (string, error)
	ValidateToken(token, movieID string) (bool, error)
}

type Service struct {
	repo     *repository.MovieRepository
	signer   TokenSigner
	tokenTTL time.Duration
	now      func() time.Time
}

type StreamAccess struct {
	URL          string
	AllowedHosts []string
}

func NewService(repo *repository.MovieRepository, signer TokenSigner, tokenTTL time.Duration) *Service {
	return &Service{
		repo:     repo,
		signer:   signer,
		tokenTTL: tokenTTL,
		now:      time.Now,
	}
}

func (s *Service) GetMovie(ctx context.Context, slug string) (movies.Movie, error) {
	return s.repo.GetMovieWithStreams(ctx, slug)
}

func (s *Service) CreatePlaybackToken(ctx context.Context, movie movies.Movie, viewerID string) (string, error) {
	if !movie.IsAvailable(s.now()) {
		return "", ErrMovieUnavailable
	}
	if s.signer == nil {
		return "", errors.New("token signer not configured")
	}
	return s.signer.SignToken(movie.ID, viewerID, s.tokenTTL)
}

func (s *Service) ResolveStream(ctx context.Context, slug, token string) (StreamAccess, error) {
	movie, err := s.repo.GetMovieWithStreams(ctx, slug)
	if err != nil {
		return StreamAccess{}, err
	}

	if !movie.IsAvailable(s.now()) {
		return StreamAccess{}, ErrMovieUnavailable
	}

	ok, err := s.signer.ValidateToken(token, movie.ID)
	if err != nil {
		return StreamAccess{}, err
	}
	if !ok {
		return StreamAccess{}, errors.New("invalid token")
	}

	allowed := append([]string{}, movie.AllowedStreamHosts...)
	if parsed, err := url.Parse(movie.StreamURL); err == nil {
		host := parsed.Hostname()
		if host != "" {
			allowed = append(allowed, host)
		}
	}

	return StreamAccess{
		URL:          movie.StreamURL,
		AllowedHosts: allowed,
	}, nil
}
