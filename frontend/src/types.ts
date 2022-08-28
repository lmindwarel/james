
export interface Account {
    id: string,
    name: string,
    icon: string
}

export interface AccountPatch{
    name?: string,
    icon?: string
}

export interface SpotifyPagination {
    total: number,
    limit: number,
    offset: number,
}
export interface SpotifyPlaylistsResult extends SpotifyPagination{
    items: SpotifyPlaylist[]
}

export interface SpotifyPlaylist {
    id: string,
    name: string,
    uri: string,
}

export interface SpotifyTrack {
    id: string,
    name: string,
    uri: string,
}