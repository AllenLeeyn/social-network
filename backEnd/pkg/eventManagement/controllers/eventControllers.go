package controller

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"social-network/pkg/middleware"
	"social-network/pkg/utils"

	errorControllers "social-network/pkg/errorManagement/controllers"
	eventModel "social-network/pkg/eventManagement/models"
	groupModel "social-network/pkg/groupManagement/models"
	userControllers "social-network/pkg/userManagement/controllers"

	_ "github.com/mattn/go-sqlite3"
)

type event = eventModel.Event
type eventRepsonse = eventModel.EventResponse

func isValidEventInfo(e *event) error {
	isOk := false
	if e.Location, isOk = utils.IsValidContent(e.Location, 3, 100); !isOk {
		return fmt.Errorf("location must be between 3 to 100 characters")
	}
	if e.Title, isOk = utils.IsValidContent(e.Title, 3, 100); !isOk {
		return fmt.Errorf("title must be between 3 to 100 characters")
	}
	if e.Description, isOk = utils.IsValidContent(e.Description, 10, 1000); !isOk {
		return fmt.Errorf("description must be between 10 to 1000 characters")
	}
	if e.StartTime == nil {
		return fmt.Errorf("start time is required")
	}
	if e.StartTime.Before(time.Now()) {
		return fmt.Errorf("start time cannot be in the past")
	}
	if e.DurationMin <= 0 {
		return fmt.Errorf("duration must be greater than 0")
	}
	return nil
}

func EventCreateHandler(w http.ResponseWriter, r *http.Request) {
	e := &event{}
	if err := utils.ReadJSON(r, e); err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	userID, isOk := middleware.GetUserID(r.Context())
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	e.CreatedBy = userID

	if err := isValidEventInfo(e); err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), http.StatusBadRequest)
		return
	}
	groupID, _, err := groupModel.SelectGroupIDcreatedByfromUUID(e.GroupUUID)
	if err != nil {
		errorControllers.CustomErrorHandler(w, r,
			"group not found", http.StatusBadRequest)
		return
	}
	if !groupModel.IsGroupMember(e.GroupUUID, userID) {
		errorControllers.CustomErrorHandler(w, r,
			"only group members can create event", http.StatusBadRequest)
		return
	}
	e.GroupID = groupID

	eventID, eventUUID, err := eventModel.InsertEvent(e)
	if err != nil {
		log.Println(err.Error())
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	err = eventModel.InsertEventResponse(&eventRepsonse{
		EventID: eventID, Response: "accepted", CreatedBy: userID})
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, "event created",
		struct {
			EventUUID string `json:"event_uuid"`
		}{eventUUID})
}

func EventUpdateHandler(w http.ResponseWriter, r *http.Request) {}

func ViewEventsHandler(w http.ResponseWriter, r *http.Request) {}

func ViewEventHandler(w http.ResponseWriter, r *http.Request) {}

func EventResponseHandler(w http.ResponseWriter, r *http.Request) {}

func EventResponsesHandler(w http.ResponseWriter, r *http.Request) {}
