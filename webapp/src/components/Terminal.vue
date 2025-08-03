<template>
  <div ref="terminal"></div>
</template>

<script setup>
import { onMounted, ref } from 'vue';

const props = defineProps({
  name: {
    required: true
  }
});

const terminal = ref(null);

onMounted(() => {
  if (terminal.value) {
    const term = new Terminal();
    term.open(terminal.value);

    const ws = new WebSocket("ws://localhost:8080/ssh?name=" + props.name);
    ws.onmessage = (e) => term.write(e.data);
    term.onData(data => ws.send(data));
  }
});

</script>