-- +goose Up
-- Initial schema for movies, streams, captions, and playback tokens

CREATE TABLE movies (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    slug VARCHAR(128) NOT NULL,
    title VARCHAR(255) NOT NULL,
    synopsis TEXT NULL,
    poster_url VARCHAR(512) NULL,
    availability_start TIMESTAMP NULL,
    availability_end TIMESTAMP NULL,
    is_visible TINYINT(1) NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uq_movies_slug (slug),
    UNIQUE KEY uq_movies_title (title),
    KEY idx_movies_availability (availability_start, availability_end),
    KEY idx_movies_visibility (is_visible)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE movie_streams (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    movie_id BIGINT UNSIGNED NOT NULL,
    stream_url VARCHAR(512) NOT NULL,
    drm_key_id VARCHAR(128) NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uq_movie_streams_movie (movie_id),
    CONSTRAINT fk_movie_streams_movie FOREIGN KEY (movie_id) REFERENCES movies(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE movie_captions (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    movie_id BIGINT UNSIGNED NOT NULL,
    language_code VARCHAR(16) NOT NULL,
    label VARCHAR(64) NOT NULL,
    caption_url VARCHAR(512) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uq_movie_captions (movie_id, language_code),
    CONSTRAINT fk_movie_captions_movie FOREIGN KEY (movie_id) REFERENCES movies(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE playback_tokens (
    token CHAR(64) NOT NULL,
    movie_id BIGINT UNSIGNED NOT NULL,
    viewer_id VARCHAR(128) NULL,
    expires_at TIMESTAMP NOT NULL,
    issued_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (token),
    KEY idx_playback_tokens_movie (movie_id),
    KEY idx_playback_tokens_expires (expires_at),
    CONSTRAINT fk_playback_tokens_movie FOREIGN KEY (movie_id) REFERENCES movies(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Ensure caption URL integrity via basic check
ALTER TABLE movie_captions
    ADD CONSTRAINT chk_movie_captions_url CHECK (caption_url LIKE 'http%' OR caption_url LIKE 'https%');

ALTER TABLE movie_streams
    ADD CONSTRAINT chk_movie_streams_url CHECK (stream_url LIKE 'http%' OR stream_url LIKE 'https%');

-- +goose Down

DROP TABLE IF EXISTS playback_tokens;
DROP TABLE IF EXISTS movie_captions;
DROP TABLE IF EXISTS movie_streams;
DROP TABLE IF EXISTS movies;
