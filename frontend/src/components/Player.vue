<template>
  <v-row
    justify="space-between"
    align="center"
  >
    <v-col cols="2">
      Musique en cours
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
import { usePlayerStore } from '@/plugins/store/player';
import { millisToMinutesAndSeconds } from '@/utils'
import { PlayerStates } from '@/types';
import api from '@/services/api';

export default {
  setup() {
    let playerStore = usePlayerStore()

    let progression = computed({
        get(){
        return playerStore.currentTrack ? playerStore.track_position / playerStore.currentTrack?.duration_ms : 0
        },
        set(newValue: number){
            if(playerStore.currentTrack){
                api.controlSpotifyPlayer({track_position_ms: playerStore.currentTrack?.duration_ms * newValue})
            }
        }
    })

    let progressionText = computed(() => millisToMinutesAndSeconds(playerStore.track_position))
    let trackDurationText = computed(() => playerStore.currentTrack ? millisToMinutesAndSeconds(playerStore.currentTrack.duration_ms) : '0:00')

    function togglePlayerState(){
        if(playerStore.currentTrack){
            api.controlSpotifyPlayer({pause: playerStore.state == PlayerStates.Playing})
        }
    }

    return {
      progression,
      playerStore,
      PlayerStates,
      togglePlayerState,
      progressionText,
      trackDurationText
    };
  },
};
</script>