DROP VIEW IF EXISTS v_posts;

CREATE TABLE users_new (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    type_id         INTEGER NOT NULL,
    first_name      TEXT NOT NULL,
    last_name       TEXT NOT NULL,
    gender          TEXT CHECK(gender IN ('Male', 'Female', 'Other')) NOT NULL,
    birthday        DATETIME NOT NULL,
    email           TEXT NOT NULL UNIQUE,
    pw_hash         TEXT NOT NULL,
    nick_name       TEXT UNIQUE,
    profile_image   TEXT UNIQUE,
    about_me        TEXT,
    visibility      TEXT NOT NULL DEFAULT 'private' CHECK(visibility IN ('private', 'public')),

    status      TEXT NOT NULL CHECK ("status" IN ('enable', 'disable', 'delete')) DEFAULT 'enable',
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by  INTEGER,
    updated_at  DATETIME,

    FOREIGN KEY (type_id)       REFERENCES account_type(id),
    FOREIGN KEY (updated_by)    REFERENCES users(id)
);

INSERT INTO users_new (
    id, type_id, first_name, last_name, nick_name, gender, birthday,
    email, pw_hash, created_at, updated_by, updated_at
)
SELECT
    id,
    type_id,
    first_name,
    last_name,
    nick_name,
    gender,
    DATE('now', '-' || age || ' years'),
    email,
    pw_hash,
    reg_date,
    0,
    CURRENT_TIMESTAMP
FROM users;

DROP TABLE users;
ALTER TABLE users_new RENAME TO users;

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
