import { reactive } from 'vue';
<template>
  <v-main>
    <div class="ma-10">
      <h1 class="text-h1">
        Paramètres
      </h1>

      <div>
        <h3 class="text-h3">
          Comptes spotify
        </h3>
        <v-list>
          <v-list-item
            v-for="credential in spotifyCredentials"
            :key="credential.id"
            :title="credential.user"
          >
            <template #prepend>
              <v-avatar :color="file.color">
                <v-btn
                  icon="mdi-pencil"
                  @click="editCredential(credential)"
                />
              </v-avatar>
            </template>
          </v-list-item>
          <v-list-item
            append-icon="mdi-account-plus"
            title="Nouveau compte Spotify"
            @click="newCredential"
          />
        </v-list>
      </div>
    </div>

    <v-dialog v-model="dialogs.editSpotifyCredential.visible">
      <v-card>
        <v-card-title>{{ dialogs.editSpotifyCredential.id ? 'Éditer' : 'Créer' }} un accès Spotify</v-card-title>
        <v-card-text>
          <v-form ref="editSpotifyCredentialForm">
            <v-text-field
              v-model="dialogs.editSpotifyCredential.user"
              label="Login"
            />
            <v-text-field
              v-model="dialogs.editSpotifyCredential.password"
              label="Mot de passe"
            />
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn
            color="red"
            @click="closeEditSpotifyCredentialDialog"
          >
            Annuler
          </v-btn>
          <v-btn
            color="success"
            @click="saveEditSpotifyCredentialDialog"
          >
            Enregistrer
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-main>
</template>

<script lang="ts">
import { reactive, ref, toRefs } from "vue";
import { CredentialPatch, SpotifyCredential } from "@/types";
import { VForm } from "vuetify/lib";
import api from "@/services/api";

export default {
  setup() {
    const state = reactive({
      spotifyCredentials: [] as SpotifyCredential[],
      dialogs: {
        editSpotifyCredential: {
          visible: false,
          saving: false,
          id: null as string | null,
          user: null as string | null,
          password: null as string | null,
        },
      },
    });

    const editSpotifyCredentialForm = ref(VForm);

    function newCredential() {
      state.dialogs.editSpotifyCredential.visible = true;
    }

    function editCredential(credential: SpotifyCredential) {
      state.dialogs.editSpotifyCredential.id = credential.id;
      state.dialogs.editSpotifyCredential.user = credential.user;
      state.dialogs.editSpotifyCredential.visible = true;
    }

    function saveCredential() {
      state.dialogs.editSpotifyCredential.saving = true;

      if (!editSpotifyCredentialForm.value?.validate()) {
        return;
      }

      let requestPromise;
      if (state.dialogs.editSpotifyCredential.id) {
        let patch = {} as CredentialPatch;
        if (state.dialogs.editSpotifyCredential.user) {
          patch.user = state.dialogs.editSpotifyCredential.user;
        }
        if (state.dialogs.editSpotifyCredential.password) {
          patch.password = state.dialogs.editSpotifyCredential.password;
        }

        requestPromise = api.patchSpotifyCredential(
          state.dialogs.editSpotifyCredential.id,
          patch
        );
      } else {
        requestPromise = api.createSpotifyCredential({
          user: state.dialogs.editSpotifyCredential.user,
          password: state.dialogs.editSpotifyCredential.password,
        });
      }
      
    }

    return {
      ...toRefs(state),
      newCredential,
      editCredential,
      editSpotifyCredentialForm,
    };
  },
};
</script>
