<template>
  <div class="settings-form">
    <div class="settings-group">
      <div class="settings-group__title">中奖限制</div>
      <SettingRow title="允许重复中奖" description="关闭后，已中奖人员不会进入后续轮次候选池">
        <a-switch v-model:checked="model.allowRepeat" />
      </SettingRow>
      <SettingRow title="同奖项重复中奖" description="控制同一人员能否重复获得相同等级奖项">
        <a-switch v-model:checked="model.allowSameLevel" />
      </SettingRow>
      <SettingRow title="单人中奖上限" description="达到上限后自动从后续抽奖候选池排除">
        <a-input-number v-model:value="model.maxWins" :min="1" :max="20" />
      </SettingRow>
    </div>

    <div class="settings-group">
      <div class="settings-group__title">阳光普照奖</div>
      <SettingRow title="排除其他奖项获奖者" description="只从尚未获得其他奖项的参与人员中产生">
        <a-switch v-model:checked="model.excludeWinners" />
      </SettingRow>
      <SettingRow title="超出奖品数量时" description="符合条件人数超过配置数量后的处理方式">
        <a-select v-model:value="model.overflowMode" class="settings-control" :options="overflowOptions" />
      </SettingRow>
    </div>

  </div>
</template>

<script setup>
import SettingRow from './setting-row.vue'

const model = defineModel({ type: Object, required: true })
const overflowOptions = [
  { label: '按配置数量随机发放', value: 'random' },
  { label: '停止发放并提示库存不足', value: 'stop' },
]
</script>

