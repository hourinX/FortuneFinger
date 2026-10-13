<template>
  <transition name="image-preview-fade">
    <div v-if="visible" class="image-preview" @click.self="close" @wheel.prevent="onWheel">
      <button type="button" class="image-preview__close" aria-label="关闭图片预览" @click="close">
        <close-outlined />
      </button>
      <img
        class="image-preview__image"
        :src="src"
        :style="{ transform: `scale(${scale})` }"
        :alt="alt"
        draggable="false"
        @dblclick="resetScale"
      />
    </div>
  </transition>
</template>

<script setup>
import { onBeforeUnmount, ref, watch } from 'vue'
import { CloseOutlined } from '@ant-design/icons-vue'

const props = defineProps({
  visible: { type: Boolean, default: false },
  src: { type: String, default: '' },
  alt: { type: String, default: '' },
})
const emit = defineEmits(['update:visible'])
const scale = ref(1)

function close() {
  emit('update:visible', false)
}
function resetScale() {
  scale.value = 1
}
function onWheel(event) {
  const step = event.deltaY < 0 ? 1.1 : 1 / 1.1
  scale.value = Math.min(4, Math.max(0.25, scale.value * step))
}
function onKeydown(event) {
  if (event.key === 'Escape') close()
}

watch(
  () => props.visible,
  (visible) => {
    scale.value = 1
    window[visible ? 'addEventListener' : 'removeEventListener']('keydown', onKeydown)
  },
  { immediate: true },
)
watch(() => props.src, resetScale)
onBeforeUnmount(() => window.removeEventListener('keydown', onKeydown))
</script>

<style scoped>
.image-preview {
  position: fixed;
  inset: 0;
  z-index: 3000;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  background: rgba(10, 18, 21, 0.78);
  cursor: zoom-out;
}
.image-preview__image {
  max-width: 90vw;
  max-height: 90vh;
  object-fit: contain;
  cursor: default;
  user-select: none;
  transition: transform 0.1s ease-out;
}
.image-preview__close {
  position: absolute;
  top: 18px;
  right: 18px;
  z-index: 1;
  display: grid;
  width: 38px;
  height: 38px;
  padding: 0;
  place-items: center;
  color: #fff;
  font-size: 18px;
  background: rgba(0, 0, 0, 0.42);
  border: 0;
  border-radius: 50%;
  cursor: pointer;
}
.image-preview__close:hover {
  background: rgba(0, 0, 0, 0.72);
}
.image-preview-fade-enter-active,
.image-preview-fade-leave-active {
  transition: opacity 0.2s;
}
.image-preview-fade-enter-from,
.image-preview-fade-leave-to {
  opacity: 0;
}
</style>
