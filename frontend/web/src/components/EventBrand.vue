<script setup lang="ts">
import wordmark from '@/assets/branding/fortune-finger-wordmark.png'
withDefaults(defineProps<{ tone?: 'light' | 'dark' }>(), { tone: 'dark' })
</script>

<template>
  <div class="event-brand" :class="`event-brand--${tone}`">
    <img :src="wordmark" alt="FortuneFinger" width="1864" height="294" />
  </div>
</template>

<style scoped>
.event-brand { display: inline-flex; align-items: center; flex-shrink: 0; line-height: 0; }
/* 保持原图宽高比：1920px 屏幕约 192 × 30px。 */
.event-brand img { display: block; width: clamp(160px, 10vw, 240px); height: auto; margin: 0; }
.event-brand--dark img { filter: brightness(0) invert(1); opacity: .88; }
@media (max-width: 600px) { .event-brand img { width: 150px; } }
</style>
