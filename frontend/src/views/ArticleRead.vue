<template>
  <div class="max-w-4xl mx-auto px-4 py-10">
    <!-- Loading -->
    <div v-if="loading" class="py-20 flex flex-col items-center">
      <div class="h-10 w-10 border-4 border-indigo-500 border-t-transparent rounded-full animate-spin mb-6"></div>
      <p class="text-sm text-gray-500">Загрузка статьи...</p>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="py-24 text-center">
      <p class="text-2xl font-semibold text-gray-800 mb-2">Ошибка</p>
      <p class="text-gray-500 mb-8">{{ error }}</p>
      <RouterLink
        to="/"
        class="inline-block px-5 py-2.5 rounded bg-indigo-600 text-white text-sm font-medium hover:bg-indigo-700 transition"
      >На главную</RouterLink>
    </div>

    <!-- Article -->
    <article v-else>
      <h1 class="text-3xl md:text-4xl font-bold tracking-tight text-gray-900 mb-6">
        {{ article.title }}
      </h1>

      <div class="flex flex-wrap items-center gap-3 text-sm text-gray-500 mb-8">
        <span class="inline-flex items-center gap-1">
          <svg class="w-4 h-4 text-indigo-500" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M5.5 21a7.5 7.5 0 0 1 13 0M12 11a4 4 0 1 0 0-8 4 4 0 0 0 0 8Z"/>
          </svg>
          {{ article.author }}
        </span>
        <span class="h-1 w-1 rounded-full bg-gray-300"></span>
        <span>ID: {{ article.id }}</span>
        <template v-if="article.media?.length">
          <span class="h-1 w-1 rounded-full bg-gray-300"></span>
          <span>{{ article.media.length }} медиа</span>
        </template>
        <RouterLink
          v-if="canEdit"
          :to="{ name: 'create', query: { edit: article?.id }}"
          class="ml-auto text-xs px-3 py-1.5 rounded border border-indigo-500 text-indigo-600 hover:bg-indigo-50 transition"
        >
          Редактировать
        </RouterLink>
      </div>

      <!-- Media Gallery -->
      <div v-if="article.media?.length" class="mb-10">
        <div class="grid gap-4 sm:grid-cols-2">
          <div
            v-for="(m,i) in article.media"
            :key="i"
            class="group relative overflow-hidden rounded-lg border bg-white shadow-sm"
          >
            <img
              :src="m"
              :alt="`media ${i+1}`"
              class="h-56 w-full object-cover transition group-hover:scale-105"
              loading="lazy"
              @error="onImgErr"
            />
            <div class="absolute inset-0 pointer-events-none shadow-inner"></div>
          </div>
        </div>
      </div>

      <!-- Markdown Content -->
      <div class="prose max-w-none bg-white/70 backdrop-blur rounded-lg border shadow-sm p-6">
        <MdPreview :id="mdId" :modelValue="article.content" />
      </div>
    </article>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import axios from 'axios'
import { MdPreview } from 'md-editor-v3'
import 'md-editor-v3/lib/preview.css'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const id = route.params.id
const mdId = 'article-view'
const loading = ref(true)
const error = ref('')
const article = ref(null)

const auth = useAuthStore?.()
const canEdit = computed(() => {
  if (!article.value || !auth?.user) return false
  return auth.user.username === article.value.author
})

const baseURL = import.meta.env.VITE_BACKEND_URL || ''
async function loadArticle() {
  loading.value = true
  error.value = ''
  try {
    const { data } = await axios.get(`${baseURL}/articles/${id}`)
    article.value = data
  } catch (e) {
    if (e.response?.status === 404) error.value = 'Статья не найдена'
    else error.value = 'Не удалось загрузить статью'
  } finally {
    loading.value = false
  }
}

function onImgErr(ev) {
  ev.target.classList.add('opacity-40')
  ev.target.alt = 'media недоступно'
}

onMounted(loadArticle)
</script>