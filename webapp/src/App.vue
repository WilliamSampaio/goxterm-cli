<template>
  <v-app>
    <v-main>
      <v-container>
        <v-select v-model="data.selectedName" label="Select" :items="data.credentials" clearable></v-select>

        <v-btn class="mb-3" @click="open">
          Connect
        </v-btn>

        <v-tabs v-model="data.tabModel">
          <v-tab v-for="tab in data.tabs" :key="tab.id">
            {{ tab.title }}
          </v-tab>
        </v-tabs>

        <v-tabs-window v-model="data.tabModel">
          <v-tabs-window-item v-for="tab in data.tabs" :key="tab.id">
            <!-- Content for {{ tab.title }} -->
            <v-card flat>
              <v-card-text>
                <Terminal :name="tab.title" />
              </v-card-text>
            </v-card>
          </v-tabs-window-item>
        </v-tabs-window>
      </v-container>
    </v-main>
  </v-app>
</template>

<script setup>
import axios from 'axios'
import { onMounted, reactive, ref, watch } from 'vue'
import Terminal from './components/Terminal.vue';

const length = ref(15)
const tab = ref(null)

watch(length, val => {
  tab.value = val - 1
})

const data = reactive({
  credentials: [],
  tabs: [],
  tabModel: null,
  selectedName: null
});

onMounted(() => {
  axios.get('http://localhost:8080/api/credentials')
    .then(response => {
      console.log('Data fetched:', response.data);
      data.credentials = response.data || [];
    })
    .catch(error => {
      console.error('Error fetching data:', error);
    });
});

const open = () => {
  const newTabId = data.tabs.length + 1;
  data.tabs.push({
    id: newTabId,
    title: data.selectedName
  });

  this.$nextTick(() => {
    data.tabModel = newTabId;
  });

  data.selectedName = null;
}
</script>
