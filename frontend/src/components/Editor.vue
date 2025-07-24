<template>
  <div class="space-y-4 p-4">

    <!-- Markdown редактор -->
    <MdEditor v-model="markdownText" height="300px" language="en-US"/>

    <!-- Просмотр markdown -->
    <div class="border p-4 bg-gray-50 rounded">
      <h2 class="font-bold text-lg mb-2">Предпросмотр</h2>
      <MdPreview :modelValue="markdownText" />
    </div>

    <!-- Загрузка изображений -->
    <div class="space-y-2">
      <h2 class="font-bold text-lg">Фотографии</h2>
      <input type="file" multiple accept="image/*" @change="handleFiles" />

      <ul>
        <li v-for="(file, idx) in imageFiles" :key="idx" class="flex items-center gap-2">
          <img :src="file.preview" alt="" class="w-16 h-16 object-cover rounded border" />
          <span>{{ file.name }}</span>
        </li>
      </ul>
    </div>

    <!-- Кнопка отправки -->
    <button
      class="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
      @click="submitForm"
    >
      Отправить
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { MdEditor, MdPreview } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'

const markdownText = ref('# Пример markdown')

interface FileWithPreview {
  file: File
  name: string
  preview: string
}

const imageFiles = ref<FileWithPreview[]>([])

function handleFiles(event: Event) {
  const input = event.target as HTMLInputElement
  if (!input.files) return

  for (const file of Array.from(input.files)) {
    if (!file.type.startsWith('image/')) continue

    const reader = new FileReader()
    reader.onload = () => {
      imageFiles.value.push({
        file,
        name: file.name,
        preview: reader.result as string,
      })
    }
    reader.readAsDataURL(file)
  }
}

async function submitForm() {
  const formData = new FormData()
  formData.append('markdown', markdownText.value)

  imageFiles.value.forEach((f, idx) => {
    formData.append(`images[${idx}]`, f.file, f.name)
  })

  try {
    const res = await fetch(import.meta.env.VITE_BACKEND_URL + '/upload', {
      method: 'POST',
      body: formData,
    })

    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    alert('Успешно отправлено')
  } catch (err) {
    console.error(err)
    alert('Ошибка при отправке')
  }
}
</script>

<style scoped>
ul {
  list-style: none;
  padding: 0;
}
</style>
