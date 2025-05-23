package controller

import (
	// "forum/middlewares"

	"database/sql"
	"errors"
	"net/http"
	errorControllers "social-network/pkg/errorManagement/controllers"
	"social-network/pkg/middleware"
	"social-network/pkg/notificationManagement/models"
	"social-network/pkg/utils"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

// validators for notification id
func isValidNotificationId(notification *models.Notification) error {
	isValid := false

	if notification.ID, isValid = utils.IsValidId(notification.ID); !isValid {
		return errors.New("notification id is required and must be numeric")
	}

	return nil
}

// validators for read status
func isValidNotificationReadStatus(notification *models.Notification) error {
	if notification.IsRead != 0 && notification.IsRead != 1 {
		return errors.New("read status is required and must be numeric either 0 or 1")
	}

	return nil
}

// validators for notification data
func isValidNotificationInfo(notification *models.Notification) error {
	isValid := false

	if notification.FromUserId != 0 {
		if notification.FromUserId, isValid = utils.IsValidId(notification.FromUserId); !isValid {
			return errors.New("notification from user id is required and must be between 3 to 100 alphanumeric characters, '_' or '-'")
		}
	}
	if notification.TargetType == "" || (notification.TargetType != "follow" && notification.TargetType != "group" && notification.TargetType != "group_event") {
		return errors.New("target type is required and must be either 'follow', 'group' or 'group_event'")
	}
	if notification.TargetDetailedType == "" || (notification.TargetDetailedType != "follow_request" && notification.TargetDetailedType != "group_invite" && notification.TargetDetailedType != "group_request" && notification.TargetDetailedType != "group_event") {
		return errors.New("target detailed type is required and must be either 'follow_request', 'group_invite', 'group_request' or 'group_event'")
	}
	if notification.TargetId, isValid = utils.IsValidId(notification.TargetId); !isValid {
		return errors.New("target id is required and must be numeric")
	}
	if notification.Message, isValid = utils.IsValidContent(notification.Message, 3, 1000); !isValid {
		return errors.New("message is required and must be between 3 to 1000 alphanumeric characters, '_' or '-'")
	}
	if notification.TargetUUIDForm != "" {
		notification.TargetUUID = sql.NullString{Valid: true, String: notification.TargetUUIDForm}
	}

	return nil
}

func insertNotification(notification *models.Notification) (int, int, error) {
	if err := isValidNotificationInfo(notification); err != nil {
		return -1, http.StatusBadRequest, err
	}

	createdNotificationId, insertError := models.InsertNotification(notification)
	if insertError != nil {
		return -1, http.StatusInternalServerError, insertError
	}
	return createdNotificationId, http.StatusCreated, nil
}

func ReadAllNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value(middleware.CtxUserID)
	userID, isOk := userIDRaw.(int)
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	posts, err := models.ReadAllNotifications(userID)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	utils.ReturnJsonSuccess(w, "Notifications fetched successfully", posts)
}

func ReadNotificationByIdHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value(middleware.CtxUserID)
	userID, isOk := userIDRaw.(int)
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	notification_id_str, errUrl := utils.ExtractFromUrl(r.URL.Path, "api/notification")
	if errUrl == "not found" {
		errorControllers.ErrorHandler(w, r, errorControllers.NotFoundError)
		return
	}

	notification_id, err := strconv.Atoi(notification_id_str)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.BadRequestError)
		return
	}

	notification, err := models.ReadNotificationById(notification_id, userID)
	if err != nil {
		if err.Error() == "id not found" {
			errorControllers.CustomErrorHandler(w, r, "notification with this id not found for you", errorControllers.NotFoundError.CodeNumber)
			return
		}
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	utils.ReturnJsonSuccess(w, "Notification fetched successfully", notification)
}

func SubmitNotificationHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value(middleware.CtxUserID)
	userID, isOk := userIDRaw.(int)
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	notification := &models.Notification{}
	notification.FromUserId = userID
	if err := utils.ReadJSON(r, notification); err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	instertedNotificationId, insertHttpStatusCode, insertErr := insertNotification(notification)
	if insertErr != nil {
		errorControllers.CustomErrorHandler(w, r, insertErr.Error(), insertHttpStatusCode)
		return
	}

	utils.ReturnJsonSuccess(w, "Notification submitted successfully", instertedNotificationId)
}

func UpdateNotificationReadStatusHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value(middleware.CtxUserID)
	userID, isOk := userIDRaw.(int)
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	notification := &models.Notification{}
	if err := utils.ReadJSON(r, notification); err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	if err := isValidNotificationId(notification); err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	if err := isValidNotificationReadStatus(notification); err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	updateError := models.UpdateNotificationReadStatus(notification.ID, notification.IsRead, userID)
	if updateError != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	utils.ReturnJsonSuccess(w, "Notification updated successfully", nil)
}

func DeleteNotificationHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value(middleware.CtxUserID)
	userID, isOk := userIDRaw.(int)
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	notification := &models.Notification{}
	if err := utils.ReadJSON(r, notification); err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	if err := isValidNotificationId(notification); err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	updateError := models.DeleteNotification(notification.ID, userID)
	if updateError != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	utils.ReturnJsonSuccess(w, "Notification removed successfully", nil)
}
