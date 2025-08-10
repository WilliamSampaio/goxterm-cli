import { defineStore } from "pinia";
import { ref } from "vue";

export const useTerminalsStore = defineStore('terminals', () => {
  const items = ref([]);
  const current = ref(null);

  const nextIndex = ref(1);

  function add(sshSessionId, shellPath, name, lock = false) {
    const terminal = {
      id: nextIndex.value++,
      sshSessionId: sshSessionId || null,
      shellPath: shellPath || null,
      name: name,
      lock: lock
    }
    items.value.push(terminal);
    current.value = terminal;
  }

  function remove(id) {
    const index = items.value.findIndex(t => t.id === id);
    items.value = items.value.filter(t => t.id !== id);
    if (current.value.id === id) {
      current.value = items.value[index - 1] || null;
    }
  }

  function toggleLock(id) {
    const index = items.value.findIndex(t => t.id === id);
    items.value[index].lock = !items.value[index].lock;
  }

  return {
    items,
    current,
    add,
    remove,
    toggleLock
  }
});