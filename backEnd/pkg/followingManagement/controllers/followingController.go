package controller

import (
	"fmt"
	"log"
	"net/http"

	errorControllers "social-network/pkg/errorManagement/controllers"
	followingModel "social-network/pkg/followingManagement/models"
	middleware "social-network/pkg/middleware"
	userControllers "social-network/pkg/userManagement/controllers"
	userModel "social-network/pkg/userManagement/models"
	"social-network/pkg/utils"
)

func processFollowingRequest(r *http.Request, userType string) (*followingModel.Following, error) {
	f := &followingModel.Following{}
	if err := utils.ReadJSON(r, f); err != nil {
		return nil, err
	}

	userID, isOk := middleware.GetUserID(r.Context())
	if !isOk {
		return nil, fmt.Errorf("error getting userID")
	}
	switch userType {
	case "follower":
		f.FollowerID = userID
	case "leader":
		f.LeaderID = userID
	}

	if err := followingModel.SelectIDsFromUUIDs(f); err != nil {
		return nil, err
	}
	return f, nil
}

func FollowingRequestHandler(w http.ResponseWriter, r *http.Request) {
	f, err := processFollowingRequest(r, "follower")
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	f.Status = "requested"

	operation := followingModel.InsertFollowing
	followingStatus, err := followingModel.SelectStatus(f)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	switch followingStatus {
	case "accepted":
		errorControllers.CustomErrorHandler(w, r, "You are a follower", http.StatusBadRequest)
		return
	case "requested", "invited":
		errorControllers.CustomErrorHandler(w, r, "Response pending. Please wait", http.StatusBadRequest)
		return
	case "declined", "inactive":
		operation = followingModel.UpdateFollowing
		f.UpdatedBy = f.FollowerID
	default:
		f.CreatedBy = f.FollowerID
	}

	// check if leader profile is public, and accept request. Need to test
	if userModel.IsPublic(f.LeaderUUID) {
		f.Status = "accepted"
	}
	if err := operation(f); err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, "request send", nil)
}

func UnfollowHandler(w http.ResponseWriter, r *http.Request) {
	f, err := processFollowingRequest(r, "follower")
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	f.UpdatedBy, f.Status = f.FollowerID, "inactive"

	message := "not following user"
	followingStatus, err := followingModel.SelectStatus(f)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	switch followingStatus {
	case "accepted":
		message = "unfollow user"
	case "requested":
		message = "cancel follow request"
	default:
		errorControllers.CustomErrorHandler(w, r, message, http.StatusBadRequest)
		return
	}

	if err := followingModel.UpdateFollowing(f); err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, message, nil)
}

func RemoveFollowerHandler(w http.ResponseWriter, r *http.Request) {
	f, err := processFollowingRequest(r, "leader")
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	f.UpdatedBy, f.Status = f.LeaderID, "declined"

	followingStatus, err := followingModel.SelectStatus(f)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	if followingStatus != "accepted" {
		errorControllers.CustomErrorHandler(w, r, "User is not a current follower", http.StatusBadRequest)
		return
	}

	if err := followingModel.UpdateFollowing(f); err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, "follower removed", nil)
}

func FollowingResponseHandler(w http.ResponseWriter, r *http.Request) {
	f, err := processFollowingRequest(r, "leader")
	if err != nil {
		log.Println(err.Error())
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	f.UpdatedBy = f.LeaderID

	message := "request accepted"
	if f.Status == "declined" {
		message = "request declined"
	}

	followingStatus, _ := followingModel.SelectStatus(f)
	if followingStatus != "requested" || (f.Status != "accepted" && f.Status != "declined") {
		errorControllers.ErrorHandler(w, r, errorControllers.BadRequestError)
		return
	}

	if err := followingModel.UpdateFollowing(f); err != nil {
		log.Println(err.Error())
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, message, nil)
}

func ViewFollowingRequestsHandler(w http.ResponseWriter, r *http.Request) {
	userID, isOk := middleware.GetUserID(r.Context())
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	pendingFollowing, err := followingModel.SelectFollowings(userID, "requested")
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	utils.ReturnJsonSuccess(w, "Data OK", pendingFollowing)
}

// show data if is follower or public
func ViewFollowingsHandler(w http.ResponseWriter, r *http.Request) {
	userID, isOk := middleware.GetUserID(r.Context())
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	pendingFollowing, err := followingModel.SelectFollowings(userID, "accepted")
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	utils.ReturnJsonSuccess(w, "Data OK", pendingFollowing)
}
