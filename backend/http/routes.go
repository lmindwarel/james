package http

import (
	"github.com/gin-gonic/gin"
)

func (a *API) setupRoutes(e *gin.Engine) {

	// public router
	e.GET("/ws", a.wshandler)
	e.GET("/accounts", a.GetAccounts)
	e.POST("/accounts", a.PostAccount)
	spotify := e.Group("/spotify")
	spotify.GET("/credentials", a.GetSpotifyCredentials)
	spotify.POST("/credentials", a.CreateSpotifyCredential)
	spotify.PATCH("/credentials/:id", a.PatchSpotifyCredential)
	spotify.GET("/playlists", a.GetSpotifyPlaylists)
	spotify.GET("/playlists/:id", a.GetSpotifyPlaylist)
	spotify.GET("/playlists/:id/tracks", a.GetSpotifyPlaylistTracks)
	spotify.GET("/tracks/:id", a.GetSpotifyTrack)
	spotify.PUT("/player/play/:id", a.PlaySpotifyTrack)
	spotify.PUT("/player/control", a.ControlSpotifyPlayer)
	spotify.POST("/player/queue/:trackID", a.AddToPlayerQueue)
	spotify.GET("/player/queue", a.GetPlayerQueue)
	spotify.DELETE("/player/queue/:trackID", a.DeleteTrackFromPlayerQueue)

	// authenticated only
	// Authorization group
	authenticated := e.Group("/")
	// per group middleware! in this case we use the custom created
	// AuthRequired() middleware just in the "authorized" group.
	authenticated.Use(a.AuthenticatedMiddleware())
}
