package login

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/lucasgjanot/go-http-server/internal/api/v1/users"
	"github.com/lucasgjanot/go-http-server/internal/auth"
	"github.com/lucasgjanot/go-http-server/internal/database"
	httperrors "github.com/lucasgjanot/go-http-server/internal/errors"
)

func PostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Parameters struct {
			Password string `json:"password"`
			Email string `json:"Email"`
		}

		params := Parameters{}
		decoder := json.NewDecoder(r.Body)
		
		if err := decoder.Decode(&params); err != nil {
			httperrors.Write(w, httperrors.BadRequestErr) 
			return
		}

		if len(params.Password) == 0 || len(params.Email) == 0 {
			httperrors.Write(w, httperrors.BadRequestErr) 
			return
		}

		user, err := database.Users.GetUserByEmail(r.Context(), params.Email)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				httperrors.Write(w, httperrors.UnauthorizedErr)
				return
			}
			httperrors.Write(w, httperrors.ServiceUnavailableErr)
			return
		}

		check, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
		if err != nil {
			log.Printf("Password check failed: %s", err)
			httperrors.Write(w, httperrors.InternalServerErr)
			return
		}

		if !check {
			httperrors.Write(w, httperrors.UnauthorizedErr)
			return 
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(users.UserResponse{
			Id: user.ID,
			Email: user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
		
	}
}