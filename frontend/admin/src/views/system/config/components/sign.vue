<template>
  <div class="settings-form">
    <div class="settings-group">
      <div class="settings-group__title">01 · 签到模式</div>
      <a-radio-group v-model:value="model.mode" class="sign-mode">
        <a-radio v-for="option in modeOptions" :key="option.value" :value="option.value" class="sign-mode__item">
          <span class="sign-mode__content">
            <span class="sign-mode__name">{{ option.label }}</span>
            <span class="sign-mode__description">{{ modeDescriptions[option.value] }}</span>
          </span>
        </a-radio>
      </a-radio-group>
    </div>

    <div class="settings-group">
      <div class="settings-group__title">02 · 签到信息</div>
      <SettingRow title="要求填写真实姓名" description="用于签到名单展示和人员匹配；名单模式下必须开启">
        <a-switch v-model:checked="model.requireName" :disabled="model.mode === 2" />
      </SettingRow>
      <SettingRow title="要求选择所属公司" description="支持集团年会、多家子公司或联合活动">
        <a-switch v-model:checked="model.requireCompany" />
      </SettingRow>
      <SettingRow title="允许上传个人头像" description="参与者可自行上传头像，用于照片墙和开奖展示">
        <a-switch v-model:checked="model.allowAvatar" />
      </SettingRow>
    </div>

    <div class="settings-group">
      <div class="settings-group__title">03 · 签到行为</div>
      <SettingRow title="重复签到处理" description="系统识别到同一参与者再次扫码时的处理方式">
        <a-select v-model:value="model.duplicateMode" class="settings-control" :options="duplicateOptions" />
      </SettingRow>
      <div class="sign-note">
        <strong>关于重复签到识别</strong>
        <p>
          自由签到下，如果没有账号、手机号验证或其他可靠身份凭据，只能识别同一签到凭证或设备会话，无法保证同一个人换设备后不会重复签到。名单模式也需要额外凭据才能防止冒用他人姓名。
        </p>
      </div>
    </div>

    <div class="settings-group">
      <div class="settings-group__title">04 · 签到展示</div>
      <SettingRow title="签到后进入照片墙" description="完成签到后，头像进入活动照片墙">
        <a-switch v-model:checked="model.enterPhotoWall" />
      </SettingRow>
      <SettingRow title="照片墙名称展示" description="设置照片墙中参与人员名称的展示方式">
        <a-select v-model:value="model.photoWallNameMode" class="settings-control" :options="nameOptions" />
      </SettingRow>
    </div>

    <div class="settings-group">
      <div class="settings-group__title">05 · 签到大屏模板</div>
      <TemplateSelector v-model="model.screen.template" />
    </div>
  </div>
</template>

<script setup>
import { watch } from 'vue'
import { useDictionaryOptions } from '@/hooks/useDictionary'
import SettingRow from './setting-row.vue'
import TemplateSelector from './template-selector.vue'

const model = defineModel({ type: Object, required: true })
const modeOptions = useDictionaryOptions('checkin_mode', true)
const modeDescriptions = {
  1: '无需提前导入人员名单，参与者填写信息后直接完成签到。',
  2: '管理员提前导入人员名单，参与者签到时需要匹配对应的公司及姓名。',
}
watch(
  () => model.value.mode,
  (mode) => {
    if (mode === 'free') {
      model.value.mode = 1
      return
    }
    if (mode === 'list') {
      model.value.mode = 2
      return
    }
    if (mode === 2) model.value.requireName = true
  },
  { immediate: true },
)
const duplicateOptions = [
  { label: '显示已有签到结果（推荐）', value: 'show' },
  { label: '允许更新签到资料，不新增参与记录', value: 'update' },
  { label: '提示已签到，不允许重复提交', value: 'reject' },
]
const nameOptions = [
  { label: '显示参与者姓名或昵称', value: 'name' },
  { label: '显示姓名或昵称及所属公司', value: 'name_company' },
  { label: '仅显示头像', value: 'avatar' },
]
</script>

<style scoped lang="less">
.sign-mode {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px;
  padding-top: 20px;
}

.sign-mode__item {
  position: relative;
  display: flex;
  min-height: 96px;
  margin: 0;
  padding: 18px 52px 18px 20px;
  align-items: center;
  background: #fff;
  border: 1px solid #e5e5e5;
  border-radius: 8px;
  cursor: pointer;
  transition: border-color 0.2s, background-color 0.2s;

  &:hover {
    border-color: #f2a5a2;
  }

  &:has(.ant-radio-checked) {
    background: #fff8f7;
    border-color: var(--color-primary);
  }

  :deep(.ant-radio) {
    position: absolute;
    top: 18px;
    right: 18px;
  }
}

.sign-mode__content {
  display: flex;
  flex-direction: column;
  gap: 6px;
  white-space: normal;
}

.sign-mode__name {
  color: #333;
  font-size: 14px;
  font-weight: 600;
  line-height: 22px;
}

.sign-mode__description {
  color: #8c8c8c;
  font-size: 13px;
  line-height: 20px;
}

.sign-note {
  margin-top: 16px;
  padding: 16px 18px;
  color: #666;
  font-size: 13px;
  line-height: 21px;
  background: #fafafa;
  border-left: 3px solid #d9d9d9;
  border-radius: 4px;

  strong {
    display: block;
    margin-bottom: 4px;
    color: #444;
    font-size: 14px;
  }

  p {
    margin: 0;
  }
}

@media (max-width: 720px) {
  .sign-mode {
    grid-template-columns: 1fr;
  }
}
</style>

