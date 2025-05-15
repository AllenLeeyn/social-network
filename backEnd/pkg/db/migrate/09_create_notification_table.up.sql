CREATE TABLE notifications (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    receiver_id  INTEGER NOT NULL,
    sender_id     INTEGER NOT NULL,
    group_id      INTEGER,
    event_id      INTEGER,
    
    type          TEXT NOT NULL CHECK (type IN (
                      'follow_request',
                      'group_invite',
                      'group_join_request',
                      'group_event',
                      'request_accepted'
                  )),
    
    message       TEXT,
    read_at       DATETIME,
    created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_by    INTEGER,
    updated_at    DATETIME,

    FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (sender_id)    REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (group_id)     REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY (event_id)     REFERENCES events(id) ON DELETE CASCADE
);
