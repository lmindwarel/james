<template>
  <v-app>
    <v-main>
      <v-app-bar
        color="primary"
        app
      >
        <div
          class="d-flex"
          @click="$router.push({name: $constants.ROUTE_NAMES.HOME})"
        >
          <v-app-bar-title class="font-weight-bold">
            James
          </v-app-bar-title>
          <v-icon end>
            mdi-comment-question-outline
          </v-icon>
        </div>
        <v-spacer />

        <div
          class="d-flex align-center"
        >
          <v-btn
          >
            <v-icon start>
              mdi-account
            </v-icon>
            {{ authStore.connectedAccount?.name }}
          </v-btn>
        </div>
      </v-app-bar>

      <router-view />

      <v-snackbar
        v-model="snackbar.visible"
        :color="snackbar.color"
      >
        {{ snackbar.message }}

        <template #actions>
          <v-btn
            icon="mdi-close"
            @click="snackbar.visible = false"
          />
        </template>
      </v-snackbar>

      <v-footer
        app
      >
        <span class="white--text">&copy; {{ currentYear }}</span>
        <v-spacer />
        <span
          v-if="buildType == 'development'"
          class="white--text"
        >{ development build }</span>
      </v-footer>
    </v-main>
  </v-app>
</template>

<script>
import moment from 'moment'
import { computed, onMounted, reactive, toRefs } from 'vue'

import { useRouter } from 'vue-router'
import { useAuthStore } from '@/plugins/store/auth'

import eventbus from '@/services/eventbus'

export default {
  setup () {
    const authStore = useAuthStore()
    const router = useRouter()

    const state = reactive({
      snackbar: {
        visible: false,
        color: '',
        message: ''
      }
    })

    // Manage events
    eventbus.on('notify', function ({ message, timeout, color }) {
      state.snackbar.message = message
      state.snackbar.color = color
      state.snackbar.visible = true

      if (timeout) {
        setTimeout(() => {
          state.snackbar.visible = false
        }, timeout)
      }
    })

    eventbus.on('unhandled-api-error', function (err) {
      console.error(err)
      eventbus.notifyError("J'ai piscine ont piscine. Réessayez dans un instant.")
    })

    const currentYear = computed(() => moment().year())

    const buildType = computed(() => process.env.NODE_ENV)

    return {
      ...toRefs(state),
      currentYear,
      buildType,
      authStore,
      router
    }
  }
}
</script>
