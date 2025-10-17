package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/leak-streaming/leak-streaming/backend/internal/domain/movies"
)

type MovieRepository struct {
	db *sql.DB
}

func NewMovieRepository(db *sql.DB) *MovieRepository {
	return &MovieRepository{db: db}
}

func (r *MovieRepository) CreateMovie(ctx context.Context, params CreateMovieParams) (movies.Movie, error) {
	if r.db == nil {
		return r.createMovieInMemory(params)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return movies.Movie{}, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	synopsis := sql.NullString{}
	if strings.TrimSpace(params.Synopsis) != "" {
		synopsis.Valid = true
		synopsis.String = params.Synopsis
	}

	posterURL := sql.NullString{}
	if strings.TrimSpace(params.PosterURL) != "" {
		posterURL.Valid = true
		posterURL.String = params.PosterURL
	}

	availabilityStart := sql.NullTime{}
	if params.AvailabilityStart != nil {
		availabilityStart.Valid = true
		availabilityStart.Time = params.AvailabilityStart.UTC()
	}

	availabilityEnd := sql.NullTime{}
	if params.AvailabilityEnd != nil {
		availabilityEnd.Valid = true
		availabilityEnd.Time = params.AvailabilityEnd.UTC()
	}

	var movieID int64
	insertMovieErr := tx.QueryRowContext(
		ctx,
		`INSERT INTO movies (slug, title, synopsis, poster_url, availability_start, availability_end, is_visible)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id`,
		params.Slug,
		params.Title,
		synopsis,
		posterURL,
		availabilityStart,
		availabilityEnd,
		params.IsVisible,
	).Scan(&movieID)
	if insertMovieErr != nil {
		return movies.Movie{}, translateCreateMovieError(insertMovieErr)
	}

	drmKey := sql.NullString{}
	if strings.TrimSpace(params.DRMKeyID) != "" {
		drmKey.Valid = true
		drmKey.String = params.DRMKeyID
	}

	allowedHosts := params.AllowedHosts
	if allowedHosts == nil {
		allowedHosts = []string{}
	}
	allowedHostsJSON, err := json.Marshal(allowedHosts)
	if err != nil {
		return movies.Movie{}, err
	}

	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO movie_streams (movie_id, stream_url, drm_key_id, allowed_hosts)
		 VALUES ($1, $2, $3, $4)`,
		movieID,
		params.StreamURL,
		drmKey,
		allowedHostsJSON,
	); err != nil {
		return movies.Movie{}, translateCreateMovieError(err)
	}

	if len(params.Captions) > 0 {
		for _, caption := range params.Captions {
			if _, err := tx.ExecContext(
				ctx,
				`INSERT INTO movie_captions (movie_id, language_code, label, caption_url)
				 VALUES ($1, $2, $3, $4)`,
				movieID,
				caption.LanguageCode,
				caption.Label,
				caption.CaptionURL,
			); err != nil {
				return movies.Movie{}, translateCreateMovieError(err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return movies.Movie{}, err
	}

	return r.GetMovieWithStreams(ctx, params.Slug)
}

func (r *MovieRepository) createMovieInMemory(params CreateMovieParams) (movies.Movie, error) {
	for _, existing := range sampleMovies {
		if existing.Slug == params.Slug {
			return movies.Movie{}, ErrDuplicateSlug
		}
		if strings.EqualFold(existing.Title, params.Title) {
			return movies.Movie{}, ErrDuplicateTitle
		}
	}

	var maxID int64
	for _, existing := range sampleMovies {
		if id, err := strconv.ParseInt(existing.ID, 10, 64); err == nil {
			if id > maxID {
				maxID = id
			}
		}
	}

	newID := strconv.FormatInt(maxID+1, 10)
	movie := movies.Movie{
		ID:                 newID,
		Slug:               params.Slug,
		Title:              params.Title,
		Synopsis:           params.Synopsis,
		PosterURL:          params.PosterURL,
		IsVisible:          params.IsVisible,
		StreamURL:          params.StreamURL,
		DRMKeyID:           params.DRMKeyID,
		Captions:           append([]movies.Caption(nil), params.Captions...),
		AllowedStreamHosts: append([]string(nil), params.AllowedHosts...),
	}
	if movie.Captions == nil {
		movie.Captions = []movies.Caption{}
	}
	if movie.AllowedStreamHosts == nil {
		movie.AllowedStreamHosts = []string{}
	}
	if params.AvailabilityStart != nil {
		movie.AvailabilityStart = params.AvailabilityStart.UTC()
	}
	if params.AvailabilityEnd != nil {
		movie.AvailabilityEnd = params.AvailabilityEnd.UTC()
	}

	sampleMovies[movie.Slug] = movie

	return movie, nil
}

func translateCreateMovieError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.ConstraintName {
		case "movies_slug_key":
			return ErrDuplicateSlug
		case "movies_title_key":
			return ErrDuplicateTitle
		}
	}
	return err
}

var (
	ErrDuplicateSlug  = errors.New("duplicate movie slug")
	ErrDuplicateTitle = errors.New("duplicate movie title")
)

type CreateMovieParams struct {
	Slug              string
	Title             string
	Synopsis          string
	PosterURL         string
	AvailabilityStart *time.Time
	AvailabilityEnd   *time.Time
	IsVisible         bool
	StreamURL         string
	DRMKeyID          string
	AllowedHosts      []string
	Captions          []movies.Caption
}

func (r *MovieRepository) ListMovies(ctx context.Context) ([]movies.Movie, error) {
	if r.db == nil {
		items := make([]movies.Movie, 0, len(sampleMovies))
		for _, movie := range sampleMovies {
			items = append(items, movie)
		}
		sort.Slice(items, func(i, j int) bool {
			return items[i].Slug < items[j].Slug
		})
		return items, nil
	}

	const query = `
SELECT id,
       slug,
       title,
       synopsis,
       poster_url,
       availability_start,
       availability_end,
       is_visible
FROM movies
WHERE is_visible = TRUE
ORDER BY COALESCE(availability_start, NOW()) ASC, slug ASC;
`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	moviesList := make([]movies.Movie, 0, 8)

	for rows.Next() {
		var (
			movieID           int64
			slug              string
			title             string
			synopsis          sql.NullString
			posterURL         sql.NullString
			availabilityStart sql.NullTime
			availabilityEnd   sql.NullTime
			isVisible         bool
		)

		if err := rows.Scan(&movieID, &slug, &title, &synopsis, &posterURL, &availabilityStart, &availabilityEnd, &isVisible); err != nil {
			return nil, err
		}

		movie := movies.Movie{
			ID:        strconv.FormatInt(movieID, 10),
			Slug:      slug,
			Title:     title,
			IsVisible: isVisible,
			Captions:  []movies.Caption{},
		}

		if synopsis.Valid {
			movie.Synopsis = synopsis.String
		}
		if posterURL.Valid {
			movie.PosterURL = posterURL.String
		}
		if availabilityStart.Valid {
			movie.AvailabilityStart = availabilityStart.Time.UTC()
		}
		if availabilityEnd.Valid {
			movie.AvailabilityEnd = availabilityEnd.Time.UTC()
		}

		moviesList = append(moviesList, movie)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return moviesList, nil
}

func (r *MovieRepository) GetMovieWithStreams(ctx context.Context, slug string) (movies.Movie, error) {
	if r.db == nil {
		movie, ok := sampleMovies[slug]
		if !ok {
			return movies.Movie{}, sql.ErrNoRows
		}
		return movie, nil
	}

	const movieQuery = `
SELECT m.id,
       m.slug,
       m.title,
       m.synopsis,
       m.poster_url,
       m.availability_start,
       m.availability_end,
       m.is_visible,
       s.stream_url,
       s.drm_key_id,
       COALESCE(s.allowed_hosts, '[]'::jsonb)
FROM movies m
LEFT JOIN movie_streams s ON s.movie_id = m.id
WHERE m.slug = $1
LIMIT 1;
`

	var (
		movieID           int64
		title             string
		synopsis          sql.NullString
		posterURL         sql.NullString
		availabilityStart sql.NullTime
		availabilityEnd   sql.NullTime
		isVisible         bool
		streamURL         sql.NullString
		drmKeyID          sql.NullString
		allowedHostsRaw   []byte
	)

	row := r.db.QueryRowContext(ctx, movieQuery, slug)
	if err := row.Scan(
		&movieID,
		&slug,
		&title,
		&synopsis,
		&posterURL,
		&availabilityStart,
		&availabilityEnd,
		&isVisible,
		&streamURL,
		&drmKeyID,
		&allowedHostsRaw,
	); err != nil {
		return movies.Movie{}, err
	}

	movie := movies.Movie{
		ID:        strconv.FormatInt(movieID, 10),
		Slug:      slug,
		Title:     title,
		IsVisible: isVisible,
	}
	if synopsis.Valid {
		movie.Synopsis = synopsis.String
	}
	if posterURL.Valid {
		movie.PosterURL = posterURL.String
	}
	if availabilityStart.Valid {
		movie.AvailabilityStart = availabilityStart.Time.UTC()
	}
	if availabilityEnd.Valid {
		movie.AvailabilityEnd = availabilityEnd.Time.UTC()
	}
	if streamURL.Valid {
		movie.StreamURL = streamURL.String
	}
	if drmKeyID.Valid {
		movie.DRMKeyID = drmKeyID.String
	}
	if len(allowedHostsRaw) > 0 {
		var hosts []string
		if err := json.Unmarshal(allowedHostsRaw, &hosts); err == nil {
			movie.AllowedStreamHosts = hosts
		}
	}

	const captionsQuery = `
SELECT language_code, label, caption_url
FROM movie_captions
WHERE movie_id = $1
ORDER BY language_code;
`

	rows, err := r.db.QueryContext(ctx, captionsQuery, movieID)
	if err != nil {
		return movies.Movie{}, err
	}
	defer rows.Close()

	captions := make([]movies.Caption, 0)
	for rows.Next() {
		var (
			languageCode string
			label        string
			captionURL   string
		)
		if err := rows.Scan(&languageCode, &label, &captionURL); err != nil {
			return movies.Movie{}, err
		}
		captions = append(captions, movies.Caption{
			LanguageCode: languageCode,
			Label:        label,
			CaptionURL:   captionURL,
		})
	}
	if err := rows.Err(); err != nil {
		return movies.Movie{}, err
	}

	movie.Captions = captions

	if movie.AllowedStreamHosts == nil {
		movie.AllowedStreamHosts = []string{}
	}
	if len(movie.AllowedStreamHosts) == 0 && movie.StreamURL != "" {
		if host := extractHost(movie.StreamURL); host != "" {
			movie.AllowedStreamHosts = []string{host}
		}
	}

	return movie, nil
}

func extractHost(raw string) string {
	parsed, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	return parsed.Hostname()
}

func (r *MovieRepository) UpsertSampleMovie(movie movies.Movie) {
	sampleMovies[movie.Slug] = movie
}

var sampleMovies = map[string]movies.Movie{
	"sample-movie": {
		ID:                "1",
		Slug:              "sample-movie",
		Title:             "ตัวอย่างภาพยนตร์",
		Synopsis:          "เรื่องราวของระบบสตรีมมิ่งที่พร้อมให้ทดสอบประสบการณ์การเล่นแบบ HLS",
		PosterURL:         "https://images.unsplash.com/photo-1524985069026-dd778a71c7b4?w=800",
		AvailabilityStart: time.Now().Add(-24 * time.Hour).UTC(),
		AvailabilityEnd:   time.Now().Add(24 * time.Hour).UTC(),
		IsVisible:         true,
		StreamURL:         "https://main.24playerhd.com/m3u8/0378b65549cda348e910faf0/0378b65549cda348e910faf0168.m3u8", // m3u8 URL ของสตรีมมิ่ง
		AllowedStreamHosts: []string{
			"main.24playerhd.com",
			"m42.winplay4.com",
			"winplay4.com",
		},
		Captions: []movies.Caption{
			{
				LanguageCode: "en",
				Label:        "English",
				CaptionURL:   "/captions/sample-en.vtt",
			},
		},
	},
	"demo-movie-2": {
		ID:                "2",
		Slug:              "demo-movie-2",
		Title:             "ตัวอย่างภาพยนตร์ลำดับที่สอง",
		Synopsis:          "ภาคต่อของการทดสอบระบบด้วยลิงก์ HLS ใหม่ พร้อมสตรีมต่อเนื่อง",
		PosterURL:         "https://images.unsplash.com/photo-1497032205916-ac775f0649ae?w=800",
		AvailabilityStart: time.Now().Add(-12 * time.Hour).UTC(),
		AvailabilityEnd:   time.Now().Add(48 * time.Hour).UTC(),
		IsVisible:         true,
		StreamURL:         "https://main.24playerhd.com/m3u8/f87ff8ffe0151aec3f5d55bc/f87ff8ffe0151aec3f5d55bc168.m3u8",
		AllowedStreamHosts: []string{
			"main.24playerhd.com",
			"m42.winplay4.com",
			"winplay4.com",
			"m37.upplay4.com",
			"upplay4.com",
		},
		Captions: []movies.Caption{
			{
				LanguageCode: "en",
				Label:        "English",
				CaptionURL:   "/captions/sample-en.vtt",
			},
		},
	},
}
