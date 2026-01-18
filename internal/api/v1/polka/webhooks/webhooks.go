package webhooks

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/lucasgjanot/go-http-server/internal/auth"
	"github.com/lucasgjanot/go-http-server/internal/database"
	httperrors "github.com/lucasgjanot/go-http-server/internal/errors"
)

func PostHandler(polkaAPIKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Parameters struct {
			Event string `json:"event"`
			Data struct {
				UserID uuid.UUID `json:"user_id"`
			} `json:"data"`
		}

		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil || apiKey != polkaAPIKey{
			httperrors.Write(w, httperrors.UnauthorizedErr)
			return
		}
		var params Parameters
		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&params); err != nil {
			httperrors.Write(w, httperrors.BadRequestErr)
			return
		}

		if params.Event != "user.upgraded" {
			w.WriteHeader(http.StatusNoContent)
		}

		if _, err := database.Users.GetUserById(
			r.Context(),
			params.Data.UserID,
		); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				httperrors.Write(w, httperrors.NotFoundErr)
				return
			}
			httperrors.Write(w, httperrors.ServiceUnavailableErr)
				return
		}

		_, err = database.Users.UpgradeUser(
			r.Context(),
			params.Data.UserID,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				httperrors.Write(w, httperrors.NotFoundErr)
				return
			}
			httperrors.Write(w, httperrors.ServiceUnavailableErr)
				return
		}

		w.WriteHeader(http.StatusNoContent)

	}
}