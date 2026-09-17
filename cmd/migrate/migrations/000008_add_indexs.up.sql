CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Full-text-ish substring/fuzzy search
CREATE INDEX idx_comments_content ON comments USING GIN (content gin_trgm_ops);

CREATE INDEX idx_posts_title ON posts USING GIN (title gin_trgm_ops);

-- Array membership / overlap queries
CREATE INDEX idx_posts_tags ON posts USING GIN (tags);

-- Username search
CREATE INDEX idx_users_username ON users USING GIN (username gin_trgm_ops);

-- Foreign-key lookup
CREATE INDEX idx_posts_user_id ON posts (user_id);

CREATE INDEX idx_comments_post_id ON comments (post_id);

SELECT

    p.id
    , p.user_id
    , p.title
    , p.content
    , p.created_at
    , p.version
    , p.tags
    , u.username
    , COUNT(c.id) AS comments_count 

FROM posts p 

LEFT JOIN comments c
    ON c.post_id = p.id 

LEFT JOIN users u
    ON p.user_id = u.id 

LEFT JOIN followers f
    ON f.follower_id = p.user_id
    AND f.user_id = $1 

WHERE
    (
        p.user_id = $1
        OR f.user_id IS NOT NULL
    )
    AND (
        p.title ILIKE '%' || $4 || '%'
        OR p.content ILIKE '%' || $4 || '%'
    ) 

GROUP BY
    p.id
    , u.username 

ORDER BY
    p.created_at ` + sort + ` 

LIMIT
    $2
OFFSET
    $3 
;
