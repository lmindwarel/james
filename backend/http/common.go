package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lmindwarel/james/backend/models"
	"github.com/lmindwarel/james/backend/utils"
)

// func (a *API) GetStatus(c *gin.Context) {

// 	if a.ctrl.IsSpotifyConnected() {
// 		spotifySession, err := a.ctrl.GetSpotifySession()
// 		if err != nil {
// 			c.AbortWithError(http.StatusInternalServerError, err)
// 			return
// 		}

// 		fullState.SpotifyPlayerStatus = spotifySession.GetPlayer().PlayerStatus
// 	}

// 	c.JSON(http.StatusOK, fullState)
// }

func (a *API) GetParameters(c *gin.Context) {
	parameters, err := a.ds.GetParameters()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, parameters)
}

func (a *API) PatchParameter(c *gin.Context) {
	id := models.UUID(c.Param("id"))

	if !utils.InArray([]models.UUID{
		models.ParamCurrentSpotifyCredential,
	}, id) {
		log.Errorf("Invalid paramter: %s", id)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var patch models.ParameterPatch
	err := c.BindJSON(&patch)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// format value
	switch id {
	case models.ParamCurrentSpotifyCredential:
		patch.Value = models.UUID(patch.Value.(string))
	}

	param, err := a.ctrl.PatchParameter(id, patch)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, param)
}
