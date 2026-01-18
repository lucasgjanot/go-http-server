package login

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lucasgjanot/go-http-server/internal/auth"
	"github.com/lucasgjanot/go-http-server/internal/database"
	httperrors "github.com/lucasgjanot/go-http-server/internal/errors"
)

type LoginResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	RefreshToken string `json:"refresh_token"`
	IsChirpyRed bool `json:"is_chirpy_red"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const (
	JWTDefaultExpire = time.Hour
	RefreshTokenExpires = 24 * 60 * time.Hour
)


func PostHandler(JWTSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Parameters struct {
			Password         string `json:"password"`
			Email            string `json:"email"`
		}

		var params Parameters
		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&params); err != nil {
			httperrors.Write(w, httperrors.BadRequestErr)
			return
		}

		if params.Password == "" || params.Email == "" {
			httperrors.Write(w, httperrors.BadRequestErr)
			return
		}

		// expiration handling (seconds -> duration)
		expire := JWTDefaultExpire


		user, err := database.Users.GetUserByEmail(r.Context(), params.Email)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				httperrors.Write(w, httperrors.UnauthorizedErr)
				return
			}
			httperrors.Write(w, httperrors.ServiceUnavailableErr)
			return
		}

		ok, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
		if err != nil {
			log.Printf("password check failed: %v", err)
			httperrors.Write(w, httperrors.InternalServerErr)
			return
		}

		if !ok {
			httperrors.Write(w, httperrors.UnauthorizedErr)
			return
		}

		jwtToken, err := auth.MakeJWT(user.ID, JWTSecret, expire)
		if err != nil {
			log.Printf("error creating jwt: %v", err)
			httperrors.Write(w, httperrors.InternalServerErr)
			return
		}

		refreshTokenString, _ := auth.MakeRefreshToken()
		now := time.Now().UTC()
		refreshToken, err := database.Auth.CreateRefreshToken(
			r.Context(),
			database.CreateRefreshTokenParams{
				Token: refreshTokenString,
				UserID: user.ID,
				ExpiresAt: now.Add(RefreshTokenExpires),
				CreatedAt: now,
				UpdatedAt: now,
			},
		)
		if err != nil {
			log.Printf("error creating refreshtoken: %v", err)
			httperrors.Write(w, httperrors.InternalServerErr)
			return 
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(LoginResponse{
			ID:        user.ID,
			Email:     user.Email,
			Token:     jwtToken,
			RefreshToken: refreshToken.Token,
			IsChirpyRed: user.IsChirpyRed,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}
}
