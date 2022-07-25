package controller

import (
	"github.com/lmindwarel/james/backend/datastore"
	"github.com/lmindwarel/james/backend/utils"
	"github.com/zmb3/spotify/v2"
)

var log = utils.GetLogger("james-controller")

// Controller is the struct for main project controller
type Controller struct {
	ds            *datastore.Datastore
	SpotifyClient *spotify.Client
}

// New create new controller with datastore
func New(ds *datastore.Datastore) *Controller {
	return &Controller{
		ds: ds,
	}
}
