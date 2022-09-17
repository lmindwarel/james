import { defineStore } from 'pinia'
import { PlayerStates, PlayerStatus, QueuedTrack, SpotifyTrack } from '@/types'
import api from '@/services/api'

export const usePlayerStore = defineStore('player', {
  state: () => ({
    queue: [] as QueuedTrack[],
    state: PlayerStates.Stopped as PlayerStates,
    track_position: 0 as number,
    currentTrack: null as SpotifyTrack | null
  }),

  getters: {
  },

  actions: {
    updateFromPlayerStatus(status: PlayerStatus) {
      console.log("loading from player status")
      if (!this.currentTrack || status.current_track_id != this.currentTrack.id) {
        // load the new curent track
        api.getSpotifyTrack(status.current_track_id).then(res => {
          this.currentTrack = res.data
        })
      }

      this.state = status.state
      this.track_position = status.track_position
    }
  }
})