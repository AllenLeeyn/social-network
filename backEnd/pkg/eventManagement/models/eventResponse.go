package models

import "time"

type EventResponse struct {
	EventID   int        `json:"event_id"`
	Response  string     `json:"response"`
	CreatedBy int        `json:"created_by"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedBy int        `json:"updated_by"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func InsertEventResponse(eReponse *EventResponse) error {

	qry := `INSERT INTO group_event_responses (
				event_id, response, created_by
			) VALUES (?, ?, ?);`
	_, err := sqlDB.Exec(qry,
		&eReponse.EventID,
		&eReponse.Response,
		&eReponse.CreatedBy)
	if err != nil {
		return err
	}
	return nil
}

func SelectEventResponse() {}
