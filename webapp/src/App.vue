<template>
  <v-app id="inspire">
    <v-toolbar color="primary">
      <v-navigation-drawer v-model="drawer">
        <v-list lines="one">
          <v-list-subheader>SESSIONS:</v-list-subheader>
          <v-divider></v-divider>
          <v-list-item v-for="item in data.sessions" :key="item.id" @click="connect(item)">
            <v-list-item-title>{{ item.name }}</v-list-item-title>
            <v-list-item-subtitle>{{ item.host }}:{{ item.port }}</v-list-item-subtitle>
            <template v-slot:append>
              <v-icon icon="mdi mdi-open-in-new"></v-icon>
            </template>
          </v-list-item>
        </v-list>
      </v-navigation-drawer>

      <v-app-bar>
        <v-app-bar-nav-icon @click="drawer = !drawer"></v-app-bar-nav-icon>
        <v-app-bar-title>GoXterm</v-app-bar-title>
        <v-btn icon="mdi-magnify"></v-btn>
        <template v-slot:extension>
          <v-tabs v-model="data.tab" align-tabs="title">
            <v-tab v-for="t in data.tabs" :key="t.id" :text="t.name" :value="t"></v-tab>
          </v-tabs>
        </template>
      </v-app-bar>
    </v-toolbar>

    <v-main>
      <v-card>
        <v-tabs-window v-model="data.tab">
          <v-tabs-window-item v-for="t in data.tabs" :key="t.id" :value="t">
            <Terminal :id="t.sessionId" />
          </v-tabs-window-item>
        </v-tabs-window>
      </v-card>
    </v-main>
  </v-app>
</template>

<script setup>
import axios from 'axios'
import { onMounted, reactive, ref, watch } from 'vue'
import Terminal from './components/Terminal.vue';
import { BACKEND_HOST } from './utils';

const drawer = ref(null);

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
      data.sessions = response.data || [];
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
</script>
