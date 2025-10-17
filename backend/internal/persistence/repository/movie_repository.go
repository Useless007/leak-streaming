package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/url"
	"sort"
	"strconv"
	"time"

	"github.com/leak-streaming/leak-streaming/backend/internal/domain/movies"
)

type MovieRepository struct {
	db *sql.DB
}

func NewMovieRepository(db *sql.DB) *MovieRepository {
	return &MovieRepository{db: db}
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
