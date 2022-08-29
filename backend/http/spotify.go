package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lmindwarel/james/backend/models"
	"github.com/lmindwarel/james/backend/spotify"
)

func (a *API) GetSpotifyPlaylists(c *gin.Context) {
	spotifySession, err := a.ctrl.GetSpotifySession()
	if err != nil {
		c.AbortWithError(http.StatusNetworkAuthenticationRequired, err)
	}

	playlistsResult, err := spotifySession.GetCurrentUserPlaylists(c.Request.Context())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, playlistsResult)
}

func (a *API) GetSpotifyPlaylist(c *gin.Context) {
	spotifySession, err := a.ctrl.GetSpotifySession()
	if err != nil {
		c.AbortWithError(http.StatusNetworkAuthenticationRequired, err)
	}

	playlistResult, err := spotifySession.GetPlaylist(c.Request.Context(), spotify.ID(models.SpotifyURI(c.Param("id"))))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, playlistResult)
}

func (a *API) GetSpotifyPlaylistTracks(c *gin.Context) {
	spotifySession, err := a.ctrl.GetSpotifySession()
	if err != nil {
		c.AbortWithError(http.StatusNetworkAuthenticationRequired, err)
	}

	tracksResult, err := spotifySession.GetPlaylistTracks(c.Request.Context(), spotify.ID(models.SpotifyURI(c.Param("id"))))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, tracksResult)
}

func (a *API) PlaySpotifyTrack(c *gin.Context) {
	spotifySession, err := a.ctrl.GetSpotifySession()
	if err != nil {
		c.AbortWithError(http.StatusNetworkAuthenticationRequired, err)
	}

	err = spotifySession.PlayTrack(spotify.ID(c.Param("id")))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
