package spotify

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
	if s.listeners.OnQueueChange != nil {
		s.listeners.OnQueueChange(s.player.queue)
	}
}

func (s *Session) RemoveTrackFromQueue(trackID ID) {
	// TODO
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
