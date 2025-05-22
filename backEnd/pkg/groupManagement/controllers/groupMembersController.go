package controller

import (
	"fmt"
	"log"
	"net/http"
	errorControllers "social-network/pkg/errorManagement/controllers"
	followingControllers "social-network/pkg/followingManagement/controllers"
	followingModel "social-network/pkg/followingManagement/models"
	groupModel "social-network/pkg/groupManagement/models"
	middleware "social-network/pkg/middleware"
	userControllers "social-network/pkg/userManagement/controllers"
	"social-network/pkg/utils"
)

func processMembersIDs(r *http.Request) (*followingModel.Following, int, string, error) {
	m, userID, _, err := followingControllers.ProcessFollowingIDs(r)
	if err != nil {
		return nil, -1, "", err
	}
	groupID, createdBy, err := groupModel.SelectGroupIDcreatedByfromUUID(m.GroupUUID)
	if err != nil {
		return nil, -1, "", err

	} else if groupID == 0 {
		return nil, -1, "", fmt.Errorf("public forum chosen as group")
	}
	m.GroupID, m.LeaderID = groupID, createdBy

	memberStatus, err := followingModel.SelectStatus(m)
	if err != nil {
		return nil, -1, "", err
	}
	return m, userID, memberStatus, nil
}

func GroupInviteRequestHandler(w http.ResponseWriter, r *http.Request) {
	m, userID, memberStatus, err := processMembersIDs(r)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	m.Status = "invited"
	operation, message := groupModel.InsertGroupMember, "invitation send"

	if !groupModel.IsGroupMember(m.GroupUUID, userID) {
		errorControllers.CustomErrorHandler(w, r,
			"only members can invite", http.StatusBadRequest)
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
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, message, nil)
}

func GroupJoinRequestHandler(w http.ResponseWriter, r *http.Request) {
	m, userID, memberStatus, err := processMembersIDs(r)
	if err != nil {
		log.Println(err.Error())
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
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
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, message, nil)
}

func GroupQuitHandler(w http.ResponseWriter, r *http.Request) {
	m, userID, memberStatus, err := processMembersIDs(r)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	m.UpdatedBy, m.Status = userID, "inactive"
	message := "not a member/ no request or invitation"

	switch memberStatus {
	case "accepted":
		message = "you have left group"
	case "requested", "invited":
		message = "cancel/decline request"
	default:
		errorControllers.CustomErrorHandler(w, r, message, http.StatusBadRequest)
		return
	}

	if err := followingModel.UpdateFollowing(m); err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, message, nil)
}

func GroupMemberRemoveHandler(w http.ResponseWriter, r *http.Request) {
	m, userID, memberStatus, err := processMembersIDs(r)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	m.UpdatedBy, m.Status = userID, "declined"

	if userID != m.LeaderID {
		errorControllers.CustomErrorHandler(w, r,
			"only leader can remove member", http.StatusBadRequest)
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
	m, userID, memberStatus, err := processMembersIDs(r)
	if err != nil {
		log.Println(err.Error())
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	message := "requested "
	if memberStatus == "invited" {
		message = "invitation "
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
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, message, nil)
}

func ViewGroupMembersHandle(w http.ResponseWriter, r *http.Request) {
	tgtUUID, err := utils.ExtractUUIDFromUrl(r.URL.Path, "api/group/members")
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.BadRequestError)
		return
	}

	groupMembers, err := groupModel.SelectGroupMembers(tgtUUID, "accepted")
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, "group members retrevied", groupMembers)
}

func ViewGroupMemberRequestsHandle(w http.ResponseWriter, r *http.Request) {
	userID, isOk := middleware.GetUserID(r.Context())
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	tgtUUID, err := utils.ExtractUUIDFromUrl(r.URL.Path, "api/group/member/requests")
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.BadRequestError)
		return
	}

	if !groupModel.IsGroupMember(tgtUUID, userID) {
		errorControllers.CustomErrorHandler(w, r,
			"only members can view", http.StatusBadRequest)
		return
	}

	groupMembers, err := groupModel.SelectGroupMembers(tgtUUID, "requested")
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, "group member requests retrevied", groupMembers)
}
