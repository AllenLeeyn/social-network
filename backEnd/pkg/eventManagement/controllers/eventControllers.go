package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"social-network/pkg/middleware"
	"social-network/pkg/utils"

	eventModel "social-network/pkg/eventManagement/models"
	groupModel "social-network/pkg/groupManagement/models"
	notificationModel "social-network/pkg/notificationManagement/models"

	errorControllers "social-network/pkg/errorManagement/controllers"
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

func parseEventRequest(r *http.Request) (*event, int, int, error) {
	e := &event{}
	if err := utils.ReadJSON(r, e); err != nil {
		return nil, -1, http.StatusBadRequest, fmt.Errorf("invalid input")
	}
	userID, isOk := middleware.GetUserID(r.Context())
	if !isOk {
		return nil, -1, http.StatusUnauthorized, fmt.Errorf("unauthorized")
	}

	if err := isValidEventInfo(e); err != nil {
		return nil, -1, http.StatusBadRequest, err
	}

	groupID, _, err := groupModel.SelectGroupIDcreatedByfromUUID(e.GroupUUID)
	if err != nil {
		return nil, -1, http.StatusNotFound, fmt.Errorf("group not found")
	}
	if !groupModel.IsGroupMember(e.GroupUUID, userID) {
		return nil, -1, http.StatusForbidden, fmt.Errorf("only group members allowed")
	}
	e.GroupID = groupID

	return e, userID, http.StatusOK, nil
}

func EventCreateHandler(w http.ResponseWriter, r *http.Request) {
	e, userID, statusCode, err := parseEventRequest(r)
	if err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), statusCode)
		return
	}
	e.CreatedBy = userID

	eventID, eventUUID, err := eventModel.InsertEvent(e)
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
	notificationModel.InsertNotificationForEvent(&notificationModel.Notification{
		FromUserId: userID, TargetId: e.GroupID, TargetUUID: e.GroupUUID,
		Message: strconv.Itoa(eventID),
	}, e.GroupID, userID)
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, "event created",
		struct {
			EventUUID string `json:"event_uuid"`
		}{eventUUID})
}

func EventUpdateHandler(w http.ResponseWriter, r *http.Request) {
	e, userID, statusCode, err := parseEventRequest(r)
	if err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), statusCode)
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
		errorControllers.ErrorHandler(w, r, errorControllers.UnauthorizedError)
		return
	}
	groupUUID, err := utils.ExtractUUIDFromUrl(r.URL.Path, "api/group/events")
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.NotFoundError)
		return
	}
	if !groupModel.IsGroupMember(groupUUID, userID) {
		errorControllers.ErrorHandler(w, r, errorControllers.NotFoundError)
		return
	}

	events, err := eventModel.SelectEvents(groupUUID, userID)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.NotFoundError)
		return
	}
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, "event list retrieved", events)
}

func ViewEventHandler(w http.ResponseWriter, r *http.Request) {
	userID, isOk := middleware.GetUserID(r.Context())
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.UnauthorizedError)
		return
	}
	eventUUID, err := utils.ExtractUUIDFromUrl(r.URL.Path, "api/group/event")
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.NotFoundError)
		return
	}
	if !eventModel.IsGroupMemberFromEventUUID(eventUUID, userID) {
		errorControllers.ErrorHandler(w, r, errorControllers.NotFoundError)
		return
	}

	events, err := eventModel.SelectEvent(eventUUID, userID)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.NotFoundError)
		return
	}
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, "event detail retrieved", events)
}
