package chirpsid

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/lucasgjanot/go-http-server/internal/api/v1/chirps"
	"github.com/lucasgjanot/go-http-server/internal/auth"
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

func DeleteHandler(JWTSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			httperrors.Write(w, httperrors.UnauthorizedErr)
			return
		}
		userID, err := auth.ValidateJWT(token, JWTSecret)
		if err != nil {
			httperrors.Write(w, httperrors.UnauthorizedErr)
			return
		}
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

		if chirp.UserID != userID {
			httperrors.Write(w, httperrors.ForbiddenErr)
			return
		}

		_, err = database.Chirps.DeleteChirp(
			r.Context(),
			chirp.ID,
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