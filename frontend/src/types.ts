export interface BaseModel{
    id: string
    date_created: string
    date_updated: string
}

export interface JamesBasics {
    player_status: PlayerStatus
    player_queue: QueuedTrack[]
}

export interface Parameters {
    current_spotify_credential: string | null
}

export interface Parameter {
    id: PARAMETERS_IDS,
    value: any
}

export enum PARAMETERS_IDS {
    CURRENT_SPOTIFY_CREDENTIAL = 'current_spotify_credential',
}

export interface Status {
    player_status: PlayerStatus
}

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

export interface SpotifyCredential extends BaseModel {
	user: string
}

export interface CredentialPatch {
	user?: string
    password?: string
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
    artists: SpotifyArtist[]
}

export interface SpotifyArtist {
    id: string
    name: string
    uri: string
    href: string
}

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
    Loading = 'loading',
}

export interface WebsocketMessage {
    topic: string,
    data: any,
}

export interface PlayerStatus {
	state: PlayerStates
    current_queue_index: number
	track_duration: number
	track_position: number
    volume: number
}

export interface SpotifyPlayerControl {
	volume?: number
	track_position_ms?: number,     
	pause?: boolean    
}

export interface QueuedTrack {
    track_id: string
    manually_added: boolean
}