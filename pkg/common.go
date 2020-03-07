package pkg

import (
	"errors"
	"net/http"
	"time"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type config struct {
	GoogleClientID     string `envconfig:"GOOGLE_CLIENT_ID" required:"true"`
	GoogleClientSecret string `envconfig:"GOOGLE_CLIENT_SECRET" required:"true"`
}

var (
	logger *zap.Logger

	cfg config

	oauthConfig = &oauth2.Config{
		ClientID:     "", // filled by config
		ClientSecret: "", // filled by config
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
		RedirectURL: "urn:ietf:wg:oauth:2.0:oob",
		Scopes:      []string{youtube.YoutubeReadonlyScope},
	}
)

func init() {
	var err error

	logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}

	err = envconfig.Process("", &cfg)
	if err != nil {
		logger.Fatal("failure reading config", zap.Error(err))
	}

	oauthConfig.ClientID = cfg.GoogleClientID
	oauthConfig.ClientSecret = cfg.GoogleClientSecret
}

func ytServiceViaBasicAuth(r *http.Request) (*youtube.Service, error) {
	_, password, ok := r.BasicAuth()
	if !ok || password == "" {
		return nil, errors.New("please provide refresh token")
	}

	oauthToken := &oauth2.Token{
		Expiry:       time.Now(),
		TokenType:    "Bearer",
		RefreshToken: password,
	}

	return youtube.NewService(
		r.Context(),
		option.WithTokenSource(oauthConfig.TokenSource(r.Context(), oauthToken)),
	)
}
