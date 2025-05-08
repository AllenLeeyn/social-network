CREATE TABLE message_files (
    id                 INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id           INTEGER,
    message_id         INTEGER NOT NULL,
    file_uploaded_name TEXT NOT NULL,
    file_real_name     TEXT NOT NULL,

    status      TEXT NOT NULL CHECK (status IN ('enable', 'delete')) DEFAULT 'enable',
    created_by  INTEGER NOT NULL,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by  INTEGER,
    updated_at  DATETIME,

    FOREIGN KEY (group_id)     REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY (message_id)  REFERENCES messages(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by)  REFERENCES users(id),
    FOREIGN KEY (updated_by)  REFERENCES users(id)
);

CREATE TABLE post_files (
    id                 INTEGER PRIMARY KEY,
    post_id            INTEGER NOT NULL,
    file_uploaded_name TEXT NOT NULL,
    file_real_name     TEXT NOT NULL,

    status      TEXT NOT NULL CHECK (status IN ('enable', 'delete')) DEFAULT 'enable',
    created_by  INTEGER NOT NULL,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by  INTEGER,
    updated_at  DATETIME,

    FOREIGN KEY (post_id)    REFERENCES posts(id),
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (updated_by) REFERENCES users(id)
);
