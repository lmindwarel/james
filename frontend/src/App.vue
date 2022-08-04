<template>
  <v-app>
      <v-app-bar
        color="primary"
        app
        
      >
        <div @click="$router.push({name: $constants.ROUTE_NAMES.HOME})" class="d-flex">
<v-app-bar-title class="font-weight-bold">
            James
          </v-app-bar-title>
          <v-icon end>
            mdi-comment-question-outline
          </v-icon>
        </div>
        <v-sheet width="500" color="transparent">
<v-text-field variant="solo" append-inner-icon="mdi-magnify" class="ml-6" bg-color="white" density="compact" hide-details single-line
            placeholder="Rechercher..." v-model="search"></v-text-field>
            </v-sheet>
 <div
          class="d-flex align-center"
        >
        </div>

        <v-spacer></v-spacer>

        <v-btn
          >
            <v-icon start>
              mdi-account
            </v-icon>
            {{ authStore.connectedAccount?.name }}
          </v-btn>
      </v-app-bar>

              <v-navigation-drawer permanent>
      <v-list nav>
       <v-list-item prepend-icon="mdi-heart" title="Titres likés" to="/liked-tracks"></v-list-item>

      <v-divider></v-divider>

      <v-sheet v-if="loadingPlaylists" height="200" class="d-flex justify-center align-center">
        <v-progress-circular indeterminate></v-progress-circular>
</v-sheet>

       <v-list-item v-for="(playlist, i) in spotifyPlaylists" :key="`playlist-${i}`" :to="`/playlist/${playlist.uri}`" :title="playlist.name"></v-list-item>
      </v-list>
    </v-navigation-drawer>

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
  </v-app>
</template>

<script lang="ts">
import moment from "moment";
import { computed, onMounted, reactive, ref, toRefs } from "vue";

import { useRouter } from "vue-router";
import { useAuthStore } from "@/plugins/store/auth";

import eventbus from "@/services/eventbus";
import api from "@/services/api";

import { SpotifyPlaylist } from "@/types";

export default {
  setup() {
    const authStore = useAuthStore();
    const router = useRouter();

    const search = ref("");

    const state = reactive({
      snackbar: {
        visible: false,
        color: "",
        message: "",
      },
      loadingPlaylists: false,
      spotifyPlaylists: [] as SpotifyPlaylist[],
    });

    function fetchSpotifyPlaylists() {
      state.loadingPlaylists = true;
      api.getSpotifyPlaylists().then(({ data }) => {
        state.spotifyPlaylists = data;
      }).finally(() => {
        state.loadingPlaylists = false;
      });
    }

    function showPlaylist(id: string) {
      router.push({
        name: "playlist",
        params: {
          id,
        },
      });
    }

    onMounted(fetchSpotifyPlaylists);

    // Manage events
    eventbus.on(
      "notify",
      function ({ message = "", timeout = 0, color = "primary" }) {
        state.snackbar.message = message;
        state.snackbar.color = color;
        state.snackbar.visible = true;

        if (timeout) {
          setTimeout(() => {
            state.snackbar.visible = false;
          }, timeout);
        }
      }
    );

    eventbus.on("unhandled-api-error", function (err) {
      console.error(err);
      eventbus.notifyError(
        "J'ai piscine ont piscine. Réessayez dans un instant."
      );
    });

    const currentYear = computed(() => moment().year());

    const buildType = computed(() => process.env.NODE_ENV);

    return {
      ...toRefs(state),
      search,
      currentYear,
      buildType,
      authStore,
      router,
    };
  },
};
</script>
