<template>
  <a-sub-menu v-if="item.types === 1" :key="String(item.id)">
    <template #icon><component :is="resolveIcon(item.icon)" /></template>
    <template #title>{{ item.title }}</template>
    <AppSiderItem v-for="child in item.children" :key="child.id" :item="child" />
  </a-sub-menu>
  <a-menu-item v-else :key="item.route">
    <template #icon><component :is="resolveIcon(item.icon)" /></template>
    {{ item.title }}
  </a-menu-item>
</template>

<script setup>
import * as AntIcons from '@ant-design/icons-vue'

defineOptions({ name: 'AppSiderItem' })
defineProps({
  item: { type: Object, required: true },
})

function resolveIcon(code) {
  if (code === 'company') return AntIcons.BankOutlined

  const name = String(code || '')
    .split('-')
    .filter(Boolean)
    .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
    .join('')
  return AntIcons[`${name}Outlined`] || AntIcons.AppstoreOutlined
}
</script>
