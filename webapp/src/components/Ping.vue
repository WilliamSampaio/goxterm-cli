<template>
  <v-overlay v-model="loading" class="align-center justify-center">
    <v-progress-circular color="error" size="64" indeterminate></v-progress-circular>
  </v-overlay>
  <span v-if="data.external && data.backend">
    <v-icon icon="mdi mdi-circle" size="x-small" color="rgb(0,255,0)"></v-icon>
    Online
  </span>
  <span v-else-if="!data.external && data.backend">
    <v-icon icon="mdi mdi-circle" size="x-small" color="rgb(255,0,0)"></v-icon>
    Offline
  </span>
  <v-overlay v-else-if="!data.backend" v-model="offline" class="align-center justify-center">
    <v-sheet :elevation="24" rounded>
      <v-empty-state>
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
import { getPing } from '@/services/api';
import { computed, onMounted, reactive, watch } from 'vue';

const data = reactive({
  external: null,
  backend: null
});

const emit = defineEmits(['reconnect']);

const ping = () => {
  getPing()
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

watch(offline, (newv, oldv) => {
  if (newv != oldv && newv === false) emit('reconnect');
});
</script>
