package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lmindwarel/james/backend/models"
)

func (a *API) GetSpotifyPlaylists(c *gin.Context) {
	playlists, err := a.ctrl.GetSpotifyPlaylists()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, playlists)
}

func (a *API) GetSpotifyPlaylistTracks(c *gin.Context) {
	playlists, err := a.ctrl.GetSpotifyTracks(models.SpotifyURI(c.Param("uri")))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, playlists)
}
