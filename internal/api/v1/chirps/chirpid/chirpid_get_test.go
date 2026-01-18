package chirpsid_test

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/lucasgjanot/go-http-server/internal/api/v1/chirps"
	"github.com/lucasgjanot/go-http-server/internal/database"
	"github.com/lucasgjanot/go-http-server/internal/testhelpers"
)

func TestGetChirpID(t *testing.T) {
	test_url, _ := testhelpers.InitTest(t)
	user, _ := testhelpers.CreateUser(testhelpers.NewUser{})
	chirp, _ := testhelpers.CreateChirp(database.CreateChirpParams{
		Body: "existentuuid",
		UserID: user.ID,
	})

	t.Run("Anonymous user", func(t *testing.T) {
		t.Run("using a existent uuid", func(t *testing.T) {
			resp, err := http.Get(test_url + "/api/chirps/" + chirp.ID.String())
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
			
			respBody := chirps.ChirpsResponse{}
			if err := json.Unmarshal(data, &respBody); err != nil {
				t.Fatalf("Error while decoding json: %s", err)
			}

			expected := chirps.ChirpsResponse{
				Id: chirp.ID,
				Body: chirp.Body,
				UserId: chirp.UserID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
			}

			if diff := cmp.Diff(expected, respBody); diff != "" {
				t.Fatalf("mismatch (-want +got):\n%s", diff)
			}
		})
	})
}

