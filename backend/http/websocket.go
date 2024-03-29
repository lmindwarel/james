package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/lmindwarel/james/backend/controller"
	"github.com/lmindwarel/james/backend/spotify"
	"github.com/pkg/errors"
)

type InWebsocketMessage struct {
	Topic string `json:"topic"`
	Data  []byte `json:"data"`
}

type OutWebsocketMessage struct {
	Topic string      `json:"topic"`
	Data  interface{} `json:"data"`
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (a *API) wshandler(c *gin.Context) {
	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Errorf("Failed to set upgrade request to websocket: %+v", err)
		return
	}

	a.websocketConnection = conn

	subscribedEvents := []string{
		controller.EventJamesStatusChange,
	}

	for _, e := range subscribedEvents {
		a.ctrl.AddEventListener(e, func(data interface{}) {
			err = conn.WriteJSON(OutWebsocketMessage{
				Topic: e,
				Data:  data,
			})
			if err != nil {
				log.Errorf("Failed to write websocket event: %s", err)
			}
		})
	}

	// spotify events
	spotifySession, notFatalErr := a.ctrl.GetSpotifySession()
	if notFatalErr == nil {
		spotifySession.ListenOnPlayerStatusChange(func(s spotify.PlayerStatus) {
			err := a.SendWebsocketMessage(OutWebsocketMessage{
				Topic: "player-status",
				Data:  s,
			})
			if err != nil {
				log.Errorf("Failed to send websocket message: %s", err)
			}
		})

		spotifySession.ListenOnPlayerQueueChange(func(q []spotify.QueuedTrack) {
			// TODO only send from -10 in queue position
			err := a.SendWebsocketMessage(OutWebsocketMessage{
				Topic: "player-queue",
				Data:  q,
			})
			if err != nil {
				log.Errorf("Failed to send websocket message: %s", err)
			}
		})
	} else {
		log.Errorf("Failed to get spotify session: %s", notFatalErr)
	}

	ticker := time.NewTicker(100 * time.Millisecond)

	for {
		select {
		case <-ticker.C:
			var message InWebsocketMessage
			err := conn.ReadJSON(&message)
			if err != nil {
				return
			}

			a.HandleWebsocketMessage(message)
		}
	}
}

func (a *API) SendWebsocketMessage(message OutWebsocketMessage) error {
	if a.websocketConnection == nil {
		return errors.New("websocket connection closed")
	}

	err := a.websocketConnection.WriteJSON(message)
	if err != nil {
		return errors.Wrap(err, "failed to write websocket message")
	}

	return nil
}

func (a *API) HandleWebsocketMessage(msg InWebsocketMessage) {
	switch msg.Topic {
	default:
		log.Warning("Unhandled websocket message: %v", msg.Topic)
	}
}
