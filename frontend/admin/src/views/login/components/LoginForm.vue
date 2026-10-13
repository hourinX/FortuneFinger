<template>
  <a-form class="login-form" :model="form" @finish="submit">
    <a-form-item name="username" :rules="[{ required: true, message: '请输入账号' }]">
      <a-input
        v-model:value="form.username"
        size="large"
        placeholder="请输入账号"
        autocomplete="username"
      >
        <template #prefix><user-outlined /></template>
      </a-input>
    </a-form-item>
    <a-form-item name="password" :rules="[{ required: true, message: '请输入密码' }]">
      <a-input-password
        v-model:value="form.password"
        size="large"
        placeholder="请输入密码"
        autocomplete="current-password"
      >
        <template #prefix><lock-outlined /></template>
      </a-input-password>
    </a-form-item>
    <a-form-item name="captcha" :rules="[{ required: true, message: '请输入验证码' }]">
      <div class="captcha-row">
        <a-input
          v-model:value="form.captcha"
          size="large"
          placeholder="请输入验证码"
          autocomplete="off"
        >
          <template #prefix><safety-certificate-outlined /></template>
        </a-input>
        <button type="button" class="captcha-image" title="点击刷新验证码" @click="loadCaptcha">
          <img v-if="captchaImage" :src="captchaImage" alt="验证码" />
          <span v-else><reload-outlined /> 验证码</span>
        </button>
      </div>
    </a-form-item>
    <div class="form-options">
      <a-checkbox v-model:checked="remember">记住我</a-checkbox>
      <a class="forgot" @click.prevent="message.info('密码找回功能正在建设中')">忘记密码?</a>
    </div>
    <a-button class="submit" type="primary" html-type="submit" size="large" block :loading="loading"
      >登录</a-button
    >
  </a-form>
  <EnhancedVerification
    :open="verificationOpen"
	scene="login"
    :subject="form.username"
    @close="verificationOpen = false"
    @verified="onVerificationPassed"
  />
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { message } from 'ant-design-vue'
import {
  LockOutlined,
  ReloadOutlined,
  SafetyCertificateOutlined,
  UserOutlined,
} from '@ant-design/icons-vue'
import { useRoute, useRouter } from 'vue-router'
import { getCaptcha } from '@/api/auth'
import { useAdminStore, useAppStore } from '@/store'
import EnhancedVerification from '@/components/EnhancedVerification/index.vue'

const savedUsername = localStorage.getItem('fortune_admin_username') || ''
const form = reactive({
  username: savedUsername,
  password: '',
  captcha: '',
  captcha_id: '',
  verification_token: '',
})
const captchaImage = ref('')
const remember = ref(Boolean(savedUsername))
const loading = ref(false)
const verificationOpen = ref(false)
const admin = useAdminStore()
const app = useAppStore()
const route = useRoute()
const router = useRouter()

async function loadCaptcha() {
  try {
    const data = await getCaptcha()
    form.captcha_id = data?.captcha_id || data?.captchaId || data?.id || ''
    const image = data?.captcha || data?.captcha_img || data?.captchaImg || data?.image || ''
    captchaImage.value =
      image && !image.startsWith('data:') ? `data:image/png;base64,${image}` : image
  } catch {
    form.captcha_id = ''
    captchaImage.value = ''
  }
}

async function submit() {
  loading.value = true
  try {
    await admin.login({ ...form, remember: remember.value })
    try {
      await app.loadRouters(router)
    } catch (error) {
      await admin.logout()
      app.resetRouters()
      throw error
    }
    if (remember.value) localStorage.setItem('fortune_admin_username', form.username)
    else localStorage.removeItem('fortune_admin_username')
    await router.replace(String(route.query.redirect || '/dashboard'))
  } catch (error) {
	if (Number(error.code) === 10) {
		verificationOpen.value = true
		return
	}
    message.error(error.message || '登录失败')
	form.verification_token = ''
    form.captcha = ''
    await loadCaptcha()
  } finally {
    loading.value = false
  }
}

async function onVerificationPassed(token) {
	verificationOpen.value = false
	form.verification_token = token
	await submit()
}

onMounted(loadCaptcha)
</script>

<style scoped>
.login-form {
  width: 100%;
  max-width: 410px;
}
:deep(.ant-form-item) {
  margin-bottom: 18px;
}
:deep(.ant-input-affix-wrapper) {
  height: 52px;
  padding: 0 18px;
  color: #8d8d8d;
  background: #fff;
  border-color: #e5e5e5;
  border-radius: 10px;
}
:deep(.ant-input-affix-wrapper:hover),
:deep(.ant-input-affix-wrapper-focused) {
  border-color: var(--color-primary);
  box-shadow: 0 0 0 2px rgba(var(--color-primary-rgb), 0.08);
}
:deep(.ant-input-prefix) {
  margin-right: 17px;
  font-size: 20px;
}
.captcha-row {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 156px;
  gap: 12px;
  width: 100%;
}
.captcha-image {
  height: 52px;
  padding: 0;
  overflow: hidden;
  color: #a0a0a0;
  background: #fafafa;
  border: 1px solid #e5e5e5;
  border-radius: 10px;
  cursor: pointer;
}
.captcha-image:hover {
  color: var(--color-primary);
  border-color: var(--color-primary);
}
.captcha-image img {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: contain;
}
.captcha-image span {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
}
.form-options {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin: 8px 1px 34px;
  color: #8f8f8f;
  font-size: 13px;
}
.forgot {
  color: #8f8f8f;
}
.forgot:hover {
  color: var(--color-primary);
}
.submit {
  height: 48px;
  font-size: 17px;
  font-weight: 600;
  background: var(--color-primary);
  border-color: var(--color-primary);
  border-radius: 10px;
  box-shadow: none;
}
.submit:hover,
.submit:focus {
  background: var(--color-primary-hover) !important;
  border-color: var(--color-primary-hover) !important;
}
:deep(.submit:active) {
  background: var(--color-primary-active) !important;
  border-color: var(--color-primary-active) !important;
}
:deep(.ant-checkbox-checked .ant-checkbox-inner) {
  background: var(--color-primary);
  border-color: var(--color-primary);
}
</style>
