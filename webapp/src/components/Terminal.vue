<template>
  <TerminalTopBar :terminalId="terminals.current.id" :locked="terminals.current.lock" @term-action="treatAction" />

  <div class="px-1 pt-1 pb-6 bg-black" ref="terminal" style="height: 90%;"></div>

  <v-dialog v-model="data.reconnect" max-width="500" contained persistent>
    <v-card v-if="data.message !== null" :title="data.message.title">
      <template v-slot:prepend>
        <v-icon :icon="'$' + data.message.type" :color="data.message.type"></v-icon>
      </template>
      <v-card-text>
        {{ data.message.text }}
      </v-card-text>
      <v-card-actions>
        <v-btn variant="tonal" color="error" text="close" rounded
          @click="terminals.remove(terminals.current.id)"></v-btn>
        <v-spacer></v-spacer>
        <v-btn variant="tonal" color="success" text="refresh" prepend-icon="mdi mdi-reload" rounded
          @click="refresh"></v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
  <v-fab v-if="terminals.current.lock" color="infor" variant="plain" :size="128" location="center center" absolute
    offset icon :disabled="true">
    <v-icon icon="mdi mdi-lock" :size="96" color="rgba(255,255,255,1)"></v-icon>
  </v-fab>
</template>

<script setup>
import { shell, ssh } from '@/services/websocket';
import { Terminal } from '@xterm/xterm';
import { onMounted, reactive, ref, watch } from 'vue';
import { FitAddon } from 'xterm-addon-fit';
import { WebLinksAddon } from 'xterm-addon-web-links';
import TerminalTopBar from './TerminalTopBar.vue';
import { ACTIONS } from '@/constants/terminalActions';
import { useTerminalsStore } from '@/stores/terminals';

const terminals = useTerminalsStore();

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

  if (terminals.current.sshSessionId) {
    ws.value = ssh(terminals.current.sshSessionId);
  } else if (terminals.current.shellPath) {
    ws.value = shell(terminals.current.shellPath);
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
    // xTerm.value.write('\r\n\x1b[31m*** Connection closed ***\x1b[0m\r\n');
  };

  ws.value.onerror = error => {
    data.reconnect = true;
    data.message = {
      title: 'Communication error',
      text: error.message || '...',
      type: 'error'
    }
    // xTerm.value.write('\r\n\x1b[31m*** Communication error ***\x1b[0m\r\n');
  };

  data.reconnect = false;
}

const clear = () => {
  xTerm.value.clear(); // clear only frontend
  ws.value.send('\x0C'); // send Ctrl+L to backend
  xTerm.value.focus();
}

const treatAction = (action) => {
  switch (action.type) {
    case ACTIONS.CLOSE:
      terminals.remove(action.terminalId);
      break;
    case ACTIONS.LOCK:
      terminals.toggleLock(action.terminalId);
      break;
    case ACTIONS.CLEAR:
      clear();
      break;
  }
}

watch(() => terminals.current.lock, (locked) => {
  if (xTerm.value) {
    xTerm.value.options.disableStdin = locked;
    xTerm.value.options.cursorBlink = !locked;
    data.locked = locked;

    if (locked === false) xTerm.value.focus();
    // if (locked) {
    //   ws.value.send("# LOCKED\r");
    // } else {
    //   ws.value.send("# UNLOCKED\r");
    // }
  }
});

onMounted(() => {
  if (terminal.value) {
    initXterm();
    initWebSocket();

    xTerm.value.focus();
  }
});
</script>
