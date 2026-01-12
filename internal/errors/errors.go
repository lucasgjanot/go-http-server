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

var InternalServerErr ErrorResponse = ErrorResponse{
	Name: "InternalServerErr",
	StatusCode: 500,
	Message: "Something went wrong",
	Action: "Try again later",
}
var BadRequestErr ErrorResponse = ErrorResponse{
	Name: "BadRequestErr",
	StatusCode: 400,
	Message: "Bad Request Error",
	Action: "Verify sent data",
}



func Write(w http.ResponseWriter, err ErrorResponse) {
	w.Header().Set("Content-Type", "applicaton/json")
	w.WriteHeader(err.StatusCode)
	_ = json.NewEncoder(w).Encode(err)
}