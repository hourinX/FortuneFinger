<template>
  <a-form class="settings-form" layout="vertical">
    <div class="settings-form__grid">
      <a-form-item label="系统名称">
        <a-input v-model:value="model.name" :maxlength="50" />
      </a-form-item>
      <a-form-item label="浏览器标题">
        <a-input v-model:value="model.title" :maxlength="50" />
      </a-form-item>
      <a-form-item label="系统简称">
        <a-input v-model:value="model.shortName" :maxlength="20" />
      </a-form-item>
      <a-form-item label="默认语言">
        <a-select v-model:value="model.language" :options="languageOptions" />
      </a-form-item>
    </div>
    <a-form-item label="页脚内容">
      <a-input v-model:value="model.footer" :maxlength="100" />
    </a-form-item>
    <SettingRow title="维护模式" description="启用后仅超级管理员可以进入管理端">
      <a-switch v-model:checked="model.maintenance" />
    </SettingRow>
  </a-form>
</template>

<script setup>
import SettingRow from './setting-row.vue'

const model = defineModel({ type: Object, required: true })
const languageOptions = [{ label: '简体中文', value: 'zh-CN' }]
</script>

