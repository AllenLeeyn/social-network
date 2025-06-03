DROP TRIGGER IF EXISTS update_like_count_after_insert_in_comment_feedback;
DROP TRIGGER IF EXISTS update_like_count_after_update_in_comment_feedback;
DROP TRIGGER IF EXISTS update_dislike_count_after_insert_in_comment_feedback;
DROP TRIGGER IF EXISTS update_dislike_count_after_update_in_comment_feedback;

CREATE TABLE comments_new (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id         INTEGER NOT NULL,
    post_id         INTEGER NOT NULL,
    parent_id       INTEGER,
    content         TEXT NOT NULL DEFAULT '',   
    attached_image  TEXT DEFAULT "",
    like_count      INTEGER NOT NULL DEFAULT 0,
    dislike_count   INTEGER NOT NULL DEFAULT 0,

    status          TEXT NOT NULL CHECK ("status" IN ('enable', 'disable', 'delete')) DEFAULT 'enable',
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by      INTEGER,
    updated_at      DATETIME,

    FOREIGN KEY (user_id)       REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id)       REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id)     REFERENCES comments(id),
    FOREIGN KEY (updated_by)    REFERENCES users(id)
);

INSERT INTO comments_new (
    id, user_id, post_id, parent_id, content,
    like_count, dislike_count, created_at, updated_by, updated_at
)
SELECT
    id, user_id, post_id, parent_id, content,
    like_count, dislike_count, created_at, 0, CURRENT_TIMESTAMP
FROM comments;

DROP TABLE comments;
ALTER TABLE comments_new RENAME TO comments;

CREATE TABLE comment_feedback_new (
    user_id     INTEGER NOT NULL,
    parent_id   INTEGER NOT NULL,
    rating      INTEGER NOT NULL CHECK (rating IN (-1, 0, 1)),
    
    status      TEXT NOT NULL CHECK ("status" IN ('enable', 'delete')) DEFAULT 'enable',
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by  INTEGER,
    updated_at  DATETIME,

    PRIMARY KEY (user_id, parent_id),
    FOREIGN KEY (user_id)       REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id)     REFERENCES comments(id) ON DELETE CASCADE,
    FOREIGN KEY (updated_by)    REFERENCES users(id)
);

INSERT INTO comment_feedback_new (
    user_id, parent_id, rating, created_at, updated_by, updated_at
)
SELECT
    user_id, parent_id, rating, created_at, 0, CURRENT_TIMESTAMP
FROM comment_feedback;

DROP TABLE comment_feedback;
ALTER TABLE comment_feedback_new RENAME TO comment_feedback;

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
