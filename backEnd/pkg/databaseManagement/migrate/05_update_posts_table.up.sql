DROP VIEW IF EXISTS v_posts;
DROP TRIGGER IF EXISTS update_like_count_after_insert_in_post_feedback;
DROP TRIGGER IF EXISTS update_like_count_after_update_in_post_feedback;
DROP TRIGGER IF EXISTS update_dislike_count_after_insert_in_post_feedback;
DROP TRIGGER IF EXISTS update_dislike_count_after_update_in_post_feedback;
DROP TRIGGER IF EXISTS update_comment_count_after_insert;

CREATE TABLE posts_new (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    uuid            TEXT NOT NULL UNIQUE,
    user_id         INTEGER NOT NULL,
    group_id        INTEGER,
    title           TEXT NOT NULL,
    content         TEXT NOT NULL,
    visibility      TEXT NOT NULL DEFAULT 'public' CHECK (visibility IN ('public', 'private', 'selected')),
    like_count      INTEGER NOT NULL DEFAULT 0,
    dislike_count   INTEGER NOT NULL DEFAULT 0,
    comment_count   INTEGER NOT NULL DEFAULT 0,

    status      TEXT NOT NULL CHECK ("status" IN ('enable', 'disable', 'delete')) DEFAULT 'enable',
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by  INTEGER,
    updated_at  DATETIME,

    FOREIGN KEY (user_id)       REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (group_id)      REFERENCES groups(id) ON DELETE SET NULL,
    FOREIGN KEY (updated_by)    REFERENCES users(id)
);

INSERT INTO posts_new (
    id, uuid, user_id, title, content, created_at, visibility, group_id,
    like_count, dislike_count, comment_count, updated_by, updated_at
)
SELECT
    p.id, pu.uuid, p.user_id, p.title, p.content, p.created_at, 'public', NULL,
    p.like_count, p.dislike_count, p.comment_count, 0, CURRENT_TIMESTAMP
FROM posts p
JOIN post_uuids pu ON pu.post_id = p.id;;

DROP TABLE post_uuids;
DROP TABLE posts;
ALTER TABLE posts_new RENAME TO posts;

CREATE TABLE post_feedback_new (
    user_id     INTEGER NOT NULL,
    parent_id   INTEGER NOT NULL,
    rating      INTEGER NOT NULL CHECK (rating IN (-1, 0, 1)),

    status      TEXT NOT NULL CHECK ("status" IN ('enable', 'delete')) DEFAULT 'enable',
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by  INTEGER,
    updated_at  DATETIME,

    PRIMARY KEY (user_id, parent_id),
    FOREIGN KEY (user_id)       REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id)     REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (updated_by)    REFERENCES users(id)
);

INSERT INTO post_feedback_new (
    user_id, parent_id, rating, created_at, updated_by, updated_at
)
SELECT
    user_id, parent_id, rating, created_at, 0, CURRENT_TIMESTAMP
FROM post_feedback;

DROP TABLE post_feedback;
ALTER TABLE post_feedback_new RENAME TO post_feedback;

CREATE TABLE post_categories_new (
    post_id     INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    
    status      TEXT NOT NULL CHECK (status IN ('enable', 'disable', 'delete')) DEFAULT 'enable',
    created_by  INTEGER NOT NULL,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by  INTEGER,
    updated_at  DATETIME,

    PRIMARY KEY (post_id, category_id),
    FOREIGN KEY (post_id)       REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id)   REFERENCES categories(id),
    FOREIGN KEY (created_by)    REFERENCES users(id),
    FOREIGN KEY (updated_by)    REFERENCES users(id)
);

INSERT INTO post_categories_new (
    post_id, category_id, status, created_by, created_at, updated_by, updated_at
)
SELECT
    post_id, category_id, 'enable', 0, CURRENT_TIMESTAMP, 0, CURRENT_TIMESTAMP
FROM post_categories;

DROP TABLE post_categories;
ALTER TABLE post_categories_new RENAME TO post_categories;

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
    p.visibility,
    p.group_id,
    GROUP_CONCAT(pc.category_id) AS category_ids
FROM posts p
INNER JOIN users u ON p.user_id = u.id
LEFT JOIN post_categories pc ON p.id = pc.post_id
GROUP BY 
    p.id, p.user_id, u.nick_name,
    p.comment_count, p.like_count, p.dislike_count,
    p.title, p.content, p.created_at, p.visibility, p.group_id;

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
