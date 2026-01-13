package reset_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lucasgjanot/go-http-server/internal/config"
	"github.com/lucasgjanot/go-http-server/internal/router"
	"github.com/lucasgjanot/go-http-server/internal/testhelpers"
)

func TestPostReset(t *testing.T) {
	t.Setenv(
		"DB_URL",
		"postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable",
	)
	cfg := config.NewConfig()
	testhelpers.InitTestDB(t)
	srv := router.New(cfg)
	ts := httptest.NewServer(srv.Handler)
	defer ts.Close()
	t.Run("Anonymous user", func(t *testing.T) {

		resp, err := http.Post(ts.URL + "/admin/reset", "", nil)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		respBody := string(data)
		if respBody != "Hits reset to 0" {
			t.Fatalf("Expected 'Hits reset to 0' got %v", respBody)
		}
	})
}

