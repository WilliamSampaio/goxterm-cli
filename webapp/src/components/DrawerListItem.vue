<template>
  <v-menu location="end">
    <template v-slot:activator="{ props: menu }">
      <v-list-item :key="session.id" v-bind="mergeProps(menu, tooltip)">
        <v-list-item-title>{{ session.name }}</v-list-item-title>
        <v-list-item-subtitle>{{ session.host }}:{{ session.port }}</v-list-item-subtitle>
        <template v-slot:append>
          <v-icon icon="mdi mdi-open-in-new"></v-icon>
        </template>
      </v-list-item>
    </template>
    <v-list lines="one">
      <v-list-subheader line="one">Users:</v-list-subheader>
      <v-list-item @click="$emit('selected', session)">
        <v-list-item-title>{{ session.user }}</v-list-item-title>
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script setup>
import { mergeProps } from 'vue'

const props = defineProps({
  session: {
    required: true
  }
});
</script>