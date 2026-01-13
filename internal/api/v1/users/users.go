package users

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lucasgjanot/go-http-server/internal/database"
	httperrors "github.com/lucasgjanot/go-http-server/internal/errors"
)

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func PostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if database.Users == nil {
			httperrors.Write(w, httperrors.ServiceUnavailableErr)
		}

		type parameters struct {
			Email string `json:"email"`
		}

		params := parameters{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&params); err != nil {
			log.Printf("Error decoding parameters: %s", err)
			httperrors.Write(w, httperrors.BadRequestErr) 
			return
		}

		newUser, err := database.Users.CreateUser(r.Context(), params.Email)
		if err != nil {
			log.Printf("Error creating user: %v", err)
			httperrors.Write(w, httperrors.InternalServerErr)
			return 
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(UserResponse{
			ID: newUser.ID,
			Email: newUser.Email,
			CreatedAt: newUser.CreatedAt,
			UpdatedAt: newUser.UpdatedAt,
		})
	}
}