<template>
  <v-list>
    <v-list-item
      v-for="track of tracks"
      :key="track.id"
      icon="mdi-music"
      :title="track.name"
    />
  </v-list>
</template>

<script lang="ts">
import { reactive, onMounted, watch } from "vue";
import { SpotifyTrack } from "@/types";
import { useRoute } from "vue-router";
import api from "@/services/api";
export default {
  setup() {
    const state = reactive({
      tracks: [] as SpotifyTrack[],
      loading: false,
      playlistURI: "",
    });

    const route = useRoute();

    function fetchPlaylistTracks() {
      state.loading = true;
      api
        .getSpotifyPlaylistTracks(state.playlistURI)
        .then(({ data }) => {
          state.tracks = data;
        })
        .finally(() => {
          state.loading = false;
        });
    }

    onMounted(() => {
      state.playlistURI = route.params.uri as string;

      if (!state.playlistURI) {
        console.warn("playlist id not found");
        return;
      }

      fetchPlaylistTracks();
    });

    watch(
      () => route.params.uri,
      () => {
        state.playlistURI = route.params.uri as string;

        if (!state.playlistURI) {
          console.warn("playlist id not found");
          return;
        } else {
          fetchPlaylistTracks();
        }
      }
    );
  },
};
</script>
