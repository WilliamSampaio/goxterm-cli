<template>
  <v-dialog max-width="500">
    <v-card title="SSH Quick Access">
      <v-card-text>
        <v-form ref="form">
          <v-text-field v-model="data.connection" label="Connection (user@host:port)" :rules="data.connectionRules"
            clearable></v-text-field>
          <v-text-field v-model="data.password" label="Password"
            :append-inner-icon="data.passwordVisible ? 'mdi-eye-off' : 'mdi-eye'"
            :type="data.passwordVisible ? 'text' : 'password'" placeholder="Enter your password"
            @click:append-inner="data.passwordVisible = !data.passwordVisible"
            :rules="data.passwordRules"></v-text-field>
        </v-form>
      </v-card-text>
      <v-card-actions>
        <v-btn class="px-5" variant="tonal" color="error" rounded @click="emit('close')">
          Cancel
        </v-btn>
        <v-spacer></v-spacer>
        <v-btn class="px-5" variant="tonal" color="success" rounded append-icon="mdi mdi-television-shimmer"
          @click="submit">
          Access
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup>
import { CONNECTION_STRING_REGEX } from '@/constants/validation';
import { useTerminalsStore } from '@/stores/terminals';
import { reactive, ref } from 'vue';

const terminals = useTerminalsStore();

const form = ref(null);

const emit = defineEmits(['close']);

const data = reactive({
  connection: null,
  password: null,
  connectionRules: [
    v => { return v ? true : 'Connection string is required.' },
    v => { return String(v).match(CONNECTION_STRING_REGEX) ? true : 'Connection string is invalid.' },
  ],
  passwordRules: [
    v => { return v ? true : 'Password is required.' },
  ],
  passwordVisible: false
});

const submit = () => {
  form.value.validate().then((v) => {
    if (v.valid) {
      terminals.add(null, null, { connection: data.connection, password: data.password }, '(quick)', false);
      emit('close');
    }
  });
}
</script>
