import { defineStore } from 'pinia'
import { Parameters, Parameter, PARAMETERS_IDS } from '@/types'
import api from '@/services/api'

const default_parameters: Parameters = {
  current_spotify_credential: null
}

export const useCommonStore = defineStore('common', {
  state: () => ({
    // default values
    parameters: default_parameters
  }),

  getters: {

  },

  actions: {
    loadParameters(parameters: Parameter[]) {
      for (const parameter of parameters) {
        this.loadParameter(parameter)
      }
    },

    loadParameter(parameter: Parameter) {
        console.log("loading parameter ", parameter)
        let value
        switch (parameter.id) {
          // strings
          case PARAMETERS_IDS.CURRENT_SPOTIFY_CREDENTIAL:
            value = parameter.value
            break;

          default:
            console.warn("Unknown parameter: " + parameter.id)
            return
        }

        this.parameters[parameter.id] = value
    },

    saveParameter(id: PARAMETERS_IDS){
      return api.patchParameter({id: id, value: this.parameters[id]}).then(response => {
        this.loadParameter(response.data)
      })
    }
  }
})