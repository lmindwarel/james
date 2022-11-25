package controller

import (
	"github.com/lmindwarel/james/backend/models"
	"github.com/lmindwarel/james/backend/spotify"
	"github.com/lmindwarel/james/backend/utils"
	"github.com/pkg/errors"
)

func (ctrl *Controller) CreateSpotifyCredential(payload models.CredentialPatch) (credential models.SpotifyCredential, err error) {
	if payload.User == nil || payload.Password == nil {
		return credential, models.ErrBadParameter
	}

	hashedPassword, err := utils.Encrypt(models.SpotifyPasswordHashKey, *payload.Password)
	if err != nil {
		return credential, errors.Wrap(err, "failed to encrypt password")
	}

	credential = models.SpotifyCredential{
		User:           *payload.User,
		HashedPassword: hashedPassword,
	}

	credential, err = ctrl.ds.UpsertSpotifyCredential(credential)
	if err != nil {
		return credential, errors.Wrap(err, "failed to upsert spotify credential")
	}

	return
}

func (ctrl *Controller) PatchSpotifyCredential(credential models.SpotifyCredential, patch models.CredentialPatch) (models.SpotifyCredential, error) {
	var err error
	if patch.User != nil {
		credential.User = *patch.User
	}

	if patch.Password != nil {
		hashedPassword, err := utils.Encrypt(models.SpotifyPasswordHashKey, *patch.Password)
		if err != nil {
			return credential, errors.Wrap(err, "failed to encrypt password")
		}

		credential.HashedPassword = hashedPassword
	}

	return credential, err
}

func (ctrl *Controller) AuthenticateSpotify(credentials models.SpotifyCredential) (err error) {
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

func (ctrl *Controller) ConnectSelectedSpotifyCredentials() error {
	usedSpotifyCredentialsID, err := ctrl.ds.GetUUIDParameter(models.ParamSpotifyCredentials)
	if err != nil {
		return errors.Wrap(err, "failed to get used spotify credentials parameter")
	}

	credentials, err := ctrl.ds.GetSpotifyCredential(usedSpotifyCredentialsID)
	if err != nil {
		return errors.Wrapf(err, "failed to get spotify credentials %s", usedSpotifyCredentialsID)
	}

	err = ctrl.AuthenticateSpotify(credentials)
	if err != nil {
		return errors.Wrapf(err, "failed to authenticate to spotify")
	}

	return nil
}
