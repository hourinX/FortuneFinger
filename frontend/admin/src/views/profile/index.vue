<template>
  <div class="profile-page">
    <div class="profile-page__toolbar">
      <div></div>
      <div class="profile-page__actions">
        <a-button :disabled="submitting || uploading" @click="resetForm">重置</a-button>
        <a-button type="primary" :loading="submitting" :disabled="uploading" @click="submit">
          保存修改
        </a-button>
      </div>
    </div>
    <a-form ref="formRef" class="profile-page__form" :model="form" :rules="rules" layout="vertical">
      <section class="profile-page__avatar-panel">
        <a-popover
          v-model:open="avatarMenuOpen"
          trigger="click"
          placement="bottom"
          :overlay-inner-style="{ padding: '6px' }"
        >
          <template #content>
            <div class="profile-page__avatar-menu">
              <button type="button" @click="generateDefaultAvatar">
                <picture-outlined />
                <span>生成头像</span>
              </button>
              <a-divider />
              <button type="button" @click="openFilePicker">
                <upload-outlined />
                <span>本地上传</span>
              </button>
            </div>
          </template>
          <button
            class="profile-page__avatar-button"
            type="button"
            :disabled="uploading"
            aria-label="更换头像"
          >
            <a-avatar :size="104" :src="form.avatar">{{ accountInitial }}</a-avatar>
            <span class="profile-page__avatar-mask">
              <loading-outlined v-if="uploading" />
              <camera-outlined v-else />
              <span>{{ uploading ? '处理中' : '更换头像' }}</span>
            </span>
          </button>
        </a-popover>
        <input
          ref="fileInputRef"
          type="file"
          accept="image/*"
          hidden
          @change="uploadAvatar"
        />
        <strong>{{ form.nickname || admin.profile?.username || '用户' }}</strong>
        <span class="profile-page__username">{{ admin.profile?.username }}</span>
        <span class="profile-page__role">{{ roleName }}</span>
        <p>支持 JPG、PNG、WEBP 等图片格式，大小不超过 5MB</p>
      </section>

      <section class="profile-page__basic-panel">
        <div class="profile-page__section-title">
          <strong>基本资料</strong>
          <span>修改当前登录账号的展示信息与联系方式</span>
        </div>
        <div class="profile-page__fields">
          <a-form-item label="昵称" name="nickname">
            <a-input v-model:value="form.nickname" :maxlength="64" placeholder="请输入昵称" />
          </a-form-item>
          <a-form-item label="手机号码" name="mobile">
            <a-input v-model:value="form.mobile" :maxlength="32" placeholder="请输入手机号码" />
          </a-form-item>
          <a-form-item label="电子邮箱" name="email">
            <a-input v-model:value="form.email" :maxlength="254" placeholder="请输入电子邮箱" />
          </a-form-item>
          <a-form-item label="登录密码" name="password">
            <a-input-password
              v-model:value="form.password"
              :maxlength="64"
              placeholder="留空表示不修改登录密码"
            />
          </a-form-item>
        </div>
      </section>
    </a-form>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, reactive, ref } from 'vue'
import {
  CameraOutlined,
  LoadingOutlined,
  PictureOutlined,
  UploadOutlined,
} from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { generateAvatar, uploadAsset } from '@/api/asset'
import { updateUserInfo } from '@/api/user'
import { useAdminStore } from '@/store'

const admin = useAdminStore()
const formRef = ref()
const fileInputRef = ref()
const submitting = ref(false)
const uploading = ref(false)
const avatarMenuOpen = ref(false)
const form = reactive({
  avatarId: '0',
  avatar: '',
  nickname: '',
  email: '',
  mobile: '',
  password: '',
})
const rules = {
  nickname: [{ required: true, message: '请输入昵称' }],
  email: [{ type: 'email', message: '请输入正确的电子邮箱' }],
  password: [{ min: 6, message: '登录密码不能少于6位' }],
}
const accountInitial = computed(() =>
  String(form.nickname || admin.profile?.username || '用').slice(0, 1),
)
const roleName = computed(() =>
  admin.profile?.role_code === 'super_admin'
    ? '超级管理员'
    : admin.profile?.role_code || '普通用户',
)
let previewURL = ''

function revokePreviewURL() {
  if (!previewURL) return
  URL.revokeObjectURL(previewURL)
  previewURL = ''
}

function resetForm() {
  revokePreviewURL()
  Object.assign(form, {
    avatarId: admin.profile?.avatar_id || '0',
    avatar: admin.profile?.avatar || '',
    nickname: admin.profile?.nickname || '',
    email: admin.profile?.email || '',
    mobile: admin.profile?.mobile || '',
    password: '',
  })
  formRef.value?.clearValidate()
}

function openFilePicker() {
  avatarMenuOpen.value = false
  if (!uploading.value) fileInputRef.value?.click()
}

