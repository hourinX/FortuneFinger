<template>
  <a-form class="settings-form" layout="vertical">
    <div class="settings-form__grid">
      <a-form-item label="图片大小上限">
        <a-input-number v-model:value="model.imageMaxSize" :min="1" :max="100" addon-after="MB" />
      </a-form-item>
      <a-form-item label="视频大小上限">
        <a-input-number v-model:value="model.videoMaxSize" :min="1" :max="2048" addon-after="MB" />
      </a-form-item>
    </div>
    <a-form-item label="允许的图片格式">
      <a-select v-model:value="model.imageFormats" mode="multiple" :options="imageFormatOptions" />
    </a-form-item>
    <a-form-item label="允许的视频格式">
      <a-select v-model:value="model.videoFormats" mode="multiple" :options="videoFormatOptions" />
    </a-form-item>
    <SettingRow title="保留原始文件名" description="关闭后仅使用系统生成的文件名保存对象">
      <a-switch v-model:checked="model.keepOriginalName" />
    </SettingRow>
  </a-form>
</template>

<script setup>
import SettingRow from './setting-row.vue'

const model = defineModel({ type: Object, required: true })
const imageFormatOptions = ['JPG', 'PNG', 'WEBP', 'GIF'].map((value) => ({ label: value, value }))
const videoFormatOptions = ['MP4', 'WEBM', 'MOV'].map((value) => ({ label: value, value }))
</script>
