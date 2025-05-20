package controller

import (
	"fmt"
	"log"
	"net/http"

	errorControllers "social-network/pkg/errorManagement/controllers"
	followingModel "social-network/pkg/followingManagement/models"
	groupModel "social-network/pkg/groupManagement/models"
	middleware "social-network/pkg/middleware"
	userControllers "social-network/pkg/userManagement/controllers"
	"social-network/pkg/utils"

	_ "github.com/mattn/go-sqlite3"
)

type group = groupModel.Group
type following = followingModel.Following

func isValidGroupInfo(g *group) error {
	isOk := false
	if g.Title, isOk = utils.IsValidContent(g.Title, 3, 100); !isOk {
		return fmt.Errorf("title must be between 3 to 100 characters")
	}
	if g.Description, isOk = utils.IsValidContent(g.Description, 10, 1000); !isOk {
		return fmt.Errorf("description must be between 10 to 1000 characters")
	}
	return nil
}

func CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	g := &group{}
	if err := utils.ReadJSON(r, g); err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	if err := isValidGroupInfo(g); err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), http.StatusBadRequest)
	}

	userId, isOk := middleware.GetUserID(r.Context())
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	g.CreatedBy = userId

	_, groupUUID, err := groupModel.InsertGroup(g)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	err = groupModel.InsertGroupMember(&following{
		LeaderID: userId, FollowerID: userId, GroupUUID: groupUUID,
		Status: "accepted", CreatedBy: userId,
	})
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, "group created",
		struct {
			GroupUUID string `json:"group_uuid"`
		}{groupUUID})
}

func UpdateGroupHandler(w http.ResponseWriter, r *http.Request) {
	g := &group{}
	if err := utils.ReadJSON(r, g); err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	if err := isValidGroupInfo(g); err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), http.StatusBadRequest)
	}

	userID, isOk := middleware.GetUserID(r.Context())
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	if !groupModel.IsGroupMember(g.UUID, userID) {
		errorControllers.ErrorHandler(w, r, errorControllers.ForbiddenError)
		return
	}
	g.UpdatedBy = userID

	if err := groupModel.UpdateGroup(g); err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, "group info updated", nil)
}

func ViewGroupsHandler(w http.ResponseWriter, r *http.Request) {
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

	joinedOnly := false
	if val := r.URL.Query().Get("joined"); val == "true" {
		joinedOnly = true
	}

	groups, err := groupModel.SelectGroups(tgtUUID, joinedOnly)
	if err != nil {
		log.Println(err.Error())
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, "group list retrieved", groups)
}

func ViewGroupHandler(w http.ResponseWriter, r *http.Request) {
	userID, isOk := middleware.GetUserID(r.Context())
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	groupUUID := r.URL.Query().Get("id")
	groups, err := groupModel.SelectGroup(userID, groupUUID)
	if err != nil {
		log.Println(err.Error())
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, "group info retrieved", groups)
}
