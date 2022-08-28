package controller

import (
	"errors"

	"github.com/lmindwarel/james/backend/datastore"
	"github.com/lmindwarel/james/backend/spotify"
	"github.com/lmindwarel/james/backend/utils"
)

var log = utils.GetLogger("james-controller")

type Config struct {
	SpotifyClientID     string `json:"spotifyClientID"`
	SpotifyClientSecret string `json:"spotifyClientSecret"`
}

// Controller is the struct for main project controller
type Controller struct {
	ds             *datastore.Datastore
	config         Config
	spotifySession *spotify.Session
}

// New create new controller with datastore
func New(ds *datastore.Datastore, config Config) *Controller {
	return &Controller{
		ds:     ds,
		config: config,
	}
}

func (ctrl *Controller) GetSpotifySession() (*spotify.Session, error) {
	if ctrl.spotifySession == nil {
		return nil, errors.New("spotify not connected")
	}

	return ctrl.spotifySession, nil
}
