<!-- eslint-disable vue/max-attributes-per-line -->
<template>
  <v-list subheader two-line flat>
    <TodoItem v-for="todo in todos" :key="todo.id" :todo="todo" @delete="onDelete" @complete="onComplete" />
  </v-list>
</template>

<script setup>
import TodoItem from './TodoItem.vue';

defineProps({
  todos: {
    type: Array,
    required: true
  }
});

const emit = defineEmits(['delete', 'complete'])

const onDelete = (id) => {
  emit('delete', id);
};

const onComplete = (id, completed) => {
  emit('complete', id, completed);
};
</script>
