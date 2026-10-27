<template>
  <a-config-provider :theme="theme" :locale="locale">
    <router-view />
    <LoadingScreen v-if="app.fullscreenLoading" />
  </a-config-provider>
</template>

<script setup>
import zhCN from 'ant-design-vue/es/locale/zh_CN'
import LoadingScreen from '@/components/loadingscreen/index.vue'
import { theme } from '@/config/theme'
import { useAppStore } from '@/store'

const app = useAppStore()

const locale = {
  ...zhCN,
  Pagination: {
    ...zhCN.Pagination,
    jump_to: '跳转至',
  },
}
</script>
