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
	spotify.GET("/playlists", a.GetSpotifyPlaylists)
	spotify.GET("/playlists/:id", a.GetSpotifyPlaylist)
	spotify.GET("/playlists/:id/tracks", a.GetSpotifyPlaylistTracks)
	spotify.GET("/tracks/:id", a.GetSpotifyTrack)
	spotify.PUT("/player/play/:id", a.PlaySpotifyTrack)
	spotify.PUT("/player/control", a.ControlSpotifyPlayer)

	// authenticated only
	// Authorization group
	authenticated := e.Group("/")
	// per group middleware! in this case we use the custom created
	// AuthRequired() middleware just in the "authorized" group.
	authenticated.Use(a.AuthenticatedMiddleware())
}
