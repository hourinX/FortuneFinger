<script setup lang="ts">
defineProps<{
  active: 'chat' | 'lottery' | 'profile'
}>()

const tabs = [
  { key: 'chat', label: '弹幕', to: { name: 'participant-home' } },
  { key: 'lottery', label: '抽奖', to: { name: 'participant-lottery' } },
  { key: 'profile', label: '我的', to: { name: 'participant-profile' } },
] as const
</script>

<template>
  <nav class="app-tabbar" aria-label="参与者端导航">
    <router-link
      v-for="tab in tabs"
      :key="tab.key"
      :to="tab.to"
      class="app-tabbar__item"
      :class="{ 'app-tabbar__item--active': active === tab.key }"
      :aria-current="active === tab.key ? 'page' : undefined"
    >
      <span class="app-tabbar__icon" :class="`app-tabbar__icon--${tab.key}`">
        <i></i>
      </span>
      <small>{{ tab.label }}</small>
    </router-link>
  </nav>
</template>

<style scoped>
.app-tabbar { display: grid; min-height: 62px; padding: 6px 22px max(8px, env(safe-area-inset-bottom)); box-sizing: border-box; grid-template-columns: repeat(3, 1fr); gap: 8px; background: rgb(255 255 255 / 97%); border-top: 1px solid #eee7de; box-shadow: 0 -8px 28px rgb(68 49 36 / 5%); }
.app-tabbar__item { position: relative; display: flex; min-height: 45px; padding: 4px 0 2px; flex-direction: column; align-items: center; justify-content: center; gap: 3px; color: #aaa19a; text-decoration: none; border-radius: 13px; transition: color .2s ease, background-color .2s ease, transform .2s ease; }
.app-tabbar small { font-size: 9px; font-weight: 500; line-height: 13px; }
.app-tabbar__icon { position: relative; display: block; width: 22px; height: 20px; }
.app-tabbar__icon i, .app-tabbar__icon::before, .app-tabbar__icon::after { position: absolute; box-sizing: border-box; content: ''; }
.app-tabbar__icon--chat::before { top: 2px; left: 2px; width: 18px; height: 14px; border: 1.5px solid currentColor; border-radius: 6px; }
.app-tabbar__icon--chat::after { bottom: 1px; left: 6px; width: 6px; height: 6px; border-bottom: 1.5px solid currentColor; transform: rotate(-35deg); }
.app-tabbar__icon--chat i { top: 8px; left: 7px; width: 2px; height: 2px; background: currentColor; border-radius: 50%; box-shadow: 4px 0 currentColor, 8px 0 currentColor; }
.app-tabbar__icon--lottery::before { inset: 3px 2px 2px; border: 1.5px solid currentColor; border-radius: 4px; }
.app-tabbar__icon--lottery::after { top: 0; left: 10px; width: 2px; height: 20px; background: currentColor; }
.app-tabbar__icon--lottery i { top: 5px; left: 5px; width: 5px; height: 5px; border: 1.5px solid currentColor; border-radius: 50%; }
.app-tabbar__icon--profile::before { top: 1px; left: 7px; width: 8px; height: 8px; border: 1.5px solid currentColor; border-radius: 50%; }
.app-tabbar__icon--profile::after { right: 3px; bottom: 1px; left: 3px; height: 9px; border: 1.5px solid currentColor; border-radius: 9px 9px 4px 4px; }
.app-tabbar__item--active { color: #df433c !important; }
.app-tabbar__item:active { transform: scale(.96); }
</style>
