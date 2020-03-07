package pkg

import (
	"net/http"

	"golang.org/x/oauth2"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	authURL := oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	http.Redirect(w, r, authURL, http.StatusFound)
}
