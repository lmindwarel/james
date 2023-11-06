import axios from 'axios'
import { useAuthStore } from '@/plugins/store/auth'
import { Account, AccountPatch, BasicsData, CredentialPatch, JamesBasics, Parameter, PlayerStatus, QueuedTrack, SpotifyCredential, SpotifyPlayerControl, SpotifyPlaylist, SpotifyPlaylistsResult, SpotifyPlaylistTracksResult, SpotifyTrack } from '@/types'

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
  getBasics:()=> apiClient.get<JamesBasics>('/basics'),
  getParameters: () => apiClient.get<Parameter[]>('/parameters'),
  patchParameter: (parameter: Parameter) => apiClient.patch<Parameter>(`/parameters/${parameter.id}`, parameter),
  getAccounts: () => apiClient.get<Account[]>('/accounts'),
  postAccount: (account: AccountPatch) => apiClient.post<Account>('/accounts', account),
  getSpotifyCredentials: () => apiClient.get<SpotifyCredential[]>('/spotify/credentials'),
  createSpotifyCredential: (credential: CredentialPatch) => apiClient.post<SpotifyCredential>('/spotify/credentials', credential),
  patchSpotifyCredential: (id: string, credential: CredentialPatch) => apiClient.patch<SpotifyCredential>(`/spotify/credentials/${id}`, credential),
  getSpotifySavedTracks: () => apiClient.get<SpotifyPlaylistTracksResult>('/spotify/saved-tracks'),
  getSpotifyPlaylists: () => apiClient.get<SpotifyPlaylistsResult>('/spotify/playlists'),
  getSpotifyPlaylist:(id: string)=> apiClient.get<SpotifyPlaylist>(`/spotify/playlists/${id}`),
  getSpotifyPlaylistTracks: (playlistID: string) => apiClient.get<SpotifyPlaylistTracksResult>(`/spotify/playlists/${playlistID}/tracks`),
  getSpotifyTrack: (trackID: string) => apiClient.get<SpotifyTrack>(`/spotify/tracks/${trackID}`),
  playSpotifyTrack: (id:string) => apiClient.put(`/spotify/player/play/${id}`),
  controlSpotifyPlayer: (control: SpotifyPlayerControl)=> apiClient.put<PlayerStatus>('/spotify/player/control', control),
  addToPlayerQueue: (trackID: string)=> apiClient.post(`/spotify/player/queue/${trackID}`),
  getPlayerQueue: ()=> apiClient.get<QueuedTrack[]>('/spotify/player/queue'),
  removeFromPlayerQueue: (id: string)=> apiClient.delete(`/spotify/player/queue/${id}`)
}