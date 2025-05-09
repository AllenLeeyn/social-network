DROP VIEW IF EXISTS v_posts;

CREATE TABLE users_old (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    type_id     INTEGER NOT NULL,
    first_name  TEXT NOT NULL,
    last_name   TEXT NOT NULL,
    nick_name   TEXT NOT NULL UNIQUE,
    gender      TEXT CHECK(gender IN ('Male', 'Female', 'Other')) NOT NULL,
    birthday    DATETIME NOT NULL,
    email       TEXT NOT NULL UNIQUE,
    pw_hash     TEXT NOT NULL,
    reg_date    DATE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_login  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (type_id) REFERENCES account_type(id)
);

INSERT INTO users_old (
    id, type_id, first_name, last_name, nick_name, gender, birthday,
    email, pw_hash, reg_date
)
SELECT
    id,
    type_id,
    first_name,
    last_name,
    COALESCE(nick_name, email) AS nick_name,
    gender,
    birthday,
    email,
    pw_hash,
    created_at
FROM users;

DROP TABLE users;
ALTER TABLE users_old RENAME TO users;

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
