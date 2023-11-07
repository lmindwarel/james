<template>
  <v-main>
    <v-alert
      v-if="!tracks.length"
      border="top"
      colored-border
      type="info"
      class="ma-2"
    >
      Aucun titre en file d'attente
    </v-alert>
    <v-table v-else>
      <thead>
        <tr>
          <th class="text-right">
            #
          </th>
          <th class="text-left">
            Titre
          </th>
          <th class="text-left">
            Album
          </th>
          <th class="text-right">
            <v-icon>mdi-clock-outline</v-icon>
          </th>
          <th class="text-right" />
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="queuedTrack of queuedTracks"
          :key="queuedTrack.id"
        >
          <td class="text-right">
            <v-sheet
              v-if="playerStore.current_track?.id == queuedTrack.id"
              width="50px"
            >
              <vue-3-lottie
                :animation-data="audioPlayingAnimation"
              />
            </v-sheet>
            
            <v-btn
              v-else
              variant="text"
              icon="mdi-play"
              @click="playTrack(queuedTrack.id)"
            />
          </td>
          <td>
            <v-list-item-title>{{ queuedTrack.track.name }}</v-list-item-title>
            <v-list-item-subtitle>{{ queuedTrack.track.name }}</v-list-item-subtitle>
          </td>
          <td>{{ queuedTrack.track.album.name }}</td>
          <td class="text-right">
            {{ millisToMinutesAndSeconds(queuedTrack.track.duration_ms) }}
          </td>
          <td>
            <v-btn
              variant="flat"
              icon="mdi-delete"
              @click="removeTrackFromQueue(queuedTrack.id)"
            />
          </td>
        </tr>
      </tbody>
    </v-table>
  </v-main>
</template>

<script lang="ts">
import { reactive, onMounted, toRefs, watch, computed } from "vue";
import { SpotifyTrack } from "@/types";
import api from "@/services/api";
import { millisToMinutesAndSeconds } from "@/utils";
import moment from 'moment'
import { usePlayerStore } from "@/plugins/store/player";

import audioPlayingAnimation from '@/assets/audio-playing-animation.json'

export default {
  setup() {
    const state = reactive({
      tracks: [] as SpotifyTrack[],
      loading: {
        tracks: false,
      },
    });

    const playerStore = usePlayerStore();

    function refresh(){
      fetchTracks()
    }

    function fetchTracks() {
      const tracksIDs = playerStore.queue.map((queuedTrack) => queuedTrack.track_id);

      if (!tracksIDs.length){
        state.tracks = []
        return
      }

      state.loading.tracks = true;
      api
        .getSpotifyTracks(tracksIDs)
        .then(({ data }) => {
          state.tracks = data;
        })
        .finally(() => {
          state.loading.tracks = false;
        });
    }

    function playTrack(id: string){
      api.playSpotifyTrack(id)
    }

    function removeTrackFromQueue(id: string){
      api.removeFromPlayerQueue(id)
    }

    watch(() => playerStore.queue, refresh);

    onMounted(() => {
      refresh();
    });

  const queuedTracks = computed(() => playerStore.queue.map((queuedTrack) => Object.assign({}, queuedTrack, {track: state.tracks.find((track) => track.id == queuedTrack.track_id)})))

    return {
      ...toRefs(state),
      moment,
      playTrack,
      playerStore,
      removeTrackFromQueue,
      millisToMinutesAndSeconds,
      audioPlayingAnimation,
      queuedTracks
    }
  },
};
</script>
