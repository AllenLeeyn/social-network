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
	GroupID     string     `json:"group_id"`
	Location    string     `json:"location"`
	StartTime   *time.Time `json:"start_time"`
	DurationMin int        `json:"duration_minutes"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	EventImage  string     `json:"event_image"`
	CreatorName string     `json:"creator_name"`
	CreatorUUID string     `json:"creator_uuid"`
	Status      string     `json:"status"`
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

func UpdateEvent() {}

func SelectEvents() {}

func SelectEvent() {}
