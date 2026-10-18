<template>
  <div class="settings-row">
    <div class="settings-row__text">
      <strong>{{ title }}</strong>
      <span>{{ description }}</span>
    </div>
    <div class="settings-row__control">
      <slot />
    </div>
  </div>
</template>

<script setup>
defineProps({
  title: { type: String, required: true },
  description: { type: String, required: true },
})
</script>

