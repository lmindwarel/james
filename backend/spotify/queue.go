package spotify

import (
	"context"

	"github.com/lmindwarel/james/backend/models"
	"github.com/pkg/errors"
	"github.com/zmb3/spotify/v2"
)

func (s *Session) AddTracksToQueue(tracksIDs []ID, manual bool) {
	for _, tID := range tracksIDs {
		s.player.queue = append(s.player.queue, QueuedTrack{TrackID: tID, ManuallyAdded: manual})
	}
}

func (s *Session) RemoveTracksFromQueue(tracksIDs []ID) {

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
