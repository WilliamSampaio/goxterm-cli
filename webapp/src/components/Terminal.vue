<template>
  <div class="pa-2 bg-black" ref="terminal"></div>
</template>

<script setup>
import { BACKEND_HOST } from '@/utils';
import { Terminal } from '@xterm/xterm';
import { onMounted, ref, watch } from 'vue';
import { FitAddon } from 'xterm-addon-fit';
import { WebLinksAddon } from 'xterm-addon-web-links';

const props = defineProps({
  id: {
    required: true
  },
  lock: false
});

const terminal = ref(null);
const xTerm = ref(null);
const ws = ref(null);

watch(() => props.lock, (locked) => {
  if (xTerm.value) {
    xTerm.value.options.disableStdin = locked;
    xTerm.value.options.cursorBlink = !locked;
    if (locked) {
      ws.value.send("# LOCKED\r");
    } else {
      ws.value.send("# UNLOCKED\r");
    }
  }
});

onMounted(() => {
  if (terminal.value) {
    xTerm.value = new Terminal({
      cursorBlink: true,
      fontFamily: 'monospace',
    });

    const fitAddon = new FitAddon();
    xTerm.value.loadAddon(fitAddon);

    const linkAddon = new WebLinksAddon();
    xTerm.value.loadAddon(linkAddon);

    xTerm.value.open(terminal.value);

    fitAddon.fit();

    // xTerm.value.write('<< WELCOME TO GOXTERM! >>\r\n');

    ws.value = new WebSocket(`ws://${BACKEND_HOST}/ws/ssh?id=${props.id}`);
    ws.value.onmessage = (e) => xTerm.value.write(e.data);
    xTerm.value.onData(data => ws.value.send(data));
  }
});
</script>
