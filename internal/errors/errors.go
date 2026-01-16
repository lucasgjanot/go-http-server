package httperrors

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Name string `json:"name"`
	StatusCode int `json:"status_code"`
	Message string `json:"error"`
	Action string `json:"action"`
}


var (
	InternalServerErr = ErrorResponse{
		Name: "InternalServerErr",
		StatusCode: http.StatusInternalServerError,
		Message: "Something went wrong",
		Action: "Try again later",
	}

	BadRequestErr = ErrorResponse{
		Name: "BadRequestErr",
		StatusCode: http.StatusBadRequest,
		Message: "Bad Request Error",
		Action: "Verify sent data",
	}

	ServiceUnavailableErr = ErrorResponse{
		Name: "ServiceUnavailableErr",
		StatusCode: http.StatusServiceUnavailable,
		Message: "Bad Request Error",
		Action: "Verify sent data",
	}

	NotFoundErr = ErrorResponse{
		Name: "NotFoundErr",
		StatusCode: http.StatusNotFound,
		Message: "Not Found",
		Action: "Check search parameter",
	}
	UnauthorizedErr = ErrorResponse{
		Name: "UnauthorizedErr",
		StatusCode: http.StatusUnauthorized,
		Message: "Incorrect email or password",
		Action: "Check credentials",
	}
)
	
func Write(w http.ResponseWriter, err ErrorResponse) {
	w.Header().Set("Content-Type", "applicaton/json")
	w.WriteHeader(err.StatusCode)
	_ = json.NewEncoder(w).Encode(err)
}