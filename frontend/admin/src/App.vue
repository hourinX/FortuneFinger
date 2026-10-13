<template>
  <a-config-provider :theme="theme" :locale="locale">
    <router-view />
  </a-config-provider>
</template>

<script setup>
import zhCN from 'ant-design-vue/es/locale/zh_CN'
import { theme } from '@/config/theme'

const locale = {
  ...zhCN,
  Pagination: {
    ...zhCN.Pagination,
    jump_to: '跳转至',
  },
}
</script>
