<template>
  <v-app>
    <v-footer height="40" app>
      <Ping @reconnect="initialize" />
      <v-spacer></v-spacer>
      <p class="text-medium-emphasis">v{{ data.info?.version }}</p>
      <v-btn href="https://github.com/WilliamSampaio/goxterm-cli" target="_blank" rel="noopener noreferrer"
        icon="mdi mdi-github" size="small" variant="plain">
      </v-btn>
    </v-footer>

    <v-navigation-drawer v-model="drawer">
      <DrawerListItem v-for="item in data.sessions" :session="item" @selected="connect(item)" />
      <template v-slot:prepend>
        <div class="pa-2">
          <v-btn color="primary" variant="tonal" block>
            Add New
          </v-btn>
        </div>
        <v-divider></v-divider>
      </template>
      <template v-slot:append>
        <v-divider></v-divider>
      </template>
    </v-navigation-drawer>

    <v-app-bar elevation="0">
      <v-app-bar-nav-icon @click="drawer = !drawer"></v-app-bar-nav-icon>
      <v-app-bar-title>GoXterm</v-app-bar-title>
      <v-btn class="mx-2" variant="tonal" color="primary" rounded="xl" prepend-icon="mdi mdi-lightning-bolt"
        @click="data.sshQuickAccessDialog = true">
        SSH
      </v-btn>
      <v-btn v-for="shell in data.info?.shells" class="mx-2" variant="tonal" color="success" rounded="xl"
        prepend-icon="mdi mdi-plus" @click="connect(shell)">
        {{ shell.bin }}{{ shell.default ? ' (default)' : '' }}
      </v-btn>
      <v-btn class="mx-2" icon>
        <v-icon icon="mdi-magnify"></v-icon>
        <v-tooltip activator="parent" location="bottom">{{ '[F1]' }}</v-tooltip>
      </v-btn>
      <!-- <v-divider vertical inset></v-divider>
      <v-btn class="mx-1" icon="mdi-cog"></v-btn> -->
    </v-app-bar>

    <v-main>
      <v-card v-if="terminals.items.length > 0" :rounded="false" style="height: 100%;">
        <v-tabs v-model="terminals.current" density="compact">
          <v-tab v-for="(t, i) in terminals.items" :key="i" :value="t">
            {{ t.name }}
          </v-tab>
        </v-tabs>
        <v-divider></v-divider>
        <Terminals />
      </v-card>
      <v-empty-state v-else headline="Whoops, 404" title="Page not found"
        text="The page you were looking for does not exist" icon="mdi mdi-console"></v-empty-state>

      <SshQuickAccessDialog v-model="data.sshQuickAccessDialog" @close="data.sshQuickAccessDialog = false" />
    </v-main>
  </v-app>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import Terminals from './components/Terminals.vue';
import DrawerListItem from './components/DrawerListItem.vue';
import Ping from './components/Ping.vue';
import { getInfo, getSshSessions } from './services/api';
import { useTerminalsStore } from './stores/terminals';
import SshQuickAccessDialog from './components/SshQuickAccessDialog.vue';

const terminals = useTerminalsStore();

const drawer = ref(null);

const data = reactive({
  drawer: null,
  info: null,
  sessions: [],
  sshQuickAccessDialog: false
});

onMounted(() => {
  initialize();
});

const initialize = () => {
  getInfo()
    .then(response => {
      if (response.headers['content-type'] === "application/json") {
        data.info = response.data;
      }
    })
    .catch(error => {
      console.error('Error fetching data:', error);
    });

  getSshSessions()
    .then(response => {
      if (response.headers['content-type'] === "application/json") {
        data.sessions = response.data || [];
      }
    })
    .catch(error => {
      console.error('Error fetching data:', error);
    });
}

const connect = (item) => {
  terminals.add(item.id, item.path, null, item.name || item.bin);
}
</script>
