package controller

import (
	"fmt"
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
	if err := utils.ReadJSON(w, r, g); err != nil {
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

	groupID, groupUUID, err := groupModel.InsertGroup(g)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	err = groupModel.InsertGroupMember(&following{
		LeaderID: userId, FollowerID: userId, GroupID: groupID,
		Status: "accepted", CreatedBy: userId,
	})
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	sessionId, _ := middleware.GetSessionID(r.Context())
	userControllers.ExtendSession(w, sessionId)
	utils.ReturnJsonSuccess(w, "Group created successfully",
		struct {
			GroupUUID string `json:"group_uuid"`
		}{groupUUID})
}

func UpdateGroupHandler(w http.ResponseWriter, r *http.Request) {
	g := &group{}
	if err := utils.ReadJSON(w, r, g); err != nil {
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

	sessionId, _ := middleware.GetSessionID(r.Context())
	userControllers.ExtendSession(w, sessionId)
	utils.ReturnJsonSuccess(w, "Updated successfully", nil)
}

// return list of all groups and an array of id for groups joined
func ViewGroupsHandler(w http.ResponseWriter, r *http.Request) {}

func ViewUserGroupHandler(w http.ResponseWriter, r *http.Request) {}

// return detail of groups with
// list of members uuid and username
// list of events
func ViewGroupHandler(w http.ResponseWriter, r *http.Request) {}
