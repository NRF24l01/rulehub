<template>
  <div class="space-y-4 p-4">
    <MdEditor
      v-model="localMarkdown"
      height="300px"
      language="en-US"
      @uploadImg="handleEditorImageInsert"
    />
  </div>
</template>

<script setup>
import { ref, watch, defineProps, defineEmits } from 'vue'
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import api from '../axios'

const props = defineProps({
  content: {
    type: String,
    default: ''
  },
  images: {
    type: Array,
    default: () => []
  }
})


const emit = defineEmits(['update:content', 'update:images'])

const localMarkdown = ref(props.content)
const localImages = ref([...props.images])

watch(() => props.content, val => localMarkdown.value = val)

function arraysShallowEqual(a, b) {
  if (a.length !== b.length) return false
  for (let i = 0; i < a.length; i++) {
    if (a[i] !== b[i]) return false
  }
  return true
}

watch(() => props.images, val => {
  if (!arraysShallowEqual(val, localImages.value)) {
    localImages.value = [...val]
  }
})

watch(localMarkdown, val => {
  if (val !== props.content) emit('update:content', val)
})
watch(localImages, val => {
  if (!arraysShallowEqual(val, props.images)) emit('update:images', val)
})

async function handleEditorImageInsert(files, callback) {
  try {
    const results = await Promise.all(
      files.map(async (file) => {
        if (!file.type.startsWith('image/')) throw new Error('Not an image')

        const form = new FormData()
        form.append('file', file)

        const response = await api.post(
          import.meta.env.VITE_BACKEND_URL + '/media/upload-temp',
          form,
          { headers: { 'Content-Type': 'multipart/form-data' } }
        )

        const tempUrl = response.data.temp_url
        const fileId = response.data.file_id

        await fetch(tempUrl, {
          method: 'PUT',
          body: file,
          headers: { 'Content-Type': file.type }
        })

        let staticUrl = tempUrl
        try {
          const res = await api.get(
            import.meta.env.VITE_BACKEND_URL + '/media/gen_static_get',
            { params: { uuid: fileId } }
          )
          staticUrl = res.data.static_url
        } catch (_) {
          console.warn('Fallback to temp URL')
        }

        localImages.value = [
          ...localImages.value,
          {
            file,
            name: file.name,
            blobUrl: staticUrl
          }
        ]
        return staticUrl
      })
    )

    callback(results)
  } catch (err) {
    console.error('Image upload failed:', err)
  }
}
</script>
