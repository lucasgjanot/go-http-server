package validatechirp

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"

	httperrors "github.com/lucasgjanot/go-http-server/internal/errors"
)

const (
	MAX_CHARACTERS int = 140
)

var BAD_WORDS []string = []string{"kerfuffle", "sharbert", "fornax"}

type ValidateChirpResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

func GetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Body string `json:"body"`
		}

		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err := decoder.Decode(&params)
		if err != nil {
			log.Printf("Error decoding parameters: %s", err)
			httperrors.Write(w, httperrors.BadRequestErr)
			return
		}
		if verr := validate(params.Body); verr != nil {
			httperrors.Write(w,*verr)
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(ValidateChirpResponse{
			CleanedBody: replaceBadWords(params.Body),
		})
		
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
