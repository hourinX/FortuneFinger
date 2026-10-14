<template>
  <span class="status-pulse" :class="`status-pulse--${status}`">
    <i />
    <span v-if="showText">{{ label }}</span>
  </span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  status: { type: Number, required: true },
  activeText: { type: String, default: '启用' },
  inactiveText: { type: String, default: '禁用' },
  showText: { type: Boolean, default: true },
})

const label = computed(() => (props.status === 1 ? props.activeText : props.inactiveText))
</script>

<style scoped>
.status-pulse {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: #3f3f3f;
  white-space: nowrap;
}

.status-pulse i {
  position: relative;
  width: 7px;
  height: 7px;
  background: #3fa765;
  border-radius: 50%;
}

.status-pulse i::after {
  position: absolute;
  inset: -5px;
  content: '';
  background: rgba(63, 167, 101, 0.3);
  border-radius: 50%;
  animation: status-breathe 1.65s ease-out infinite;
}

.status-pulse--2 i {
  background: #d7473f;
}

.status-pulse--2 i::after {
  background: rgba(215, 71, 63, 0.28);
}

@keyframes status-breathe {
  0% {
    opacity: 1;
    transform: scale(0.25);
  }
  75%,
  100% {
    opacity: 0;
    transform: scale(1.45);
  }
}
</style>
