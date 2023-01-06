package main

import (
	"encoding/json"
	"os"

	"github.com/lmindwarel/james/backend/controller"
	"github.com/lmindwarel/james/backend/datastore"
	"github.com/lmindwarel/james/backend/http"
	"github.com/lmindwarel/james/backend/spotify"
	"github.com/lmindwarel/james/backend/utils"
	"github.com/xlab/portaudio-go/portaudio"
)

var log = utils.GetLogger("james")

// Config is the core configuration
type Config struct {
	controller.Config
	LogPath         string           `json:"logPath"`
	Datastore       datastore.Config `json:"datastore"`
	API             http.Config      `json:"http"`
	SpotifyUser     string           `json:"spotifyUser"`
	SpotifyPassword string           `json:"spotifyPassword"`
}

func main() {
	var err error

	configFileName := "config.json"
	if len(os.Args) > 1 {
		configFileName = os.Args[1]
	}

	_, err = os.Stat(configFileName)
	if os.IsNotExist(err) {
		panic("please provide config.json or give the path in arg")
	}

	log.Noticef("Loading config file...")
	configFile, err := os.Open(configFileName)
	if err != nil {
		panic(err)
	}

	var config Config
	parser := json.NewDecoder(configFile)
	if err = parser.Decode(&config); err != nil {
		panic(err)
	}

	configJSON, _ := json.MarshalIndent(config, "    ", "    ")
	log.Debugf("Loaded config: %s", configJSON)

	log.Noticef("Initialize logger...")
	utils.InitLogger(config.LogPath)
	log.Noticef("Logger initialized to %s", config.LogPath)

	log.Noticef("Initialize datastore...")
	ds, err := datastore.New(config.Datastore)
	if err != nil {
		panic(err)
	}
	log.Noticef("Datastore initialized successfully")

	log.Noticef("Initialize controller...")
	ctrl := controller.New(ds, config.Config)
	log.Noticef("Contorller initialized")

	// log.Noticef("Connecting spotify account...")
	// err = ctrl.ConnectSpotifyAccount("4153cca5cb4544ad8973eb94a7de36e1", os.Getenv("SPOTIFY_SECRET"))
	// if err != nil {
	// 	panic(err)
	// }
	// log.Noticef("ok\n")

	// log.Noticef("Get liked titles...")
	// err = ctrl.GetLikedTitles()
	// if err != nil {
	// 	panic(err)
	// }
	// log.Noticef("ok\n")

	log.Noticef("Initialize api...")
	a := http.New(config.API, ctrl, ds)
	go func() {
		err = a.StartServer()
		if err != nil {
			panic(err)
		}
	}()
	log.Noticef("API initialized")

	log.Noticef("Initialize port audio...")
	if err := portaudio.Initialize(); spotify.PAError(err) {
		panic("PortAudio init error: " + spotify.PAErrorText(err))
	}
	log.Noticef("Port audio initialized")

	notFatalErr := ctrl.AuthenticateCurrentSpotifyCredential()
	if notFatalErr != nil {
		log.Warningf("Failed to authenticate with current spotify credential: %s", notFatalErr)
	}

	log.Noticef("Every services initialized, at your service.")

	utils.Standby()
}
