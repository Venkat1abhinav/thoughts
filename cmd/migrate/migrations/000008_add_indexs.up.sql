CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Full-text-ish substring/fuzzy search
CREATE INDEX idx_comments_content
ON comments USING GIN (content gin_trgm_ops);

CREATE INDEX idx_posts_title
ON posts USING GIN (title gin_trgm_ops);

-- Array membership / overlap queries
CREATE INDEX idx_posts_tags
ON posts USING GIN (tags);

-- Username search
CREATE INDEX idx_users_username
ON users USING GIN (username gin_trgm_ops);

-- Foreign-key lookup
CREATE INDEX idx_posts_user_id
ON posts (user_id);

CREATE INDEX idx_comments_post_id
ON comments (post_id);
