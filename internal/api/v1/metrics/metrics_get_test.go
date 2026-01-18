package metrics_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/lucasgjanot/go-http-server/internal/testhelpers"
)

func TestGetMetrics(t *testing.T) {
	test_url, _ := testhelpers.InitTest(t)
	t.Run("Anonymous user", func(t *testing.T) {
		resp, err := http.Get(test_url + "/admin/metrics")
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