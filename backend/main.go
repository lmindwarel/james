package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/lmindwarel/james/backend/controller"
	"github.com/lmindwarel/james/backend/datastore"
	"github.com/lmindwarel/james/backend/http"
	"github.com/lmindwarel/james/backend/models"
	"github.com/lmindwarel/james/backend/utils"
)

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

	fmt.Printf("Reading config file...")
	configFile, err := os.Open(configFileName)
	if err != nil {
		panic(err)
	}

	var config Config
	parser := json.NewDecoder(configFile)
	if err = parser.Decode(&config); err != nil {
		panic(err)
	}
	fmt.Printf("ok\n")

	fmt.Printf("Initialize logger...")
	utils.InitLogger(config.LogPath)
	fmt.Printf("ok\n")

	fmt.Printf("Initialize datastore...")
	ds, err := datastore.New(config.Datastore)
	if err != nil {
		panic(err)
	}
	fmt.Printf("ok\n")

	fmt.Printf("Initialize controller...")
	ctrl := controller.New(ds, config.Config)
	fmt.Printf("ok\n")

	// fmt.Printf("Connecting spotify account...")
	// err = ctrl.ConnectSpotifyAccount("4153cca5cb4544ad8973eb94a7de36e1", os.Getenv("SPOTIFY_SECRET"))
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("ok\n")

	// fmt.Printf("Get liked titles...")
	// err = ctrl.GetLikedTitles()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("ok\n")

	fmt.Printf("Initialize api...")
	a := http.New(config.API, ctrl, ds)
	go func() {
		err = a.StartServer()
		if err != nil {
			panic(err)
		}
	}()
	fmt.Printf("ok\n")

	pass, err := utils.Encrypt(models.SpotifyPasswordHashKey, config.SpotifyPassword)

	err = ctrl.AuthenticateSpotify(models.SpotifyCredentials{
		User:           config.SpotifyUser,
		HashedPassword: pass,
	})
	if err != nil {
		panic(err)
	}

	utils.Standby()
}
