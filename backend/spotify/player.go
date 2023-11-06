package spotify

import (
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
	"github.com/librespot-org/librespot-golang/Spotify"
	"github.com/librespot-org/librespot-golang/librespot/utils"
	"github.com/lmindwarel/james/backend/models"
	jamesUtils "github.com/lmindwarel/james/backend/utils"
	"github.com/pkg/errors"
	"github.com/xlab/portaudio-go/portaudio"
)

const (
	samplesPerChannel = 2048
	// The samples bit depth
	bitDepth = 16
	// The samples format
	sampleFormat = portaudio.PaFloat32
)

func (s *Session) PlayTrack(id ID) error {
	speaker.Clear()

	s.player.QueueTrackAtIndex(s.player.CurrentQueueIndex+1, QueuedTrack{TrackID: id, ManuallyAdded: true})

	if s.listeners.OnQueueChange != nil {
		s.listeners.OnQueueChange(s.player.queue)
	}

	return s.PlayNextQueuedTrack()
}

func (s *Session) PlayNextQueuedTrack() error {
	if len(s.player.queue) <= s.player.CurrentQueueIndex+1 {
		return errors.New("no queued tracks")
	}

	s.player.CurrentQueueIndex += 1

	queuedTrack := s.player.queue[s.player.CurrentQueueIndex]

	// Get the track metadata: it holds information about which files and encodings are available
	track, err := s.librespotSession.Mercury().GetTrack(utils.Base62ToHex(string(queuedTrack.TrackID)))
	if err != nil {
		return errors.Wrap(err, "failed to load track")
	}

	// Updating player status
	s.player.TrackPosition = 0
	s.player.TrackDuration = models.DurationMs(time.Duration(*track.Duration) * time.Millisecond)
	s.player.State = PlayerState(PlayerStateLoading)
	if s.listeners.OnPlayerStatusChange != nil {
		log.Debugf("Change player statusTrack: %s", track.GetName())
		s.listeners.OnPlayerStatusChange(s.player.PlayerStatus)
	}

	log.Debugf("Track: %s", track.GetName())
	availableFiles := track.GetFile()
	if len(availableFiles) == 0 {
		for _, a := range track.GetAlternative() {
			availableFiles = append(availableFiles, a.GetFile()...)
		}
	}

	if len(availableFiles) == 0 {
		return errors.New("No file in track")
	}

	// As a demo, select the OGG 160kbps variant of the track. The "high quality" setting in the official Spotify
	// app is the OGG 320kbps variant.
	var selectedFile *Spotify.AudioFile
	for _, file := range availableFiles {
		if jamesUtils.InArray([]Spotify.AudioFile_Format{Spotify.AudioFile_OGG_VORBIS_320, Spotify.AudioFile_OGG_VORBIS_160, Spotify.AudioFile_OGG_VORBIS_96}, file.GetFormat()) {
			selectedFile = file
		} else {
			log.Errorf("Unsupported file format: %s", file.GetFormat())
		}
	}

	if selectedFile == nil {
		return errors.New("Unsupported track format: %s")
	}

	// Synchronously load the track
	audioFile, err := s.librespotSession.Player().LoadTrack(selectedFile, track.GetGid())
	if err != nil {
		return errors.Wrap(err, "failed to load track")
	}

	time.Sleep(200 * time.Millisecond) // to recover header and audio file length before

	// We have the track audio, let's play it! Initialize the OGG decoder, and start a PortAudio stream.
	// Note that we skip the first 167 bytes as it is a Spotify-specific header. You can decode it by
	// using this: https://sourceforge.net/p/despotify/code/HEAD/tree/java/trunk/src/main/java/se/despotify/client/player/SpotifyOggHeader.java
	log.Debug("Decoding stream...")
	// fmt.Printf("PortAudio channels: %d / SampleRate: %f\n", info.Channels, info.SampleRate)

	streamer, format, err := vorbis.Decode(audioFile)
	if err != nil {
		return errors.Wrap(err, "failed to decode audio file")
	}

	log.Debug("Setting up  stream...")
	s.player.streamer = streamer
	s.player.ctrl = &beep.Ctrl{Streamer: beep.Loop(1, streamer)}
	s.player.resampler = beep.ResampleRatio(4, 1, s.player.ctrl)
	s.player.volume = &effects.Volume{Streamer: s.player.resampler, Base: 2}
	s.player.sampleRate = format.SampleRate

	// Updating player status
	s.player.TrackPosition = 0
	s.player.TrackDuration = models.DurationMs(time.Duration(*track.Duration) * time.Millisecond)
	s.player.State = PlayerState(PlayerStatePlaying)
	if s.listeners.OnPlayerStatusChange != nil {
		s.listeners.OnPlayerStatusChange(s.player.PlayerStatus)
	}

	log.Debug("Initializing speaker...")
	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		return errors.Wrap(err, "failed to initialize speaker")
	}

	log.Debug("Playing track")

	speaker.Play(s.player.volume)
	s.ensureTickerStarted()

	log.Debug("Done!")

	return nil
}

