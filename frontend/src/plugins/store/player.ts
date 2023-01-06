import { defineStore } from 'pinia'
import {PlayerStates, PlayerStatus, QueuedTrack, SpotifyTrack } from '@/types'
import api from '@/services/api'

export const usePlayerStore = defineStore('player', {
  state: () => ({
    queue: [] as QueuedTrack[],
    state: PlayerStates.Stopped as PlayerStates,
    track_position: 0 as number,
    current_track: null as SpotifyTrack | null,
    current_queue_index: null as number | null,
  }),

  getters: {},

  actions: {
    updateFromPlayerStatus(status: PlayerStatus) {
      console.log("loading from player status")
      if (!this.current_track && !!this.queue[status.current_queue_index]) {
        // load the new curent track
        api.getSpotifyTrack(this.queue[status.current_queue_index].track_id).then(res => {
          this.current_track = res.data
        })
      }

      this.state = status.state
      this.track_position = status.track_position
      this.current_queue_index = status.current_queue_index
    }
  }
})