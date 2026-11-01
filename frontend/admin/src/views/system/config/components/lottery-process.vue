<template>
  <div class="settings-form lottery-process">
    <div class="settings-group">
      <div class="settings-group__title">开奖触发</div>
      <div class="process-mode">
        <button
          v-for="option in triggerModeOptions"
          :key="option.value"
          type="button"
          class="process-mode__item"
          :class="{ 'process-mode__item--active': model.trigger_mode === option.value }"
          @click="model.trigger_mode = option.value"
        >
          <span class="process-mode__icon">
            <component :is="option.icon" />
          </span>
          <span class="process-mode__body">
            <strong>{{ option.label }}</strong>
            <small>{{ option.description }}</small>
          </span>
          <span class="process-mode__check"></span>
        </button>
      </div>
    </div>

    <div class="settings-group">
      <div class="settings-group__title">现场操作</div>
      <SettingRow
        title="触发交互"
        description="手动开奖时，大屏提供给主持人的启动方式"
      >
        <a-select
          v-model:value="model.trigger_interaction"
          class="settings-control"
          :options="interactionOptions"
          :disabled="model.trigger_mode === 'automatic'"
        />
      </SettingRow>
      <SettingRow
        v-if="model.trigger_interaction === 'hold'"
        title="长按蓄力时长"
        description="按住操作区域达到设定时间后才会触发开奖"
      >
        <a-input-number
          v-model:value="model.hold_duration"
          class="settings-control"
          :min="1"
          :max="10"
          :precision="0"
          addon-after="秒"
        />
      </SettingRow>
      <SettingRow
        title="允许提前开奖"
        description="开启后，轮次到达计划时间前也可以由现场人员触发"
      >
        <a-switch v-model:checked="model.allow_early" />
      </SettingRow>
    </div>

    <div class="settings-group">
      <div class="settings-group__title">动画衔接</div>
      <SettingRow
        title="开奖预热时长"
        description="触发成功后，先播放预热动画，再进入正式开奖表现"
      >
        <a-input-number
          v-model:value="model.warmup_duration"
          class="settings-control"
          :min="0"
          :max="30"
          :precision="0"
          addon-after="秒"
        />
      </SettingRow>
      <SettingRow
        title="中奖结果揭晓"
        description="同一轮抽出多名中奖者时的默认展示方式"
      >
        <a-select
          v-model:value="model.reveal_mode"
          class="settings-control"
          :options="revealOptions"
        />
      </SettingRow>
      <SettingRow
        title="自动开奖失败后"
        description="自动模式未能按时执行时，保留轮次并等待现场处理"
      >
        <a-select
          v-model:value="model.auto_failure_mode"
          class="settings-control"
          :options="failureOptions"
          :disabled="model.trigger_mode !== 'automatic'"
        />
      </SettingRow>
    </div>

    <div class="process-note">
      <info-circle-outlined />
      <span>这里保存的是全局默认流程。后续为单个抽奖轮次设置不同流程时，以轮次配置为准。</span>
    </div>
  </div>
</template>

<script setup>
import {
  ClockCircleOutlined,
  InfoCircleOutlined,
  PlayCircleOutlined,
} from '@ant-design/icons-vue'
import SettingRow from './setting-row.vue'

const model = defineModel({ type: Object, required: true })

const triggerModeOptions = [
  {
    label: '手动触发',
    value: 'manual',
    description: '到达轮次时间后解锁操作，由主持人或现场人员主动开始',
    icon: PlayCircleOutlined,
  },
  {
    label: '自动触发',
    value: 'automatic',
    description: '到达轮次时间后由后端自动执行，不依赖大屏本地计时',
    icon: ClockCircleOutlined,
  },
]

const interactionOptions = [
  { label: '长按蓄力', value: 'hold' },
  { label: '普通按钮', value: 'button' },
  { label: '老虎机拉杆', value: 'lever' },
]

const revealOptions = [
  { label: '全部抽取后统一揭晓', value: 'batch' },
  { label: '按中奖顺序逐个揭晓', value: 'sequence' },
]

const failureOptions = [
  { label: '切换为手动等待触发', value: 'manual' },
  { label: '保持失败状态等待管理端处理', value: 'blocked' },
]
</script>

<style scoped lang="less">
.lottery-process {
  padding-bottom: 28px;
}

.process-mode {
  display: grid;
  padding-top: 18px;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.process-mode__item {
  position: relative;
  display: grid;
  min-height: 112px;
  padding: 20px;
  grid-template-columns: 42px minmax(0, 1fr) 16px;
  align-items: center;
  gap: 14px;
  color: #555;
  text-align: left;
  background: #fafafa;
  border: 1px solid #ededed;
  border-radius: 8px;
  cursor: pointer;
  transition: border-color 0.2s, background 0.2s, box-shadow 0.2s;

  &:hover {
    border-color: #efb4b2;
    background: #fff;
  }
}

.process-mode__item--active {
  background: #fffafa;
  border-color: var(--color-primary);
  box-shadow: 0 4px 16px rgb(238 62 58 / 8%);
}

.process-mode__icon {
  display: inline-flex;
  width: 42px;
  height: 42px;
  align-items: center;
  justify-content: center;
  color: var(--color-primary);
  font-size: 20px;
  background: var(--color-primary-light);
  border-radius: 50%;
}

.process-mode__body {
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 7px;

  strong {
    color: #333;
    font-size: 15px;
    font-weight: 600;
  }

  small {
    color: #8c8c8c;
    font-size: 12px;
    line-height: 19px;
  }
}

.process-mode__check {
  width: 14px;
  height: 14px;
  border: 1px solid #cfcfcf;
  border-radius: 50%;
  box-shadow: inset 0 0 0 3px #fff;
}

.process-mode__item--active .process-mode__check {
  background: var(--color-primary);
  border-color: var(--color-primary);
}

.process-note {
  display: flex;
  margin-top: 28px;
  padding: 13px 16px;
  align-items: flex-start;
  gap: 9px;
  color: #777;
  font-size: 13px;
  line-height: 20px;
  background: #fafafa;
  border-radius: 6px;

  > :first-child {
    margin-top: 3px;
    color: #999;
  }
}

@media (max-width: 760px) {
  .process-mode {
    grid-template-columns: 1fr;
  }
}
</style>
