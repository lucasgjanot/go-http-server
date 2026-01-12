package healthz

import (
	"encoding/json"
	"net/http"
)

func GetHandler() http.HandlerFunc {
	type response struct {
		Status string `json:"status"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(response{
			Status: "ok",
		})
	}
}
