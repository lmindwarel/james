<template>
  <v-row
    justify="space-between"
    align="center"
  >
    <v-col cols="3">
      <v-row align="center">
        <v-col cols="3">
          <v-img
            v-if="playerStore.current_track"
            :aspect-ratio="1"
            width="100"
            cover
            :src="playerStore.current_track.album.images[0].url"
          />
          <v-icon
            v-else
            size="100"
          >
            mdi-music
          </v-icon>
        </v-col>
        <v-col>
          <span v-if="playerStore.current_track">{{ playerStore.current_track.name }}</span>
        </v-col>
      </v-row>
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
      <div class="d-flex align-center">
        <v-btn
          variant="flat"
          icon="mdi-menu"
          @click="$emit('open-details')"
        />
        <v-slider
          density="compact"
          hide-details
          prepend-icon="mdi-volume-high"
        />
      </div>
    </v-col>
  </v-row>
</template>

<script lang="ts">
import { computed, reactive, toRefs, } from "vue";
import { usePlayerStore } from "@/plugins/store/player";
import { millisToMinutesAndSeconds } from "@/utils";
import { PlayerStates, SpotifyPlayerControl } from "@/types";
import api from "@/services/api";
import _ from "lodash";

export default {
  setup() {
    const playerStore = usePlayerStore();

    const controlDebouced = _.debounce((control: SpotifyPlayerControl) => {
      api.controlSpotifyPlayer(control).then(res => {
        playerStore.updateFromPlayerStatus(res.data)
      });
    }, 400);

    const progression = computed({
      get() {
        return !!playerStore.current_track
          ? (playerStore.track_position / playerStore.current_track?.duration_ms) * 100
          : 0;
      },
      set(newValue: number) {
        if (playerStore.current_track) {
          controlDebouced({
            track_position_ms: Math.floor(
              playerStore.current_track?.duration_ms * (newValue / 100)
            ),
          });
        }
      },
    });

    const progressionText = computed(() =>
      millisToMinutesAndSeconds(playerStore.track_position)
    );
    const trackDurationText = computed(() =>
      playerStore.current_track
        ? millisToMinutesAndSeconds(playerStore.current_track.duration_ms)
        : "0:00"
    );

    function togglePlayerState() {
      if (playerStore.current_track) {
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