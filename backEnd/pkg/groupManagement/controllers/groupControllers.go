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

func GroupCreateHandler(w http.ResponseWriter, r *http.Request) {
	g := &group{}
	if err := utils.ReadJSON(r, g); err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	userId, isOk := middleware.GetUserID(r.Context())
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	g.CreatedBy = userId

	if err := isValidGroupInfo(g); err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), http.StatusBadRequest)
	}

	groupId, groupUUID, err := groupModel.InsertGroup(g)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	err = groupModel.InsertGroupMember(&following{
		LeaderID: userId, FollowerID: userId, GroupID: groupId, GroupUUID: groupUUID,
		Type: "group", Status: "accepted", CreatedBy: userId,
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

func GroupUpdateHandler(w http.ResponseWriter, r *http.Request) {
	g := &group{}
	if err := utils.ReadJSON(r, g); err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	userID, isOk := middleware.GetUserID(r.Context())
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	g.UpdatedBy = userID

	if err := isValidGroupInfo(g); err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), http.StatusBadRequest)
	}
	if !groupModel.IsGroupMember(g.UUID, userID) {
		errorControllers.ErrorHandler(w, r, errorControllers.ForbiddenError)
		return
	}

	if err := groupModel.UpdateGroup(g); err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, "group info updated", nil)
}

func ViewGroupsHandler(w http.ResponseWriter, r *http.Request) {
	tgtUUID, statusCode := userControllers.GetTgtUUID(r, "api/groups")
	if statusCode != http.StatusOK {
		errorControllers.CustomErrorHandler(w, r, tgtUUID, statusCode)
		return
	}

	joined := r.URL.Query().Get("joined")

	groups, err := groupModel.SelectGroups(tgtUUID, joined == "true")
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

	groupUUID, err := utils.ExtractUUIDFromUrl(r.URL.Path, "api/group")
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.NotFoundError)
		return
	}

	groups, err := groupModel.SelectGroup(userID, groupUUID)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.NotFoundError)
		return
	}
	userControllers.ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, "group info retrieved", groups)
}
