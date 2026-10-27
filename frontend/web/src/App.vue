<template>
  <router-view v-slot="{ Component }">
    <transition name="screen" mode="out-in">
      <component :is="Component" />
    </transition>
  </router-view>
</template>
