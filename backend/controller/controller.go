package controller

import (
	"github.com/lmindwarel/james/backend/datastore"
	"github.com/lmindwarel/james/backend/models"
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
	listeners      map[string][]func(data interface{})
	state          models.JamesStatus
}

// New create new controller with datastore
func New(ds *datastore.Datastore, config Config) *Controller {
	return &Controller{
		ds:        ds,
		config:    config,
		listeners: map[string][]func(data interface{}){},
	}
}

func (ctrl *Controller) GetJamesStatus() models.JamesStatus {
	return ctrl.state
}

func (ctrl *Controller) updateJamesStatus(patch models.JamesStatusPatch) (err error) {
	hasChanged := false
	if patch.AuthenticatedSpotifyCredentialID != nil {
		ctrl.state.AuthenticatedSpotifyCredentialID = *patch.AuthenticatedSpotifyCredentialID
		hasChanged = true
	}

	if hasChanged {
		ctrl.triggerListeners(EventJamesStatusChange, patch)
	} else {
		log.Warningf("Nothing has changed")
	}

	return
}
