<template>
  <a-form class="settings-form" layout="vertical">
    <div class="settings-group">
      <div class="settings-group__title">01 · 邮件服务</div>
      <SettingRow title="启用邮件服务" description="开启后，系统可以发送验证码和业务通知邮件">
        <a-switch v-model:checked="model.enabled" :checked-value="1" :un-checked-value="2" />
      </SettingRow>
      <div class="mail-grid">
        <a-form-item label="SMTP 服务器地址">
          <a-input v-model:value="model.host" :maxlength="255" placeholder="smtp-relay.brevo.com" />
        </a-form-item>
        <a-form-item label="SMTP 端口">
          <a-input-number v-model:value="model.port" :min="1" :max="65535" />
        </a-form-item>
        <a-form-item label="连接安全">
          <a-select v-model:value="model.security" :options="securityOptions" />
        </a-form-item>
        <a-form-item label="连接超时">
          <a-input-number v-model:value="model.timeout" :min="1" :max="120" addon-after="秒" />
        </a-form-item>
      </div>
      <p class="mail-tip">Brevo 推荐使用 smtp-relay.brevo.com、587 端口和 STARTTLS。</p>
    </div>

    <div class="settings-group">
      <div class="settings-group__title mail-section-title">
        <span>02 · 认证信息</span>
        <a-popover
          content="认证信息保存后由服务端加密存储，重新加载页面时仅展示脱敏内容。"
          :trigger="['hover', 'click']"
          placement="right"
        >
          <button type="button" class="mail-section-title__help" aria-label="认证信息安全说明">
            <exclamation-circle-outlined />
          </button>
        </a-popover>
      </div>
      <div class="mail-grid mail-grid--section">
        <a-form-item label="SMTP Login">
          <a-input
            v-model:value="model.username"
            :maxlength="255"
            autocomplete="off"
            placeholder="请输入 Brevo SMTP & API 页面显示的 Login"
          />
        </a-form-item>
        <a-form-item label="SMTP Key">
          <a-input-password
            v-model:value="model.password"
            :maxlength="255"
            autocomplete="new-password"
            placeholder="请输入 Brevo 生成的 SMTP Key"
          />
        </a-form-item>
      </div>
      <p class="mail-tip">
        SMTP Login 以 Brevo 后台显示值为准；SMTP Key 不是 Brevo、Google 或邮箱账户的登录密码。
      </p>
    </div>

    <div class="settings-group">
      <div class="settings-group__title">03 · 发件人信息</div>
      <div class="mail-grid mail-grid--section">
        <a-form-item label="发件人名称">
          <a-input
            v-model:value="model.senderName"
            :maxlength="100"
            placeholder="例如 FortuneFinger"
          />
        </a-form-item>
        <a-form-item label="已验证发件人邮箱">
          <a-input
            v-model:value="model.senderEmail"
            :maxlength="255"
            placeholder="请输入已在 Brevo 验证的发件人邮箱"
          />
        </a-form-item>
        <a-form-item label="回复邮箱（Reply-To）">
          <a-input
            v-model:value="model.replyTo"
            :maxlength="255"
            placeholder="留空时使用发件人邮箱"
          />
        </a-form-item>
        <a-form-item label="测试收件邮箱">
          <div class="mail-test-field">
            <a-input
              :value="testEmail"
              :maxlength="254"
              placeholder="请输入接收测试邮件的邮箱"
              @update:value="emit('update:test-email', $event)"
            />
            <a-button
              v-permission="'system:config:mail'"
              type="primary"
              :loading="testLoading"
              :disabled="!testEmail"
              @click="emit('send-test', testEmail)"
            >
              发送测试邮件
            </a-button>
          </div>
        </a-form-item>
      </div>
      <p class="mail-tip">发件人邮箱必须已在 Brevo 中完成验证；回复邮箱留空时使用发件人邮箱。</p>
    </div>
  </a-form>
</template>

<script setup>
import { ExclamationCircleOutlined } from '@ant-design/icons-vue'
import { useDictionaryOptions } from '@/hooks/useDictionary'
import SettingRow from './setting-row.vue'

const model = defineModel({ type: Object, required: true })
defineProps({
  testEmail: { type: String, default: '' },
  testLoading: { type: Boolean, default: false },
})
const emit = defineEmits(['send-test', 'update:test-email'])
const securityOptions = useDictionaryOptions('mail_secury')
</script>

<style scoped lang="less">
.mail-grid {
  display: grid;
  padding-top: 22px;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0 20px;
}

.mail-grid--section {
  padding-top: 20px;
}

:deep(.ant-input-number),
:deep(.ant-select) {
  width: 100%;
}

.mail-tip {
  margin: -4px 0 0;
  color: #8c8c8c;
  font-size: 12px;
  line-height: 20px;
}

.mail-test-field {
  display: flex;
  gap: 10px;

  .ant-input {
    min-width: 0;
  }

  .ant-btn {
    flex: none;
  }
}

.mail-section-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.mail-section-title__help {
  display: inline-flex;
  padding: 0;
  align-items: center;
  color: #999;
  font-size: 14px;
  line-height: 1;
  background: transparent;
  border: 0;
  cursor: pointer;

  &:hover,
  &:focus-visible {
    color: var(--color-primary);
  }
}

@media (max-width: 720px) {
  .mail-grid {
    grid-template-columns: 1fr;
  }

  .mail-test-field {
    align-items: stretch;
  }
}
</style>
