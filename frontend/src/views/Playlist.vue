<template>
  <v-main>
    <v-row
      class="ma-2"
      align="center"
    >
      <v-col cols="2">
        <v-img :src="playlist?.images[0].url" />
      </v-col>
      <v-col>
        <h1 class="text-h1">
          {{ playlist?.name }}
        </h1>
      </v-col>
    </v-row>

    <v-table>
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
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="playlistTrack of tracks"
          :key="playlistTrack.track.id"
        >
          <td class="text-right">
            <v-btn
              variant="text"
              icon="mdi-play"
              @click="playTrack(playlistTrack.track.id)"
            />
          </td>
          <td>
            <v-list-item-title>{{ playlistTrack.track.name }}</v-list-item-title>
            <v-list-item-subtitle>{{ playlistTrack.track.name }}</v-list-item-subtitle>
          </td>
          <td>{{ playlistTrack.track.album.name }}</td>
          <td>{{ moment(playlistTrack.added_at).format("DD MMM. YYYY") }}</td>
          <td class="text-right">
            {{ playlistTrack.track.duration_ms }} ms
          </td>
        </tr>
      </tbody>
    </v-table>
  </v-main>
</template>

<script lang="ts">
import { reactive, onMounted, watch, toRefs } from "vue";
import { SpotifyPlaylistTrack, SpotifyPlaylist } from "@/types";
import { useRoute } from "vue-router";
import api from "@/services/api";
import moment from 'moment'

export default {
  setup() {
    const state = reactive({
      tracks: [] as SpotifyPlaylistTrack[],
      loading: {
        playlist: false,
        tracks: false,
      },
      playlistID: "",
      playlist: null as SpotifyPlaylist|null
    });

    const route = useRoute();

    function refresh(){
      fetchPlaylist()
      fetchPlaylistTracks()
    }

    function fetchPlaylist() {
      state.loading.playlist = true;
      api
        .getSpotifyPlaylist(state.playlistID)
        .then(({ data }) => {
          state.playlist = data;
        })
        .finally(() => {
          state.loading.playlist = false;
        });
    }

    function fetchPlaylistTracks() {
      state.loading.tracks = true;
      api
        .getSpotifyPlaylistTracks(state.playlistID)
        .then(({ data }) => {
          state.tracks = data.items;
        })
        .finally(() => {
          state.loading.tracks = false;
        });
    }

    function playTrack(id){
      api.playSpotifyTrack(id)
    }

    onMounted(() => {
      state.playlistID = route.params.id as string;

      if (!state.playlistID) {
        console.warn("playlist id not found");
        return;
      }

      refresh();
    });

    watch(
      () => route.params.uri,
      () => {
        state.playlistID = route.params.id as string;

        if (!state.playlistID) {
          console.warn("playlist id not found");
          return;
        } else {
          refresh();
        }
      }
    );

    return {
      ...toRefs(state),
      moment,
      playTrack
    }
  },
};
</script>
