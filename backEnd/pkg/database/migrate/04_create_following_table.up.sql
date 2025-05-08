CREATE TABLE following (
    leader_id INTEGER NOT NULL,
    follower_id INTEGER NOT NULL,
    group_id INTEGER,
    status TEXT NOT NULL CHECK(status IN ('requested', 'invited', 'accepted', 'declined', 'inactive')),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER NOT NULL,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER,
    PRIMARY KEY (leader_id, follower_id, group_id),
    FOREIGN KEY (leader_id) REFERENCES users(id),
    FOREIGN KEY (follower_id) REFERENCES users(id),
    FOREIGN KEY (group_id) REFERENCES groups(id),
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (updated_by) REFERENCES users(id) 
);
