package metrics_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lucasgjanot/go-http-server/internal/config"
	"github.com/lucasgjanot/go-http-server/internal/router"
)

func TestGetMetrics(t *testing.T) {
	cfg := config.NewConfig()
	srv := router.New(cfg)

	ts := httptest.NewServer(srv.Handler)
	defer ts.Close()
	t.Run("Anonymous user", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/admin/metrics")
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
		html := string(data)
		if !strings.Contains(html,"Welcome, Chirpy Admin") {
			t.Fatalf("Invalid html response")
		}
	})

}