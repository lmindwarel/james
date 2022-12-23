import { defineStore } from 'pinia'
import {PlayerStates, PlayerStatus, QueuedTrack, SpotifyTrack } from '@/types'
import api from '@/services/api'

export const usePlayerStore = defineStore('player', {
  state: () => ({
    queue: [] as QueuedTrack[],
    state: PlayerStates.Stopped as PlayerStates,
    track_position: 0 as number,
    current_track: null as SpotifyTrack | null,
    authenticated_crendential_id: null as string | null
  }),

  getters: {},

  actions: {
    updateFromPlayerStatus(status: PlayerStatus) {
      console.log("loading from player status")
      if (!this.current_track || status.current_track_id != this.current_track.id) {
        // load the new curent track
        api.getSpotifyTrack(status.current_track_id).then(res => {
          this.current_track = res.data
        })
      }

      this.state = status.state
      this.track_position = status.track_position
    }
  }
})