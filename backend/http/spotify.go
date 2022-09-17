package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lmindwarel/james/backend/models"
	"github.com/lmindwarel/james/backend/spotify"
	"github.com/pkg/errors"
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

func (a *API) GetSpotifyTrack(c *gin.Context) {
	spotifySession, err := a.ctrl.GetSpotifySession()
	if err != nil {
		c.AbortWithError(http.StatusNetworkAuthenticationRequired, err)
	}

	track, err := spotifySession.GetTrack(c.Request.Context(), spotify.ID(c.Param("id")))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, track)
}

func (a *API) ControlSpotifyPlayer(c *gin.Context) {
	var control models.SpotifyPlayerControl
	err := c.BindJSON(&control)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.Wrap(err, "failed to parse json"))
		return
	}

	spotifySession, err := a.ctrl.GetSpotifySession()
	if err != nil {
		c.AbortWithError(http.StatusNetworkAuthenticationRequired, err)
	}

	log.Debugf("control: %v", control)

	if control.Pause != nil {
		if len(spotifySession.GetPlayer().GetQueue()) == 0 {
			c.AbortWithError(http.StatusNotAcceptable, errors.New("no current track"))
			return
		}
		if *control.Pause {
			spotifySession.Pause()
		} else {
			spotifySession.Resume()
		}
	}

	if control.TrackPositionMs != nil {
		if len(spotifySession.GetPlayer().GetQueue()) == 0 {
			c.AbortWithError(http.StatusNotAcceptable, errors.New("no current track"))
			return
		}

		err = spotifySession.SetTrackPosition(time.Duration(*control.TrackPositionMs) * time.Millisecond)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	c.JSON(http.StatusOK, spotifySession.GetPlayer().GetStatus())
}

func (a *API) AddToPlayerQueue(c *gin.Context) {
	spotifySession, err := a.ctrl.GetSpotifySession()
	if err != nil {
		c.AbortWithError(http.StatusNetworkAuthenticationRequired, err)
	}

	spotifySession.AddTrackToQueue(spotify.ID(c.Param("trackID")))

	c.Status(http.StatusOK)
}

func (a *API) GetPlayerQueue(c *gin.Context) {
	spotifySession, err := a.ctrl.GetSpotifySession()
	if err != nil {
		c.AbortWithError(http.StatusNetworkAuthenticationRequired, err)
	}

	queuedTracks, err := spotifySession.GetPlayerQueue(c.Request.Context())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, queuedTracks)
}

func (a *API) DeleteTrackFromPlayerQueue(c *gin.Context) {
	spotifySession, err := a.ctrl.GetSpotifySession()
	if err != nil {
		c.AbortWithError(http.StatusNetworkAuthenticationRequired, err)
	}

	spotifySession.RemoveTrackFromQueue(spotify.ID(c.Param("trackID")))
}
