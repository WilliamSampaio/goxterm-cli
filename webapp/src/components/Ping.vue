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
  <v-dialog v-else-if="!data.backend" v-model="offline" max-width="320" persistent>
    <v-list class="py-2" elevation="12" rounded="lg">
      <v-list-item prepend-icon="mdi mdi-connection" title="Backend is Offline...">
        <template v-slot:prepend>
          <div class="pe-4">
            <v-icon color="error" size="x-large"></v-icon>
          </div>
        </template>
        <small class="text-caption text-medium-emphasis">
          Check if GoXterm is running
        </small>
        <template v-slot:append>
          <v-progress-circular color="error" indeterminate="disable-shrink" size="16" width="2"></v-progress-circular>
        </template>
      </v-list-item>
    </v-list>
  </v-dialog>
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
