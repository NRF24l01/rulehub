<template>
  <div class="space-y-4 p-4">
    <MdEditor v-model="markdownText" height="300px" language="en-US" @uploadImg="handleEditorImageInsert" />
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
</script>