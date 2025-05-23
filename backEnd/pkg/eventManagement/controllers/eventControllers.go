package controller

import (
	"fmt"
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
type eventResponse = eventModel.EventResponse

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

func parseEventRequest(r *http.Request) (*event, int, error) {
	e := &event{}
	if err := utils.ReadJSON(r, e); err != nil {
		return nil, -1, fmt.Errorf("invalid input")
	}
	userID, isOk := middleware.GetUserID(r.Context())
	if !isOk {
		return nil, -1, fmt.Errorf("InternalServerError")
	}

	if err := isValidEventInfo(e); err != nil {
		return nil, -1, err
	}

	groupID, _, err := groupModel.SelectGroupIDcreatedByfromUUID(e.GroupUUID)
	if err != nil {
		return nil, -1, fmt.Errorf("group not found")
	}
	if !groupModel.IsGroupMember(e.GroupUUID, userID) {
		return nil, -1, fmt.Errorf("only group members allowed")
	}
	e.GroupID = groupID

	return e, userID, nil
}

func EventCreateHandler(w http.ResponseWriter, r *http.Request) {
	e, userID, err := parseEventRequest(r)
	if err != nil {
		if err.Error() == "InternalServerError" {
			errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		} else {
			errorControllers.CustomErrorHandler(w, r, err.Error(), http.StatusBadRequest)
		}
		return
	}
	e.CreatedBy = userID

	_, eventUUID, err := eventModel.InsertEvent(e)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	err = eventModel.InsertEventResponse(&eventResponse{
		EventUUID: eventUUID, Response: "accepted", CreatedBy: userID})
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

func EventUpdateHandler(w http.ResponseWriter, r *http.Request) {
	e, userID, err := parseEventRequest(r)
	if err != nil {
		if err.Error() == "InternalServerError" {
			errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		} else {
			errorControllers.CustomErrorHandler(w, r, err.Error(), http.StatusBadRequest)
		}
		return
	}
	e.UpdatedBy = userID

	if !eventModel.IsEventCreator(e.UUID, userID) {
		errorControllers.ErrorHandler(w, r, errorControllers.ForbiddenError)
		return
	}

	if err := eventModel.UpdateEvent(e); err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, "group event updated", nil)
}

func ViewEventsHandler(w http.ResponseWriter, r *http.Request) {
	userID, isOk := middleware.GetUserID(r.Context())
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	groupUUID, err := utils.ExtractUUIDFromUrl(r.URL.Path, "api/group/events")
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.BadRequestError)
		return
	}
	if !groupModel.IsGroupMember(groupUUID, userID) {
		errorControllers.CustomErrorHandler(w, r,
			"only members can view", http.StatusBadRequest)
		return
	}

	events, err := eventModel.SelectEvents(groupUUID, userID)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, "event list retrieved", events)
}

func ViewEventHandler(w http.ResponseWriter, r *http.Request) {
	userID, isOk := middleware.GetUserID(r.Context())
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	eventUUID, err := utils.ExtractUUIDFromUrl(r.URL.Path, "api/group/event")
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.BadRequestError)
		return
	}
	if !eventModel.IsGroupMemberFromEventUUID(eventUUID, userID) {
		errorControllers.CustomErrorHandler(w, r,
			"only members can view", http.StatusBadRequest)
		return
	}

	events, err := eventModel.SelectEvent(eventUUID, userID)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, "event detail retrieved", events)
}
