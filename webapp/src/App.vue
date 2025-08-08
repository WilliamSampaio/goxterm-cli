<template>
  <v-app>
    <v-navigation-drawer permanent>
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
      <v-app-bar-title>GoXterm</v-app-bar-title>
      <v-btn icon="mdi-magnify"></v-btn>
    </v-app-bar>

    <v-footer height="40" app>
      <v-spacer></v-spacer>
      2025 <v-divider class="mx-1" vertical></v-divider> v0.2.0
      <v-btn href="https://github.com/WilliamSampaio/goxterm-cli" target="_blank" rel="noopener noreferrer"
        icon="mdi mdi-github" variant="plain">
      </v-btn>
    </v-footer>

    <v-main>
      <v-card v-if="data.tabs.length > 0" :rounded="false">
        <v-tabs v-model="data.tab" density="compact">
          <v-tab v-for="(t, i) in data.tabs" :key="i" :value="t">
            {{ t.name }}
          </v-tab>
        </v-tabs>
        <v-divider></v-divider>
        <v-tabs-window v-model="data.tab">
          <v-tabs-window-item v-for="(t, i) in data.tabs" :key="i" :value="t">
            <v-sheet class="text-center">
              <v-btn class="ma-1" variant="tonal" color="error" size="x-small" @dblclick="closeTab()" icon>
                <v-icon icon="mdi-close"></v-icon>
                <v-tooltip activator="parent" location="bottom">Double Click</v-tooltip>
              </v-btn>
              <v-btn class="ma-1" variant="tonal" color="warning" :icon="lockIcon(t.lock)" size="x-small"
                @click="lockTab(i)"></v-btn>
            </v-sheet>
            <Terminal :id="t.sessionId" :lock="t.lock" />
          </v-tabs-window-item>
        </v-tabs-window>
      </v-card>
      <v-empty-state v-else headline="Whoops, 404" title="Page not found"
        text="The page you were looking for does not exist" icon="mdi mdi-console"></v-empty-state>
    </v-main>
  </v-app>
</template>

<script setup>
import axios from 'axios'
import { onMounted, reactive } from 'vue'
import { BACKEND_HOST } from './utils';
import Terminal from './components/Terminal.vue';
import DrawerListItem from './components/DrawerListItem.vue';

const data = reactive({
  drawer: null,
  sessions: [],
  tabs: [],
  tab: null,
});

onMounted(() => {
  axios.get(`http://${BACKEND_HOST}/api/ssh/sessions`)
    .then(response => {
      if (response.headers['content-type'] === "application/json") {
        data.sessions = response.data || [];
      }
    })
    .catch(error => {
      console.error('Error fetching data:', error);
    });
});

const connect = (item) => {
  const len = data.tabs.push({
    id: data.tabs.length + 1,
    sessionId: item.id,
    name: item.name,
    lock: false
  });
  data.tab = data.tabs[len - 1];
}

const closeTab = () => {
  const index = data.tabs.indexOf(data.tab);
  if (data.tabs[index] !== undefined) {
    data.tabs.splice(index, 1);
    if (index > data.tabs.length - 1) {
      data.tab = data.tabs[data.tabs.length - 1];
    } else {
      data.tab = data.tabs[index];
    }
  };
}

const lockTab = (index) => {
  data.tabs[index].lock = !data.tabs[index].lock;
}

const lockIcon = (lock) => {
  if (lock === true) {
    return 'mdi mdi-lock';
  }
  return 'mdi mdi-lock-open-variant';
};
</script>
