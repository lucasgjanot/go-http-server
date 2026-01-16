package chirpsid

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/lucasgjanot/go-http-server/internal/api/v1/chirps"
	"github.com/lucasgjanot/go-http-server/internal/database"
	httperrors "github.com/lucasgjanot/go-http-server/internal/errors"
)

func GetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chirpId := r.PathValue("chirpID")
		chirpIdParsed, err := uuid.Parse(chirpId)
		if err != nil {
			httperrors.Write(w,httperrors.BadRequestErr)
			return
		}

		chirp, err := database.Chirps.GetChirp(r.Context(), chirpIdParsed)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				httperrors.Write(w, httperrors.NotFoundErr)
				return
			}
			httperrors.Write(w, httperrors.ServiceUnavailableErr)
			return 
		}

		response := chirps.ChirpsResponse{
			Id: chirp.ID,
			Body: chirp.Body,
			UserId: chirp.UserID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	}
}