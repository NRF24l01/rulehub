<template>
  <div class="flex flex-col space-y-4 mt-4 max-w-xl mx-auto">
    <input
      v-model="articleName"
      type="text"
      placeholder="Название статьи"
      class="border rounded px-3 py-2 border-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
    />
  </div>
  <Editor v-model:content="content" v-model:images="images" />
  <div class="flex flex-col space-y-4 mt-4 max-w-xl mx-auto">
    <button
      @click="handleClick"
      class="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 transition"
    >
      Сохранить
    </button>
  </div>
  <div v-if="error !== ''" class="text-red-500 text-sm mt-2">{{ error }}</div>
</template>

<script setup>
import Editor from '@/components/Editor.vue';
import { ref, watch } from 'vue';
import api from '../axios'

const articleName = ref('');
const content = ref('# Пример markdown')
const images = ref([]);
const error = ref('');

async function handleClick() {
  error.value = '';
  try {
    const res = await api.post(
      import.meta.env.VITE_BACKEND_URL + "/articles/",
      {
        name: articleName.value,
        content: content.value,
        images: images.value
      }
    );
    if (res.status !== 201) {
      error.value = res.data.error || res.data.message || 'Ошибка при создании статьи';
      return;
    }
    error.value = res.data.id;
  } catch (e) {
    // axios error: e.response содержит ответ сервера
    error.value =
      e?.response?.data?.error ||
      e?.response?.data?.message ||
      'Ошибка при создании статьи';
  }
}
</script>