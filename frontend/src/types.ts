import { Moment } from 'moment'

export interface Account {
    id: string,
    name: string,
    icon: string
}

export interface AccountPatch {
    name?: string,
    icon?: string
}

export interface SpotifyPagination {
    total: number,
    limit: number,
    offset: number,
}
export interface SpotifyPlaylistsResult extends SpotifyPagination {
    items: SpotifyPlaylist[]
}

export interface SpotifyPlaylist {
    id: string,
    name: string,
    uri: string,
    images: SpotifyImage[]
}

export interface SpotifyTrack {
    id: string,
    name: string,
    uri: string,
    album: SpotifyAlbum
    duration_ms: number
}

export interface PlayerQueuedTrack {
    track_id: string,
    manually_added: boolean,}

export interface SpotifyPlaylistTracksResult extends SpotifyPagination {
    items: SpotifyPlaylistTrack[]
}

export interface SpotifyImage {
    url: string
}

export interface SpotifyAlbum {
    id: string
    name: string
    images: SpotifyImage[]
}
export interface SpotifyPlaylistTrack {
    added_at: string,
    track: SpotifyTrack
}

export enum PlayerStates {
    Playing = 'playing',
    Paused = 'paused',
    Stopped = 'stopped',
}

export interface WebsocketMessage {
    topic: string,
    data: any,
}

export interface PlayerStatus {
	state: PlayerStates
    current_track_id: string
	track_duration: number
	track_position: number
}

export interface SpotifyPlayerControl {
	volume?: number
	track_position_ms?: number,     
	pause?: boolean    
}

export interface QueuedTrack {
    track: SpotifyTrack
    manually_added: boolean
}