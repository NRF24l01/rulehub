<template>
  <Editor v-model:content="content" v-model:images="images" />

  <div class="mt-4">
    <h2>Markdown Preview</h2>
    <div class="markdown-preview">{{ content }}</div>
  </div>
  <div class="mt-4">
    <h2>Uploaded Images</h2>
    <ul>
      <li v-for="image in images" :key="image.blobUrl">
        <img :src="image.blobUrl" :alt="image.name" class="w-32 h-32 object-cover" />
      </li>
    </ul>
  </div>
</template>

<script setup>
import Editor from '@/components/Editor.vue';
import { ref, watch } from 'vue';

const content = ref('# Пример markdown')
const images = ref([])
</script>