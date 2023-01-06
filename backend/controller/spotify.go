package controller

import (
	"github.com/lmindwarel/james/backend/models"
	"github.com/lmindwarel/james/backend/spotify"
	"github.com/lmindwarel/james/backend/utils"
	"github.com/pkg/errors"
)

func (ctrl *Controller) GetSpotifySession() (*spotify.Session, error) {
	if ctrl.spotifySession == nil {
		return nil, errors.New("spotify not connected")
	}

	return ctrl.spotifySession, nil
}

func (ctrl *Controller) IsSpotifyConnected() bool {
	return ctrl.spotifySession != nil
}

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

func (ctrl *Controller) AuthenticateSpotify(credential models.SpotifyCredential) (err error) {
	// decrypt password
	password, err := utils.Decrypt(models.SpotifyPasswordHashKey, credential.HashedPassword)
	if err != nil {
		return errors.Wrap(err, "failed to decrypt password")
	}

	log.Infof("authenticating user %s to Spotify with password %s", credential.User, password)

	ctrl.spotifySession, err = spotify.Authenticate(ctrl.config.SpotifyClientID, ctrl.config.SpotifyClientSecret, credential.User, password)
	if err != nil {
		return errors.Wrap(err, "failed to authenticate to spotify")
	}

	err = ctrl.updateJamesStatus(models.JamesStatusPatch{
		AuthenticatedSpotifyCredentialID: &credential.ID,
	})
	if err != nil {
		return errors.Wrap(err, "failed to update james state")
	}

	return
}

func (ctrl *Controller) GetCurrentSpotifyCredential() (credential *models.SpotifyCredential, err error) {
	selectedCredentialID, notFatalErr := ctrl.ds.GetUUIDParameter(models.ParamCurrentSpotifyCredential)
	if notFatalErr != nil {
		if !ctrl.ds.IsNotFoundError(notFatalErr) {
			return credential, errors.Wrap(err, "failed to get selected spotify credential parameter")
		}

		return nil, nil // no credentials selected
	}

	credentialTmp, err := ctrl.ds.GetSpotifyCredential(selectedCredentialID)
	if err != nil {
		return credential, errors.Wrapf(err, "failed to get selected spotify credential %s", selectedCredentialID)
	}

	credential = &credentialTmp

	return
}

func (ctrl *Controller) AuthenticateCurrentSpotifyCredential() (err error) {
	currentCredential, err := ctrl.GetCurrentSpotifyCredential()
	if err != nil {
		return errors.Wrap(err, "failed to get current spotify credential")
	}

	if currentCredential == nil {
		return
	}

	err = ctrl.AuthenticateSpotify(*currentCredential)
	if err != nil {
		return errors.Wrap(err, "failed to authenticate to spotify")
	}

	return
}
