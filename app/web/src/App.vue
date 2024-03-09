<!-- eslint-disable vue/max-attributes-per-line -->
<template>
  <v-app>
    <v-app-bar color="primary" dark>
      <v-icon icon="mdi-checkbox-marked-circle-plus-outline" end />
      <v-app-bar-title>TODO Application</v-app-bar-title>
      <v-spacer />
      <v-btn icon :href="docsURL" target="_blank">
        <v-icon icon="mdi-help" />
      </v-btn>
    </v-app-bar>
    <v-container>
      <v-row justify="center">
        <v-col>
          <v-card>
            <AddTodoForm @add="onAdd" />
            <TodoList :todos="todos" @delete="onDelete" @complete="onComplete" />
          </v-card>
        </v-col>
      </v-row>
      <v-snackbar v-model="errorSnackbar.show" color="red" top>
        {{ errorSnackbar.message }}
        <template #actions>
          <v-btn @click="errorSnackbar.show = false">
            <v-icon icon="mdi-close" />
          </v-btn>
        </template>
      </v-snackbar>
    </v-container>
  </v-app>
</template>

<script setup>
import TodoList from './components/TodoList.vue';
import AddTodoForm from './components/AddTodoForm.vue';

import { ref, onMounted } from 'vue';
import axios from 'axios';

const todos = ref([]);

const errorSnackbar = ref({ show: false, message: '' });

const apiBaseURL = '/api/v1/todos';
const docsURL = ref('/swagger');

axios.interceptors.response.use(
  (response) => response,
  (error) => {
    showErrorSnackbar(error.message || 'A server-side error occurred');
    return Promise.reject(error);
  }
);

const showErrorSnackbar = (message) => {
  errorSnackbar.value = { show: true, message };
}

const fetchTodos = async () => {
  try {
    const response = await axios.get(apiBaseURL);
    todos.value = response.data;
  } catch (error) {
    console.error('Error fetching todos:', error);
  }
}

const onAdd = async (newTodo) => {
  try {
    const response = await axios.post(apiBaseURL, newTodo);
    todos.value.push(response.data);
  } catch (error) {
    console.error('Error adding todo:', error);
  }
}

const onDelete = async (id) => {
  try {
    await axios.delete(`${apiBaseURL}/${id}`);
    todos.value = todos.value.filter((todo) => todo.id !== id);
  } catch (error) {
    console.error('Error deleting todo:', error);
  }
}

const onComplete = async (id, completed) => {
  try {
    await axios.put(`${apiBaseURL}/${id}`, {completed});
    todos.value.find((todo) => todo.id == id).completed = completed
  } catch (error) {
    console.error('Error marking todo as complete:', error);
  }
}

onMounted(fetchTodos);
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
