<template>
  <v-sheet class="text-center">
    <!-- close button -->
    <v-btn class="ma-1" variant="flat" color="error" size="x-small" @dblclick="$emit('term-action', {
      type: ACTIONS.CLOSE,
      terminalId: terminalId
    })" icon>
      <v-icon icon="mdi-close"></v-icon>
      <v-tooltip activator="parent" location="bottom">Close (Double Click)</v-tooltip>
    </v-btn>

    <!-- lock button -->
    <v-btn class="ma-1" variant="flat" color="warning" size="x-small" @click="emitLock" icon>
      <v-icon :icon="lockIcon"></v-icon>
      <v-tooltip activator="parent" location="bottom">{{ locked ? 'Unlock' : 'Lock' }}</v-tooltip>
    </v-btn>

    <!-- clear button -->
    <v-btn class="ma-1" variant="flat" color="success" size="x-small" @click="$emit('term-action', {
      type: ACTIONS.CLEAR,
      terminalId: terminalId
    })" icon>
      <v-icon icon="mdi mdi-broom"></v-icon>
      <v-tooltip activator="parent" location="bottom">Clear</v-tooltip>
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

const emit = defineEmits(['term-action']);

const emitLock = (e) => {
  emit('term-action', {
    type: ACTIONS.LOCK,
    terminalId: props.terminalId
  });

  e.target.blur();
}

const lockIcon = computed(() => {
  if (props.locked === true) {
    return 'mdi mdi-lock';
  }
  return 'mdi mdi-lock-open-variant';
});
</script>
