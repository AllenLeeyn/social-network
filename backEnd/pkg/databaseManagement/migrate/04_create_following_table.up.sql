CREATE TABLE following (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    leader_id       INTEGER NOT NULL,
    follower_id     INTEGER NOT NULL,
    group_id        INTEGER NOT NULL DEFAULT 0,

    type        TEXT NOT NULL CHECK(type IN ('group', 'user')) DEFAULT 'user',
    status      TEXT NOT NULL CHECK(status IN ('requested', 'invited', 'accepted', 'declined', 'inactive')),
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by  INTEGER NOT NULL,
    updated_at  DATETIME,
    updated_by  INTEGER,

    UNIQUE (leader_id, follower_id, group_id),
    FOREIGN KEY (leader_id)     REFERENCES users(id),
    FOREIGN KEY (follower_id)   REFERENCES users(id),
    FOREIGN KEY (group_id)      REFERENCES groups(id),
    FOREIGN KEY (created_by)    REFERENCES users(id),
    FOREIGN KEY (updated_by)    REFERENCES users(id)
);

CREATE TRIGGER update_group_members_count_on_accepted
AFTER UPDATE OF status ON following
WHEN NEW.status = 'accepted' AND OLD.status != 'accepted'
BEGIN
    UPDATE groups
    SET members_count = members_count + 1
    WHERE id = NEW.group_id;
END;

CREATE TRIGGER update_group_members_count_on_not_accepted
AFTER UPDATE OF status ON following
WHEN OLD.status = 'accepted' AND NEW.status != 'accepted'
BEGIN
    UPDATE groups
    SET members_count = members_count - 1
    WHERE id = NEW.group_id;
END;
