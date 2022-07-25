package controller

import (
	"github.com/lmindwarel/james/backend/models"
	"github.com/pkg/errors"
)

func (ctrl *Controller) PatchAccount(account models.Account, patch models.AccountPatch) (models.Account, error) {
	if patch.Name != nil {
		account.Name = *patch.Name
	}

	if patch.Icon != nil {
		account.Icon = *patch.Icon
	}

	account, err := ctrl.ds.UpsertAccount(account)
	if err != nil {
		return account, errors.Wrap(err, "failed to upsert account")
	}

	return account, nil
}
