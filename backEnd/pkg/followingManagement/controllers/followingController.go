package controller

import (
	"fmt"
	"net/http"

	middleware "social-network/pkg/middleware"
	"social-network/pkg/utils"

	followingModel "social-network/pkg/followingManagement/models"
	notificationModel "social-network/pkg/notificationManagement/models"
	userModel "social-network/pkg/userManagement/models"

	errorControllers "social-network/pkg/errorManagement/controllers"
	userControllers "social-network/pkg/userManagement/controllers"
)

func ParseFollowingRequest(r *http.Request) (*followingModel.Following, int, string, int, error) {
	_, userID, userUUID, isOk := middleware.GetSessionCredentials(r.Context())
	if !isOk {
		return nil, -1, "", http.StatusUnauthorized, fmt.Errorf("unauthorized")
	}

	f := &followingModel.Following{}
	if err := utils.ReadJSON(r, f); err != nil {
		return nil, -1, "", http.StatusBadRequest, fmt.Errorf("invalid input")
	}
	if f.FollowerUUID == "" {
		f.FollowerID, f.FollowerUUID = userID, userUUID
	} else {
		f.LeaderID, f.LeaderUUID = userID, userUUID
	}
	f.Type = "user"

	if err := followingModel.SelectIDsFromUUIDs(f); err != nil {
		return nil, -1, "", http.StatusBadRequest, err
	}

	followingStatus, err := followingModel.SelectStatus(f)
	if err != nil {
		return nil, -1, "", http.StatusBadRequest, fmt.Errorf("bad request")
	}
	return f, userID, followingStatus, http.StatusOK, nil
}

func FollowingRequestHandler(w http.ResponseWriter, r *http.Request) {
	f, userID, followingStatus, statusCode, err := ParseFollowingRequest(r)
	if err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), statusCode)
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
	notificationModel.InsertNotification(&notificationModel.Notification{
		ToUserId: f.LeaderID, FromUserId: f.FollowerID,
		TargetType: "following", TargetDetailedType: "follow_request",
		TargetId: f.FollowerID, TargetUUID: f.FollowerUUID,
		Message: f.Status,
	})
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, message, nil)
}

func FollowingResponseHandler(w http.ResponseWriter, r *http.Request) {
	f, userID, followingStatus, statusCode, err := ParseFollowingRequest(r)
	if err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), statusCode)
		return
	}
	f.UpdatedBy = userID
	message := "request " + f.Status

	if followingStatus != "requested" || (f.Status != "accepted" && f.Status != "declined") {
		errorControllers.ErrorHandler(w, r, errorControllers.BadRequestError)
		return
	}

	if err := followingModel.UpdateFollowing(f); err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	notificationModel.InsertNotification(&notificationModel.Notification{
		ToUserId: f.FollowerID, FromUserId: f.LeaderID,
		TargetType: "following", TargetDetailedType: "follow_request",
		TargetId: f.LeaderID, TargetUUID: f.LeaderUUID,
		Message: f.Status,
	})
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, message, nil)
}

func UnfollowHandler(w http.ResponseWriter, r *http.Request) {
	f, userID, followingStatus, statusCode, err := ParseFollowingRequest(r)
	if err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), statusCode)
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
	if followingStatus == "requested" {
		notificationModel.UpdateNotificationOnCancel(
			&notificationModel.Notification{
				ToUserId: f.LeaderID, FromUserId: f.FollowerID,
				Message: "inactive", UpdatedBy: &userID,
			}, "following", followingStatus)
	}
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, message, nil)
}

func FollowingRemoveHandler(w http.ResponseWriter, r *http.Request) {
	f, userID, followingStatus, statusCode, err := ParseFollowingRequest(r)
	if err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), statusCode)
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

func ViewFollowingsHandler(w http.ResponseWriter, r *http.Request) {
	tgtUUID, statusCode, err := userControllers.GetTgtUUID(r, "api/followers")
	if err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), statusCode)
		return
	}

	followings, err := followingModel.SelectFollowings(tgtUUID, "accepted")
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, "follow data retrieved", followings)
}

func ViewFollowingRequestsHandler(w http.ResponseWriter, r *http.Request) {
	userUUID, isOk := middleware.GetUserUUID(r.Context())
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.UnauthorizedError)
		return
	}

	pendingFollowing, err := followingModel.SelectFollowings(userUUID, "requested")
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, "follow requests retrieved", pendingFollowing)
}
