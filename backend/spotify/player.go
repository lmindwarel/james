package spotify

import (
	"fmt"
	"io"
	"time"

	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
	"github.com/librespot-org/librespot-golang/Spotify"
	"github.com/librespot-org/librespot-golang/librespot/utils"
	"github.com/pkg/errors"
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
	// Get the track metadata: it holds information about which files and encodings are available
	track, err := s.librespotSession.Mercury().GetTrack(utils.Base62ToHex(string(id)))
	if err != nil {
		return errors.Wrap(err, "failed to load track")
	}

	log.Debugf("Track: %s", track.GetName())

	// As a demo, select the OGG 160kbps variant of the track. The "high quality" setting in the official Spotify
	// app is the OGG 320kbps variant.
	var selectedFile *Spotify.AudioFile
	for _, file := range track.GetFile() {
		if file.GetFormat() == Spotify.AudioFile_OGG_VORBIS_160 {
			selectedFile = file
		}
	}

	// Synchronously load the track
	audioFile, err := s.librespotSession.Player().LoadTrack(selectedFile, track.GetGid())

	// TODO: channel to be notified of chunks downloaded (or reader?)

	if err != nil {
		fmt.Printf("Error while loading track: %s\n", err)
	} else {
		// We have the track audio, let's play it! Initialize the OGG decoder, and start a PortAudio stream.
		// Note that we skip the first 167 bytes as it is a Spotify-specific header. You can decode it by
		// using this: https://sourceforge.net/p/despotify/code/HEAD/tree/java/trunk/src/main/java/se/despotify/client/player/SpotifyOggHeader.java
		fmt.Println("Setting up OGG decoder...")
		// dec, err := decoder.New(audioFile, samplesPerChannel)
		// if err != nil {
		// 	return errors.Wrap(err, "failed to load decoder")
		// }

		fmt.Println("Setting up PortAudio stream...")
		// fmt.Printf("PortAudio channels: %d / SampleRate: %f\n", info.Channels, info.SampleRate)

		streamer, format, err := vorbis.Decode(io.NopCloser(audioFile))
		if err != nil {
			log.Fatal(err)
		}

		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		speaker.Play(streamer)

	}

	return nil
}
