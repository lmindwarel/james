package spotify

import (
	"context"

	"github.com/librespot-org/librespot-golang/librespot"
	"github.com/lmindwarel/james/backend/models"
	"github.com/pkg/errors"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

func Authenticate(clientID, clientSecret, userID, userSecret string) (session *Session, err error) {
	ctx := context.Background()

	log.Debugf("userID: %s, userSecret: %s", userID, userSecret)

	// web api
	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     spotifyauth.TokenURL,
	}
	token, err := config.Token(ctx)
	if err != nil {
		return session, errors.Wrap(err, "failed to get token")
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	webapiClient := spotify.New(httpClient)

	// librespot
	librespotSession, err := librespot.Login(userID, userSecret, models.SpotifyDeviceName)
	if err != nil {
		return session, errors.Wrap(err, "failed to authenticate to Spotify")
	}

	return &Session{
		librespotSession: librespotSession,
		webapiClient:     webapiClient,
		userID:           userID,
	}, err
}