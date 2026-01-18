package users

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

type UserResponse struct {
	Id        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsChirpyRed bool 	`json:"is_chirpy_red"`
	
}

func PostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if database.Users == nil {
			httperrors.Write(w, httperrors.ServiceUnavailableErr)
		}

		type parameters struct {
			Email string `json:"email"`
			Password string `json:"password"`

		}
		
		params := parameters{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&params); err != nil {
			httperrors.Write(w, httperrors.BadRequestErr) 
			return
		}
	
		if len(params.Email) == 0 || len(params.Password) == 0 {
			httperrors.Write(w, httperrors.BadRequestErr)
			return
		}
		hashedPassword, err := auth.HashPassword(params.Password)
		if err != nil {
			log.Printf("Error in hashing password: %s", err)
			httperrors.Write(w, httperrors.InternalServerErr)
		}
		newUser, err := database.Users.CreateUser(
			r.Context(),
			database.CreateUserParams{
				Email: params.Email,
				HashedPassword: hashedPassword,
			},
		)
		if err != nil {
			log.Printf("Error creating user: %v", err)
			httperrors.Write(w, httperrors.InternalServerErr)
			return 
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(UserResponse{
			Id: newUser.ID,
			Email: newUser.Email,
			IsChirpyRed: newUser.IsChirpyRed,
			CreatedAt: newUser.CreatedAt,
			UpdatedAt: newUser.UpdatedAt,
			
		})
	}
}

func PutHandler(JWTSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Email string `json:"email"`
			Password string `json:"password"`
		}

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
		
		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err = decoder.Decode(&params)
		if err != nil {
			log.Printf("Error decoding parameters: %s", err)
			httperrors.Write(w, httperrors.BadRequestErr)
			return
		}

		hashedPassword, err := auth.HashPassword(params.Password)
		if err != nil {
			log.Printf("Error on hashing password: %s", err)
		}

		updatedUser, err := database.Users.UpdateUser(
			r.Context(),
			database.UpdateUserParams{
				ID: userID,
				Email: params.Email,
				HashedPassword: hashedPassword,
			},
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				httperrors.Write(w, httperrors.BadRequestErr)
				return 
			}
			httperrors.Write(w, httperrors.ServiceUnavailableErr)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(UserResponse{
			Id: updatedUser.ID,
			Email: updatedUser.Email,
			IsChirpyRed: updatedUser.IsChirpyRed,
			CreatedAt: updatedUser.CreatedAt,
			UpdatedAt: updatedUser.UpdatedAt,
		})


	}
}