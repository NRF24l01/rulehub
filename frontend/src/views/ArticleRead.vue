<template>
  <div
    class="w-full mx-auto
           px-3 sm:px-4 lg:px-8 2xl:px-14
           py-8 sm:py-10
           max-w-3xl md:max-w-4xl lg:max-w-5xl xl:max-w-6xl 2xl:max-w-7xl"
  >
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
    <div v-else class="flex flex-col lg:flex-row gap-8 min-h-screen">
      <article class="flex-1">
        <!-- Markdown Content -->
        <div>
          <MdPreview :id="mdId" :modelValue="article.content" />
        </div>
      </article>

      <!-- Каталог справа -->
      <aside
        class="hidden lg:block w-72 flex-shrink-0"
        style="position:sticky;top:2rem;align-self:flex-start;"
      >
        <MdCatalog :editorId="mdId" :scrollElement="scrollElement" />
      </aside>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import axios from 'axios'
import { MdPreview, MdCatalog } from 'md-editor-v3'
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

// Добавим ссылку на прокручиваемый элемент для MdCatalog
const scrollElement = document.documentElement

onMounted(loadArticle)

// Обновляем заголовок вкладки при изменении article.title
watch(
  () => article.value?.title,
  (title) => {
    if (title) {
      document.title = title
    }
  }
)
</script>

<style>
/* Применяем стили ко всем спискам внутри компонента */
ul, ol {
  list-style-position: inside !important;
  margin-left: 1.5em !important;
  padding-left: 0 !important;
  list-style: revert !important;
}

ul {
  list-style-type: disc !important;
}

ol {
  list-style-type: decimal !important;
}

li {
  display: list-item !important;
  margin-bottom: 0.25em !important;
}
</style>