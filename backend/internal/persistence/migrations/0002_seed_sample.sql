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
        'sample-movie',
        'ตัวอย่างภาพยนตร์',
        'เรื่องราวของระบบสตรีมมิ่งที่พร้อมให้ทดสอบประสบการณ์การเล่นแบบ HLS',
        'https://images.unsplash.com/photo-1524985069026-dd778a71c7b4?w=800',
        NOW() - INTERVAL '1 day',
        NOW() + INTERVAL '1 day',
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
       'https://main.24playerhd.com/m3u8/0378b65549cda348e910faf0/0378b65549cda348e910faf0168.m3u8',
       NULL,
       '["main.24playerhd.com","m42.winplay4.com","winplay4.com"]'::jsonb
FROM movie
ON CONFLICT (movie_id) DO UPDATE SET
    stream_url = EXCLUDED.stream_url,
    drm_key_id = EXCLUDED.drm_key_id,
    allowed_hosts = EXCLUDED.allowed_hosts;

WITH movie AS (
    SELECT id FROM movies WHERE slug = 'sample-movie'
)
INSERT INTO movie_captions (movie_id, language_code, label, caption_url)
SELECT id, 'en', 'English', '/captions/sample-en.vtt' FROM movie
ON CONFLICT (movie_id, language_code) DO UPDATE SET
    label = EXCLUDED.label,
    caption_url = EXCLUDED.caption_url;

-- +goose Down
DELETE FROM movie_captions WHERE movie_id IN (SELECT id FROM movies WHERE slug = 'sample-movie');
DELETE FROM movie_streams WHERE movie_id IN (SELECT id FROM movies WHERE slug = 'sample-movie');
DELETE FROM movies WHERE slug = 'sample-movie';
