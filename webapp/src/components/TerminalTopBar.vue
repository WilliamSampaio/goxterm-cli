<template>
  <v-sheet class="text-center">
    <!-- close button -->
    <v-btn class="ma-1" variant="tonal" color="error" size="x-small" @dblclick="$emit('term-action', {
      type: ACTIONS.CLOSE,
      terminalId: terminalId
    })" icon>
      <v-icon icon="mdi-close"></v-icon>
      <v-tooltip activator="parent" location="bottom">Double Click</v-tooltip>
    </v-btn>

    <!-- lock button -->
    <v-btn class="ma-1" variant="tonal" color="warning" :icon="lockIcon" size="x-small" @click="$emit('term-action', {
      type: ACTIONS.LOCK,
      terminalId: terminalId
    })">
    </v-btn>
  </v-sheet>
</template>

<script setup>
import { ACTIONS } from '@/constants/terminalActions';
import { computed } from 'vue';

const props = defineProps({
  terminalId: null,
  locked: false
});

const lockIcon = computed(() => {
  if (props.locked === true) {
    return 'mdi mdi-lock';
  }
  return 'mdi mdi-lock-open-variant';
});
</script>
