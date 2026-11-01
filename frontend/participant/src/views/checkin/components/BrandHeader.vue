<script setup lang="ts">
import brandWordmark from '@/assets/branding/fortune-finger-wordmark.png'

defineProps<{
  company: string
  caption?: string
}>()
</script>

<template>
  <header class="brand-header">
    <div>
      <img class="brand-header__wordmark" :src="brandWordmark" alt="FortuneFinger" />
      <small>{{ caption || `${company} · 活动现场移动签到` }}</small>
    </div>
  </header>
</template>

<style scoped>
.brand-header { display: flex; align-items: center; }
.brand-header > div:last-child { display: flex; min-width: 0; flex-direction: column; }
.brand-header__wordmark { width: 154px; height: auto; object-fit: contain; }
.brand-header small { margin-top: 4px; color: #a39b94; font-size: 8px; letter-spacing: .12em; }
</style>
