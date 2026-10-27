<script setup lang="ts">
import type { CheckinTemplateProps } from '../useCheckin'

defineProps<Pick<CheckinTemplateProps, 'qr' | 'deadline'>>()
</script>
<template>
  <div class="scan-block">
    <div
      class="placeholder"
      :class="{ flash: qr.refreshing }"
    >
      <svg
        viewBox="0 0 64 64"
        fill="none"
        aria-hidden="true"
      >
        <path
          d="M8 24V8h16M40 8h16v16M56 40v16H40M24 56H8V40M25 25h14v14H25Z"
          stroke="currentColor"
          stroke-width="3"
        />
      </svg>
      <strong>{{ deadline.closed ? '签到已结束' : qr.refreshing ? '二维码已刷新' : '扫码加入现场' }}</strong>
      <small>演示占位 · 不可扫描</small>
    </div>
    <p>{{ deadline.closed ? '已停止刷新' : `${qr.refreshSeconds}s 后刷新 · 打开手机扫一扫` }}</p>
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
</style>
