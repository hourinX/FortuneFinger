<script setup lang="ts">
import brandIcon from '@/assets/branding/fortunefinger.png'

defineProps<{
  value: string
}>()
</script>

<template>
  <span
    v-if="value"
    class="qr-code-image"
  >
    <img
      class="qr-code-image__code"
      :src="value"
      alt="签到二维码"
    />
    <span class="qr-code-image__brand">
      <img
        :src="brandIcon"
        alt=""
      />
    </span>
  </span>
</template>

<style scoped>
.qr-code-image {
  position: relative;
  display: block;
  width: 100%;
  aspect-ratio: 1;
  overflow: hidden;
  background: #fff;
}
.qr-code-image__code {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: contain;
}
.qr-code-image__brand {
  position: absolute;
  top: 50%;
  left: 50%;
  display: grid;
  box-sizing: border-box;
  width: 20%;
  aspect-ratio: 1;
  place-items: center;
  padding: 2.5%;
  border-radius: 20%;
  background: #fff;
  transform: translate(-50%, -50%);
  filter: drop-shadow(0 2px 5px rgb(26 18 18 / 22%));
}
.qr-code-image__brand img {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: contain;
}
</style>
