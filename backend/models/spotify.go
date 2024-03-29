package models

import (
	"strings"
)

type SpotifyURI string

func (uri SpotifyURI) GetID() string {
	uriParts := strings.Split(string(uri), ":")
	return uriParts[len(uriParts)-1]
}

func (uri SpotifyURI) GetURL() string {
	url := strings.TrimPrefix(string(uri), "spotify:")
	return strings.ReplaceAll(url, ":", "/")
}

type SpotifyPlaylist struct {
	ID   string     `json:"id"`
	Name string     `json:"name"`
	URI  SpotifyURI `json:"uri"`
}

type SpotifyTrack struct {
	ID   string     `json:"id"`
	Name string     `json:"name"`
	URI  SpotifyURI `json:"uri"`
}

type SpotifyPlayerControl struct {
	Volume          *float64 `json:"volume"` // 0-1
	TrackPositionMs *int     `json:"track_position_ms"`
	Pause           *bool    `json:"pause"`
}

type AddToPlayerQueueOptions struct {
	PlaylistID *string `json:"playlist_id"`
}
