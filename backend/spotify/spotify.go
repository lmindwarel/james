package spotify

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/librespot-org/librespot-golang/librespot/core"
	"github.com/lmindwarel/james/backend/models"
	"github.com/lmindwarel/james/backend/utils"
	"github.com/zmb3/spotify/v2"
)

var log = utils.GetLogger("james-spotify")

type ID spotify.ID

type PlayerState string

const (
	PlayerStatePlaying ID = "playing"
	PlayerStatePaused  ID = "paused"
	PlayerStateStopped ID = "stopped"
)

type PlayerStatus struct {
	State          PlayerState       `json:"state"`
	CurrentTrackID *ID               `json:"current_track_id"`
	TrackDuration  models.DurationMs `json:"track_duration"`
	TrackPosition  models.DurationMs `json:"track_position"`
}

type Player struct {
	PlayerStatus
	sampleRate beep.SampleRate
	streamer   beep.StreamSeeker
	ctrl       *beep.Ctrl
	resampler  *beep.Resampler
	volume     *effects.Volume
}

type Listeners struct {
	OnPlayerStatusChange func(s PlayerStatus)
}

type Session struct {
	userID           string
	librespotSession *core.Session
	webapiClient     *spotify.Client
	player           Player
	listeners        Listeners
}

func (s *Session) ListenOnPlayerStatusChange(listener func(s PlayerStatus)) {
	s.listeners.OnPlayerStatusChange = listener
}

func (s *Session) GetPlayerStatus() PlayerStatus {
	return s.player.PlayerStatus
}
