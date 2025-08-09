<template>
  <v-overlay v-model="loading" class="align-center justify-center">
    <v-progress-circular color="error" size="64" indeterminate></v-progress-circular>
  </v-overlay>
  <v-overlay v-model="offline" class="align-center justify-center" contained>
    <v-sheet :elevation="24" rounded>
      <v-empty-state v-if="!data.external && data.backend">
        <template v-slot:media>
          <v-icon icon="mdi mdi-cloud-off-outline" color="error"></v-icon>
        </template>
        <template v-slot:headline>
          <div class="text-h5">
            You are offline!
          </div>
        </template>
        <template v-slot:title>
          <div class="text-h6">
            Check your internet connection
          </div>
        </template>
      </v-empty-state>
      <v-empty-state v-else-if="!data.backend">
        <template v-slot:media>
          <v-icon icon="mdi mdi-connection" color="error"></v-icon>
        </template>
        <template v-slot:headline>
          <div class="text-h5">
            Backend is offline!
          </div>
        </template>
        <template v-slot:title>
          <div class="text-h6">
            check if GoXterm is running!
          </div>
        </template>
      </v-empty-state>
    </v-sheet>
  </v-overlay>
</template>

<script setup>
import { BACKEND_HOST } from '@/utils';
import axios from 'axios';
import { computed, onMounted, reactive } from 'vue';

const data = reactive({
  external: null,
  backend: null
});

const ping = () => {
  axios.get(`http://${BACKEND_HOST}/api/ping`)
    .then(response => {
      if (response.headers['content-type'] === "application/json") {
        data.external = response.data?.alive || false;
        data.backend = true;
      }
    })
    .catch(() => {
      data.backend = false;
    });
}

onMounted(() => {
  ping();
  setInterval(ping, 3000);
});

const loading = computed(() => {
  return (data.external === null && data.backend !== null) || data.backend === null ? true : false;
});

const offline = computed(() => {
  return !data.backend || !data.external ? true : false;
});
</script>