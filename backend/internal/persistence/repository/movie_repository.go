package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/leak-streaming/leak-streaming/backend/internal/domain/movies"
)

type MovieRepository struct {
	db *sql.DB
}

func NewMovieRepository(db *sql.DB) *MovieRepository {
	return &MovieRepository{db: db}
}

func (r *MovieRepository) GetMovieWithStreams(ctx context.Context, slug string) (movies.Movie, error) {
	movie, ok := sampleMovies[slug]
	if !ok {
		return movies.Movie{}, sql.ErrNoRows
	}
	return movie, nil
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
		StreamURL:         "https://main.24playerhd.com/m3u8/0378b65549cda348e910faf0/0378b65549cda348e910faf0168.m3u8",
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
}
