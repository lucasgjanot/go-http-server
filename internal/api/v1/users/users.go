package users

import (
	"encoding/json"
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
			CreatedAt: newUser.CreatedAt,
			UpdatedAt: newUser.UpdatedAt,
		})
	}
}