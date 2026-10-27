<template>
  <ScreenShell accent="#eab308">
    <div class="screen-intro">
      <p class="screen-intro__eyebrow">LOTTERY SCREEN</p>
      <h1>抽奖大屏</h1>
      <p>抽奖预热、开奖动画与中奖结果将在这里呈现。</p>
      <span>活动编号：{{ eventId }}</span>
    </div>
  </ScreenShell>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import ScreenShell from '@/components/ScreenShell.vue'

const route = useRoute()
const eventId = computed(() => String(route.params.eventId || 'preview'))
</script>
