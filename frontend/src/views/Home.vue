<template>
  <v-container>
    <v-sheet
      height="500"
      class="overflow-y-auto"
    >
      <v-row>
        <v-col>
          <h2>Playlists</h2>
          <v-list />
        </v-col>
      </v-row>
    </v-sheet>
  </v-container>
</template>

<script lang="ts">
import { reactive, toRefs } from "vue";
import { useAuthStore } from "@/plugins/store/auth";
import { useRouter } from "vue-router";

export default {
  setup() {
    const authStore = useAuthStore();
    const router = useRouter();

    const state = reactive({
    });

    return {
      ...toRefs(state),
      authStore,
      router,
    };
  },
};
</script>
