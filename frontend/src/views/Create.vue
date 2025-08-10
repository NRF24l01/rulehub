<template>
  <div class="flex flex-col space-y-4 mt-4 max-w-xl mx-auto">
    <input
      v-model="articleName"
      type="text"
      placeholder="Название статьи"
      class="border rounded px-3 py-2 border-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
    />
  </div>
  <div class="flex flex-col h-[80vh] min-h-[300px] w-full">
    <Editor v-model:content="content" v-model:images="images" class="flex-1 min-h-0" />
  </div>
  <div v-if="error !== ''" class="flex justify-center mt-4">
    <div class="shadow-md rounded px-6 py-4 text-red-600 text-sm max-w-md w-full text-center border border-red-500 bg-red-100">
      {{ error }}
    </div>
  </div>
  <div class="flex flex-col space-y-4 mt-4 max-w-xl mx-auto">
    <button
      @click="handleClick"
      class="bg-blue-600 text-white px-4 py-2 rounded transition"
      :class="[
      (error !== '' || articleName.trim() === '' || content.trim() === '')
        ? 'opacity-50 cursor-not-allowed'
        : 'hover:bg-blue-700'
      ]"
      :disabled="error !== '' || articleName.trim() === '' || content.trim() === ''"
    >
      Сохранить
    </button>
  </div>
</template>

<script setup>
import Editor from '@/components/Editor.vue';
import { ref, watch } from 'vue';
import api from '../axios'
import { useRouter } from 'vue-router'

const articleName = ref('');
const content = ref('# Пример markdown')
const images = ref([]);
const error = ref('');

const router = useRouter();

watch([articleName, content, images], () => {
  if (error.value !== '') {
    error.value = '';
  }
});

function validate() {
  if (articleName.value.trim() === '') {
    error.value = 'Название статьи не может быть пустым';
    return false;
  }
  if (content.value.trim() === '') {
    error.value = 'Содержание статьи не может быть пустым';
    return false;
  }
  error.value = '';
  return true;
}

async function handleClick() {
  error.value = '';
  if (!validate()) {
    return;
  }
  try {
    const res = await api.post(
      import.meta.env.VITE_BACKEND_URL + "/articles/",
      {
        title: articleName.value,
        content: content.value,
        media: images.value.map(img => img.fileID)
      }
    );
    if (res.status !== 201) {
      error.value = res.data.error || res.data.message || 'Ошибка при создании статьи';
      return;
    }
    error.value = res.data.id;
    router.push({ name: 'article', params: { id: res.data.id } });
  } catch (e) {
    error.value =
      e?.response?.data?.error ||
      e?.response?.data?.message ||
      'Ошибка при создании статьи';
    console.log(e);
  }
}
</script>