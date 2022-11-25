package controller

import (
	"github.com/lmindwarel/james/backend/models"
	"github.com/pkg/errors"
)

func (ctrl *Controller) UpdateParam(param models.Parameter) (err error) {
	switch param.ID {
	case models.ParamSpotifyCredentials:
		credentials, err := ctrl.ds.GetSpotifyCredential(param.Value.(models.UUID))
		if err != nil {
			return errors.Wrapf(err, "failed to get spotify credentials %s", param.Value)
		}

		err = ctrl.AuthenticateSpotify(credentials)
		if err != nil {
			return errors.Wrapf(err, "failed to authenticate to spotify %s", param.ID)
		}
	}

	_, err = ctrl.ds.UpsertParameter(param)
	if err != nil {
		return errors.Wrapf(err, "failed to upsert parameter %s", param.ID)
	}

	return nil
}
