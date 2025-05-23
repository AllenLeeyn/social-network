CREATE TABLE post_selected_audience (
    post_id     INTEGER NOT NULL,
    user_id     INTEGER NOT NULL,
    
    status      TEXT NOT NULL CHECK ("status" IN ('enable', 'disable', 'delete')) DEFAULT 'enable',
    created_by  INTEGER NOT NULL,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by  INTEGER,
    updated_at  DATETIME,

    PRIMARY KEY (post_id, user_id),
    FOREIGN KEY (post_id)       REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id)       REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by)    REFERENCES users(id),
    FOREIGN KEY (updated_by)    REFERENCES users(id)
);
