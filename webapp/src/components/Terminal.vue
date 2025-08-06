<template>
  <div class="pa-2 bg-black" ref="terminal"></div>
</template>

<script setup>
import { BACKEND_HOST } from '@/utils';
import { Terminal } from '@xterm/xterm';
import { onMounted, ref } from 'vue';
import { FitAddon } from 'xterm-addon-fit';
import { WebLinksAddon } from 'xterm-addon-web-links';

const props = defineProps({
  id: {
    required: true
  }
});

const terminal = ref(null);

onMounted(() => {
  if (terminal.value) {
    const term = new Terminal({
      cursorBlink: true,
      fontFamily: 'monospace',
    });

    const fitAddon = new FitAddon();
    term.loadAddon(fitAddon);

    const linkAddon = new WebLinksAddon();
    term.loadAddon(linkAddon);

    term.open(terminal.value);

    fitAddon.fit();

    term.write('<< WELCOME TO GOXTERM! >>\r\n');

    const ws = new WebSocket(`ws://${BACKEND_HOST}/ws/ssh?id=${props.id}`);
    ws.onmessage = (e) => term.write(e.data);
    term.onData(data => ws.send(data));
  }
});
</script>
