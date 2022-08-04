package controller

import (
	"github.com/librespot-org/librespot-golang/librespot/core"
	"github.com/lmindwarel/james/backend/datastore"
	"github.com/lmindwarel/james/backend/utils"
)

var log = utils.GetLogger("james-controller")

type Config struct {
}

// Controller is the struct for main project controller
type Controller struct {
	ds             *datastore.Datastore
	config         Config
	SpotifySession *core.Session
}

// New create new controller with datastore
func New(ds *datastore.Datastore, config Config) *Controller {
	return &Controller{
		ds:     ds,
		config: config,
	}
}
