package refresh

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/lucasgjanot/go-http-server/internal/api/v1/login"
	"github.com/lucasgjanot/go-http-server/internal/auth"
	"github.com/lucasgjanot/go-http-server/internal/database"
	httperrors "github.com/lucasgjanot/go-http-server/internal/errors"
)

type LoginResponse struct {
	Token string `json:"token"`
}

func PostHandler(JWTSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authToken, err := auth.GetBearerToken(r.Header)
		if err != nil {
			httperrors.Write(w, httperrors.UnauthorizedErr)
		}

		refreshToken, err := database.Auth.GetRefreshToken(
			r.Context(),
			authToken,
		)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				httperrors.Write(w, httperrors.UnauthorizedErr)
				return 
			}
			log.Printf("Error validating refresh token: %s", err)
			httperrors.Write(w,httperrors.ServiceUnavailableErr)
			return 
		}

		if refreshToken.ExpiresAt.Before(time.Now()) || refreshToken.RevokedAt.Valid {
			httperrors.Write(w, httperrors.UnauthorizedErr)
			return
		}

		newJWT, err := auth.MakeJWT(refreshToken.UserID, JWTSecret, login.JWTDefaultExpire)
		if err != nil {
			log.Printf("Error creating JWT token: %s", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(LoginResponse{
			Token: newJWT,
		})
	}
}