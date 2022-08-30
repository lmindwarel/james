package spotify

import (
	"context"

	"github.com/pkg/errors"
	"github.com/xlab/portaudio-go/portaudio"
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
