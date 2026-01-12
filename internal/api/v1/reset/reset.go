package reset

import (
	"net/http"

	"github.com/lucasgjanot/go-http-server/internal/config"
)

func PostHandler(m *config.Metrics) http.HandlerFunc {
	return  func(w http.ResponseWriter, r *http.Request) {
		m.Reset()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hits reset to 0"))
	}
}