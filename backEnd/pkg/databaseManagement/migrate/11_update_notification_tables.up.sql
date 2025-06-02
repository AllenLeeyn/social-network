ALTER TABLE notifications RENAME TO notifications_old;

CREATE TABLE notifications (
    id                      INTEGER PRIMARY KEY,
    to_user_id              INTEGER NOT NULL,
    from_user_id            INTEGER NOT NULL DEFAULT 0,
    target_id               INTEGER NOT NULL,
    target_uuid             TEXT NULL,
    target_type             TEXT NOT NULL CHECK (target_type IN ('following', 'groups', 'group_event')),
    target_detailed_type    TEXT NOT NULL CHECK (target_detailed_type IN (
                                'follow_request',
                                'follow_request_responded',
                                'group_invite',
                                'group_invite_responded',
                                'group_request',
                                'group_request_responded',
                                'group_event')),
    message                 TEXT NOT NULL,
    is_read                 INTEGER NOT NULL CHECK (is_read IN (0, 1)) DEFAULT 0,
    data                    TEXT CHECK (json_valid(data)),
    status                  TEXT NOT NULL CHECK (status IN ('enable', 'delete')) DEFAULT 'enable',
    created_at              DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by              INTEGER,
    updated_at              DATETIME,

    FOREIGN KEY (to_user_id)   REFERENCES users(id),
    FOREIGN KEY (from_user_id)   REFERENCES users(id),
    FOREIGN KEY (updated_by)      REFERENCES users(id)
);

INSERT INTO notifications (
    id, to_user_id, from_user_id, target_id, target_uuid, target_type,
    target_detailed_type, message, is_read, data, status,
    created_at, updated_by, updated_at
)
SELECT
    id,
    to_user_id,
    from_user_id,
    target_id,
    target_uuid,
    target_type,
    CASE
        WHEN target_detailed_type = 'follow_request_accepted' THEN 'follow_request_responded'
        ELSE target_detailed_type
    END,
    message,
    is_read,
    data,
    status,
    created_at,
    updated_by,
    updated_at
FROM notifications_old;

DROP TABLE notifications_old;
