// Styles
import '@mdi/font/css/materialdesignicons.css'
import 'vuetify/styles'

// Labs components
import { VSkeletonLoader } from 'vuetify/labs/VSkeletonLoader'

// Vuetify
import { createVuetify } from 'vuetify'

export default createVuetify(
  {
    components: {
      VSkeletonLoader
    },
    theme: {
      themes: {
        light: {
          colors: {
            primary: '#ff7f11',
            secondary: '#beb7a4',
            background: '#fffffc'
            // #d7c0d0, #d33f49, 262730
          }
        }
      }
    }
  }
)
