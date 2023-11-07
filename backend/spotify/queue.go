package spotify

import "github.com/lmindwarel/james/backend/models"

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

	s.player.QueueTrackAtIndex(insertAtIndex, QueuedTrack{ID: models.NewUUID(), TrackID: trackID, ManuallyAdded: true})
	if s.listeners.OnQueueChange != nil {
		s.listeners.OnQueueChange(s.player.queue)
	}
}

func (s *Session) RemoveTrackFromQueue(id models.UUID) {
	index := -1
	for queueIndex, queuedTrack := range s.player.queue {
		if queuedTrack.ID == id {
			index = queueIndex
			break
		}
	}

	if index == -1 {
		return
	}

	s.player.queue = append(s.player.queue[:index], s.player.queue[index+1:]...)
	if s.listeners.OnQueueChange != nil {
		s.listeners.OnQueueChange(s.player.queue)
	}
}

func (p *Player) QueueTrackAtIndex(index int, track QueuedTrack) {
	if index > len(p.queue)-1 || index < 0 {
		p.queue = append(p.queue, track)
		return
	}
	nextQueue := p.queue[index:]
	p.queue = p.queue[:index]
	p.queue = append(p.queue, track)
	p.queue = append(p.queue, nextQueue...)
}
