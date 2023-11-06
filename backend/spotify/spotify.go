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
	PlayerStatePlaying PlayerState = "playing"
	PlayerStatePaused  PlayerState = "paused"
	PlayerStateStopped PlayerState = "stopped"
	PlayerStateLoading PlayerState = "loading"
)

type PlayerStatus struct {
	State             PlayerState       `json:"state"`
	CurrentQueueIndex int               `json:"current_queue_index"`
	TrackDuration     models.DurationMs `json:"track_duration"`
	TrackPosition     models.DurationMs `json:"track_position"`
	Volume            int               `json:"volume"` // in % (0-100)
}

type Player struct {
	PlayerStatus
	sampleRate beep.SampleRate
	streamer   beep.StreamSeeker
	ctrl       *beep.Ctrl
	resampler  *beep.Resampler
	volume     *effects.Volume
	queue      []QueuedTrack // sorted array
}

type QueuedTrack struct {
	TrackID       ID   `json:"track_id"`
	ManuallyAdded bool `json:"manually_added"`
}

type Listeners struct {
	OnPlayerStatusChange func(s PlayerStatus)
	OnQueueChange        func(queue []QueuedTrack)
}

type Session struct {
	userID           string
	librespotSession *core.Session
	webapiClient     *spotify.Client
	player           Player
	listeners        Listeners
	ticking          bool
}

func (s *Session) ListenOnPlayerStatusChange(listener func(s PlayerStatus)) {
	s.listeners.OnPlayerStatusChange = listener
}

func (s *Session) ListenOnPlayerQueueChange(listener func(queue []QueuedTrack)) {
	s.listeners.OnQueueChange = listener
}

func (s *Session) GetPlayer() *Player {
	return &s.player
}

func (p *Player) GetQueue() []QueuedTrack {
	return p.queue
}
