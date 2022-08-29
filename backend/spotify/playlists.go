package spotify

import (
	"context"
	"fmt"
	"sync"

	"github.com/librespot-org/librespot-golang/Spotify"
	"github.com/librespot-org/librespot-golang/librespot/utils"
	"github.com/pkg/errors"
	"github.com/xlab/portaudio-go/portaudio"
	"github.com/xlab/vorbis-go/decoder"
	"github.com/zmb3/spotify/v2"
)

const (
	samplesPerChannel = 2048
	// The samples bit depth
	bitDepth = 16
	// The samples format
	sampleFormat = portaudio.PaFloat32
)

func (s *Session) GetCurrentUserPlaylists(ctx context.Context) (result *spotify.SimplePlaylistPage, err error) {
	result, err = s.webapiClient.GetPlaylistsForUser(ctx, s.userID)
	if err != nil || result == nil {
		return result, errors.Wrapf(err, "failed to get playlists")
	}

	return

	// playlistResult, err := s.SpotifySession.Mercury().GetRootPlaylist(ctrl.SpotifySession.Username())
	// if err != nil || playlistResult.Contents == nil {
	// 	return playlists, errors.Wrapf(err, "failed to get playlists")
	// }

	// for _, item := range playlistResult.Contents.Items {
	// 	spotifyURI := models.SpotifyURI(item.GetUri())
	// 	list, _ := ctrl.SpotifySession.Mercury().GetPlaylist(spotifyURI.GetURL())
	// 	playlists = append(playlists, models.SpotifyPlaylist{
	// 		ID:   spotifyURI.GetID(),
	// 		Name: *list.Attributes.Name,
	// 		URI:  spotifyURI,
	// 	})
	// }

	// return
}

func (s *Session) GetPlaylist(ctx context.Context, id ID) (result *spotify.FullPlaylist, err error) {
	result, err = s.webapiClient.GetPlaylist(ctx, spotify.ID(id))
	if err != nil || result == nil {
		return result, errors.Wrapf(err, "failed to get playlist %s tracks", id)
	}

	return
}

func (s *Session) GetPlaylistTracks(ctx context.Context, id ID) (result *spotify.PlaylistTrackPage, err error) {
	result, err = s.webapiClient.GetPlaylistTracks(ctx, spotify.ID(id))
	if err != nil || result == nil {
		return result, errors.Wrapf(err, "failed to get playlist %s tracks", id)
	}

	return

	// playlistResult, err := ctrl.SpotifySession.Mercury().GetPlaylist(uri.GetURL())
	// if err != nil || playlistResult.Contents == nil {
	// 	return tracks, errors.Wrapf(err, "failed to get tracks")
	// }

	// log.Debugf("Got %d tracks", len(playlistResult.Contents.Items))
	// log.Debugf("%+v", playlistResult.Contents.Items)

	// for _, item := range playlistResult.Contents.Items {
	// 	spotifyURI := models.SpotifyURI(item.GetUri())

	// 	uri := item.GetUri()
	// 	parts := strings.Split(uri, ":")
	// 	id := parts[len(parts)-1]
	// 	log.Debugf("Getting track %s (%s)", spotifyURI.GetURL(), id)

	// 	track, err := ctrl.SpotifySession.Mercury().GetTrack(id)
	// 	if err != nil {
	// 		return tracks, errors.Wrapf(err, "failed to get track")
	// 	}

	// 	if track == nil {
	// 		return tracks, errors.Errorf("track is nil")
	// 	}

	// 	log.Debugf("%+v", *track)

	// 	tracks = append(tracks, models.SpotifyTrack{
	// 		ID:   spotifyURI.GetID(),
	// 		Name: track.GetName(),
	// 		URI:  spotifyURI,
	// 	})
	// }

	// return
}

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
		dec, err := decoder.New(audioFile, samplesPerChannel)
		if err != nil {
			return errors.Wrap(err, "failed to load decoder")
		}

		info := dec.Info()

		go func() {
			dec.Decode()
			dec.Close()
		}()

		fmt.Println("Setting up PortAudio stream...")
		fmt.Printf("PortAudio channels: %d / SampleRate: %f\n", info.Channels, info.SampleRate)

		var wg sync.WaitGroup
		var stream *portaudio.Stream
		callback := paCallback(&wg, int(info.Channels), dec.SamplesOut())

		defaultDevice := portaudio.GetDefaultInputDevice()
		log.Debugf("device count: %d", portaudio.GetDeviceCount())
		deviceInfos := portaudio.GetDeviceInfo(defaultDevice)
		if deviceInfos == nil {
			return errors.New("could not find default input device")
		}
		log.Debugf("Default device: %s, is free %t", deviceInfos.Name)

		if err := portaudio.OpenDefaultStream(&stream, 0, info.Channels, sampleFormat, info.SampleRate,
			samplesPerChannel, callback, nil); paError(err) {
			return fmt.Errorf("failed to open PortAudio stream: %s", paErrorText(err))
		}

		fmt.Println("Starting playback...")
		if err := portaudio.StartStream(stream); paError(err) {
			return fmt.Errorf("failed to start PortAudio stream: %s", paErrorText(err))
		}

		wg.Wait()
	}

	return nil
}
