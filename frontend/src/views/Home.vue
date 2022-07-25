<template>
  <v-carousel>
    <v-carousel-item v-if="authStore.connectedAccount">
      <v-sheet
        color="orange"
        height="100%"
        tile
      >
        <v-row
          class="fill-height d-flex flex-column"
          align="center"
          justify="center"
        >
          <div class="text-h2">
            Bon retour {{ authStore.connectedAccount.name }}!
          </div>
          <div class="text-body-1">
            Jouez avec vos amis, challengez les autres joueurs en créant vos propre quizz !
          </div>
        </v-row>
      </v-sheet>
    </v-carousel-item>
    <v-carousel-item v-else>
      <v-sheet
        color="#6699cc"
        height="100%"
        tile
      >
        <v-row
          class="fill-height d-flex flex-column"
          align="center"
          justify="center"
        >
          <div class="text-h2">
            Yo moi c'est James
          </div>
          <div class="d-flex flex-wrap my-6">
            <v-icon
              size="150"
            >
              mdi-controller-classic
            </v-icon>
            <v-icon
              size="150"
              class="mx-2"
            >
              mdi-help
            </v-icon>
            <v-icon size="150">
              mdi-play
            </v-icon>
          </div>
          <div class="text-body-1">
            Controllez votre musiques avec vos collègues grâce à moi.
          </div>
        </v-row>
      </v-sheet>
    </v-carousel-item>
  </v-carousel>
</template>

<script>
import { ref } from 'vue'
import { useAuthStore } from '@/plugins/store/auth'
import { useRouter } from 'vue-router'

export default {
  setup () {
    const authStore = useAuthStore()
    const router = useRouter()

    return {
      authStore,
      router,
    }
  }
}
</script>
