<template>
  <aside class="settings-nav">
    <div class="settings-nav__mark" aria-hidden="true">
      <img :src="brandWordmark" alt="" />
    </div>
    <button
      v-for="item in items"
      :key="item.code"
      v-permission="item.permission"
      type="button"
      class="settings-nav__item"
      :class="{ 'settings-nav__item--active': modelValue === item.code }"
      @click="emit('update:modelValue', item.code)"
    >
      <component :is="item.icon" />
      <span>{{ item.label }}</span>
      <right-outlined />
    </button>
  </aside>
</template>

<script setup>
import { RightOutlined } from '@ant-design/icons-vue'
import brandWordmark from '@/assets/branding/fortune-finger-wordmark.png'

defineProps({
  modelValue: { type: String, required: true },
  items: { type: Array, required: true },
})
const emit = defineEmits(['update:modelValue'])
</script>

<style scoped lang="less">
.settings-nav {
  min-height: 0;
  padding: 0 12px 16px;
  overflow-y: auto;
  scrollbar-width: none;
  background: #fafafa;
  border-right: 1px solid #efeded;

  &::-webkit-scrollbar {
    display: none;
  }
}

.settings-nav__mark {
  display: flex;
  height: 64px;
  margin-bottom: 10px;
  align-items: center;
  justify-content: center;
  border-bottom: 1px solid #ececec;

  img {
    width: 148px;
    height: auto;
    object-fit: contain;
  }
}

.settings-nav__item {
  display: grid;
  width: 100%;
  height: 42px;
  margin-bottom: 4px;
  padding: 0 12px;
  grid-template-columns: 20px minmax(0, 1fr) 14px;
  align-items: center;
  gap: 8px;
  color: #555;
  text-align: left;
  background: transparent;
  border: 0;
  border-radius: 6px;
  cursor: pointer;

  > :last-child {
    color: #bbb;
    font-size: 10px;
  }

  &:hover {
    color: var(--color-primary);
    background: #fff;
  }
}

.settings-nav__item--active {
  color: var(--color-primary);
  background: var(--color-primary-light);

  > :last-child {
    color: var(--color-primary);
  }
}

@media (max-width: 900px) {
  .settings-nav {
    display: flex;
    gap: 6px;
    overflow-y: visible;
    overflow-x: auto;
    border-right: 0;
    border-bottom: 1px solid #efeded;
  }

  .settings-nav__mark {
    display: none;
  }

  .settings-nav__item {
    width: auto;
    min-width: max-content;
    margin: 0;
    grid-template-columns: 18px auto;

    > :last-child {
      display: none;
    }
  }
}
</style>
