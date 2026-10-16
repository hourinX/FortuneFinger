<template>
  <a-layout-header class="header">
    <div class="header__left">
      <menu-unfold-outlined v-if="app.collapsed" class="trigger" @click="app.toggleCollapsed" />
      <menu-fold-outlined v-else class="trigger" @click="app.toggleCollapsed" />
      <span class="header__separator"></span>
      <a-breadcrumb class="header__breadcrumb" separator=">">
        <a-breadcrumb-item v-for="(item, index) in breadcrumbs" :key="item.key">
          <router-link v-if="item.route && index < breadcrumbs.length - 1" :to="item.route">
            {{ item.title }}
          </router-link>
          <template v-else>{{ item.title }}</template>
        </a-breadcrumb-item>
      </a-breadcrumb>
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
        <a-popover
          v-model:open="settingsOpen"
          :trigger="['click']"
          placement="bottomRight"
          overlay-class-name="settings-popover"
        >
          <button
            class="header__action"
            type="button"
            aria-label="打开界面设置"
            :aria-expanded="settingsOpen"
          >
            <setting-outlined />
          </button>
          <template #content>
            <div class="settings-menu">
              <button class="settings-menu__item" type="button" @click="onThemeSettings">
                <skin-outlined />
                <span>主题设置</span>
              </button>
              <div class="settings-menu__divider"></div>
              <button class="settings-menu__item" type="button" @click="onRestoreDefaults">
                <reload-outlined />
                <span>恢复默认</span>
              </button>
            </div>
          </template>
        </a-popover>
      </div>
      <a-dropdown
        v-model:open="accountMenuOpen"
        :trigger="['click']"
        placement="bottomRight"
        overlay-class-name="account-dropdown"
      >
        <button class="account" type="button" :aria-expanded="accountMenuOpen">
          <a-avatar :src="admin.profile?.avatar">{{ accountInitial }}</a-avatar>
          <span class="account__name">{{ accountName }}</span>
          <down-outlined
            class="account__arrow"
            :class="{ 'account__arrow--open': accountMenuOpen }"
          />
        </button>
        <template #overlay>
          <div class="account-menu">
            <div class="account-menu__profile">
              <button
                class="account-menu__avatar"
                type="button"
                :disabled="!admin.profile?.avatar"
                aria-label="预览当前用户头像"
                @click.stop="openAvatarPreview"
              >
                <a-avatar :size="38" :src="admin.profile?.avatar">{{ accountInitial }}</a-avatar>
              </button>
              <div class="account-menu__identity">
                <strong>{{ accountName }}</strong>
                <span>{{ admin.profile?.username || '未设置账号' }}</span>
              </div>
            </div>
            <div class="account-menu__actions">
              <button class="account-menu__item" type="button" @click="onProfile">
                <user-outlined />
                <span>个人信息</span>
              </button>
              <button
                class="account-menu__item account-menu__item--danger"
                type="button"
                @click="onLogout"
              >
                <logout-outlined />
                <span>退出登录</span>
              </button>
            </div>
          </div>
        </template>
      </a-dropdown>
    </div>
    <ImagePreview
      v-model:visible="avatarPreviewVisible"
      :src="admin.profile?.avatar || ''"
      :alt="`${accountName}的头像`"
    />
  </a-layout-header>
</template>

<script setup>
import {
  DownOutlined,
  GithubOutlined,
  LogoutOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  ReloadOutlined,
  SettingOutlined,
  SkinOutlined,
  UserOutlined,
} from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ImagePreview from '@/components/ImagePreview.vue'
import { useAdminStore, useAppStore } from '@/store'

const app = useAppStore()
const admin = useAdminStore()
const route = useRoute()
const router = useRouter()
const accountMenuOpen = ref(false)
const settingsOpen = ref(false)
const avatarPreviewVisible = ref(false)
const accountName = computed(() => admin.profile?.nickname || admin.profile?.username || '用户')
const accountInitial = computed(() => String(accountName.value).slice(0, 1))
const breadcrumbs = computed(() => {
  const home = { key: 'home', title: '首页', route: '/dashboard' }
  if (route.path === home.route) return [home]
  const menuBreadcrumbs = findMenuBreadcrumbs(app.routers, route.path)
  if (menuBreadcrumbs.length) return [home, ...menuBreadcrumbs]
  const matched = route.matched
    .filter((item) => item.meta?.title)
    .map((item) => ({ key: String(item.name || item.path), title: item.meta.title }))
  return [home, ...matched]
})

function findMenuBreadcrumbs(items, path, parents = []) {
  for (const item of items || []) {
    const current = [
      ...parents,
      { key: String(item.id), title: item.title, route: item.route || '' },
    ]
    if (item.route === path) return current
    const children = findMenuBreadcrumbs(item.children, path, current)
    if (children.length) return children
  }
  return []
}

