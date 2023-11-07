<template>
  <v-main>
    <v-row
      class="ma-2"
      align="center"
    >
      <v-col cols="2">
        <v-icon
          size="100"
        >
          mdi-music
        </v-icon>
      </v-col>
      <v-col>
        <h1
          class="text-h1"
        >
          Mes titres sauvegardés
        </h1>
      </v-col>
    </v-row>

    <v-skeleton-loader
      v-if="loading.tracks"
      type="table"
    />
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
          <th class="text-left">
            Ajouté le
          </th>
          <th class="text-right">
            <v-icon>mdi-clock-outline</v-icon>
          </th>
          <th class="text-right" />
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="playlistTrack of tracks"
          :key="playlistTrack.track.id"
        >
          <td class="d-flex justify-end">
            <v-sheet
              v-if="playerStore.current_track?.id == playlistTrack.track.id"
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
              @click="playTrack(playlistTrack.track.id)"
            />
          </td>
          <td>
            <v-list-item-title>{{ playlistTrack.track.name }}</v-list-item-title>
            <v-list-item-subtitle>{{ playlistTrack.track.artists[0].name }}</v-list-item-subtitle>
          </td>
          <td>{{ playlistTrack.track.album.name }}</td>
          <td>{{ moment(playlistTrack.added_at).format("DD MMM. YYYY") }}</td>
          <td class="text-right">
            {{ millisToMinutesAndSeconds(playlistTrack.track.duration_ms) }}
          </td>
          <td class="text-center">
            <v-icon v-if="playerStore.queue.some(t => t.track_id == playlistTrack.track.id)">
              mdi-playlist-check
            </v-icon>
            <v-btn
              v-else
              variant="flat"
              icon="mdi-playlist-plus"
              @click="addToQueue(playlistTrack.track.id)"
            />
          </td>
        </tr>
      </tbody>
    </v-table>
  </v-main>
</template>

<script lang="ts">
import { reactive, onMounted, watch, toRefs } from "vue";
import { SpotifyPlaylistTrack } from "@/types";
import { useRoute } from "vue-router";
import api from "@/services/api";
import { millisToMinutesAndSeconds } from "@/utils";
import moment from 'moment'
import { usePlayerStore } from '@/plugins/store/player';

import audioPlayingAnimation from '@/assets/audio-playing-animation.json'

export default {
  setup() {
    const playerStore = usePlayerStore()

    const state = reactive({
      tracks: [] as SpotifyPlaylistTrack[],
      loading: {
        tracks: false,
      },
    });

    function fetchSavedTracks() {
      state.loading.tracks = true;
      api
        .getSpotifySavedTracks()
        .then(({ data }) => {
          state.tracks = data.items;
        })
        .finally(() => {
          state.loading.tracks = false;
        });
    }

    function playTrack(id: string){
      api.playSpotifyTrack(id)
    }

    function addToQueue(trackID: string){
      api.addToPlayerQueue(trackID)
    }

    onMounted(() => {
      fetchSavedTracks();
    });

    return {
      ...toRefs(state),
      moment,
      playTrack,
      playerStore,
      addToQueue,
      millisToMinutesAndSeconds,
      audioPlayingAnimation
    }
  },
};
</script>
