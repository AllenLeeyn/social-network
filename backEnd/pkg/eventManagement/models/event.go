package models

import (
	"social-network/pkg/utils"
	"time"
)

type Event struct {
	ID          int        `json:"id"`
	UUID        string     `json:"uuid"`
	GroupID     int        `json:"group_id"`
	GroupUUID   string     `json:"group_uuid"`
	Location    string     `json:"location"`
	StartTime   *time.Time `json:"start_time"`
	DurationMin int        `json:"duration_minutes"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	EventImage  string     `json:"event_image"`
	Status      string     `json:"status"`
	CreatedBy   int        `json:"created_by"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedBy   int        `json:"updated_by"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type EventView struct {
	UUID        string     `json:"uuid"`
	GroupUUID   string     `json:"group_uuid"`
	GroupTitle  string     `json:"group_title"`
	Location    string     `json:"location"`
	StartTime   *time.Time `json:"start_time"`
	DurationMin int        `json:"duration_minutes"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	EventImage  string     `json:"event_image"`
	AttendCount int        `json:"attend_count"`
	CreatorName string     `json:"creator_name"`
	CreatorUUID string     `json:"creator_uuid"`
	Status      string     `json:"status"`
}

func IsEventCreator(eventUUID string, userID int) bool {
	qry := `SELECT 1
			FROM group_events
			WHERE uuid = ? AND created_by = ?
			LIMIT 1;`

	var exists int
	err := sqlDB.QueryRow(qry, eventUUID, userID).Scan(&exists)
	return err == nil
}

func IsGroupMemberFromEventUUID(eventUUID string, userID int) bool {
	qry := `SELECT 1
			FROM following f
			JOIN groups g ON f.group_id = g.id
			JOIN group_events ge ON ge.group_id = g.id
			WHERE ge.uuid = ? AND f.follower_id = ? AND f.status = 'accepted'
			LIMIT 1;`

	var exists int
	err := sqlDB.QueryRow(qry, eventUUID, userID).Scan(&exists)
	return err == nil
}

func InsertEvent(event *Event) (int, string, error) {
	uuid, err := utils.GenerateUuid()
	if err != nil {
		return -1, "", err
	}
	event.UUID = uuid

	qry := `INSERT INTO group_events (
				uuid, group_id, location,
				start_time, duration_minutes,
				title, description, event_image, 
				created_by
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`
	result, err := sqlDB.Exec(qry,
		event.UUID, event.GroupID, event.Location,
		event.StartTime, event.DurationMin,
		event.Title, event.Description, event.EventImage,
		event.CreatedBy)
	if err != nil {
		return -1, "", err
	}

	eventID, err := result.LastInsertId()
	if err != nil {
		return -1, "", err
	}
	return int(eventID), event.UUID, err
}

func UpdateEvent(event *Event) error {
	updateQuery := `
		UPDATE group_events
		SET title = ?,	description = ?, event_image =?,
			start_time = ?, duration_minutes = ?, location = ?,
			updated_by = ?, updated_at = CURRENT_TIMESTAMP
		WHERE uuid = ?;`

	_, err := sqlDB.Exec(updateQuery,
		event.Title, event.Description, event.EventImage,
		event.StartTime, event.DurationMin, event.Location,
		event.UpdatedBy,
		event.UUID,
	)
	return err
}

func SelectEvents(groupUUID string, userID int) (*[]EventView, error) {
	qry := `SELECT
				ge.uuid, g.uuid, g.title,
				ge.location, ge.start_time, ge.duration_minutes,
				ge.title, ge.description, ge.event_image,
				( SELECT COUNT(*) FROM group_event_responses ger
				  WHERE ger.event_id = ge.id AND ger.response = 'accepted'
				) AS attend_count,
				u.nick_name,
				u.uuid,
    			COALESCE(ger.response, '') AS response
			FROM group_events ge
			JOIN groups g ON ge.group_id = g.id
			JOIN users u ON g.created_by = u.id
			LEFT JOIN group_event_responses ger ON ger.event_id = ge.id
				AND ger.created_by = ?
			WHERE g.uuid = ?
			AND g.status = 'enable';`

	rows, err := sqlDB.Query(qry, userID, groupUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []EventView
	for rows.Next() {
		var e EventView
		err := rows.Scan(
			&e.UUID, &e.GroupUUID, &e.GroupTitle,
			&e.Location, &e.StartTime, &e.DurationMin,
			&e.Title, &e.Description, &e.EventImage,
			&e.AttendCount,
			&e.CreatorName, &e.CreatorUUID,
			&e.Status)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &events, nil
}

func SelectEvent(eventUUID string, userID int) (*EventView, error) {
	qry := `SELECT
				ge.uuid, g.uuid, g.title,
				ge.location, ge.start_time, ge.duration_minutes,
				ge.title, ge.description, ge.event_image,
				( SELECT COUNT(*) FROM group_event_responses ger
				  WHERE ger.event_id = ge.id AND ger.response = 'accepted'
				) AS attend_count,
				u.nick_name,
				u.uuid,
    			COALESCE(ger.response, '') AS response
			FROM group_events ge
			JOIN groups g ON ge.group_id = g.id
			JOIN users u ON g.created_by = u.id
			LEFT JOIN group_event_responses ger ON ger.event_id = ge.id
				AND ger.created_by = ?
			WHERE ge.uuid = ?
			AND g.status = 'enable';`

	var e EventView
	err := sqlDB.QueryRow(qry, userID, eventUUID).Scan(
		&e.UUID, &e.GroupUUID, &e.GroupTitle,
		&e.Location, &e.StartTime, &e.DurationMin,
		&e.Title, &e.Description, &e.EventImage,
		&e.AttendCount,
		&e.CreatorName, &e.CreatorUUID,
		&e.Status)
	if err != nil {
		return nil, err
	}
	return &e, nil
}
