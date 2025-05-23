DROP VIEW IF EXISTS v_posts;
DROP TRIGGER IF EXISTS update_like_count_after_insert_in_post_feedback;
DROP TRIGGER IF EXISTS update_like_count_after_update_in_post_feedback;
DROP TRIGGER IF EXISTS update_dislike_count_after_insert_in_post_feedback;
DROP TRIGGER IF EXISTS update_dislike_count_after_update_in_post_feedback;
DROP TRIGGER IF EXISTS update_comment_count_after_insert;

CREATE TABLE posts_old (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id         INTEGER NOT NULL,
    comment_count   INTEGER NOT NULL DEFAULT 0,
    like_count      INTEGER NOT NULL DEFAULT 0,
    dislike_count   INTEGER NOT NULL DEFAULT 0,
    title           TEXT NOT NULL,
    content         TEXT NOT NULL,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

INSERT INTO posts_old (
    id, user_id, title, content, created_at,
    like_count, dislike_count, comment_count
)
SELECT
    id, user_id, title, content, created_at,
    like_count, dislike_count, comment_count
FROM posts;

DROP TABLE posts;
ALTER TABLE posts_old RENAME TO posts;

CREATE TABLE post_feedback_old (
    user_id     INTEGER NOT NULL,
    parent_id   INTEGER NOT NULL,
    rating      INTEGER NOT NULL CHECK (rating IN (-1, 0, 1)),
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (user_id, parent_id),
    FOREIGN KEY (user_id)   REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES posts(id) ON DELETE CASCADE
);

INSERT INTO post_feedback_old (
    user_id, parent_id, rating, created_at
)
SELECT
    user_id, parent_id, rating, created_at
FROM post_feedback;

DROP TABLE post_feedback;
ALTER TABLE post_feedback_old RENAME TO post_feedback;

CREATE TABLE post_categories_old (
    post_id     INTEGER NOT NULL,
    category_id INTEGER NOT NULL,

    PRIMARY KEY (post_id, category_id),
    FOREIGN KEY (post_id)       REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id)   REFERENCES categories(id)
);

INSERT INTO post_categories_old (
    post_id, category_id
)
SELECT
    post_id, category_id
FROM post_categories;

DROP TABLE post_categories;
ALTER TABLE post_categories_old RENAME TO post_categories;

CREATE VIEW v_posts AS
SELECT 
    p.id,
    p.user_id AS puser_id,
    u.nick_name AS user_name,
    p.comment_count,
    p.like_count,
    p.dislike_count,
    p.title,
    p.content,
    p.created_at AS pcreated_at,
    GROUP_CONCAT(pc.category_id) AS category_ids
FROM posts p
INNER JOIN users u ON p.user_id = u.id
LEFT JOIN post_categories pc ON p.id = pc.post_id
GROUP BY 
    p.id, p.user_id, u.nick_name,
    p.comment_count, p.like_count, p.dislike_count,
    p.title, p.content, p.created_at;

CREATE TRIGGER update_like_count_after_insert_in_post_feedback
AFTER INSERT ON post_feedback
FOR EACH ROW
BEGIN
    UPDATE posts
    SET like_count = (
        SELECT COUNT(*) 
        FROM post_feedback
        WHERE parent_id = NEW.parent_id AND rating = 1)
    WHERE id = NEW.parent_id;
END;

CREATE TRIGGER update_like_count_after_update_in_post_feedback
AFTER UPDATE ON post_feedback
FOR EACH ROW
BEGIN
    UPDATE posts
    SET like_count = (
        SELECT COUNT(*) 
        FROM post_feedback
        WHERE parent_id = NEW.parent_id AND rating = 1)
    WHERE id = NEW.parent_id;
END;

CREATE TRIGGER update_dislike_count_after_insert_in_post_feedback
AFTER INSERT ON post_feedback
FOR EACH ROW
BEGIN
    UPDATE posts
    SET dislike_count = (
        SELECT COUNT(*) 
        FROM post_feedback
        WHERE parent_id = NEW.parent_id AND rating = -1)
    WHERE id = NEW.parent_id;
END;

CREATE TRIGGER update_dislike_count_after_update_in_post_feedback
AFTER UPDATE ON post_feedback
FOR EACH ROW
BEGIN
    UPDATE posts
    SET dislike_count = (
        SELECT COUNT(*) 
        FROM post_feedback
        WHERE parent_id = NEW.parent_id AND rating = -1)
    WHERE id = NEW.parent_id;
END;

CREATE TRIGGER update_comment_count_after_insert
AFTER INSERT ON comments
FOR EACH ROW
BEGIN
    UPDATE posts
    SET comment_count = (
        SELECT COUNT(*) 
        FROM comments 
        WHERE post_id = NEW.post_id)
    WHERE id = NEW.post_id;
END;
