<template>
  <v-tabs-window v-model="terminals.current" style="height: 100%;">
    <v-tabs-window-item v-for="(t, i) in terminals.items" :key="i" :value="t" style="height: 100%;">
      <v-sheet style="height: 100%;">
        <Terminal />
      </v-sheet>
    </v-tabs-window-item>
  </v-tabs-window>
</template>

<script setup>
import Terminal from './Terminal.vue';
import { useTerminalsStore } from '@/stores/terminals';

const terminals = useTerminalsStore();
</script>
