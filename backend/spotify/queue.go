package spotify

import (
	"context"

	"github.com/lmindwarel/james/backend/models"
	"github.com/pkg/errors"
	"github.com/zmb3/spotify/v2"
)

// AddTrackToQueue will add the track to the manual queue, moving all automatically queued tracks
func (s *Session) AddTrackToQueue(trackID ID) {
	insertAtIndex := -1
	for queueIndex, queuedTracks := range s.player.queue {
		if queuedTracks.ManuallyAdded {
			continue
		}

		// Here we are after all manually queued tracks. Add the track here and move others
		insertAtIndex = queueIndex
		break
	}

	s.player.QueueTrackAtIndex(insertAtIndex, QueuedTrack{TrackID: trackID, ManuallyAdded: true})
	s.listeners.OnQueueChange(s.player.queue)
}

func (s *Session) RemoveTrackFromQueue(trackID ID) {
	// TODO
}

func (s *Session) GetPlayerQueue(ctx context.Context) (queuedTracks []models.SpotifyQueuedTrack, err error) {
	var ids []spotify.ID
	queuedTracksByID := map[spotify.ID]QueuedTrack{}
	for _, queuedTrack := range s.player.queue {
		ids = append(ids, spotify.ID(queuedTrack.TrackID))
		queuedTracksByID[spotify.ID(queuedTrack.TrackID)] = queuedTrack
	}

	tracks, err := s.webapiClient.GetTracks(ctx, ids)
	if err != nil {
		return queuedTracks, errors.Wrap(err, "failed to get tracks")
	}

	for _, t := range tracks {
		queuedTracks = append(queuedTracks, models.SpotifyQueuedTrack{
			Track:         *t,
			ManuallyAdded: queuedTracksByID[t.ID].ManuallyAdded,
		})
	}

	return
}

func (p *Player) QueueTrackAtIndex(index int, track QueuedTrack) {
	if len(p.queue) < index {
		p.queue = append(p.queue, track)
		return
	}
	nextQueue := p.queue[index:]
	p.queue = p.queue[:index]
	p.queue = append(p.queue, track)
	p.queue = append(p.queue, nextQueue...)
}
