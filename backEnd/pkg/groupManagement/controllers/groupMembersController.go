package controller

import (
	"fmt"
	"net/http"

	middleware "social-network/pkg/middleware"
	"social-network/pkg/utils"

	followingModel "social-network/pkg/followingManagement/models"
	groupModel "social-network/pkg/groupManagement/models"
	notificationModel "social-network/pkg/notificationManagement/models"

	errorControllers "social-network/pkg/errorManagement/controllers"
	followingControllers "social-network/pkg/followingManagement/controllers"
	userControllers "social-network/pkg/userManagement/controllers"
)

func parseMemberRequest(r *http.Request) (*followingModel.Following, int, string, int, error) {
	m, userID, _, statusCode, err := followingControllers.ParseFollowingRequest(r)
	if err != nil {
		return nil, -1, "", statusCode, err
	}
	groupID, createdBy, err := groupModel.SelectGroupIDcreatedByfromUUID(m.GroupUUID)
	if err != nil {
		return nil, -1, "", http.StatusBadRequest, fmt.Errorf("bad request")

	} else if groupID == 0 {
		return nil, -1, "", http.StatusBadRequest, fmt.Errorf("public forum chosen as group")
	}
	m.GroupID, m.LeaderID, m.Type = groupID, createdBy, "group"

	memberStatus, err := followingModel.SelectStatus(m)
	if err != nil {
		return nil, -1, "", http.StatusBadRequest, fmt.Errorf("bad request")
	}
	return m, userID, memberStatus, http.StatusOK, nil
}

func GroupInviteRequestHandler(w http.ResponseWriter, r *http.Request) {
	m, userID, memberStatus, statusCode, err := parseMemberRequest(r)
	if err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), statusCode)
		return
	}
	m.Status = "invited"
	operation, message := groupModel.InsertGroupMember, "invitation send"

	if !groupModel.IsGroupMember(m.GroupUUID, userID) {
		errorControllers.CustomErrorHandler(w, r,
			"only members can invite", http.StatusForbidden)
		return
	}

	switch memberStatus {
	case "accepted":
		errorControllers.CustomErrorHandler(w, r,
			"user is a member", http.StatusBadRequest)
		return
	case "requested", "invited":
		errorControllers.CustomErrorHandler(w, r,
			"Response pending. Please wait", http.StatusBadRequest)
		return
	case "declined", "inactive":
		operation = followingModel.UpdateFollowing
		m.UpdatedBy = userID
	default:
		m.CreatedBy = userID
	}

	if err := operation(m); err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	notificationModel.InsertNotification(&notificationModel.Notification{
		ToUserId: m.FollowerID, FromUserId: m.LeaderID,
		TargetType: "group", TargetDetailedType: "group_invite",
		TargetId: m.GroupID, TargetUUID: m.GroupUUID,
		Message: m.Status,
	})
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, message, nil)
}

func GroupJoinRequestHandler(w http.ResponseWriter, r *http.Request) {
	m, userID, memberStatus, statusCode, err := parseMemberRequest(r)
	if err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), statusCode)
		return
	}
	m.Status = "requested"
	operation, message := groupModel.InsertGroupMember, "request send"

	switch memberStatus {
	case "accepted":
		errorControllers.CustomErrorHandler(w, r,
			"you are a member", http.StatusBadRequest)
		return
	case "requested", "invited":
		errorControllers.CustomErrorHandler(w, r,
			"Response pending. Please wait", http.StatusBadRequest)
		return
	case "declined", "inactive":
		operation = followingModel.UpdateFollowing
		m.UpdatedBy = userID
	default:
		m.CreatedBy = userID
	}

	if err := operation(m); err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	notificationModel.InsertNotification(&notificationModel.Notification{
		ToUserId: m.LeaderID, FromUserId: userID,
		TargetType: "group", TargetDetailedType: "group_request",
		TargetId: m.FollowerID, TargetUUID: m.FollowerUUID,
		Message: m.Status,
	})
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, message, nil)
}