async function generateDefaultAvatar() {
  avatarMenuOpen.value = false
  try {
    uploading.value = true
    const resp = await generateAvatar()
    revokePreviewURL()
    form.avatarId = resp.storage_id
    form.avatar = resp.url
    message.success('头像生成成功，请保存修改')
  } catch (error) {
    message.error(error.message)
  } finally {
    uploading.value = false
  }
}

async function uploadAvatar(event) {
  const file = event.target.files?.[0]
  event.target.value = ''
  if (!file) return
  if (!file.type.startsWith('image/')) {
    message.error('头像只能上传图片文件')
    return
  }
  if (file.size > 5 * 1024 * 1024) {
    message.error('头像文件不能超过5MB')
    return
  }
  try {
    uploading.value = true
    const resp = await uploadAsset(file, 'avatar')
    revokePreviewURL()
    previewURL = URL.createObjectURL(file)
    form.avatarId = resp.storage_id
    form.avatar = previewURL
    message.success('头像上传成功，请保存修改')
  } catch (error) {
    message.error(error.message)
  } finally {
    uploading.value = false
  }
}

async function submit() {
  try {
    await formRef.value.validate()
    submitting.value = true
    await updateUserInfo({
      avatar_id: form.avatarId,
      nickname: form.nickname,
      email: form.email,
      mobile: form.mobile,
      password: form.password,
    })
    await admin.loadProfile()
    resetForm()
    message.success('个人信息修改成功')
  } catch (error) {
    if (error?.errorFields) return
    message.error(error.message)
  } finally {
    submitting.value = false
  }
}

resetForm()
onBeforeUnmount(revokePreviewURL)
</script>

<style scoped lang="less">
	.profile-page {
		min-width: 0;

		&__toolbar {
			display: flex;
			align-items: center;
			justify-content: space-between;
			min-height: 34px;
			padding-bottom: 16px;
			border-bottom: 1px solid #f0f0f0;
		}

		&__actions {
			display: flex;
			gap: 10px;
		}

		&__form {
			display: grid;
			grid-template-columns: 260px minmax(0, 680px);
			gap: 28px;
			padding-top: 20px;
		}

		&__avatar-panel,
		&__basic-panel {
			min-width: 0;
			border: 1px solid #efeded;
			border-radius: 8px;
		}

		&__avatar-panel {
			display: flex;
			align-items: center;
			align-self: start;
			flex-direction: column;
			padding: 28px 20px 24px;
			background: #fafafa;

			strong {
				margin-top: 14px;
				color: var(--color-text);
				font-size: 16px;
			}

			p {
				margin: 18px 0 0;
				color: #aaa;
				font-size: 12px;
				line-height: 1.7;
				text-align: center;
			}
		}

		&__avatar-button {
			position: relative;
			display: block;
			width: 104px;
			height: 104px;
			padding: 0;
			overflow: hidden;
			background: transparent;
			border: 0;
			border-radius: 50%;
			cursor: pointer;

			&:hover .profile-page__avatar-mask,
			&:focus-visible .profile-page__avatar-mask {
				opacity: 1;
			}
		}

		&__avatar-mask {
			position: absolute;
			inset: 0;
			display: flex;
			align-items: center;
			justify-content: center;
			flex-direction: column;
			gap: 5px;
			color: #fff;
			font-size: 12px;
			background: rgb(0 0 0 / 48%);
			opacity: 0;
			transition: opacity 160ms ease;

			:deep(.anticon) {
				font-size: 20px;
			}
		}

		&__avatar-menu {
			display: flex;
			flex-direction: column;
			width: 104px;

			:deep(.ant-divider) {
				margin: 2px 0;
			}

			button {
				display: flex;
				align-items: center;
				gap: 7px;
				width: 100%;
				padding: 7px 8px;
				color: var(--color-text);
				background: transparent;
				border: 0;
				border-radius: 5px;
				cursor: pointer;

				&:hover {
					color: var(--color-primary);
					background: var(--color-primary-light);
				}
			}
		}

		&__username,
		&__role {
			color: var(--color-muted);
			font-size: 12px;
		}

		&__username {
			margin-top: 2px;
		}

		&__role {
			margin-top: 9px;
			padding: 2px 8px;
			color: var(--color-primary);
			background: var(--color-primary-light);
			border-radius: 3px;
		}

		&__basic-panel {
			padding: 22px 24px 18px;
			background: #fff;
		}

		&__section-title {
			display: flex;
			flex-direction: column;
			padding-bottom: 20px;
			border-bottom: 1px solid #f2f0f0;

			strong {
				color: var(--color-text);
				font-size: 15px;
			}

			span {
				margin-top: 2px;
				color: var(--color-muted);
				font-size: 12px;
			}
		}

		&__fields {
			display: grid;
			grid-template-columns: repeat(2, minmax(0, 1fr));
			gap: 0 20px;
			padding-top: 20px;
		}

		@media (max-width: 900px) {
			&__form,
			&__fields {
				grid-template-columns: 1fr;
			}
		}
	}
</style>
