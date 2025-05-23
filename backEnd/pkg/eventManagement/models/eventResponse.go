package models

import (
	"time"
)

type EventResponse struct {
	EventID   string     `json:"event_id"`
	EventUUID string     `json:"event_uuid"`
	Response  string     `json:"response"`
	CreatedBy int        `json:"created_by"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedBy int        `json:"updated_by"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type EventResponseView struct {
	Response      string     `json:"response"`
	CreatedByUUID string     `json:"created_by_uuid,omitempty"`
	CreatedByName string     `json:"created_by_name,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
}

func InsertEventResponse(eResponse *EventResponse) error {
	qry := `INSERT INTO group_event_responses (
				event_id, response, created_by
			) SELECT id, ?, ? FROM group_events WHERE uuid = ?;`
	_, err := sqlDB.Exec(qry,
		eResponse.Response,
		eResponse.CreatedBy,
		eResponse.EventUUID)
	return err
}

func UpdateEventResponse(eResponse *EventResponse) error {
	qry := `UPDATE group_event_responses
			SET response = ?, updated_by = ?, updated_at = CURRENT_TIMESTAMP
			WHERE event_id = (SELECT id FROM group_events WHERE uuid = ?)
			AND created_by = ?;`
	_, err := sqlDB.Exec(qry,
		eResponse.Response,
		eResponse.UpdatedBy,
		eResponse.EventUUID,
		eResponse.CreatedBy)
	return err
}

func SelectStatus(eReponse *EventResponse) (string, error) {
	qry := `SELECT ger.response
			FROM group_event_responses ger
			JOIN group_events ge ON ger.event_id = ge.id
			WHERE ger.created_by = ? AND ge.uuid = ?`

	var status string
	err := sqlDB.QueryRow(qry, eReponse.CreatedBy, eReponse.EventUUID).Scan(&status)
	return status, checkErrNoRows(err)
}

func SelectEventResponses(eventUUID string) (*[]EventResponseView, error) {
	qry := `SELECT ger.response, u.uuid, u.nick_name, ger.created_at
			FROM group_event_responses ger
			JOIN group_events ge ON ger.event_id = ge.id
			JOIN users u ON ger.created_by = u.id
			WHERE ge.uuid = ?
			ORDER BY ger.created_at ASC;`

	rows, err := sqlDB.Query(qry, eventUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var responses []EventResponseView
	for rows.Next() {
		var res EventResponseView
		err := rows.Scan(&res.Response, &res.CreatedByUUID, &res.CreatedByName, &res.CreatedAt)
		if err != nil {
			return nil, err
		}
		responses = append(responses, res)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &responses, nil
}
