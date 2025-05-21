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

func ProcessFollowingIDs(r *http.Request) (*followingModel.Following, int, string, error) {
	userID, isOk := middleware.GetUserID(r.Context())
	if !isOk {
		return nil, -1, "", fmt.Errorf("error getting userID")
	}

	f := &followingModel.Following{}
	if err := utils.ReadJSON(r, f); err != nil {
		return nil, -1, "", err
	}
	if f.FollowerUUID == "" {
		f.FollowerID = userID
	} else {
		f.LeaderID = userID
	}

	if err := followingModel.SelectIDsFromUUIDs(f); err != nil {
		return nil, -1, "", err
	}

	followingStatus, err := followingModel.SelectStatus(f)
	if err != nil {
		return nil, -1, "", err
	}
	return f, userID, followingStatus, nil
}

func FollowingRequestHandler(w http.ResponseWriter, r *http.Request) {
	f, userID, followingStatus, err := ProcessFollowingIDs(r)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	f.Status = "requested"

	operation, message := followingModel.InsertFollowing, "request send"

	switch followingStatus {
	case "accepted":
		errorControllers.CustomErrorHandler(w, r,
			"you are a follower", http.StatusBadRequest)
		return
	case "requested", "invited":
		errorControllers.CustomErrorHandler(w, r,
			"Response pending. Please wait", http.StatusBadRequest)
		return
	case "declined", "inactive":
		operation = followingModel.UpdateFollowing
		f.UpdatedBy = userID
	default:
		f.CreatedBy = userID
	}

	if userModel.IsPublic(f.LeaderUUID) {
		f.Status, message = "accepted", "request accepted"
	}
	if err := operation(f); err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, message, nil)
}

func UnfollowHandler(w http.ResponseWriter, r *http.Request) {
	f, userID, followingStatus, err := ProcessFollowingIDs(r)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	f.UpdatedBy, f.Status = userID, "inactive"

	message := "not following user"

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
	f, userID, followingStatus, err := ProcessFollowingIDs(r)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	f.UpdatedBy, f.Status = userID, "declined"

	if followingStatus != "accepted" {
		errorControllers.CustomErrorHandler(w, r,
			"User is not a current follower", http.StatusBadRequest)
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
	f, userID, followingStatus, err := ProcessFollowingIDs(r)
	if err != nil {
		log.Println(err.Error())
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	f.UpdatedBy = userID

	message := "request " + f.Status

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
	userUUID, isOk := middleware.GetUserUUID(r.Context())
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	pendingFollowing, err := followingModel.SelectFollowings(userUUID, "requested")
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, "Data OK", pendingFollowing)
}

func ViewFollowingsHandler(w http.ResponseWriter, r *http.Request) {
	tgtUUID, statusCode := userControllers.GetTgtUUID(r)
	if statusCode == http.StatusInternalServerError {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	} else if statusCode == http.StatusForbidden {
		errorControllers.CustomErrorHandler(w, r,
			"access denied: private profile and user is not follower",
			http.StatusForbidden)
		return
	}

	followings, err := followingModel.SelectFollowings(tgtUUID, "accepted")
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, "Data OK", followings)
}
