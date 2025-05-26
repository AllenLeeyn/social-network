package controller

import (
	"fmt"
	"net/http"

	"social-network/pkg/middleware"
	"social-network/pkg/utils"

	errorControllers "social-network/pkg/errorManagement/controllers"
	eventModel "social-network/pkg/eventManagement/models"
	userControllers "social-network/pkg/userManagement/controllers"
)

func parseEventResponse(r *http.Request) (*eventResponse, string, int, error) {
	userID, isOk := middleware.GetUserID(r.Context())
	if !isOk {
		return nil, "", http.StatusUnauthorized, fmt.Errorf("unauthorized")
	}
	eResponse := &eventResponse{}
	if err := utils.ReadJSON(r, eResponse); err != nil {
		return nil, "", http.StatusBadRequest, fmt.Errorf("invalid input")
	}
	eResponse.CreatedBy = userID

	responseStatus, err := eventModel.SelectStatus(eResponse)
	if err != nil {
		return nil, "", http.StatusNotFound, fmt.Errorf("not found")
	}
	return eResponse, responseStatus, http.StatusOK, nil
}

func EventResponseHandler(w http.ResponseWriter, r *http.Request) {
	eResponse, responseStatus, statusCode, err := parseEventResponse(r)
	if err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), statusCode)
		return
	}

	if !eventModel.IsGroupMemberFromEventUUID(eResponse.EventUUID, eResponse.CreatedBy) {
		errorControllers.CustomErrorHandler(w, r,
			"only members can reponse", http.StatusForbidden)
		return
	}

	operation := eventModel.InsertEventResponse
	if responseStatus != "" && responseStatus != "accepted" && responseStatus != "declined" {
		errorControllers.ErrorHandler(w, r, errorControllers.BadRequestError)
		return
	}
	if responseStatus != "" {
		operation = eventModel.UpdateEventResponse
		eResponse.UpdatedBy = eResponse.CreatedBy
	}

	if err := operation(eResponse); err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, "RSVP submitted", nil)
}

func ViewEventResponsesHandler(w http.ResponseWriter, r *http.Request) {
	eventUUID, err := utils.ExtractUUIDFromUrl(r.URL.Path, "api/group/event/responses")
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.NotFoundError)
		return
	}
	userID, isOk := middleware.GetUserID(r.Context())
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.UnauthorizedError)
		return
	}

	if !eventModel.IsGroupMemberFromEventUUID(eventUUID, userID) {
		errorControllers.CustomErrorHandler(w, r,
			"only members can reponse", http.StatusForbidden)
		return
	}
	eResponses, err := eventModel.SelectEventResponses(eventUUID)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.NotFoundError)
		return
	}
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, "event response retrieved", eResponses)
}
