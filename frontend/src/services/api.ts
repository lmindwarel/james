import axios from 'axios'
import { useAuthStore } from '@/plugins/store/auth'
import { Account, AccountPatch, SpotifyPlaylist } from '@/types'

let apiClient = axios.create({
  // @ts-ignore
  baseURL: import.meta.env.VITE_JAMES_API_ADDRESS,
  headers: {
    'Content-Type': 'application/json',
  },
  transformRequest: [function (req, headers) {
    if (headers){
      let accountID = useAuthStore().connectedAccount?.id
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
  // public
  getAccounts: () => apiClient.get<Account[]>('/accounts'),
  postAccount: (account: AccountPatch) => apiClient.post<Account>('/accounts', account),
  getSpotifyPlaylists: () => apiClient.get<SpotifyPlaylist[]>('/spotify/playlists'),
}