func (s *Session) Pause() {
	s.player.ctrl.Paused = true
	s.player.State = PlayerState(PlayerStatePaused)
	if s.listeners.OnPlayerStatusChange != nil {
		s.listeners.OnPlayerStatusChange(s.player.PlayerStatus)
	}
}

func (s *Session) Resume() {
	s.player.ctrl.Paused = false
	s.player.State = PlayerState(PlayerStatePlaying)
	s.listeners.OnPlayerStatusChange(s.player.PlayerStatus)
	s.ensureTickerStarted()
}

func (s *Session) ensureTickerStarted() {
	if s.ticking {
		return
	}

	s.ticking = true

	go func() {
		for {
			if !s.ticking {
				return
			}

			s.player.TrackPosition = models.DurationMs(s.player.sampleRate.D(s.player.streamer.Position()))
			s.player.Volume = int(s.player.volume.Volume)
			// log.Debugf("track position: %dms, duration: %dms", time.Duration(s.player.TrackPosition).Milliseconds(), time.Duration(s.player.TrackDuration).Milliseconds())
			if s.player.TrackPosition >= s.player.TrackDuration {
				s.onCurrentTrackEnd()
			}

			if s.listeners.OnPlayerStatusChange != nil {
				s.listeners.OnPlayerStatusChange(s.player.PlayerStatus)
			}

			if s.player.State != PlayerStatePlaying {
				return // stop ticking
			}

			time.Sleep(time.Second)
		}
	}()
}

func (s *Session) onCurrentTrackEnd() {
	if len(s.player.queue) > s.player.CurrentQueueIndex {
		// play next track in queue
		err := s.PlayNextQueuedTrack()
		if err != nil {
			log.Errorf("failed to play next track: %s", err)
		}
	} else {
		s.player.State = PlayerStatePaused
	}
}

func (s *Session) SetTrackPosition(pos time.Duration) error {
	speaker.Lock()
	defer speaker.Unlock()

	log.Debugf("SetTrackPosition: %fs", pos.Seconds())

	newPos := s.player.sampleRate.N(pos)
	log.Debugf("to sample: %d, duration: %d", newPos, s.player.streamer.Len())

	if newPos < 0 {
		newPos = 0
	}

	// trackDurationSamplesLen := s.player.streamer.Len()
	// log.Debugf("max track duration: %fs", s.player.sampleRate.D(trackDurationSamplesLen).Seconds())
	// if newPos >= trackDurationSamplesLen {
	// 	return fmt.Errorf("want to set track position to %fs but track duration is %fs", pos.Seconds(), s.player.sampleRate.D(trackDurationSamplesLen).Seconds())
	// }

	log.Debugf("SetTrackPosition to %f seconds (%d)", s.player.sampleRate.D(newPos).Seconds(), newPos)

	err := s.player.streamer.Seek(newPos)
	if err != nil {
		return errors.Wrap(err, "failed to seek to position")
	}

	s.player.TrackPosition = models.DurationMs(pos)

	return nil
}

func (s *Session) SetVolume(vol float64) error {
	s.player.volume.Volume = vol
	s.player.PlayerStatus.Volume = int(vol)
	return nil
}
