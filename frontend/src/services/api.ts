import axios from 'axios'
import { useAuthStore } from '@/plugins/store/auth'
import { Account, AccountPatch, SpotifyPlaylist, SpotifyPlaylistsResult, SpotifyPlaylistTracksResult, SpotifyTrack } from '@/types'

const apiClient = axios.create({
  // @ts-ignore
  baseURL: import.meta.env.VITE_JAMES_API_ADDRESS,
  headers: {
    'Content-Type': 'application/json',
  },
  transformRequest: [function (req, headers) {
    if (headers){
      const accountID = useAuthStore().connectedAccount?.id
      if (accountID) {
        headers['X-Doer'] = accountID
      }
    }
    return JSON.stringify(req)
  }],
})

apiClient.interceptors.response.use((res) => res, (res) => {
  if (res.response) {
    res.httpStatus = res.response.status
    if (res.response.data && res.response.data.code) {
      res.apiCode = res.response.data.code
    }
    if (res.httpStatus < 200 && res.httpStatus >= 300) {
      throw res
    }
  }
  return Promise.reject(res)
})


export default {
  getAccounts: () => apiClient.get<Account[]>('/accounts'),
  postAccount: (account: AccountPatch) => apiClient.post<Account>('/accounts', account),
  getSpotifyPlaylists: () => apiClient.get<SpotifyPlaylistsResult>('/spotify/playlists'),
  getSpotifyPlaylist:(id: string)=> apiClient.get<SpotifyPlaylist>(`/spotify/playlists/${id}`),
  getSpotifyPlaylistTracks: (playlistID: string) => apiClient.get<SpotifyPlaylistTracksResult>(`/spotify/playlists/${playlistID}/tracks`),
  playSpotifyTrack: (id:string) => apiClient.put(`/spotify/play-track/${id}`)
}