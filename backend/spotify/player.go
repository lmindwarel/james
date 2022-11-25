package spotify

import (
	"fmt"
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

// func (s *Session) PlayTrack(id ID) error {
// 	// Get the track metadata: it holds information about which files and encodings are available
// 	track, err := s.librespotSession.Mercury().GetTrack(utils.Base62ToHex(string(id)))
// 	if err != nil {
// 		return errors.Wrap(err, "failed to load track")
// 	}

// 	log.Debugf("Track: %s", track.GetName())

// 	// As a demo, select the OGG 160kbps variant of the track. The "high quality" setting in the official Spotify
// 	// app is the OGG 320kbps variant.
// 	var selectedFile *Spotify.AudioFile
// 	for _, file := range track.GetFile() {
// 		if file.GetFormat() == Spotify.AudioFile_OGG_VORBIS_160 {
// 			selectedFile = file
// 		}
// 	}

// 	// Synchronously load the track
// 	audioFile, err := s.librespotSession.Player().LoadTrack(selectedFile, track.GetGid())

// 	// TODO: channel to be notified of chunks downloaded (or reader?)

// 	if err != nil {
// 		fmt.Printf("Error while loading track: %s\n", err)
// 	} else {
// 		// We have the track audio, let's play it! Initialize the OGG decoder, and start a PortAudio stream.
// 		// Note that we skip the first 167 bytes as it is a Spotify-specific header. You can decode it by
// 		// using this: https://sourceforge.net/p/despotify/code/HEAD/tree/java/trunk/src/main/java/se/despotify/client/player/SpotifyOggHeader.java
// 		fmt.Println("Setting up OGG decoder...")
// 		dec, err := decoder.New(audioFile, samplesPerChannel)
// 		if err != nil {
// 			return errors.Wrap(err, "failed to load decoder")
// 		}

// 		info := dec.Info()

// 		go func() {
// 			dec.Decode()
// 			dec.Close()
// 		}()

// 		fmt.Println("Setting up PortAudio stream...")
// 		fmt.Printf("PortAudio channels: %d / SampleRate: %f\n", info.Channels, info.SampleRate)

// 		var wg sync.WaitGroup
// 		var stream *portaudio.Stream
// 		callback := paCallback(&wg, int(info.Channels), dec.SamplesOut())

// 		defaultDevice := portaudio.GetDefaultInputDevice()
// 		log.Debugf("device count: %d", portaudio.GetDeviceCount())
// 		deviceInfos := portaudio.GetDeviceInfo(defaultDevice)
// 		if deviceInfos == nil {
// 			return errors.New("could not find default input device")
// 		}
// 		log.Debugf("Default device: %s, is free %t", deviceInfos.Name)

// 		if err := portaudio.OpenDefaultStream(&stream, 0, info.Channels, sampleFormat, info.SampleRate,
// 			samplesPerChannel, callback, nil); PAError(err) {
// 			return fmt.Errorf("failed to open PortAudio stream: %s", PAErrorText(err))
// 		}

// 		fmt.Println("Starting playback...")
// 		if err := portaudio.StartStream(stream); PAError(err) {
// 			return fmt.Errorf("failed to start PortAudio stream: %s", PAErrorText(err))
// 		}

// 		wg.Wait()
// 	}

// 	return nil
// }

func (s *Session) PlayTrack(id ID) error {
	s.player.QueueTrackAtIndex(s.player.CurrentQueueIndex+1, QueuedTrack{TrackID: id, ManuallyAdded: true})
	s.listeners.OnQueueChange(s.player.queue)
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

	log.Debugf("Track: %s", track.GetName())

	// As a demo, select the OGG 160kbps variant of the track. The "high quality" setting in the official Spotify
	// app is the OGG 320kbps variant.
	var selectedFile *Spotify.AudioFile
	for _, file := range track.GetFile() {
		if jamesUtils.InArray([]Spotify.AudioFile_Format{Spotify.AudioFile_OGG_VORBIS_320, Spotify.AudioFile_OGG_VORBIS_160, Spotify.AudioFile_OGG_VORBIS_96}, file.GetFormat()) {
			selectedFile = file
		}
	}

	// Synchronously load the track
	if selectedFile == nil {
		return errors.New("Unsupported track format")
	}

	audioFile, err := s.librespotSession.Player().LoadTrack(selectedFile, track.GetGid())
	if err != nil {
		return errors.Wrap(err, "failed to load track")
	}

	// We have the track audio, let's play it! Initialize the OGG decoder, and start a PortAudio stream.
	// Note that we skip the first 167 bytes as it is a Spotify-specific header. You can decode it by
	// using this: https://sourceforge.net/p/despotify/code/HEAD/tree/java/trunk/src/main/java/se/despotify/client/player/SpotifyOggHeader.java
	fmt.Println("Setting up OGG decoder...")
	// dec, err := decoder.New(audioFile, samplesPerChannel)
	// if err != nil {
	// 	return errors.Wrap(err, "failed to load decoder")
	// }

	fmt.Println("Decoding stream...")
	// fmt.Printf("PortAudio channels: %d / SampleRate: %f\n", info.Channels, info.SampleRate)

	streamer, format, err := vorbis.Decode(audioFile)
	if err != nil {
		return errors.Wrap(err, "failed to decode audio file")
	}

	speaker.Lock()
	fmt.Println("Setting up  stream...")
	s.player.streamer = streamer
	s.player.ctrl = &beep.Ctrl{Streamer: beep.Loop(1, streamer)}
	s.player.resampler = beep.ResampleRatio(4, 1, s.player.ctrl)
	s.player.volume = &effects.Volume{Streamer: s.player.resampler, Base: 2}
	s.player.sampleRate = format.SampleRate

	log.Debug("Initializing speaker...")
	// speaker.Clear()
	// log.Debug("Cleared")

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		return errors.Wrap(err, "failed to initialize speaker")
	}

	log.Debug("ok")

	speaker.Play(beep.Seq(s.player.volume, beep.Callback(func() {
		if len(s.player.queue) > s.player.CurrentQueueIndex {
			// play next track in queue
			err := s.PlayNextQueuedTrack()
			if err != nil {
				log.Errorf("failed to play next track: %s", err)
			}
		}
	})))

	s.player.TrackDuration = models.DurationMs(time.Duration(*track.Duration) * time.Millisecond)
	s.player.State = PlayerState(PlayerStatePlaying)
	s.listeners.OnPlayerStatusChange(s.player.PlayerStatus)
	s.StartPlayerListener()
	speaker.Unlock()

	log.Debug("Done!")

	return nil
}

