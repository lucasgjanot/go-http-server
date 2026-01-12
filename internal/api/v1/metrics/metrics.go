package metrics

import (
	"fmt"
	"net/http"

	"github.com/lucasgjanot/go-http-server/internal/config"
)

func GetHandler(m *config.Metrics) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write(
			fmt.Appendf(
				[]byte{},
`<html>
	<body>
    	<h1>Welcome, Chirpy Admin</h1>
    	<p>Chirpy has been visited %d times!</p>
  	</body>
</html>`,
				m.Hits(),
			),
		)
	}
}