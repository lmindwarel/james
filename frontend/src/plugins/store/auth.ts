import { defineStore } from 'pinia'
import { Account } from '@/types'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    connectedAccount: null as Account | null
  }),

  getters: {
    isConnected: (state) => state.connectedAccount != null,
  },

  actions: {
  }
})