func GroupQuitHandler(w http.ResponseWriter, r *http.Request) {
	m, userID, memberStatus, statusCode, err := parseMemberRequest(r)
	if err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), statusCode)
		return
	}
	m.UpdatedBy, m.Status = userID, "inactive"
	message := "not a member/ no request or invitation"

	toUserID, fromUserID := m.LeaderID, userID
	switch memberStatus {
	case "accepted":
		message = "you have left group"
	case "requested":
		message = "cancel request"
	case "invited":
		message, m.Status = "decline invite", "declined"
		toUserID, fromUserID = userID, m.LeaderID

	default:
		errorControllers.CustomErrorHandler(w, r, message, http.StatusBadRequest)
		return
	}

	if err := followingModel.UpdateFollowing(m); err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	if memberStatus != "accepted" {
		notificationModel.UpdateNotificationOnCancel(
			&notificationModel.Notification{
				ToUserId: toUserID, FromUserId: fromUserID,
				Message: m.Status, UpdatedBy: &userID,
			}, "group", memberStatus)
	}
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, message, nil)
}

func GroupMemberRemoveHandler(w http.ResponseWriter, r *http.Request) {
	m, userID, memberStatus, statusCode, err := parseMemberRequest(r)
	if err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), statusCode)
		return
	}
	m.UpdatedBy, m.Status = userID, "declined"

	if userID != m.LeaderID {
		errorControllers.CustomErrorHandler(w, r,
			"only leader can remove member", http.StatusForbidden)
		return
	}
	if memberStatus != "accepted" {
		errorControllers.CustomErrorHandler(w, r,
			"user is not a member", http.StatusBadRequest)
		return
	}

	if err := followingModel.UpdateFollowing(m); err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, "member removed", nil)
}

func GroupMemberResponseHandler(w http.ResponseWriter, r *http.Request) {
	m, userID, memberStatus, statusCode, err := parseMemberRequest(r)
	if err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), statusCode)
		return
	}
	message := "request "
	toUserID, fromUserID, detailedType := m.FollowerID, userID, "group_request"
	if memberStatus == "invited" {
		message = "invitation "
		toUserID, fromUserID, detailedType = m.LeaderID, userID, "group_invite"
	}
	message += m.Status

	if (memberStatus != "requested" && memberStatus != "invited") ||
		(m.Status != "accepted" && m.Status != "declined") ||
		(memberStatus == "invited" && userID != m.FollowerID) ||
		(memberStatus == "requested" && userID != m.LeaderID) {
		errorControllers.ErrorHandler(w, r, errorControllers.BadRequestError)
		return
	}
	if err := followingModel.UpdateFollowing(m); err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	notificationModel.InsertNotification(
		&notificationModel.Notification{
			ToUserId: toUserID, FromUserId: fromUserID,
			TargetType: "group", TargetDetailedType: detailedType,
			TargetId: m.GroupID, TargetUUID: m.GroupUUID,
			Message: m.Status, UpdatedBy: &userID,
		})
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, message, nil)
}

func ViewGroupMembersHandle(w http.ResponseWriter, r *http.Request) {
	tgtUUID, err := utils.ExtractUUIDFromUrl(r.URL.Path, "api/group/members")
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.NotFoundError)
		return
	}

	groupMembers, err := groupModel.SelectGroupMembers(tgtUUID, "accepted")
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.NotFoundError)
		return
	}
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, "group members retrevied", groupMembers)
}

func ViewGroupMemberRequestsHandle(w http.ResponseWriter, r *http.Request) {
	userID, isOk := middleware.GetUserID(r.Context())
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.UnauthorizedError)
		return
	}
	tgtUUID, err := utils.ExtractUUIDFromUrl(r.URL.Path, "api/group/member/requests")
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.NotFoundError)
		return
	}

	if !groupModel.IsGroupMember(tgtUUID, userID) {
		errorControllers.CustomErrorHandler(w, r,
			"only members can view", http.StatusNotFound)
		return
	}

	groupMembers, err := groupModel.SelectGroupMembers(tgtUUID, "requested")
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.NotFoundError)
		return
	}
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, "group member requests retrevied", groupMembers)
}
