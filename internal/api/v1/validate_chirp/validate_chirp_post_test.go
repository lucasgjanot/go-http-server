package validatechirp_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	validatechirp "github.com/lucasgjanot/go-http-server/internal/api/v1/validate_chirp"
	"github.com/lucasgjanot/go-http-server/internal/config"
	httperrors "github.com/lucasgjanot/go-http-server/internal/errors"
	"github.com/lucasgjanot/go-http-server/internal/router"
)

func TestGetValidateChirp(t *testing.T) {
	cfg := config.NewConfig()
	r := router.New(cfg)

	ts := httptest.NewServer(r.Handler)
	defer ts.Close()

	t.Run("Anonymous user", func(t *testing.T) {
		t.Run("Valid chirpy", func(t *testing.T) {
			payload := map[string]string{
				"body": "valid chirp",
			}

			body, err := json.Marshal(payload)
			if err != nil {
				t.Fatalf("failed to marshal payload: %v", err)
			}

			resp, err := http.Post(
				ts.URL + "/api/validate_chirp",
				"application/json",
				bytes.NewBuffer(body),
			)
			if err != nil {
				t.Fatalf("Request failed: %s", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Fatalf("Expected Status 200 got: %d", resp.StatusCode)
			}

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Error while reading body: %s", err)
			}
			
			respBody := validatechirp.ValidateChirpResponse{}
			if err := json.Unmarshal(data, &respBody); err != nil {
				t.Fatalf("Error while decoding json: %s", err)
			}

			if respBody.CleanedBody != "valid chirp" {
				t.Fatalf("Expected respBody.CleanedBody to be 'valid chirp', got: %v", respBody.CleanedBody)
			}
		}) 
		t.Run("Chirp is to long", func(t *testing.T) {
			payload := map[string]string{
				"body": strings.Repeat("a", validatechirp.MAX_CHARACTERS+1),
			}

			body, err := json.Marshal(payload)
			if err != nil {
				t.Fatalf("failed to marshal payload: %v", err)
			}

			resp, err := http.Post(
				ts.URL + "/api/validate_chirp",
				"application/json",
				bytes.NewBuffer(body),
			)
			if err != nil {
				t.Fatalf("Request failed: %s", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != 400 {
				t.Fatalf("Expected Status 400 got: %d", resp.StatusCode)
			}

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Error while reading body: %s", err)
			}
			
			respBody := httperrors.ErrorResponse{}
			if err := json.Unmarshal(data, &respBody); err != nil {
				t.Fatalf("Error while decoding json: %s", err)
			}

			if respBody.Message != "Chirp is too long" {
				t.Fatalf("Expected respBody.Message to be 'Chirp is too long', got: %v", respBody.Message)
			}
		}) 
		t.Run("Chirp is empty", func(t *testing.T) {
			payload := map[string]string{
				"body": "",
			}

			body, err := json.Marshal(payload)
			if err != nil {
				t.Fatalf("failed to marshal payload: %v", err)
			}

			resp, err := http.Post(
				ts.URL + "/api/validate_chirp",
				"application/json",
				bytes.NewBuffer(body),
			)
			if err != nil {
				t.Fatalf("Request failed: %s", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusBadRequest {
				t.Fatalf("Expected Status 400 got: %d", resp.StatusCode)
			}

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Error while reading body: %s", err)
			}
			
			respBody := httperrors.ErrorResponse{}
			if err := json.Unmarshal(data, &respBody); err != nil {
				t.Fatalf("Error while decoding json: %s", err)
			}

			if respBody.Message != "Chirp body cannot be empty" {
				t.Fatalf("Expected respBody.Message to be 'Chirp body cannot be empty', got: %v", respBody.Message)
			}
		}) 
		t.Run("Total match bad word", func(t *testing.T) {
			payload := map[string]string{
				"body": "totalMatch, kerfuffle",
			}

			body, err := json.Marshal(payload)
			if err != nil {
				t.Fatalf("failed to marshal payload: %v", err)
			}

			resp, err := http.Post(
				ts.URL + "/api/validate_chirp",
				"application/json",
				bytes.NewBuffer(body),
			)
			if err != nil {
				t.Fatalf("Request failed: %s", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Fatalf("Expected Status 200 got: %d", resp.StatusCode)
			}

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Error while reading body: %s", err)
			}
			
			respBody := validatechirp.ValidateChirpResponse{}
			if err := json.Unmarshal(data, &respBody); err != nil {
				t.Fatalf("Error while decoding json: %s", err)
			}

			if respBody.CleanedBody != "totalMatch, ****" {
				t.Fatalf("Expected respBody.CleanedBody to be 'totalMatch, ****', got: %v", respBody.CleanedBody)
			}
		}) 
		t.Run("Different case bad word", func(t *testing.T) {
			payload := map[string]string{
				"body": "totalMatch, kerFuffle",
			}

			body, err := json.Marshal(payload)
			if err != nil {
				t.Fatalf("failed to marshal payload: %v", err)
			}

			resp, err := http.Post(
				ts.URL + "/api/validate_chirp",
				"application/json",
				bytes.NewBuffer(body),
			)
			if err != nil {
				t.Fatalf("Request failed: %s", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Fatalf("Expected Status 200 got: %d", resp.StatusCode)
			}

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Error while reading body: %s", err)
			}
			
			respBody := validatechirp.ValidateChirpResponse{}
			if err := json.Unmarshal(data, &respBody); err != nil {
				t.Fatalf("Error while decoding json: %s", err)
			}

			if respBody.CleanedBody != "totalMatch, ****" {
				t.Fatalf("Expected respBody.CleanedBody to be 'totalMatch, ****', got: %v", respBody.CleanedBody)
			}
		}) 
	})
}
