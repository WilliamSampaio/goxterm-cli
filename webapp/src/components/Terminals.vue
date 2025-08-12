<template>
  <v-tabs-window v-model="terminals.current">
    <v-tabs-window-item v-for="(t, i) in terminals.items" :key="i" :value="t">
      <Terminal />
    </v-tabs-window-item>
  </v-tabs-window>
</template>

<script setup>
import Terminal from './Terminal.vue';
import { useTerminalsStore } from '@/stores/terminals';

const terminals = useTerminalsStore();
</script>