func (s *Session) Pause() {
	speaker.Lock()
	defer speaker.Unlock()
	s.player.ctrl.Paused = true
	s.player.State = PlayerState(PlayerStatePaused)
	s.listeners.OnPlayerStatusChange(s.player.PlayerStatus)
}

func (s *Session) Resume() {
	speaker.Lock()
	defer speaker.Unlock()
	s.player.ctrl.Paused = false
	s.player.State = PlayerState(PlayerStatePlaying)
	s.listeners.OnPlayerStatusChange(s.player.PlayerStatus)
	s.StartPlayerListener()
}

func (s *Session) StartPlayerListener() {
	ticker := time.NewTicker(time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				speaker.Lock()
				s.player.TrackPosition = models.DurationMs(s.player.sampleRate.D(s.player.streamer.Position()))
				s.listeners.OnPlayerStatusChange(s.player.PlayerStatus)
				log.Debugf("track position: %dms, duration: %dms", time.Duration(s.player.TrackPosition).Milliseconds(), time.Duration(s.player.TrackDuration).Milliseconds())
				speaker.Unlock()

				if s.player.State != PlayerState(PlayerStatePlaying) {
					log.Debugf("player state: %s", s.player.State)
					return
				}
			}
		}
	}()
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
	trackDurationSamplesLen := s.player.streamer.Len()
	log.Debugf("max track duration: %fs", s.player.sampleRate.D(trackDurationSamplesLen).Seconds())
	if newPos >= trackDurationSamplesLen {
		newPos = trackDurationSamplesLen - 1
	}

	log.Debugf("SetTrackPosition to %d seconds (%d)", s.player.sampleRate.D(newPos).Seconds(), newPos)

	err := s.player.streamer.Seek(newPos)
	if err != nil {
		return errors.Wrap(err, "failed to seek to position")
	}

	s.player.TrackPosition = models.DurationMs(pos)

	return nil
}
