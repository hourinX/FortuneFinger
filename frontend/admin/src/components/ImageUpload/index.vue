<template>
  <div class="image-upload" :class="{ 'image-upload--filled': hasImage }">
    <input ref="fileInputRef" type="file" :accept="accept" hidden @change="upload" />
    <slot :uploading="uploading" :preview="displayPreview" :has-image="hasImage" />

    <div v-if="hasImage" class="image-upload__mask image-upload__actions">
      <button
        type="button"
        aria-label="预览图片"
        :disabled="!displayPreview"
        @click="previewImage"
      >
        <eye-outlined />
      </button>
      <button type="button" aria-label="移除图片" @click="remove">
        <delete-outlined />
      </button>
    </div>
    <button
      v-else
      type="button"
      class="image-upload__mask image-upload__trigger"
      :disabled="uploading"
      @click="open"
    >
      <loading-outlined v-if="uploading" />
      <camera-outlined v-else />
      <span>{{ uploading ? '上传中' : '上传图片' }}</span>
    </button>

    <ImagePreview
      v-model:visible="previewVisible"
      :src="displayPreview"
      alt="图片预览"
    />
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, ref } from 'vue'
import {
  CameraOutlined,
  DeleteOutlined,
  EyeOutlined,
  LoadingOutlined,
} from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { uploadAsset } from '@/api/asset'
import ImagePreview from '@/components/ImagePreview.vue'

const props = defineProps({
  modelValue: { type: [String, Number], default: '0' },
  scene: { type: String, required: true },
  accept: { type: String, default: 'image/*' },
  src: { type: String, default: '' },
})
const emit = defineEmits(['update:modelValue', 'uploaded', 'removed'])

const fileInputRef = ref()
const uploading = ref(false)
const preview = ref('')
const previewVisible = ref(false)
const displayPreview = computed(() => preview.value || props.src)
const hasImage = computed(() => Boolean(displayPreview.value) || Number(props.modelValue) > 0)

function open() {
  if (!uploading.value) fileInputRef.value?.click()
}

function previewImage() {
  if (displayPreview.value) previewVisible.value = true
}

function remove() {
  previewVisible.value = false
  clearPreview()
  if (fileInputRef.value) fileInputRef.value.value = ''
  emit('update:modelValue', '0')
  emit('removed')
}

async function upload(event) {
  const file = event.target.files?.[0]
  event.target.value = ''
  if (!file) return
  if (!file.type.startsWith('image/')) {
    message.warning('请选择图片文件')
    return
  }

  try {
    uploading.value = true
    const resp = await uploadAsset(file, props.scene)
    clearPreview()
    preview.value = URL.createObjectURL(file)
    emit('update:modelValue', resp.storage_id)
    emit('uploaded', resp)
  } catch (error) {
    message.error(error.message)
  } finally {
    uploading.value = false
  }
}

function clearPreview() {
  if (preview.value) URL.revokeObjectURL(preview.value)
  preview.value = ''
}

onBeforeUnmount(clearPreview)
</script>

<style scoped>
.image-upload {
  position: relative;
  display: inline-block;
  min-width: 0;
  overflow: hidden;
}

.image-upload__mask {
  position: absolute;
  inset: 0;
  color: #fff;
  background: rgb(0 0 0 / 48%);
  border: 0;
  opacity: 0;
  transition: opacity 0.2s ease;
}

.image-upload:hover .image-upload__mask,
.image-upload:focus-within .image-upload__mask,
.image-upload__trigger:disabled {
  opacity: 1;
}

.image-upload__trigger {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 7px;
  cursor: pointer;
}

.image-upload__trigger:disabled {
  cursor: wait;
}

.image-upload__trigger :deep(.anticon) {
  font-size: 24px;
}

.image-upload__actions {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 18px;
  pointer-events: none;
}

.image-upload__actions button {
  display: grid;
  width: 34px;
  height: 34px;
  padding: 0;
  place-items: center;
  color: #fff;
  font-size: 19px;
  background: transparent;
  border: 0;
  border-radius: 50%;
  cursor: pointer;
  pointer-events: auto;
}

.image-upload__actions button:hover {
  background: rgb(255 255 255 / 18%);
}

.image-upload__actions button:disabled {
  color: rgb(255 255 255 / 42%);
  cursor: not-allowed;
}
</style>
