<template>
  <v-main>
    <div class="ma-10">
      <h1 class="text-h1">
        Bonjour {{ authStore.connectedAccount?.name }} !
      </h1>

      <div class="my-8">
        <h1 class="text-h2">
          Mes playlists
        </h1>
        <v-row class="my-2">
          <v-col
            v-for="playlist in spotifyPlaylists"
            :key="playlist.id"
            cols="2"
          >
            <v-card
              class="d-flex"
              :to="`/playlist/${playlist.id}`"
            >
              <v-img
                v-if="playlist?.images.length"
                width="100"
                height="100"
                cover
                aspect-ratio="1"
                :src="playlist?.images[0].url"
              />
              <v-icon
                v-else
                size="100"
              >
                mdi-music
              </v-icon>
              <div>
                <v-card-title>{{ playlist.name }}</v-card-title>
              </div>
            </v-card>
          </v-col>
        </v-row>
      </div>
    </div>
  </v-main>
</template>

<script lang="ts">
import { onMounted, reactive, toRefs, watch } from "vue";
import { useAuthStore } from "@/plugins/store/auth";
import { useRouter } from "vue-router";
import { SpotifyPlaylist } from '@/types';
import api from '@/services/api';
import { useCommonStore } from '@/plugins/store/common';

export default {
  setup() {
    const authStore = useAuthStore();
    const router = useRouter();
    const commonStore = useCommonStore()

    const state = reactive({
      loadingPlaylists: false,
      spotifyPlaylists: [] as SpotifyPlaylist[],
    });

    function fetchSpotifyPlaylists() {
      state.loadingPlaylists = true;
      api
        .getSpotifyPlaylists()
        .then(({ data }) => {
          state.spotifyPlaylists = data ? data.items : [];
        })
        .finally(() => {
          state.loadingPlaylists = false;
        });
    }

    onMounted(fetchSpotifyPlaylists)

    watch(
      () => commonStore.parameters?.current_spotify_credential,
      fetchSpotifyPlaylists
    );

    return {
      ...toRefs(state),
      authStore,
      router,
    };
  },
};
</script>
