package chirps

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lucasgjanot/go-http-server/internal/auth"
	"github.com/lucasgjanot/go-http-server/internal/database"
	httperrors "github.com/lucasgjanot/go-http-server/internal/errors"
)

const (
	MAX_CHARACTERS int = 140
)

var BAD_WORDS []string = []string{"kerfuffle", "sharbert", "fornax"}

type ChirpsResponse struct {
	Id uuid.UUID `json:"id"`
	Body string `json:"body"`
	UserId uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			chirps []database.Chirp
			err error
		)

		authorID := r.URL.Query().Get("author_id")
		desc := r.URL.Query().Get("sort") == "desc"
		if authorID == "" {
			chirps, err = database.Chirps.GetChirps(r.Context(), desc)
		} else {
			parsedAuthorID, err := uuid.Parse(authorID)
			if err != nil {
				log.Printf("Error parsing uuid: %s", err)
				httperrors.Write(w, httperrors.InternalServerErr)
				return
			}
			chirps, err = database.Chirps.GetChirpsByUserId(
				r.Context(),
				database.GetChirpsByUserIdParams{
					UserID: parsedAuthorID,
					IsDesc: desc,
				},
			)
		}
		
		if err != nil {
			log.Printf("Error getting chirps: %s", err)
			httperrors.Write(w,httperrors.InternalServerErr)
			return
		}
		if chirps == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode([]ChirpsResponse{})
			return 
		}
		
		chirpsResponses := []ChirpsResponse{}
		for _, chirp := range chirps {
			chirpsResponses = append(chirpsResponses, ChirpsResponse{
				Id: chirp.ID,
				Body: chirp.Body,
				UserId: chirp.UserID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(chirpsResponses)
	}
}
func PostHandler(JWTSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Body string `json:"body"`
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
		if verr := validate(params.Body); verr != nil {
			httperrors.Write(w,*verr)
			return
		}
		newChirp, err := database.Chirps.CreateChirp(
			r.Context(),
			database.CreateChirpParams{
				Body: replaceBadWords(params.Body),
				UserID: userID,
			},
		)
		if err != nil {
			log.Printf("Error creating Chirp: %s", err)
			httperrors.Write(w, httperrors.InternalServerErr)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(
			ChirpsResponse{
				Id: newChirp.ID,
				Body: newChirp.Body,
				UserId: newChirp.UserID,
				CreatedAt: newChirp.CreatedAt,
				UpdatedAt: newChirp.UpdatedAt,
			},
		)
		
	}
}

func validate(body string) *httperrors.ErrorResponse {
	body = strings.TrimSpace(body)

	if body == "" {
		err := httperrors.BadRequestErr
		err.Message = "Chirp body cannot be empty"
		return &err
	}

	if len(body) > MAX_CHARACTERS {
		err := httperrors.BadRequestErr
		err.Message = "Chirp is too long"
		err.Action = "Reduce the chirp length to 140 characters"
		return &err
	}
	return nil
}

func replaceBadWords(s string) string {
	words := strings.Fields(s)

	for i, word := range words {
		trimmed := strings.Trim(word, ".,!?;:")
		lower := strings.ToLower(trimmed)

		if slices.Contains(BAD_WORDS, lower) {
			words[i] = strings.Replace(word, trimmed, "****", 1)
		}
	}

	return strings.Join(words, " ")
}
