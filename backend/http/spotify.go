package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lmindwarel/james/backend/models"
	"github.com/lmindwarel/james/backend/spotify"
	"github.com/pkg/errors"
)

func (a *API) SpotifyConnectedMiddleware(c *gin.Context) {
	if !a.ctrl.IsSpotifyConnected() {
		c.AbortWithStatus(http.StatusPreconditionRequired)
		return
	}
}

func (a *API) GetSpotifyCredentials(c *gin.Context) {
	credentials, err := a.ds.GetSpotifyCredentials()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, credentials)
}

func (a *API) CreateSpotifyCredential(c *gin.Context) {
	var credential models.CredentialPatch
	err := c.ShouldBindJSON(&credential)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	spotifyCredential, err := a.ctrl.CreateSpotifyCredential(credential)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, spotifyCredential)
}

func (a *API) PatchSpotifyCredential(c *gin.Context) {
	var patch models.CredentialPatch
	err := c.ShouldBindJSON(&patch)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	credential, err := a.ds.GetSpotifyCredential(models.UUID(c.Param("id")))
	if err != nil {
		if a.ds.IsNotFoundError(err) {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	spotifyCredential, err := a.ctrl.PatchSpotifyCredential(credential, patch)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, spotifyCredential)
}

func (a *API) AuthenticateSpotifyCredential(c *gin.Context) {
	credential, err := a.ds.GetSpotifyCredential(models.UUID(c.Param("id")))
	if err != nil {
		if a.ds.IsNotFoundError(err) {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	err = a.ctrl.AuthenticateSpotify(credential)
	if err != nil {
		if errors.Is(err, models.ErrAuthenticationFailed) {
			c.AbortWithError(http.StatusForbidden, err)
			return

		}
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

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

func (a *API) GetSpotifySavedTracks(c *gin.Context) {
	spotifySession, err := a.ctrl.GetSpotifySession()
	if err != nil {
		c.AbortWithError(http.StatusNetworkAuthenticationRequired, err)
	}

	tracksResult, err := spotifySession.GetCurrentUserSavedTracks(c.Request.Context())
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
		return
	}

	track, err := spotifySession.GetTrack(c.Request.Context(), spotify.ID(c.Param("id")))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, track)
}

func (a *API) GetSpotifyTracks(c *gin.Context) {
	spotifySession, err := a.ctrl.GetSpotifySession()
	if err != nil {
		c.AbortWithError(http.StatusNetworkAuthenticationRequired, err)
		return
	}

	track, err := spotifySession.GetTracks(c.Request.Context(), spotify.StringSliceToIDs(c.QueryArray("ids[]")))
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

	if control.Volume != nil {
		err = spotifySession.SetVolume(*control.Volume)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	c.JSON(http.StatusOK, spotifySession.GetPlayer().PlayerStatus)
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
	// spotifySession, err := a.ctrl.GetSpotifySession()
	// if err != nil {
	// 	c.AbortWithError(http.StatusNetworkAuthenticationRequired, err)
	// }

	// // queuedTracks, err := spotifySession.GetPlayerQueue(c.Request.Context())
	// if err != nil {
	// 	c.AbortWithError(http.StatusInternalServerError, err)
	// 	return
	// }

	// c.JSON(http.StatusOK, queuedTracks)
}

func (a *API) DeleteTrackFromPlayerQueue(c *gin.Context) {
	spotifySession, err := a.ctrl.GetSpotifySession()
	if err != nil {
		c.AbortWithError(http.StatusNetworkAuthenticationRequired, err)
	}

	spotifySession.RemoveTrackFromQueue(models.UUID(c.Param("id")))
}
