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
          v-for="queuedTrack of tracks"
          :key="queuedTrack.track.id"
        >
          <td class="text-right">
            <v-btn
              variant="text"
              icon="mdi-play"
              @click="playTrack(queuedTrack.track.id)"
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
              @click="removeTrackFromQueue(queuedTrack.track.id)"
            />
          </td>
        </tr>
      </tbody>
    </v-table>
  </v-main>
</template>

<script lang="ts">
import { reactive, onMounted, toRefs } from "vue";
import { QueuedTrack } from "@/types";
import api from "@/services/api";
import { millisToMinutesAndSeconds } from "@/utils";
import moment from 'moment'

export default {
  setup() {
    const state = reactive({
      tracks: [] as QueuedTrack[],
      loading: {
        tracks: false,
      },
    });

    function refresh(){
      fetchPlaylistTracks()
    }

    function fetchPlaylistTracks() {
      state.loading.tracks = true;
      api
        .getPlayerQueue()
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

    function removeTrackFromQueue(trackID: string){
      api.removeFromPlayerQueue(trackID)
    }

    onMounted(() => {
      refresh();
    });

    return {
      ...toRefs(state),
      moment,
      playTrack,
      removeTrackFromQueue,
      millisToMinutesAndSeconds,
    }
  },
};
</script>
