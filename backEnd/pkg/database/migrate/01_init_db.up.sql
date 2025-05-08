CREATE TABLE IF NOT EXISTS account_type (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    name            TEXT NOT NULL UNIQUE,
    can_create_post BOOLEAN NOT NULL,
    can_comment     BOOLEAN NOT NULL,
    can_feedback    BOOLEAN NOT NULL,
    can_moderate    BOOLEAN NOT NULL,
    can_ban_user    BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    type_id     INTEGER NOT NULL,
    first_name  TEXT NOT NULL,
    last_name   TEXT NOT NULL,
    nick_name   TEXT NOT NULL UNIQUE,
    gender      TEXT CHECK(gender IN ('Male', 'Female', 'Other')) NOT NULL,
    age         INTEGER NOT NULL,
    email       TEXT NOT NULL UNIQUE,
    pw_hash     TEXT NOT NULL,
    reg_date    DATE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_login  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (type_id) REFERENCES account_type(id)
);

CREATE TABLE IF NOT EXISTS sessions (
    id UUID     PRIMARY KEY,
    user_id     INTEGER NOT NULL,
    is_active   BOOLEAN NOT NULL,
    start_time  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expire_time DATETIME NOT NULL,
    last_access DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS categories (
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    name    TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS posts (
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

CREATE TABLE IF NOT EXISTS post_categories (
    post_id     INTEGER NOT NULL,
    category_id INTEGER NOT NULL,

    PRIMARY KEY (post_id, category_id),
    FOREIGN KEY (post_id)       REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id)   REFERENCES categories(id)
);

CREATE TABLE IF NOT EXISTS comments (
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

CREATE TABLE IF NOT EXISTS post_feedback (
    user_id     INTEGER NOT NULL,
    parent_id   INTEGER NOT NULL,
    rating      INTEGER NOT NULL CHECK (rating IN (-1, 0, 1)),
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (user_id, parent_id),
    FOREIGN KEY (user_id)   REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES posts(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS comment_feedback (
    user_id     INTEGER NOT NULL,
    parent_id   INTEGER NOT NULL,
    rating      INTEGER NOT NULL CHECK (rating IN (-1, 0, 1)),
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (user_id, parent_id),
    FOREIGN KEY (user_id)   REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES comments(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS messages (
    id          INTEGER PRIMARY KEY,
    sender_id   INTEGER NOT NULL,
    receiver_id INTEGER NOT NULL,
    content     TEXT NOT NULL,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    read_at     DATETIME NULL,

    FOREIGN KEY (sender_id)     REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (receiver_id)   REFERENCES users(id) ON DELETE CASCADE
);

INSERT INTO account_type (id, name, can_create_post, can_comment, can_feedback, can_moderate, can_ban_user)
SELECT 1, 'user', 1, 1, 1, 0, 0
WHERE NOT EXISTS (SELECT 1 FROM account_type WHERE id = 1);

INSERT INTO account_type (id, name, can_create_post, can_comment, can_feedback, can_moderate, can_ban_user)
SELECT 2, 'moderator', 1, 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM account_type WHERE id = 2);

INSERT INTO account_type (id, name, can_create_post, can_comment, can_feedback, can_moderate, can_ban_user)
SELECT 3, 'administrator', 1, 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM account_type WHERE id = 3);

INSERT INTO categories (id, name)
SELECT 0, 'General'
WHERE NOT EXISTS (SELECT 1 FROM categories WHERE id = 0);

INSERT INTO categories (id, name)
SELECT 1, 'golang'
WHERE NOT EXISTS (SELECT 1 FROM categories WHERE id = 1);

INSERT INTO categories (id, name)
SELECT 2, 'html'
WHERE NOT EXISTS (SELECT 1 FROM categories WHERE id = 2);

INSERT INTO categories (id, name)
SELECT 3, 'css'
WHERE NOT EXISTS (SELECT 1 FROM categories WHERE id = 3);

INSERT INTO categories (id, name)
SELECT 4, 'sqlite3'
WHERE NOT EXISTS (SELECT 1 FROM categories WHERE id = 4);

INSERT INTO users (id, type_id, first_name, last_name, nick_name, gender, age, email, pw_hash)
SELECT 0, 1, 'ADMIN', 'ADMIN', 'ADMIN', 'Other', 0, 'admin@somewhere.earth', 's0m3thingS3cured?'
WHERE NOT EXISTS (SELECT 1 FROM users);

INSERT INTO posts (user_id, title, content)
SELECT 
    id,
    'Welcome. Please read and respect each other.',
    'Welcome!

We are so glad to have you here! Whether you are here to share your thoughts, ask questions, or simply learn, this is the perfect place to connect with like-minded individuals.

To make sure everyone has a positive experience, here are a few simple guidelines:

1. Be Respectful: Treat everyone with kindness and respect, even if you disagree. Healthy debate is welcome, but personal attacks are not.
2. Stay on Topic: Keep your posts relevant to the discussion. If you are starting a new topic, make sure it fits in the right category.
3. No Spam: We want to keep the forum focused and valuable for all users. Please avoid unsolicited promotions, ads, or irrelevant links.
4. Help Each Other: If you know the answer to a question, feel free to jump in and share! We are all here to learn from one another.

Feel free to explore the different sections, introduce yourself, and make yourself at home! If you need any help or have any questions, do not hesitate to ask our community or team.

Happy posting!'
FROM users
WHERE email = 'admin@somewhere.earth'
  AND NOT EXISTS (
    SELECT 1 FROM posts 
    WHERE user_id = (SELECT id FROM users WHERE email = 'admin@somewhere.earth')
      AND title = 'Welcome. Please read and respect each other.'
);

INSERT INTO post_categories (post_id, category_id)
SELECT 
    p.id,
    c.id
FROM posts p
JOIN users u ON p.user_id = u.id
JOIN categories c ON c.name = 'General'
WHERE u.email = 'admin@somewhere.earth'
  AND NOT EXISTS (
    SELECT 1 FROM post_categories 
    WHERE post_id = p.id AND category_id = c.id
);

DROP VIEW IF EXISTS v_posts;

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

-- Comment count after new comment
DROP TRIGGER IF EXISTS update_comment_count_after_insert;
CREATE TRIGGER update_comment_count_after_insert
AFTER INSERT ON comments
FOR EACH ROW
BEGIN
    UPDATE posts
    SET comment_count = (
        SELECT COUNT(*) 
        FROM comments 
        WHERE post_id = NEW.post_id
    )
    WHERE id = NEW.post_id;
END;

-- Comment like count after insert
DROP TRIGGER IF EXISTS update_like_count_after_insert_in_comment_feedback;
CREATE TRIGGER update_like_count_after_insert_in_comment_feedback
AFTER INSERT ON comment_feedback
FOR EACH ROW
BEGIN
    UPDATE comments
    SET like_count = (
        SELECT COUNT(*) 
        FROM comment_feedback 
        WHERE parent_id = NEW.parent_id AND rating = 1
    )
    WHERE id = NEW.parent_id;
END;

-- Comment like count after update
DROP TRIGGER IF EXISTS update_like_count_after_update_in_comment_feeback;
CREATE TRIGGER update_like_count_after_update_in_comment_feedback
AFTER UPDATE ON comment_feedback
FOR EACH ROW
BEGIN
    UPDATE comments
    SET like_count = (
        SELECT COUNT(*) 
        FROM comment_feedback 
        WHERE parent_id = NEW.parent_id AND rating = 1
    )
    WHERE id = NEW.parent_id;
END;

-- Post like count after insert
DROP TRIGGER IF EXISTS update_like_count_after_insert_in_post_feedback;
CREATE TRIGGER update_like_count_after_insert_in_post_feedback
AFTER INSERT ON post_feedback
FOR EACH ROW
BEGIN
    UPDATE posts
    SET like_count = (
        SELECT COUNT(*) 
        FROM post_feedback 
        WHERE parent_id = NEW.parent_id AND rating = 1
    )
    WHERE id = NEW.parent_id;
END;

-- Post like count after update
DROP TRIGGER IF EXISTS update_like_count_after_update_in_post_feedback;
CREATE TRIGGER update_like_count_after_update_in_post_feedback
AFTER UPDATE ON post_feedback
FOR EACH ROW
BEGIN
    UPDATE posts
    SET like_count = (
        SELECT COUNT(*) 
        FROM post_feedback 
        WHERE parent_id = NEW.parent_id AND rating = 1
    )
    WHERE id = NEW.parent_id;
END;

-- Comment dislike count after insert
DROP TRIGGER IF EXISTS update_dislike_count_after_insert_in_comment_feedback;
CREATE TRIGGER update_dislike_count_after_insert_in_comment_feedback
AFTER INSERT ON comment_feedback
FOR EACH ROW
BEGIN
    UPDATE comments
    SET dislike_count = (
        SELECT COUNT(*) 
        FROM comment_feedback 
        WHERE parent_id = NEW.parent_id AND rating = -1
    )
    WHERE id = NEW.parent_id;
END;

-- Comment dislike count after update
DROP TRIGGER IF EXISTS update_dislike_count_after_update_in_comment_feeback;
CREATE TRIGGER update_dislike_count_after_update_in_comment_feedback
AFTER UPDATE ON comment_feedback
FOR EACH ROW
BEGIN
    UPDATE comments
    SET dislike_count = (
        SELECT COUNT(*) 
        FROM comment_feedback 
        WHERE parent_id = NEW.parent_id AND rating = -1
    )
    WHERE id = NEW.parent_id;
END;

-- Post dislike count after insert
DROP TRIGGER IF EXISTS update_dislike_count_after_insert_in_post_feedback;
CREATE TRIGGER update_dislike_count_after_insert_in_post_feedback
AFTER INSERT ON post_feedback
FOR EACH ROW
BEGIN
    UPDATE posts
    SET dislike_count = (
        SELECT COUNT(*) 
        FROM post_feedback 
        WHERE parent_id = NEW.parent_id AND rating = -1
    )
    WHERE id = NEW.parent_id;
END;

-- Post dislike count after update
DROP TRIGGER IF EXISTS update_dislike_count_after_update_in_post_feedback;
CREATE TRIGGER update_dislike_count_after_update_in_post_feedback
AFTER UPDATE ON post_feedback
FOR EACH ROW
BEGIN
    UPDATE posts
    SET dislike_count = (
        SELECT COUNT(*) 
        FROM post_feedback 
        WHERE parent_id = NEW.parent_id AND rating = -1
    )
    WHERE id = NEW.parent_id;
END;
