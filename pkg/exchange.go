package pkg

import (
	"fmt"
	"net/http"
)

func ExchangeHandler(w http.ResponseWriter, r *http.Request) {
	token, err := oauthConfig.Exchange(r.Context(), r.URL.Query().Get("token"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	fmt.Fprintf(w, "Refresh Token: %s", token.RefreshToken)
}
