DROP TRIGGER IF EXISTS update_like_count_after_insert_in_comment_feedback;
DROP TRIGGER IF EXISTS update_like_count_after_update_in_comment_feeback;
DROP TRIGGER IF EXISTS update_dislike_count_after_insert_in_comment_feedback;
DROP TRIGGER IF EXISTS update_dislike_count_after_update_in_comment_feeback;
DROP TRIGGER IF EXISTS update_comment_count_after_insert;

CREATE TABLE comments_old (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id         INTEGER NOT NULL,
    user_name       TEXT NOT NULL,
    post_id         INTEGER NOT NULL,
    parent_id       INTEGER NULL,
    content         TEXT NOT NULL DEFAULT 0,
    like_count      INTEGER NOT NULL DEFAULT 0,
    dislike_count   INTEGER NOT NULL DEFAULT 0,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id)   REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id)   REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES comments(id)
);

INSERT INTO comments_old (
    id, user_id, user_name, post_id, parent_id, content,
    like_count, dislike_count, created_at
)
SELECT
    c.id,
    c.user_id,
    u.nick_name AS user_name,
    c.post_id,
    c.parent_id,
    c.content,
    c.like_count,
    c.dislike_count,
    c.created_at
FROM comments c
JOIN users u ON c.user_id = u.id;

DROP TABLE comments;
ALTER TABLE comments_old RENAME TO comments;

CREATE TABLE comment_feedback_old (
    user_id     INTEGER NOT NULL,
    parent_id   INTEGER NOT NULL,
    rating      INTEGER NOT NULL CHECK (rating IN (-1, 0, 1)),
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (user_id, parent_id),
    FOREIGN KEY (user_id)   REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES comments(id) ON DELETE CASCADE
);

INSERT INTO comment_feedback_old (
    user_id, parent_id, rating, created_at
)
SELECT
    user_id, parent_id, rating, created_at
FROM comment_feedback;

DROP TABLE comment_feedback;
ALTER TABLE comment_feedback_old RENAME TO comment_feedback;

CREATE TRIGGER update_like_count_after_insert_in_comment_feedback
AFTER INSERT ON comment_feedback
FOR EACH ROW
BEGIN
    UPDATE comments
    SET like_count = (
        SELECT COUNT(*) 
        FROM comment_feedback 
        WHERE parent_id = NEW.parent_id AND rating = 1)
    WHERE id = NEW.parent_id;
END;

CREATE TRIGGER update_like_count_after_update_in_comment_feedback
AFTER UPDATE ON comment_feedback
FOR EACH ROW
BEGIN
    UPDATE comments
    SET like_count = (
        SELECT COUNT(*) 
        FROM comment_feedback 
        WHERE parent_id = NEW.parent_id AND rating = 1)
    WHERE id = NEW.parent_id;
END;

CREATE TRIGGER update_dislike_count_after_insert_in_comment_feedback
AFTER INSERT ON comment_feedback
FOR EACH ROW
BEGIN
    UPDATE comments
    SET dislike_count = (
        SELECT COUNT(*) 
        FROM comment_feedback 
        WHERE parent_id = NEW.parent_id AND rating = -1)
    WHERE id = NEW.parent_id;
END;

CREATE TRIGGER update_dislike_count_after_update_in_comment_feeback
AFTER UPDATE ON comment_feedback
FOR EACH ROW
BEGIN
    UPDATE comments
    SET dislike_count = (
        SELECT COUNT(*) 
        FROM comment_feedback 
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
