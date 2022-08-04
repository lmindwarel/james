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
