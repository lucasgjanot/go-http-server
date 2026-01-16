package users_test

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

func TestPostUsers(t *testing.T) {
	test_url := testhelpers.InitTest(t)
	t.Run("Anonymous user", func(t *testing.T) {
		t.Run("With valid email", func(t *testing.T) {
			 payload := map[string]string{
				"email": "valid.email@example.com",
				"password": "secure",
			}

			jsonData, err := json.Marshal(payload)
			if err != nil {
				t.Fatalf("failed to marshal payload: %v", err)
			}

			resp, err := http.Post(
				test_url + "/api/users",
				"application/json",
				bytes.NewBuffer(jsonData),
			)
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
			
			respBody := users.UserResponse{}
			if err := json.Unmarshal(data, &respBody); err != nil {
				t.Fatalf("Error while decoding json: %s", err)
			}

			expected := users.UserResponse{
				Id: respBody.Id,
				Email: "valid.email@example.com",
				CreatedAt: respBody.CreatedAt,
				UpdatedAt: respBody.UpdatedAt,
			}

			if diff := cmp.Diff(expected, respBody); diff != "" {
				t.Fatalf("mismatch (-want +got):\n%s", diff)
			}
		})
		t.Run("With empty email", func(t *testing.T) {
			 payload := map[string]string{
				"email": "",
			}

			jsonData, err := json.Marshal(payload)
			if err != nil {
				t.Fatalf("failed to marshal payload: %v", err)
			}

			resp, err := http.Post(
				test_url + "/api/users",
				"application/json",
				bytes.NewBuffer(jsonData),
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

			expected := httperrors.BadRequestErr

			if diff := cmp.Diff(expected, respBody); diff != "" {
				t.Fatalf("mismatch (-want +got):\n%s", diff)
			}
		})
		t.Run("With no body", func(t *testing.T) {
			resp, err := http.Post(
				test_url + "/api/users",
				"application/json",
				nil,
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

			expected := httperrors.BadRequestErr

			if diff := cmp.Diff(expected, respBody); diff != "" {
				t.Fatalf("mismatch (-want +got):\n%s", diff)
			}
		})
	})
}
