package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lmindwarel/james/backend/models"
	"github.com/pkg/errors"
)

func (a *API) GetAccounts(c *gin.Context) {
	accounts, err := a.ds.GetAccounts()
	if err != nil {
		c.AbortWithError(http.StatusNotFound, errors.Wrap(err, "failed to get accounts"))
		return
	}

	c.JSON(http.StatusOK, accounts)
}

func (a *API) PostAccount(c *gin.Context) {
	var patch models.AccountPatch
	if err := c.ShouldBindJSON(&patch); err != nil {
		log.Errorf("failed to decode account patch: %s", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	account, err := a.ctrl.PatchAccount(models.Account{}, patch)
	if err != nil {
		log.Errorf("failed to patch account: %s", err)
		c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "failed to patch account"))
		return
	}

	c.JSON(http.StatusNoContent, account)
}
