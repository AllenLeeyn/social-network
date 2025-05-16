package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

type Result struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	HttpStatus int         `json:"httpStatus,omitempty"`
}

func ReturnJson(w http.ResponseWriter, outputData Result) {
	w.Header().Set("Content-Type", "application/json")
	if outputData.Success {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(outputData.HttpStatus)
	}

	if err := json.NewEncoder(w).Encode(outputData); err != nil {
		errorResponse := Result{
			Success: false,
			Message: "Failed to encode JSON",
			Data:    nil,
		}
		json.NewEncoder(w).Encode(errorResponse)
	}
}

func ReturnJsonSuccess(w http.ResponseWriter, message string, data interface{}) {
	ReturnJson(w, Result{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(body, data); err != nil {
		return err
	}
	return nil
}
