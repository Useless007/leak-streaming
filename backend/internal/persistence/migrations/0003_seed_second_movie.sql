-- +goose Up
WITH movie AS (
    INSERT INTO movies (
        slug,
        title,
        synopsis,
        poster_url,
        availability_start,
        availability_end,
        is_visible
    ) VALUES (
        'demo-movie-2',
        'ตัวอย่างภาพยนตร์ลำดับที่สอง',
        'ภาคต่อของการทดสอบระบบด้วยลิงก์ HLS ใหม่ พร้อมสตรีมต่อเนื่อง',
        'https://images.unsplash.com/photo-1497032205916-ac775f0649ae?w=800',
        NOW() - INTERVAL '12 hours',
        NOW() + INTERVAL '48 hours',
        TRUE
    )
    ON CONFLICT (slug) DO UPDATE SET
        title = EXCLUDED.title,
        synopsis = EXCLUDED.synopsis,
        poster_url = EXCLUDED.poster_url,
        availability_start = EXCLUDED.availability_start,
        availability_end = EXCLUDED.availability_end,
        is_visible = EXCLUDED.is_visible
    RETURNING id
)
INSERT INTO movie_streams (movie_id, stream_url, drm_key_id, allowed_hosts)
SELECT id,
       'https://main.24playerhd.com/m3u8/f87ff8ffe0151aec3f5d55bc/f87ff8ffe0151aec3f5d55bc168.m3u8',
       NULL,
       '["main.24playerhd.com","m42.winplay4.com","winplay4.com","m37.upplay4.com","upplay4.com"]'::jsonb
FROM movie
ON CONFLICT (movie_id) DO UPDATE SET
    stream_url = EXCLUDED.stream_url,
    drm_key_id = EXCLUDED.drm_key_id,
    allowed_hosts = EXCLUDED.allowed_hosts;

WITH movie AS (
    SELECT id FROM movies WHERE slug = 'demo-movie-2'
)
INSERT INTO movie_captions (movie_id, language_code, label, caption_url)
SELECT id, 'en', 'English', '/captions/sample-en.vtt' FROM movie
ON CONFLICT (movie_id, language_code) DO UPDATE SET
    label = EXCLUDED.label,
    caption_url = EXCLUDED.caption_url;

-- +goose Down
DELETE FROM movie_captions WHERE movie_id IN (SELECT id FROM movies WHERE slug = 'demo-movie-2');
DELETE FROM movie_streams WHERE movie_id IN (SELECT id FROM movies WHERE slug = 'demo-movie-2');
DELETE FROM movies WHERE slug = 'demo-movie-2';
