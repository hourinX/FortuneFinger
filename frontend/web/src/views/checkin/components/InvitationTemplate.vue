<script setup lang="ts">
import CheckinLayout from './CheckinLayout.vue'
import type { CheckinTemplateProps } from '../useCheckin'

defineProps<CheckinTemplateProps>()
</script>

<template>
  <CheckinLayout
    :content="content"
    :qr="qr"
    :attendance="attendance"
    :deadline="deadline"
    theme="invitation"
  />
</template>
