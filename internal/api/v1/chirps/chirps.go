package chirps

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
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
		chirps, err := database.Chirps.GetChirps(r.Context())
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
				UserId: chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(chirpsResponses)
	}
}
func PostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Body string `json:"body"`
			UserId string `json:"user_id"`
		}

		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err := decoder.Decode(&params)
		if err != nil {
			log.Printf("Error decoding parameters: %s", err)
			httperrors.Write(w, httperrors.BadRequestErr)
			return
		}
		uid, err := uuid.Parse(params.UserId)
		if err != nil {
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
				UserID: uid,
			},
		)
		if err != nil {
			log.Printf("Error creating User: %s", err)
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
