package chirps_test

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	chirps "github.com/lucasgjanot/go-http-server/internal/api/v1/chirps"
	"github.com/lucasgjanot/go-http-server/internal/database"
	"github.com/lucasgjanot/go-http-server/internal/testhelpers"
)

func TestGetChirps(t *testing.T) {
	test_url, _ := testhelpers.InitTest(t)

	t.Run("Anonymous user", func(t *testing.T) {
		user, _ := testhelpers.CreateUser(testhelpers.NewUser{})
		chirp1, _ := testhelpers.CreateChirp(
			database.CreateChirpParams{
				Body: "Chirp 1 body",
				UserID: user.ID,
			},
		)
		chirp2, _ := testhelpers.CreateChirp(
			database.CreateChirpParams{
				Body: "Chirp 2 body",
				UserID: user.ID,
			},
		)

		t.Run("Get data", func(t *testing.T) {
			resp, err := http.Get(test_url + "/api/chirps")
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
			
			respBody := []chirps.ChirpsResponse{}
			if err := json.Unmarshal(data, &respBody); err != nil {
				t.Fatalf("Error while decoding json: %s", err)
			}

			expected := []chirps.ChirpsResponse{
				{
					Id: chirp1.ID,
					Body: chirp1.Body,
					UserId: chirp1.UserID,
					CreatedAt: chirp1.CreatedAt,
					UpdatedAt: chirp1.UpdatedAt,
				},
				{
					Id: chirp2.ID,
					Body: chirp2.Body,
					UserId: chirp2.UserID,
					CreatedAt: chirp2.CreatedAt,
					UpdatedAt: chirp2.UpdatedAt,
				},
			}

			if diff := cmp.Diff(expected, respBody); diff != "" {
				t.Fatalf("mismatch (-want +got):\n%s", diff)
			}
		}) 
	})
}
