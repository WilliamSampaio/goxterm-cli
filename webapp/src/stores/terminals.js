import { defineStore } from "pinia";
import { ref } from "vue";

export const useTerminalsStore = defineStore('terminals', () => {
  const items = ref([]);
  const current = ref(null);

  function add(sshSessionId, shellPath, sshConnection, name, lock = false) {
    const terminal = {
      id: Date.now(),
      sshSessionId: sshSessionId || null,
      shellPath: shellPath || null,
      sshConnection: sshConnection || { connection: null, password: null },
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
}, {
  persist: true
});
