CREATE TABLE messages_old (
    id          INTEGER PRIMARY KEY,
    sender_id   INTEGER NOT NULL,
    receiver_id INTEGER NOT NULL,
    content     TEXT NOT NULL,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    read_at     DATETIME,

    FOREIGN KEY (sender_id)     REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (receiver_id)   REFERENCES users(id) ON DELETE CASCADE
);

INSERT INTO messages_old (
    id, sender_id, receiver_id, content, created_at, read_at
)
SELECT
    id, sender_id, receiver_id, content, created_at, read_at
FROM messages;

DROP TABLE messages;
ALTER TABLE messages_old RENAME TO messages;
