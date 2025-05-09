CREATE TABLE messages_new (
    id          INTEGER PRIMARY KEY,
    sender_id   INTEGER NOT NULL,
    receiver_id INTEGER NOT NULL,
    group_id    INTEGER,
    content     TEXT NOT NULL,
    
    status      TEXT NOT NULL CHECK ("status" IN ('enable', 'disable', 'delete')) DEFAULT 'enable',
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    read_at     DATETIME,
    updated_by  INTEGER,
    updated_at  DATETIME,

    FOREIGN KEY (sender_id)     REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (receiver_id)   REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (group_id)      REFERENCES groups(id) ON DELETE SET NULL,
    FOREIGN KEY (updated_by)    REFERENCES users(id)
);

INSERT INTO messages_new (
    id, sender_id, receiver_id, content, created_at, read_at, updated_by, updated_at
)
SELECT
    id, sender_id, receiver_id, content, created_at, read_at, 0, CURRENT_TIMESTAMP
FROM messages;

DROP TABLE messages;
ALTER TABLE messages_new RENAME TO messages;
