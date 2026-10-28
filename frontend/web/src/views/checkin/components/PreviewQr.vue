<script setup lang="ts">
import type { CheckinTemplateProps } from '../useCheckin'
import QrCodeImage from './QrCodeImage.vue'

defineProps<Pick<CheckinTemplateProps, 'qr' | 'deadline'>>()
</script>
<template>
  <div class="scan-block">
    <div
      class="placeholder"
      :class="{ flash: qr.refreshing, inactive: !qr.active }"
    >
      <QrCodeImage
        v-if="qr.active && qr.value"
        :value="qr.value"
      />
      <strong v-else>
        {{ deadline.status === 'pending' ? '签到尚未开始' : deadline.closed ? '已停止签到' : qr.refreshing ? '二维码已刷新' : '扫码加入现场' }}
      </strong>
      <small v-if="qr.active && !qr.value">二维码加载中</small>
    </div>
    <p>{{ qr.active ? `${qr.refreshSeconds}s 后刷新 · 打开手机扫一扫` : '签到二维码暂不可用' }}</p>
  </div>
</template>
<style scoped>
.scan-block {
  text-align: center;
}
.placeholder {
  width: 100%;
  aspect-ratio: 1;
  background: #f3f1ea;
  color: #283e51;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  gap: 15px;
  border-radius: 12px;
  transition: box-shadow 0.3s;
}
.placeholder svg {
  width: 28%;
}
.placeholder strong {
  font-size: 14px;
  font-weight: 500;
}
.placeholder small {
  font-size: 9px;
  opacity: 0.5;
}
.scan-block p {
  font-size: 10px;
  opacity: 0.6;
  margin: 18px 0 0;
}
.flash {
  box-shadow: 0 0 35px #c5e8ff80;
}
.inactive {
  min-height: 150px;
  aspect-ratio: auto;
}
</style>
