package login_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/lucasgjanot/go-http-server/internal/api/v1/users"
	httperrors "github.com/lucasgjanot/go-http-server/internal/errors"
	"github.com/lucasgjanot/go-http-server/internal/testhelpers"
)

func TestPostLogin(t *testing.T) {
	test_url := testhelpers.InitTest(t)
	user, _ :=  testhelpers.CreateUser(testhelpers.NewUser{
		Email: "correct.email@example.com",
		Password: "correctpassword",
	})

	t.Run("Anonymous user", func(t *testing.T) {
		t.Run("correct email and password", func(t *testing.T) {
			type Payload struct {
				Email string `json:"email"`
				Password string `json:"password"`
			}
			payload := Payload{Email: "correct.email@example.com", Password: "correctpassword"}

			jsonData, err := json.Marshal(payload)
			if err != nil {
				t.Fatalf("failed to marshal payload: %v", err)
			} 

			resp, err := http.Post(test_url + "/api/login", "application/json", bytes.NewBuffer(jsonData))
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

			respBody := users.UserResponse{}
			if err := json.Unmarshal(data, &respBody); err != nil {
				t.Fatalf("Error while decoding json: %s", err)
			}

			expected := users.UserResponse{
				Id: user.ID,
				Email: user.Email,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			}

			if diff := cmp.Diff(expected, respBody); diff != "" {
				t.Fatalf("mismatch (-want +got):\n%s", diff)
			}
		})
		t.Run("invalid email and correct password", func(t *testing.T) {
			type Payload struct {
				Email string `json:"email"`
				Password string `json:"password"`
			}
			payload := Payload{Email: "invalid.email@example.com", Password: "correctpassword"}

			jsonData, err := json.Marshal(payload)
			if err != nil {
				t.Fatalf("failed to marshal payload: %v", err)
			} 

			resp, err := http.Post(test_url + "/api/login", "application/json", bytes.NewBuffer(jsonData))
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
		t.Run("valid email and invalid password", func(t *testing.T) {
			type Payload struct {
				Email string `json:"email"`
				Password string `json:"password"`
			}
			payload := Payload{Email: "correct.email@example.com", Password: "invalidpassword"}

			jsonData, err := json.Marshal(payload)
			if err != nil {
				t.Fatalf("failed to marshal payload: %v", err)
			} 

			resp, err := http.Post(test_url + "/api/login", "application/json", bytes.NewBuffer(jsonData))
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
		t.Run("invalid email and invalid password", func(t *testing.T) {
			type Payload struct {
				Email string `json:"email"`
				Password string `json:"password"`
			}
			payload := Payload{Email: "invalid.email@example.com", Password: "invalidpassword"}

			jsonData, err := json.Marshal(payload)
			if err != nil {
				t.Fatalf("failed to marshal payload: %v", err)
			} 

			resp, err := http.Post(test_url + "/api/login", "application/json", bytes.NewBuffer(jsonData))
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