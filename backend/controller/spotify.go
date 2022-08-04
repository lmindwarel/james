package controller

import (
	"github.com/librespot-org/librespot-golang/librespot"
	"github.com/lmindwarel/james/backend/models"
	"github.com/lmindwarel/james/backend/utils"
	"github.com/pkg/errors"
)

var ErrNotConnected = errors.New("spotify not connected")

func (ctrl *Controller) AuthenticateSpotify(credentials models.SpotifyCredentials) (err error) {
	// decrypt password
	password, err := utils.Decrypt(models.SpotifyPasswordHashKey, credentials.HashedPassword)
	if err != nil {
		return errors.Wrap(err, "failed to decrypt password")
	}

	// Authenticate using a regular login and password, and store it in the blob file.
	ctrl.SpotifySession, err = librespot.Login(credentials.User, password, models.SpotifyDeviceName)
	if err != nil {
		return errors.Wrap(err, "failed to authenticate to Spotify")
	}

	return
}

func (ctrl *Controller) GetSpotifyPlaylists() (playlists []models.SpotifyPlaylist, err error) {
	if ctrl.SpotifySession == nil {
		return playlists, ErrNotConnected
	}

	playlistResult, err := ctrl.SpotifySession.Mercury().GetRootPlaylist(ctrl.SpotifySession.Username())
	if err != nil || playlistResult.Contents == nil {
		return playlists, errors.Wrapf(err, "failed to get playlists")
	}

	for _, item := range playlistResult.Contents.Items {
		spotifyURI := models.SpotifyURI(item.GetUri())
		list, _ := ctrl.SpotifySession.Mercury().GetPlaylist(spotifyURI.GetURL())
		playlists = append(playlists, models.SpotifyPlaylist{
			ID:   spotifyURI.GetID(),
			Name: *list.Attributes.Name,
			URI:  spotifyURI,
		})
	}

	return
}

func (ctrl *Controller) GetSpotifyTracks(uri models.SpotifyURI) (tracks []models.SpotifyTrack, err error) {
	if ctrl.SpotifySession == nil {
		return tracks, ErrNotConnected
	}

	playlistResult, err := ctrl.SpotifySession.Mercury().GetPlaylist(uri.GetURL())
	if err != nil || playlistResult.Contents == nil {
		return tracks, errors.Wrapf(err, "failed to get tracks")
	}

	for _, item := range playlistResult.Contents.Items {
		log.Debugf("%+v", item.Uri)
		// spotifyURI := models.SpotifyURI(item.GetUri())
		// track, _ := ctrl.SpotifySession.Mercury().GetTrack(spotifyURI.GetURL())
		// tracks = append(tracks, models.SpotifyTrack{
		// 	ID:   spotifyURI.GetID(),
		// 	Name: track.GetName(),
		// 	URI:  spotifyURI,
		// })
	}

	return
}
