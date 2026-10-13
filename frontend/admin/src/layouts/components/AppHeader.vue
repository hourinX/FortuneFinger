<template>
  <a-layout-header class="header">
    <div class="header__left">
      <menu-unfold-outlined v-if="app.collapsed" class="trigger" @click="app.toggleCollapsed" />
      <menu-fold-outlined v-else class="trigger" @click="app.toggleCollapsed" />
    </div>
    <div class="header__right">
      <div class="header__actions">
        <a-tooltip title="GitHub">
          <a
            class="header__action"
            href="https://github.com/hourinX/FortuneFinger"
            target="_blank"
            rel="noopener noreferrer"
            aria-label="打开 FortuneFinger GitHub 仓库"
          >
            <github-outlined />
          </a>
        </a-tooltip>
        <a-tooltip title="设置功能建设中">
          <button class="header__action" type="button" aria-label="设置功能建设中">
            <setting-outlined />
          </button>
        </a-tooltip>
      </div>
      <a-dropdown>
        <div class="account">
          <a-avatar :src="admin.profile?.avatar">{{ accountInitial }}</a-avatar>
          <span>{{ admin.profile?.nickname || admin.profile?.username || '用户' }}</span>
        </div>
        <template #overlay>
          <a-menu @click="onAccountMenu">
            <a-menu-item key="logout">退出登录</a-menu-item>
          </a-menu>
        </template>
      </a-dropdown>
    </div>
  </a-layout-header>
</template>

<script setup>
import {
  GithubOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  SettingOutlined,
} from '@ant-design/icons-vue'
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAdminStore, useAppStore } from '@/store'

const app = useAppStore()
const admin = useAdminStore()
const router = useRouter()
const accountInitial = computed(() =>
  String(admin.profile?.nickname || admin.profile?.username || '用').slice(0, 1),
)

async function onAccountMenu({ key }) {
  if (key !== 'logout') return
  await admin.logout()
  app.resetRouters()
  router.replace('/login')
}
</script>

<style scoped>
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 64px;
  padding: 0 24px;
  background: #fff;
  border-bottom: 1px solid #edf0f4;
}
.trigger {
  font-size: 19px;
  cursor: pointer;
}
.header__left,
.header__right,
.header__actions,
.header__action {
  display: flex;
  align-items: center;
}
.header__right {
  gap: 14px;
}
.header__actions {
  gap: 4px;
}
.header__action {
  justify-content: center;
  width: 34px;
  height: 34px;
  padding: 0;
  color: var(--color-muted);
  font-size: 18px;
  background: transparent;
  border: 0;
  border-radius: 8px;
  cursor: pointer;
  transition:
    color 160ms ease,
    background-color 160ms ease;
}
.header__action:hover {
  color: var(--color-primary);
  background: var(--color-primary-light);
}
.account {
  display: flex;
  align-items: center;
  gap: 9px;
  cursor: pointer;
}
</style>
