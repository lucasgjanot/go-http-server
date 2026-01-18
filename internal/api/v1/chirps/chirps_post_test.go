package chirps_test

import (
	"bytes"
	"fmt"
	"strings"

	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	chirps "github.com/lucasgjanot/go-http-server/internal/api/v1/chirps"
	httperrors "github.com/lucasgjanot/go-http-server/internal/errors"
	"github.com/lucasgjanot/go-http-server/internal/testhelpers"
)

func TestPostChirps(t *testing.T) {
	test_url, cfg := testhelpers.InitTest(t)
	user, _ := testhelpers.CreateUser(testhelpers.NewUser{})
	token, _ := testhelpers.AuthenticateUser(user, cfg.JWTSecret)
	t.Run("Authenticated user", func(t *testing.T) {	
		t.Run("Valid chirp", func(t *testing.T) {
			payload := map[string]string{
				"body": "valid chirp",
			}

			jsonData, err := json.Marshal(payload)
			if err != nil {
				t.Fatalf("failed to marshal payload: %v", err)
			}
			req, err := http.NewRequest(
				"POST",
				test_url + "/api/chirps",
				bytes.NewBuffer(jsonData),
			)
			if err != nil {
				t.Fatalf("Request failed: %s", err)
			}

			
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			client := http.DefaultClient
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Request failed: %s", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusCreated {
				t.Fatalf("Expected Status 201 got: %d", resp.StatusCode)
			}

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Error while reading body: %s", err)
			}
			
			respBody := chirps.ChirpsResponse{}
			if err := json.Unmarshal(data, &respBody); err != nil {
				t.Fatalf("Error while decoding json: %s", err)
			}

			expected := chirps.ChirpsResponse{
				Id: respBody.Id,
				Body: "valid chirp",
				UserId: user.ID,
				CreatedAt: respBody.CreatedAt,
				UpdatedAt: respBody.UpdatedAt,
			}

			if diff := cmp.Diff(expected, respBody); diff != "" {
				t.Fatalf("mismatch (-want +got):\n%s", diff)
			}
		}) 
		t.Run("Chirp is to long", func(t *testing.T) {
			payload := map[string]string{
				"body": strings.Repeat("a", chirps.MAX_CHARACTERS+1),
			}

			jsonData, err := json.Marshal(payload)
			if err != nil {
				t.Fatalf("failed to marshal payload: %v", err)
			}
			req, err := http.NewRequest(
				"POST",
				test_url + "/api/chirps",
				bytes.NewBuffer(jsonData),
			)
			if err != nil {
				t.Fatalf("Request failed: %s", err)
			}

			
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			client := http.DefaultClient
			resp, err := client.Do(req)
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

			if respBody.Message != "Chirp is too long" {
				t.Fatalf("Expected respBody.Message to be 'Chirp is too long', got: %v", respBody.Message)
			}
		}) 
		t.Run("Chirp is empty", func(t *testing.T) {
			payload := map[string]string{
				"body": "",
				"user_id": user.ID.String(),
			}

			jsonData, err := json.Marshal(payload)
			if err != nil {
				t.Fatalf("failed to marshal payload: %v", err)
			}
			req, err := http.NewRequest(
				"POST",
				test_url + "/api/chirps",
				bytes.NewBuffer(jsonData),
			)
			if err != nil {
				t.Fatalf("Request failed: %s", err)
			}

			
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			client := http.DefaultClient
			resp, err := client.Do(req)
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

			jsonData, err := json.Marshal(payload)
			if err != nil {
				t.Fatalf("failed to marshal payload: %v", err)
			}
			req, err := http.NewRequest(
				"POST",
				test_url + "/api/chirps",
				bytes.NewBuffer(jsonData),
			)
			if err != nil {
				t.Fatalf("Request failed: %s", err)
			}

			
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			client := http.DefaultClient
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Request failed: %s", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusCreated {
				t.Fatalf("Expected Status 201 got: %d", resp.StatusCode)
			}

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Error while reading body: %s", err)
			}
			
			respBody := chirps.ChirpsResponse{}
			if err := json.Unmarshal(data, &respBody); err != nil {
				t.Fatalf("Error while decoding json: %s", err)
			}

			expected := chirps.ChirpsResponse{
				Id: respBody.Id,
				Body: "totalMatch, ****",
				UserId: user.ID,
				CreatedAt: respBody.CreatedAt,
				UpdatedAt: respBody.UpdatedAt,
			}

			if diff := cmp.Diff(expected, respBody); diff != "" {
				t.Fatalf("mismatch (-want +got):\n%s", diff)
			}
		}) 
		t.Run("Different case bad word", func(t *testing.T) {
			payload := map[string]string{
				"body": "totalMatch, kerFuffle",
			}

			jsonData, err := json.Marshal(payload)
			if err != nil {
				t.Fatalf("failed to marshal payload: %v", err)
			}
			req, err := http.NewRequest(
				"POST",
				test_url + "/api/chirps",
				bytes.NewBuffer(jsonData),
			)
			if err != nil {
				t.Fatalf("Request failed: %s", err)
			}

			
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			client := http.DefaultClient
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Request failed: %s", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusCreated {
				t.Fatalf("Expected Status 201 got: %d", resp.StatusCode)
			}

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Error while reading body: %s", err)
			}
			
			respBody := chirps.ChirpsResponse{}
			if err := json.Unmarshal(data, &respBody); err != nil {
				t.Fatalf("Error while decoding json: %s", err)
			}

			expected := chirps.ChirpsResponse{
				Id: respBody.Id,
				Body: "totalMatch, ****",
				UserId: user.ID,
				CreatedAt: respBody.CreatedAt,
				UpdatedAt: respBody.UpdatedAt,
			}

			if diff := cmp.Diff(expected, respBody); diff != "" {
				t.Fatalf("mismatch (-want +got):\n%s", diff)
			}
		}) 
	})
	t.Run("Anonymous user", func(t *testing.T) {
		t.Run("Creating Chirp", func(t *testing.T) {
			payload := map[string]string{
				"body": "valid body",
			}

			body, err := json.Marshal(payload)
			if err != nil {
				t.Fatalf("failed to marshal payload: %v", err)
			}

			resp, err := http.Post(
				test_url + "/api/chirps",
				"application/json",
				bytes.NewBuffer(body),
			)
			if err != nil {
				t.Fatalf("Request failed: %s", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusUnauthorized {
				t.Fatalf("Expected Status 401 got: %d", resp.StatusCode)
			}

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Error while reading body: %s", err)
			}
			
			respBody := httperrors.ErrorResponse{}
			if err := json.Unmarshal(data, &respBody); err != nil {
				t.Fatalf("Error while decoding json: %s", err)
			}

			expected := httperrors.UnauthorizedErr

			if diff := cmp.Diff(expected, respBody); diff != "" {
				t.Fatalf("mismatch (-want +got):\n%s", diff)
			}
		}) 
	}) 
}
