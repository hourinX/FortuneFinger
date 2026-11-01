<script setup lang="ts">
import { onBeforeUnmount, reactive, ref } from 'vue'
import type { CheckinSessionResp } from '@/api/checkin'
import type { OrgCompanyOption } from '@/api/org'
import type { CheckinFormValue } from '@/stores/checkin'

const props = defineProps<{
  session: CheckinSessionResp
  companies: OrgCompanyOption[]
}>()
const emit = defineEmits<{
  submit: [value: CheckinFormValue]
}>()
const value = reactive<CheckinFormValue>({
  name: '',
  companyId: '',
  companyName: '',
  phone: '',
})
const avatarPreview = ref('')

function selectAvatar(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file || !file.type.startsWith('image/') || file.size > 5 * 1024 * 1024) return
  if (avatarPreview.value) URL.revokeObjectURL(avatarPreview.value)
  value.avatar = file
  avatarPreview.value = URL.createObjectURL(file)
}

function confirm() {
  const selected = props.companies.find(item => item.id === value.companyId)
  if (selected) value.companyName = selected.name
  emit('submit', { ...value })
}

onBeforeUnmount(() => {
  if (avatarPreview.value) URL.revokeObjectURL(avatarPreview.value)
})
</script>

<template>
  <form class="checkin-form" @submit.prevent="confirm">
    <label v-if="session.config.allow_avatar" class="avatar-picker">
      <input type="file" accept="image/*" @change="selectAvatar" />
      <span class="avatar-picker__preview">
        <img v-if="avatarPreview" :src="avatarPreview" alt="头像预览" />
        <span v-else>+</span>
      </span>
      <span class="avatar-picker__text">
        <strong>选择个人头像</strong>
        <small>用于现场签到动态与照片墙展示</small>
      </span>
    </label>

    <div class="field-group">
      <label class="field">
        <span>真实姓名</span>
        <span class="field__control">
          <input
            v-model="value.name"
            maxlength="30"
            autocomplete="name"
            placeholder="用于抽奖开奖人姓名展示"
          />
        </span>
      </label>

      <label class="field">
        <span>手机号码</span>
        <span class="field__control">
          <input
            v-model="value.phone"
            type="tel"
            maxlength="20"
            inputmode="tel"
            autocomplete="tel"
            placeholder="用于辅助识别重复签到"
          />
        </span>
      </label>

      <label v-if="session.config.require_company" class="field">
        <span>所属企业</span>
        <select v-model="value.companyId" :class="{ 'field__select--placeholder': !value.companyId }">
          <option value="" disabled>
            {{ companies.length ? '选择所属公司' : '暂无可选企业' }}
          </option>
          <option v-for="company in companies" :key="company.id" :value="company.id">
            {{ company.name }}
          </option>
        </select>
      </label>
    </div>

    <button type="submit">确认签到</button>
  </form>
</template>

<style scoped>
.checkin-form { display: flex; padding: 12px 22px 24px; flex-direction: column; gap: 14px; }
.avatar-picker { display: flex; padding: 6px 2px 10px; flex-direction: column; align-items: center; gap: 7px; text-align: center; cursor: pointer; }
.avatar-picker input { display: none; }
.avatar-picker__preview { display: grid; width: 78px; height: 78px; flex: 0 0 78px; overflow: hidden; place-items: center; color: #d9443d; font-size: 30px; background: #f7f3ed; border-radius: 50%; box-shadow: inset 0 0 0 1px rgb(205 193 179 / 35%), 0 10px 22px rgb(84 66 52 / 7%); }
.avatar-picker__preview img { width: 100%; height: 100%; object-fit: cover; }
.avatar-picker__text { display: flex; min-width: 0; flex-direction: column; gap: 3px; }
.avatar-picker__text strong { color: #403b37; font-size: 11px; }
.avatar-picker__text small { color: #a39b94; font-size: 9px; }
.field-group { padding: 0 2px; }
.field { display: grid; min-height: 54px; grid-template-columns: 74px minmax(0, 1fr); align-items: center; gap: 12px; border-bottom: 1px solid #e6ded4; transition: border-color .2s ease; }
.field:last-child { border-bottom: 0; }
.field:focus-within { border-bottom-color: #dc5a53; }
.field > span { color: #625b55; font-size: 12px; font-weight: 600; transition: color .2s ease; }
.field:focus-within > span { color: #d5443d; }
.field__control { display: flex; min-width: 0; }
.field input, .field select { width: 100%; height: 34px; padding: 0; box-sizing: border-box; color: #332f2c; font: inherit; font-size: 16px; outline: none; background: transparent; border: 0; border-radius: 0; box-shadow: none; }
.field input:focus, .field select:focus { background: transparent; box-shadow: none; }
.field input::placeholder { color: #b7afa8; font-size: 13px; font-weight: 400; opacity: 1; }
.field__select--placeholder { color: #b7afa8 !important; font-size: 13px !important; font-weight: 400; }
.field > select, .field > input { height: 52px; }
.checkin-form button { height: 50px; margin-top: 3px; color: #fff; font-size: 15px; font-weight: 650; letter-spacing: .06em; background: linear-gradient(100deg, #ef4139, #da312f); border: 0; border-radius: 13px; box-shadow: 0 13px 26px rgb(219 49 46 / 20%); }
</style>
