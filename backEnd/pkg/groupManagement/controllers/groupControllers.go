package controller

import (
	"net/http"

	errorControllers "social-network/pkg/errorManagement/controllers"
	followingModel "social-network/pkg/followingManagement/models"
	groupModel "social-network/pkg/groupManagement/models"
	middleware "social-network/pkg/middleware"
	"social-network/pkg/utils"

	_ "github.com/mattn/go-sqlite3"
)

type group = groupModel.Group
type following = followingModel.Following

func CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	g, isOk := &group{}, false
	if err := utils.ReadJSON(w, r, g); err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	if g.Title, isOk = utils.IsValidContent(g.Title, 3, 100); !isOk {
		errorControllers.CustomErrorHandler(w, r, "Title must be between 3 to 100 characters",
			http.StatusBadRequest)
	}
	if g.Description, isOk = utils.IsValidContent(g.Description, 10, 1000); !isOk {
		errorControllers.CustomErrorHandler(w, r, "Description must be between 10 to 1000 characters",
			http.StatusBadRequest)
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

	err = followingModel.InsertFollowing(&following{
		LeaderID: userId, FollowerID: userId, GroupID: groupID,
		Status: "accepted", CreatedBy: userId,
	})
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	utils.ReturnJsonSuccess(w, "Group created successfully",
		struct {
			GroupUUID string `json:"group_uuid"`
		}{groupUUID})
}
