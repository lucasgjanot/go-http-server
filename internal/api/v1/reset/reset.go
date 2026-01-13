package reset

import (
	"log"
	"net/http"

	"github.com/lucasgjanot/go-http-server/internal/config"
	"github.com/lucasgjanot/go-http-server/internal/database"
	httperrors "github.com/lucasgjanot/go-http-server/internal/errors"
)

func PostHandler(m *config.Metrics) http.HandlerFunc {
	return  func(w http.ResponseWriter, r *http.Request) {
		if database.Users == nil {
			httperrors.Write(w, httperrors.ServiceUnavailableErr)
			return 
		}
		if _, err := database.Users.DeleteAllUsers(r.Context()); err != nil {
			log.Printf("Error reseting database: %v", err)
			httperrors.Write(w, httperrors.InternalServerErr)
			return 
		}
		log.Printf("Database reseted")
		m.Reset()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hits reset to 0"))
	}
}