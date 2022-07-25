import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from '@/App.vue'
import vuetify from '@/plugins/vuetify'
import router from '@/plugins/router'
import constants from '@/plugins/constants'
import { loadFonts } from './plugins/webfontloader'
import eventbus from '@/services/eventbus'

loadFonts()
const pinia = createPinia()

const app = createApp(App)
  .use(vuetify)
  .use(router)
  .use(constants)
  .use(pinia)

app.config.globalProperties.$eventbus = eventbus

app.mount('#app')
