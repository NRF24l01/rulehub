<template>
  <div class="space-y-4 p-4">

    <!-- Markdown редактор -->
    <MdEditor v-model="markdownText" height="300px" language="en-US" @uploadImg="handleEditorImageInsert" />

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

<script setup>
import { ref } from 'vue'
import { MdEditor, MdPreview } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import { useAuthStore } from '../stores/auth'
import api from '../axios'

const authStore = useAuthStore()
const markdownText = ref('# Пример markdown')

const images = ref([]) // { file, name, blobUrl }

function handleEditorImageInsert(files, callback) {
  const uploadImg = async (files, callback) => {
    const res = await Promise.all(
      files.map((file) => {
        return new Promise((resolve, reject) => {
          if (!file.type.startsWith('image/')) {
            reject('Not an image file');
            return;
          }
          
          const form = new FormData();
          form.append('file', file);
          
          api.post(import.meta.env.VITE_BACKEND_URL + '/media/upload-temp', form, {
            headers: {
              'Content-Type': 'multipart/form-data'
            }
          })
          .then(async response => {
            const tempUrl = response.data.temp_url;
            const fileId = response.data.file_id;

            // Upload to S3
            await fetch(tempUrl, {
              method: 'PUT',
              body: file,
              headers: {
                'Content-Type': file.type
              }
            });

            // Request static GET URL by fileId
            let staticUrl = '';
            try {
              const staticRes = await api.get(
                import.meta.env.VITE_BACKEND_URL + '/media/gen_static_get',
                { params: { uuid: fileId } }
              );
              staticUrl = staticRes.data.static_url;
            } catch (err) {
              console.error('Failed to get static GET url:', err);
              staticUrl = tempUrl; // fallback to temp url
            }

            images.value.push({
              file,
              name: file.name,
              blobUrl: staticUrl
            });

            resolve({ temp_url: staticUrl });
          })
          .catch(error => {
            console.error('Failed to upload image:', error);
            reject(error);
          });
        });
      })
    );

    callback(res.map((item) => item.temp_url));
  };

  uploadImg(files, callback);
  return;
}

const imageFiles = ref([])

function handleFiles(event) {
  const input = event.target
  if (!input.files) return

  for (const file of Array.from(input.files)) {
    if (!file.type.startsWith('image/')) continue

    const reader = new FileReader()
    reader.onload = () => {
      imageFiles.value.push({
        file,
        name: file.name,
        preview: reader.result
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
    // Use axios instance for consistency with auth
    const res = await api.post('/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })

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