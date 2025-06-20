package controller

import (
	"net/http"
	"social-network/pkg/utils"
	"strconv"
)

type ErrorPageData struct {
	Name       string
	Code       string
	CodeNumber int
	CodeSlice  []string
	Info       string
}

var PredefinedErrors = map[string]ErrorPageData{
	"BadRequestError": {
		Name:       "BadRequestError",
		Code:       strconv.Itoa(http.StatusBadRequest),
		CodeNumber: http.StatusBadRequest,
		CodeSlice:  splitString(strconv.Itoa(http.StatusBadRequest)),
		Info:       "Bad request",
	},
	"UnauthorizedError": {
		Name:       "UnauthorizedError",
		Code:       strconv.Itoa(http.StatusUnauthorized),
		CodeNumber: http.StatusUnauthorized,
		CodeSlice:  splitString(strconv.Itoa(http.StatusUnauthorized)),
		Info:       "Unauthorized",
	},
	"ForbiddenError": {
		Name:       "ForbiddenError",
		Code:       strconv.Itoa(http.StatusForbidden),
		CodeNumber: http.StatusUnauthorized,
		CodeSlice:  splitString(strconv.Itoa(http.StatusForbidden)),
		Info:       "Forbidden",
	},
	"NotFoundError": {
		Name:       "NotFoundError",
		Code:       strconv.Itoa(http.StatusNotFound),
		CodeNumber: http.StatusNotFound,
		CodeSlice:  splitString(strconv.Itoa(http.StatusNotFound)),
		Info:       "Page not found",
	},
	"MethodNotAllowedError": {
		Name:       "MethodNotAllowedError",
		Code:       strconv.Itoa(http.StatusMethodNotAllowed),
		CodeNumber: http.StatusMethodNotAllowed,
		CodeSlice:  splitString(strconv.Itoa(http.StatusMethodNotAllowed)),
		Info:       "Method not allowed",
	},
	"InternalServerError": {
		Name:       "InternalServerError",
		Code:       strconv.Itoa(http.StatusInternalServerError),
		CodeNumber: http.StatusInternalServerError,
		CodeSlice:  splitString(strconv.Itoa(http.StatusInternalServerError)),
		Info:       "Internal server error",
	},
}

func splitString(s string) []string {
	result := make([]string, len(s))
	for i, r := range s {
		result[i] = string(r)
	}
	return result
}

var (
	BadRequestError       = PredefinedErrors["BadRequestError"]
	UnauthorizedError     = PredefinedErrors["UnauthorizedError"]
	ForbiddenError        = PredefinedErrors["ForbiddenError"]
	NotFoundError         = PredefinedErrors["NotFoundError"]
	MethodNotAllowedError = PredefinedErrors["MethodNotAllowedError"]
	InternalServerError   = PredefinedErrors["InternalServerError"]
)

func ErrorHandler(w http.ResponseWriter, r *http.Request, errorPageData ErrorPageData) {
	res := utils.Result{
		Success:    false,
		Message:    errorPageData.Info,
		Data:       nil,
		HttpStatus: errorPageData.CodeNumber,
	}
	utils.ReturnJson(w, res)
}

func CustomErrorHandler(w http.ResponseWriter, r *http.Request, message string, errorCode int) {
	res := utils.Result{
		Success:    false,
		Message:    message,
		Data:       nil,
		HttpStatus: errorCode,
	}
	utils.ReturnJson(w, res)
}