async function onLogout() {
  accountMenuOpen.value = false
  await admin.logout()
  app.resetRouters()
  router.replace('/login')
}

function onProfile() {
  accountMenuOpen.value = false
  router.push('/profile')
}

function openAvatarPreview() {
  if (!admin.profile?.avatar) return
  accountMenuOpen.value = false
  avatarPreviewVisible.value = true
}

function onThemeSettings() {
  settingsOpen.value = false
  message.info('主题设置功能建设中')
}

function onRestoreDefaults() {
  settingsOpen.value = false
  app.restoreDefaults()
  message.success('已恢复默认设置')
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
.header__left {
  min-width: 0;
}
.header__separator {
  width: 1px;
  height: 16px;
  margin: 0 16px;
  background: #e8e8e8;
}
.header__breadcrumb {
  min-width: 0;
  overflow: hidden;
  white-space: nowrap;
}
.header__breadcrumb :deep(.ant-breadcrumb-link),
.header__breadcrumb :deep(.ant-breadcrumb-separator) {
  font-size: 13px;
}
.header__breadcrumb :deep(li:last-child .ant-breadcrumb-link) {
  color: var(--color-text);
  font-weight: 500;
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
.header__action[aria-expanded='true'] {
  color: var(--color-primary);
  background: var(--color-primary-light);
}
.settings-menu {
  display: flex;
  width: 104px;
  flex-direction: column;
  padding: 0;
}
.settings-menu__item {
  display: flex;
  align-items: center;
  gap: 7px;
  width: 100%;
  height: 36px;
  padding: 0 8px;
  color: var(--color-text);
  font-size: 14px;
  text-align: left;
  background: transparent;
  border: 0;
  border-radius: 5px;
  cursor: pointer;
  transition:
    color 160ms ease,
    background-color 160ms ease;
}
.settings-menu__item:hover {
  color: var(--color-primary);
  background: var(--color-primary-light);
}
.settings-menu__divider {
  height: 1px;
  margin: 3px 8px;
  background: #f0f0f0;
}
.account {
  display: flex;
  align-items: center;
  gap: 8px;
  height: 40px;
  min-width: 0;
  padding: 0 10px 0 6px;
  color: var(--color-text);
  background: transparent;
  border: 0;
  border-radius: 6px;
  cursor: pointer;
  transition:
    color 160ms ease,
    background-color 160ms ease;
}
.account:hover,
.account[aria-expanded='true'] {
  color: var(--color-primary);
  background: #f8f6f6;
}
.account__name {
  max-width: 120px;
  overflow: hidden;
  color: inherit;
  font-size: 14px;
  line-height: 1;
  white-space: nowrap;
  text-overflow: ellipsis;
}
.account__arrow {
  margin-left: 1px;
  color: #999;
  font-size: 10px;
  transition: transform 180ms ease;
}
.account__arrow--open {
  transform: rotate(180deg);
}
.account-menu {
  position: relative;
  width: 220px;
  margin-top: 5px;
  padding: 6px;
  background: #fff;
  border: 1px solid #eee;
  border-radius: 8px;
  box-shadow: 0 6px 20px rgb(0 0 0 / 8%);
}
.account-menu::before {
  position: absolute;
  top: -6px;
  right: 24px;
  width: 11px;
  height: 11px;
  background: #fff;
  border-top: 1px solid #eee;
  border-left: 1px solid #eee;
  transform: rotate(45deg);
  content: '';
}
.account-menu__profile {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
  padding: 10px;
  background: #fafafa;
  border-radius: 6px;
}
.account-menu__avatar {
  display: inline-flex;
  flex: 0 0 auto;
  padding: 0;
  background: transparent;
  border: 0;
  border-radius: 50%;
  cursor: pointer;
}
.account-menu__avatar:disabled {
  cursor: default;
}
.account-menu__identity {
  display: flex;
  flex: 1;
  flex-direction: column;
  min-width: 0;
}
.account-menu__identity strong,
.account-menu__identity span {
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}
.account-menu__identity strong {
  color: var(--color-text);
  font-size: 14px;
  font-weight: 600;
  line-height: 20px;
}
.account-menu__identity span {
  color: var(--color-muted);
  font-size: 12px;
  line-height: 18px;
}
.account-menu__actions {
  display: flex;
  flex-direction: column;
  gap: 2px;
  margin-top: 5px;
}
.account-menu__item {
  display: flex;
  align-items: center;
  gap: 9px;
  width: 100%;
  height: 36px;
  margin: 0;
  padding: 0 11px;
  color: var(--color-text);
  font-size: 14px;
  text-align: left;
  background: transparent;
  border: 0;
  border-radius: 5px;
  cursor: pointer;
  transition:
    color 160ms ease,
    background-color 160ms ease;
}
.account-menu__item--danger:hover {
  color: #d9363e;
  background: #fff2f0;
}
</style>
