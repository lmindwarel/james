package controller

import (
	"github.com/lmindwarel/james/backend/models"
	"github.com/pkg/errors"
)

func (ctrl *Controller) PatchParameter(id models.UUID, patch models.ParameterPatch) (models.Parameter, error) {
	var err error

	// get param
	param, notFatalErr := ctrl.ds.GetParameter(id)
	if notFatalErr != nil {
		if !ctrl.ds.IsNotFoundError(notFatalErr) {
			return param, errors.Wrapf(notFatalErr, "failed to get parameter %s", id)
		}

		// Param do not exist yet, create it
		param.ID = id
	}

	// set param new value
	if param.Value == patch.Value {
		log.Warning("Not changed")
		return param, nil
	}

	param.Value = patch.Value

	switch param.ID {
	case models.ParamCurrentSpotifyCredential:
		credentials, err := ctrl.ds.GetSpotifyCredential(param.Value.(models.UUID))
		if err != nil {
			return param, errors.Wrapf(err, "failed to get spotify credentials %s", param.Value)
		}

		err = ctrl.AuthenticateSpotify(credentials)
		if err != nil {
			return param, errors.Wrapf(err, "failed to authenticate to spotify %s", param.ID)
		}
	}

	log.Debugf("Upserting parameters %s", param.ID)
	param, err = ctrl.ds.UpsertParameter(param)
	if err != nil {
		return param, errors.Wrapf(err, "failed to upsert parameter %s", param.ID)
	}

	return param, err
}
