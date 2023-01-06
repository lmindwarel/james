import { reactive } from 'vue';
<template>
  <v-main>
    <div class="ma-10">
      <h1 class="text-h1">
        Paramètres
      </h1>

      <v-row>
        <v-col cols="6">
          <v-card
            class="my-8"
            variant="outlined"
          >
            <v-card-title class="d-flex justify-space-between">
              <span>Comptes Spotify</span>
              <v-btn
                append-icon="mdi-account-plus"
                variant="text"
                @click="newCredential"
              >
                Ajouter
              </v-btn>
            </v-card-title>
              
            <v-list>
              <v-list-item
                v-for="credential in spotifyCredentials"
                :key="credential.id"
                :title="credential.user"
              >
                <template #prepend>
                  <v-sheet
                    width="100px"
                    color="transparent"
                    class="d-flex justify-center"
                  >
                    <v-progress-circular
                      v-if="authenticating_credential_id == credential.id"
                      indeterminate
                    />
                    <v-icon v-else-if="commonStore.parameters.current_spotify_credential == credential.id">
                      mdi-circle
                    </v-icon>
                    <v-btn
                      v-else
                      variant="icon"
                      icon="mdi-login"
                      @click="useCredential(credential.id)"
                    />
                  </v-sheet>
                </template>

                <template #append>
                  <v-avatar>
                    <v-btn
                      icon="mdi-pencil"
                      @click="editCredential(credential)"
                    />
                  </v-avatar>
                </template>
              </v-list-item>
            </v-list>
          </v-card>
        </v-col>
      </v-row>
    </div>

    <v-dialog
      v-model="dialogs.editSpotifyCredential.visible"
      width="unset"
    >
      <v-card min-width="500">
        <v-card-title>{{ dialogs.editSpotifyCredential.id ? 'Éditer' : 'Créer' }} un accès Spotify</v-card-title>
        <v-card-text>
          <v-form ref="editSpotifyCredentialForm">
            <v-text-field
              v-model="dialogs.editSpotifyCredential.user"
              label="Login"
            />
            <v-text-field
              v-model="dialogs.editSpotifyCredential.password"
              type="password"
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
import { onMounted, reactive, ref, toRefs } from "vue";
import { CredentialPatch, PARAMETERS_IDS, SpotifyCredential } from "@/types";
import { VForm } from "vuetify/components";
import api from "../services/api";
import { usePlayerStore } from "@/plugins/store/player";
import eventbus from "@/services/eventbus";
import { useCommonStore } from '@/plugins/store/common';

export default {
  setup() {
    const commonStore = useCommonStore()
    const playerStore = usePlayerStore();

    const state = reactive({
      spotifyCredentials: [] as SpotifyCredential[],
      authenticating_credential_id: null as string | null,
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

    onMounted(fetchSpotifyCredentials);

    function fetchSpotifyCredentials() {
      api.getSpotifyCredentials().then((res) => {
        state.spotifyCredentials = res.data;
      });
    }

    function newCredential() {
      state.dialogs.editSpotifyCredential.visible = true;
    }

    function editCredential(credential: SpotifyCredential) {
      state.dialogs.editSpotifyCredential.id = credential.id;
      state.dialogs.editSpotifyCredential.user = credential.user;
      state.dialogs.editSpotifyCredential.visible = true;
    }

    function closeEditSpotifyCredentialDialog() {
      state.dialogs.editSpotifyCredential.visible = false;
      state.dialogs.editSpotifyCredential.id = null;
      state.dialogs.editSpotifyCredential.user = null;
    }

    function saveEditSpotifyCredentialDialog() {
      if (!editSpotifyCredentialForm.value?.validate()) {
        return;
      }

      state.dialogs.editSpotifyCredential.saving = true;

      let requestPromise;
      if (state.dialogs.editSpotifyCredential.id) {
        const patch = {} as CredentialPatch;
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
        if (
          !state.dialogs.editSpotifyCredential.user ||
          !state.dialogs.editSpotifyCredential.password
        ) {
          console.error("undefined user or password");
          return;
        }
        requestPromise = api.createSpotifyCredential({
          user: state.dialogs.editSpotifyCredential.user,
          password: state.dialogs.editSpotifyCredential.password,
        });
      }

      requestPromise
        .then((res) => {
          if (state.dialogs.editSpotifyCredential.id) {
            state.spotifyCredentials.splice(state.spotifyCredentials.findIndex(c=> c.id == state.dialogs.editSpotifyCredential.id), 1, res.data)
          } else {
            state.spotifyCredentials.push(res.data);
          }
          closeEditSpotifyCredentialDialog();
        })
        .finally(() => {
          state.dialogs.editSpotifyCredential.saving = false;
        });
    }

    function useCredential(id: string) {
      commonStore.parameters.current_spotify_credential = id
      state.authenticating_credential_id = id;
       commonStore.saveParameter(PARAMETERS_IDS.CURRENT_SPOTIFY_CREDENTIAL).catch((err) => {
          console.error(err);
          eventbus.notifyError(
            "Impossible de s'authentifier à Spotify. Veuillez vérifier votre identifiant et mot de passe"
          );
        })
        .finally(() => {
          state.authenticating_credential_id = null;
        });
    }

    return {
      playerStore,
      commonStore,
      ...toRefs(state),
      newCredential,
      editCredential,
      editSpotifyCredentialForm,
      saveEditSpotifyCredentialDialog,
      closeEditSpotifyCredentialDialog,
      useCredential,
    };
  },
};
</script>
