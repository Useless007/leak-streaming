-- +goose Up
-- Initial schema for movies, streams, captions, and playback tokens (PostgreSQL)

CREATE TABLE movies (
    id BIGSERIAL PRIMARY KEY,
    slug VARCHAR(128) NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL UNIQUE,
    synopsis TEXT NULL,
    poster_url VARCHAR(512) NULL,
    availability_start TIMESTAMPTZ NULL,
    availability_end TIMESTAMPTZ NULL,
    is_visible BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_movies_availability ON movies (availability_start, availability_end);
CREATE INDEX idx_movies_visibility ON movies (is_visible);

CREATE TABLE movie_streams (
    id BIGSERIAL PRIMARY KEY,
    movie_id BIGINT NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    stream_url VARCHAR(1024) NOT NULL,
    drm_key_id VARCHAR(128) NULL,
    allowed_hosts JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (movie_id)
);

CREATE TABLE movie_captions (
    id BIGSERIAL PRIMARY KEY,
    movie_id BIGINT NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    language_code VARCHAR(16) NOT NULL,
    label VARCHAR(64) NOT NULL,
    caption_url VARCHAR(1024) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (movie_id, language_code),
    CHECK (caption_url LIKE 'http%' OR caption_url LIKE 'https%' OR caption_url LIKE '/%')
);

CREATE TABLE playback_tokens (
    token VARCHAR(64) PRIMARY KEY,
    movie_id BIGINT NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    viewer_id VARCHAR(128) NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    issued_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_playback_tokens_movie ON playback_tokens (movie_id);
CREATE INDEX idx_playback_tokens_expires ON playback_tokens (expires_at);

ALTER TABLE movies
    ADD CONSTRAINT chk_movies_availability CHECK (
        availability_start IS NULL
        OR availability_end IS NULL
        OR availability_start <= availability_end
    );

ALTER TABLE movie_streams
    ADD CONSTRAINT chk_movie_streams_url CHECK (stream_url LIKE 'http%' OR stream_url LIKE 'https%');

-- +goose Down

DROP TABLE IF EXISTS playback_tokens;
DROP TABLE IF EXISTS movie_captions;
DROP TABLE IF EXISTS movie_streams;
DROP TABLE IF EXISTS movies;
