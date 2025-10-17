-- +goose Up
UPDATE movie_streams
SET allowed_hosts = '["main.24playerhd.com","m42.winplay4.com","winplay4.com","m37.upplay4.com","upplay4.com"]'::jsonb
WHERE movie_id = (SELECT id FROM movies WHERE slug = 'demo-movie-2');

-- +goose Down
UPDATE movie_streams
SET allowed_hosts = '["main.24playerhd.com","m42.winplay4.com","winplay4.com"]'::jsonb
WHERE movie_id = (SELECT id FROM movies WHERE slug = 'demo-movie-2');
