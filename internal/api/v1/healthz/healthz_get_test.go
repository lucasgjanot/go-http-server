package healthz_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/lucasgjanot/go-http-server/internal/testhelpers"
)

type healthzResponse struct {
	Status string `json:"status"`
}

func TestGetHealthz(t *testing.T) {
	test_url := testhelpers.InitTest(t)

	t.Run("GET /api/healthz", func(t *testing.T) {
		t.Run("Anonymous user", func(t *testing.T) {
			resp, err := http.Get(test_url + "/api/healthz")
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Fatalf("expected 200, got %d", resp.StatusCode)
			}

			var body healthzResponse
			if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
				t.Fatalf("failed to decode JSON: %v", err)
			}

			if body.Status != "ok" {
				t.Fatalf("expected status 'ok', got %q", body.Status)
			}
		})
	})
}
