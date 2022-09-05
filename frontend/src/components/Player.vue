<template>
  <v-row
    justify="space-between"
    align="center"
  >
    <v-col cols="2">
      <div class="d-flex">
        <v-img
          v-if="playerStore.currentTrack"
          aspect-ratio="1"
          height="50"
          :src="playerStore.currentTrack.album.images[0].url"
        />
        <v-icon size="100">
          <mdi-music />
        </v-icon>
        <span v-if="playerStore.currentTrack">{{ playerStore.currentTrack.name }}</span>
      </div>
    </v-col>

    <v-col cols="4">
      <div class="d-flex flex-column justify-center align-center">
        <div>
          <v-btn
            icon="mdi-skip-previous"
            variant="flat"
          />
          <v-btn
            :icon="playerStore.state == PlayerStates.Playing ? 'mdi-pause': 'mdi-play'"
            @click="togglePlayerState"
          />
          <v-btn
            icon="mdi-skip-next"
            variant="flat"
          />
        </div>
      </div>
      <div class="d-flex align-center">
        {{ progressionText }}
        <v-slider
          v-model="progression"
          hide-details
          class="mx-2"
          density="compact"
        />
        {{ trackDurationText }}
      </div>
    </v-col>

    <v-col cols="2">
      <v-slider
        prepend-icon="mdi-volume-high"
      />
    </v-col>
  </v-row>
</template>

<script lang="ts">
import { computed, ref } from "vue";
import { usePlayerStore } from "@/plugins/store/player";
import { millisToMinutesAndSeconds } from "@/utils";
import { PlayerStates, SpotifyPlayerControl } from "@/types";
import api from "@/services/api";
import _ from "lodash";

export default {
  setup() {
    let playerStore = usePlayerStore();

    const controlDebouced = _.debounce((control: SpotifyPlayerControl) => {
      api.controlSpotifyPlayer(control).then(res => {
        playerStore.updateFromPlayerStatus(res.data)
      });
    }, 400);

    let progression = computed({
      get() {
        return playerStore.currentTrack
          ? (playerStore.track_position / playerStore.currentTrack?.duration_ms) * 100
          : 0;
      },
      set(newValue: number) {
        if (playerStore.currentTrack) {
          controlDebouced({
            track_position_ms: Math.floor(
              playerStore.currentTrack?.duration_ms * (newValue / 100)
            ),
          });
        }
      },
    });

    let progressionText = computed(() =>
      millisToMinutesAndSeconds(playerStore.track_position)
    );
    let trackDurationText = computed(() =>
      playerStore.currentTrack
        ? millisToMinutesAndSeconds(playerStore.currentTrack.duration_ms)
        : "0:00"
    );

    function togglePlayerState() {
      if (playerStore.currentTrack) {
        controlDebouced({
          pause: playerStore.state == PlayerStates.Playing,
        });
      }
    }

    return {
      progression,
      playerStore,
      PlayerStates,
      togglePlayerState,
      progressionText,
      trackDurationText,
    };
  },
};
</script>