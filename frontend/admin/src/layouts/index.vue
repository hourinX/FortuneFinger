<template>
  <a-layout class="admin-layout">
    <AppSider />
    <a-layout class="admin-main">
      <AppHeader />
      <AppRouterTags />
      <a-layout-content class="admin-content">
        <router-view v-slot="{ Component, route }">
          <keep-alive>
            <component :is="Component" v-if="route.meta.keepAlive" :key="route.name" />
          </keep-alive>
          <component :is="Component" v-if="!route.meta.keepAlive" :key="route.fullPath" />
        </router-view>
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup>
import AppHeader from './components/AppHeader.vue'
import AppRouterTags from './components/AppRouterTags.vue'
import AppSider from './components/AppSider.vue'
</script>

<style scoped>
.admin-layout {
  width: 100%;
  max-width: 100%;
  min-height: 100vh;
}
.admin-main {
  width: 0;
  min-width: 0;
  flex: 1 1 0;
}
.admin-content {
  width: auto;
  max-width: calc(100% - 40px);
  min-width: 0;
  margin: 5px 20px 20px;
  padding: 24px;
  overflow: auto;
  background: #fff;
  border-radius: 10px;
}
</style>
