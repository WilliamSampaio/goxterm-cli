<template>
  <v-alert v-if="data.message !== null" class="ma-3" :title="data.message.title" :text="data.message.text"
    :type="data.message.type" variant="tonal" density="compact" @click:close="data.message = null" closable></v-alert>
  <div class="pa-2 bg-black" ref="terminal"></div>
  <v-fab v-if="data.reconnect" color="primary" extended text="refresh" variant="tonal" prepend-icon="mdi mdi-reload"
    location="center center" @click="refresh" absolute offset></v-fab>
</template>

<script setup>
import { shell, ssh } from '@/plugins/websocket';
import { Terminal } from '@xterm/xterm';
import { onMounted, reactive, ref, watch } from 'vue';
import { FitAddon } from 'xterm-addon-fit';
import { WebLinksAddon } from 'xterm-addon-web-links';

const props = defineProps({
  sshSessionId: null,
  shellPath: null,
  lock: false
});

const terminal = ref(null);
const xTerm = ref(null);
const ws = ref(null);

const data = reactive({
  message: null,
  reconnect: false
});

const refresh = () => {
  data.message = null;
  initXterm();
  initWebSocket();
}

const initXterm = () => {
  if (xTerm.value) {
    terminal.value.innerHTML = '';
  };

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
}

const initWebSocket = () => {

  if (props.sshSessionId) {
    ws.value = ssh(props.sshSessionId);
  } else if (props.shellPath) {
    ws.value = shell(props.shellPath);
  } else {
    return;
  }

  ws.value.onmessage = (e) => xTerm.value.write(e.data);
  xTerm.value.onData(data => ws.value.send(data));

  ws.value.onclose = event => {
    data.reconnect = true;
    data.message = {
      title: 'Connection closed',
      text: `${event.code} (${event.type}): ${event.reason || '...'}`,
      type: 'error'
    }
    xTerm.value.write('\r\n\x1b[31m*** Connection closed ***\x1b[0m\r\n');
  };

  ws.value.onerror = error => {
    data.reconnect = true;
    data.message = {
      title: 'Communication error',
      text: error.message || '...',
      type: 'error'
    }
    xTerm.value.write('\r\n\x1b[31m*** Communication error ***\x1b[0m\r\n');
  };

  data.reconnect = false;
}

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
    initXterm();
    initWebSocket();
  }
});
</script>
