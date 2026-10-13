<template>
  <a-layout-sider
    v-model:collapsed="app.collapsed"
    class="app-sider"
    collapsible
    :width="220"
    :collapsed-width="80"
    :trigger="null"
    theme="light"
  >
    <div class="brand" :class="{ 'brand--collapsed': app.collapsed }">
      <img class="brand__icon" :src="brandIcon" alt="" />
      <img class="brand__wordmark" :src="brandWordmark" :alt="DEFAULT_BRANDING.name" />
    </div>
    <a-menu
      theme="light"
      mode="inline"
      :selected-keys="selectedKeys"
      :open-keys="openKeys"
      @click="onClick"
    >
      <AppSiderItem v-for="item in app.routers" :key="item.id" :item="item" />
    </a-menu>
  </a-layout-sider>
</template>

<script setup>
import { computed, nextTick, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAppStore } from '@/store'
import { DEFAULT_BRANDING } from '@/constants/branding'
import AppSiderItem from './AppSiderItem.vue'
import brandIcon from '@/assets/branding/fortune-finger.png'
import brandWordmark from '@/assets/branding/fortune-finger-wordmark.png'

const app = useAppStore()
const route = useRoute()
const router = useRouter()
const pendingPath = ref('')
const selectedKeys = computed(() => [pendingPath.value || route.path])
const openKeys = computed(() => findAncestorIDs(app.routers, route.path) || [])

watch(
  () => route.path,
  () => {
    pendingPath.value = ''
  },
)

function findAncestorIDs(items, path, ancestors = []) {
  for (const item of items) {
    if (item.route === path) return ancestors
    const found = findAncestorIDs(item.children || [], path, [...ancestors, String(item.id)])
    if (found) return found
  }
  return null
}

async function onClick({ key }) {
  if (key === route.path) return
  pendingPath.value = key
  await nextTick()
  await new Promise((resolve) => requestAnimationFrame(resolve))
  try {
    await router.push(key)
  } finally {
    if (route.path !== key) pendingPath.value = ''
  }
}
</script>

<style scoped>
.app-sider {
  background: #fff;
  border-right: 1px solid #edf0f4;
}
.app-sider :deep(.ant-layout-sider-children),
.app-sider :deep(.ant-menu-light) {
  background: transparent;
}
.app-sider :deep(.ant-menu-light .ant-menu-sub) {
  background: #faf8f8;
}
.app-sider :deep(.ant-menu-light .ant-menu-item-selected) {
  color: var(--color-primary);
  background: var(--color-primary-light);
}
.app-sider :deep(.ant-menu-light .ant-menu-item-selected::after) {
  border-inline-end-color: var(--color-primary);
}
.brand {
  height: 64px;
  padding: 0 20px;
  overflow: hidden;
  color: var(--color-text);
  font-size: 18px;
  font-weight: 700;
  display: flex;
  align-items: center;
  gap: 10px;
  white-space: nowrap;
}
.brand__icon {
  width: 34px;
  height: 34px;
  flex: 0 0 34px;
  object-fit: contain;
}
.brand__wordmark {
  width: 132px;
  height: auto;
  flex: 0 0 132px;
  object-fit: contain;
  opacity: 1;
  transform: translateY(2px);
  transition: opacity 120ms ease 180ms;
}
.brand--collapsed .brand__wordmark {
  opacity: 0;
  transition-delay: 0ms;
}
</style>
