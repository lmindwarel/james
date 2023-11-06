package spotify

import (
	"context"

	"github.com/pkg/errors"
	"github.com/zmb3/spotify/v2"
)

func (s *Session) GetCurrentUserPlaylists(ctx context.Context) (result *spotify.SimplePlaylistPage, err error) {
	result, err = s.webapiClient.GetPlaylistsForUser(ctx, s.userID)
	if err != nil || result == nil {
		return result, errors.Wrapf(err, "failed to get playlists")
	}

	return
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
}

func (s *Session) GetTrack(ctx context.Context, id ID) (result *spotify.FullTrack, err error) {
	result, err = s.webapiClient.GetTrack(ctx, spotify.ID(id))
	if err != nil || result == nil {
		return result, errors.Wrapf(err, "failed to get track %s", id)
	}

	return
}

func (s *Session) GetCurrentUserSavedTracks(ctx context.Context) (result *spotify.SavedTrackPage, err error) {
	result, err = s.webapiClient.CurrentUsersTracks(ctx)
	if err != nil || result == nil {
		return result, errors.Wrapf(err, "failed to get saved tracks")
	}

	return
}
