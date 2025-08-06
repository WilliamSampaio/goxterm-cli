<template>
  <v-app id="inspire">
    <v-navigation-drawer permanent>
      <v-list-item v-for="item in data.sessions" :key="item.id" @click="connect(item)">
        <v-list-item-title>{{ item.name }}</v-list-item-title>
        <v-list-item-subtitle>{{ item.host }}:{{ item.port }}</v-list-item-subtitle>
        <template v-slot:append>
          <v-icon icon="mdi mdi-open-in-new"></v-icon>
        </template>
      </v-list-item>
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
        <v-btn href="https://github.com/WilliamSampaio/goxterm-cli" target="_blank" rel="noopener noreferrer"
          append-icon="mdi mdi-github" variant="text" block>
          GitHub
        </v-btn>
      </template>
    </v-navigation-drawer>
    <v-app-bar elevation="0">
      <v-app-bar-title>GoXterm</v-app-bar-title>
      <v-btn icon="mdi-magnify"></v-btn>
    </v-app-bar>
    <v-main>
      <v-card v-if="data.tabs.length > 0">
        <v-tabs v-model="data.tab">
          <v-tab v-for="t in data.tabs" :key="t.id" :text="t.name" :value="t"></v-tab>
        </v-tabs>
        <v-tabs-window v-model="data.tab">
          <v-tabs-window-item v-for="t in data.tabs" :key="t.id" :value="t">
            <v-sheet class="pa-2 text-center">
              <v-btn variant="tonal" color="error" @click="" @dblclick="closeTab()" icon>
                <v-icon icon="mdi-close"></v-icon>
                <v-tooltip activator="parent" location="bottom">Double Click</v-tooltip>
              </v-btn>
            </v-sheet>
            <Terminal :id="t.sessionId" />
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
import { onMounted, reactive, ref, watch } from 'vue'
import Terminal from './components/Terminal.vue';
import { BACKEND_HOST } from './utils';

const length = ref(15)
const tab = ref(null)

watch(length, val => {
  tab.value = val - 1
})

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
  const newTabId = data.tabs.length + 1;
  data.tabs.push({
    id: newTabId,
    sessionId: item.id,
    name: item.name
  });
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
</script>
