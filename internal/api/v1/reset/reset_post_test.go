package reset_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/lucasgjanot/go-http-server/internal/testhelpers"
)

func TestPostReset(t *testing.T) {
	test_url, _ := testhelpers.InitTest(t)
	t.Run("Anonymous user", func(t *testing.T) {

		resp, err := http.Post(test_url + "/admin/reset", "", nil)
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

