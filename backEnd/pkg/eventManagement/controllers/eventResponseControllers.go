package controller

import (
	"fmt"
	"log"
	"net/http"

	"social-network/pkg/middleware"
	"social-network/pkg/utils"

	errorControllers "social-network/pkg/errorManagement/controllers"
	eventModel "social-network/pkg/eventManagement/models"
	userControllers "social-network/pkg/userManagement/controllers"
)

func parseEventResponse(r *http.Request) (*eventResponse, string, error) {
	userID, isOk := middleware.GetUserID(r.Context())
	if !isOk {
		return nil, "", fmt.Errorf("InternalServerError")
	}
	eResponse := &eventResponse{}
	if err := utils.ReadJSON(r, eResponse); err != nil {
		return nil, "", fmt.Errorf("invalid input")
	}
	eResponse.CreatedBy = userID

	responseStatus, err := eventModel.SelectStatus(eResponse)
	if err != nil {
		return nil, "", fmt.Errorf("InternalServerError")
	}
	return eResponse, responseStatus, nil
}

func EventResponseHandler(w http.ResponseWriter, r *http.Request) {
	eResponse, responseStatus, err := parseEventResponse(r)
	if err != nil {
		if err.Error() == "InternalServerError" {
			errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		} else {
			errorControllers.CustomErrorHandler(w, r, err.Error(), http.StatusBadRequest)
		}
		return
	}

	if !eventModel.IsGroupMemberFromEventUUID(eResponse.EventUUID, eResponse.CreatedBy) {
		errorControllers.CustomErrorHandler(w, r,
			"only members can reponse", http.StatusBadRequest)
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
		errorControllers.ErrorHandler(w, r, errorControllers.BadRequestError)
		return
	}
	userID, isOk := middleware.GetUserID(r.Context())
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	log.Println(eventUUID)
	if !eventModel.IsGroupMemberFromEventUUID(eventUUID, userID) {
		errorControllers.CustomErrorHandler(w, r,
			"only members can reponse", http.StatusBadRequest)
		return
	}
	eResponses, err := eventModel.SelectEventResponses(eventUUID)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, "event response retrieved", eResponses)
}
