package controller

import (
	"github.com/lmindwarel/james/backend/models"
	"github.com/lmindwarel/james/backend/spotify"
	"github.com/lmindwarel/james/backend/utils"
	"github.com/pkg/errors"
)

func (ctrl *Controller) AuthenticateSpotify(credentials models.SpotifyCredentials) (err error) {
	// decrypt password
	password, err := utils.Decrypt(models.SpotifyPasswordHashKey, credentials.HashedPassword)
	if err != nil {
		return errors.Wrap(err, "failed to decrypt password")
	}

	ctrl.spotifySession, err = spotify.Authenticate(ctrl.config.SpotifyClientID, ctrl.config.SpotifyClientSecret, credentials.User, password)
	if err != nil {
		return errors.Wrap(err, "failed to authenticate to spotify")
	}

	return
}
