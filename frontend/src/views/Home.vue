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
            v-for="playlist in playlists"
            :key="playlist.id"
            cols="2"
          >
            <v-card class="d-flex">
              <v-img
                width="100"
                height="100"
                cover
                aspect-ratio="1"
                :src="playlist.image"
              />
              <div>
                <v-card-title>{{ playlist.title }}</v-card-title>
                <v-card-subtitle>{{ playlist.subtitle }}</v-card-subtitle>
              </div>
            </v-card>
          </v-col>
        </v-row>
      </div>
    </div>
  </v-main>
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
      playlists: [
        {
          id: 1,
          image: "https://picsum.photos/200?random=1",
          title: "Rock",
          subtitle: "56 titres"
        },
         {
          id: 2,
          image: "https://picsum.photos/200?random=2",
          title: "Pop",
          subtitle: "25 titres"
        },
         {
          id: 3,
          image: "https://picsum.photos/200?random=3",
          title: "Electro",
          subtitle: "10 titres"
        },
         {
          id: 5,
          image: "https://picsum.photos/200?random=4",
          title: "Swing",
          subtitle: "68 titres"
        },
         {
          id: 5,
          image: "https://picsum.photos/200?random=5",
          title: "Jazz",
          subtitle: "2 titres"
        }
      ]
    });

    return {
      ...toRefs(state),
      authStore,
      router,
    };
  },
};
</script>
