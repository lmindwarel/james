package controller

import (
	"context"

	"github.com/pkg/errors"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

func (ctrl *Controller) ConnectSpotifyAccount(id, secret string) (err error) {
	ctx := context.Background()
	config := &clientcredentials.Config{
		ClientID:     id,
		ClientSecret: secret,
		TokenURL:     spotifyauth.TokenURL,
	}

	token, err := config.Token(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get token")
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	ctrl.SpotifyClient = spotify.New(httpClient)

	return nil
}
