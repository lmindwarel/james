<template>
  <div class="fill-height d-flex flex-column justify-space-around align-center">
    <div class="d-flex flex-column justify-center align-center mb-12">
      <v-icon
        size="300"
        class="mb-4"
      >
        mdi-robot-excited-outline
      </v-icon>
      <span class="text-h1">Hello, moi c'est James !</span>
    </div>
    <v-card
      width="600"
      max-height="300"
      class="overflow-y-auto"
    >
      <v-card-title>
        Connexion
      </v-card-title>
      <v-window v-model="page">
        <v-window-item>
          <v-list tile>
            <v-list-item
              v-for="account in accounts"
              :key="account.id"
              :prepend-icon="account.icon || 'mdi-account'"
              :title="account.name"
              @click="login(account)"
            />
            <v-list-item
              prepend-icon="mdi-account-plus"
              title=" Créer un compte"
              @click="page = 1"
            />
          </v-list>
        </v-window-item>
        <v-window-item>
          <v-card-text>
            <div class="d-flex justify-space-between">
              <v-btn
                flat
                @click="page = 0"
              >
                <v-icon start>
                  mdi-arrow-left
                </v-icon>
                Retour
              </v-btn>
              <v-btn
                color="success"
                :disabled="!newAccountName.trim()"
                @click="createAccount"
              >
                Créer
              </v-btn>
            </div>

            <v-text-field
              v-model="newAccountName"
              label="Pseudo"
              class="my-4"
            />
          </v-card-text>
        </v-window-item>
      </v-window>
    </v-card>
  </div>
</template>

<script lang="ts">
import api from "@/services/api";
import { onMounted, reactive, toRefs } from "vue";
import { Account } from "@/types";
import { useAuthStore } from '@/plugins/store/auth';
import { useRouter } from 'vue-router';
import { ROUTE_NAMES } from '@/constants';

export default {
  setup() {
    const authStore = useAuthStore();
    const router = useRouter()

    const state = reactive({
      accounts: [] as Account[],
      newAccountName: "",
      page: 0,
    });

    function fetchAccounts() {
      api.getAccounts().then((res) => {
        state.accounts = res.data;
      });
    }

    function createAccount() {
      api.postAccount({ name: state.newAccountName }).then((res) => {
        state.accounts.push(res.data);
        state.newAccountName = "";
        state.page = 0;
      });
    }

    function login(account: Account){
      authStore.connectedAccount = account;
      router.push({ name: ROUTE_NAMES.HOME });
    }

    onMounted(() => {
      fetchAccounts();
    });

    return {
      ...toRefs(state),
      createAccount,
      login,
    };
  },
};
</script>