package spotify

import (
	"github.com/librespot-org/librespot-golang/librespot/core"
	"github.com/lmindwarel/james/backend/utils"
	"github.com/zmb3/spotify/v2"
)

var log = utils.GetLogger("james-spotify")

type ID spotify.ID

type Session struct {
	userID           string
	librespotSession *core.Session
	webapiClient     *spotify.Client
}
