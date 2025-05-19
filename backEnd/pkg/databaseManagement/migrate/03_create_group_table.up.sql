CREATE TABLE groups (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    uuid            TEXT NOT NULL UNIQUE,
    title           TEXT NOT NULL,
    description     TEXT,
    banner_image    TEXT,
    members_count   INTEGER NOT NULL DEFAULT 1,

    status      TEXT NOT NULL CHECK ("status" IN ('enable', 'disable', 'delete')) DEFAULT 'enable',
    created_by  INTEGER NOT NULL,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by  INTEGER,
    updated_at  DATETIME,

    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (updated_by) REFERENCES users(id)
);

CREATE TABLE group_events (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id            INTEGER NOT NULL,
    location            TEXT NOT NULL,
    start_time          DATETIME NOT NULL,
    duration_minutes    INTEGER NOT NULL,
    title               TEXT NOT NULL,
    description         TEXT,
    event_image         TEXT,

    status      TEXT NOT NULL DEFAULT 'scheduled' CHECK(status IN ('scheduled', 'cancelled', 'completed')),
    created_by  INTEGER NOT NULL,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by  INTEGER,
    updated_at  DATETIME,

    FOREIGN KEY (group_id)      REFERENCES groups(id),
    FOREIGN KEY (created_by)    REFERENCES users(id),
    FOREIGN KEY (updated_by)    REFERENCES users(id)
);

CREATE TABLE group_event_responses (
    event_id    INTEGER NOT NULL,
    response    TEXT NOT NULL CHECK(response IN ( 'accepted', 'declined')),
    created_by  INTEGER NOT NULL,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by  INTEGER,
    updated_at  DATETIME,

    PRIMARY KEY (event_id, created_by),
    FOREIGN KEY (event_id)      REFERENCES group_events(id),
    FOREIGN KEY (created_by)    REFERENCES users(id),
    FOREIGN KEY (updated_by)    REFERENCES users(id)
);
INSERT INTO groups (
    uuid,
    title,
    description,
    banner_image,
    status,
    created_by
) VALUES (
    '00000000-0000-0000-0000-000000000000',
    'Public',
    'Please be kind and respect one another',
    NULL,
    'enable',
    0
);
