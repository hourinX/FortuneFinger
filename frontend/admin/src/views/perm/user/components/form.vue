<template>
  <a-modal
    v-model:open="visible"
    :title="editingID ? '编辑用户' : '新增用户'"
    width="620px"
    :confirm-loading="submitting"
    ok-text="确定"
    cancel-text="取消"
    destroy-on-close
    @ok="submit"
  >
    <a-form ref="formRef" class="user-form" :model="form" :rules="rules" layout="vertical">
      <a-form-item label="用户名" name="username">
        <a-input v-model:value="form.username" :maxlength="64" placeholder="请输入用户名" />
      </a-form-item>
      <a-form-item label="昵称" name="nickname">
        <a-input v-model:value="form.nickname" :maxlength="64" placeholder="请输入昵称" />
      </a-form-item>
      <a-form-item label="状态" name="status">
        <a-radio-group v-model:value="form.status">
          <a-radio :value="1">启用</a-radio>
          <a-radio :value="2">禁用</a-radio>
        </a-radio-group>
      </a-form-item>
      <a-form-item label="电子邮箱" name="email">
        <a-input v-model:value="form.email" :maxlength="254" placeholder="请输入电子邮箱" />
      </a-form-item>
      <a-form-item label="手机号码" name="mobile">
        <a-input v-model:value="form.mobile" :maxlength="32" placeholder="请输入手机号码" />
      </a-form-item>
      <a-form-item v-if="!isSuperAdmin" label="所属部门" name="departmentId">
        <a-tree-select
          v-model:value="form.departmentId"
          :tree-data="departmentOptions"
          allow-clear
          tree-default-expand-all
          placeholder="请选择所属部门"
        />
      </a-form-item>
      <a-form-item label="登录密码" name="password">
        <a-input-password
          v-model:value="form.password"
          :maxlength="64"
          :placeholder="editingID ? '留空表示不修改密码' : '请输入登录密码'"
        />
      </a-form-item>
      <div v-if="!isSuperAdmin" class="user-form__password-tip">
        <span>新增用户时必须设置密码。</span>
      </div>
    </a-form>
  </a-modal>
</template>

<script setup>
import { computed, nextTick, reactive, ref } from 'vue'
import { message } from 'ant-design-vue'
import { createUser, getUserDetail, updateUser } from '@/api/user'

const props = defineProps({ organizations: { type: Array, default: () => [] } })
const emit = defineEmits(['success'])
const visible = ref(false)
const submitting = ref(false)
const editingID = ref()
const roleCode = ref('')
const formRef = ref()
const isSuperAdmin = computed(() => roleCode.value === 'super_admin')
const form = reactive({
  username: '',
  nickname: '',
  avatar: '',
  status: 1,
  email: '',
  mobile: '',
  departmentId: undefined,
  password: '',
})
const departmentOptions = computed(() => {
  const buildTree = (organizations, depth = 0) =>
    organizations.map((organization) => ({
      title: organization.title,
      value: organization.key,
      selectable: depth >= 2,
      children: buildTree(organization.children || [], depth + 1),
    }))
  return buildTree(props.organizations)
})
const rules = {
  username: [{ required: true, message: '请输入用户名' }],
  nickname: [{ required: true, message: '请输入昵称' }],
  status: [{ required: true, message: '请选择状态' }],
  email: [{ type: 'email', message: '请输入正确的电子邮箱' }],
  departmentId: [
    {
      validator: async (_, value) => {
        if (!editingID.value && !value) throw new Error('请选择所属部门')
      },
    },
  ],
  password: [
    {
      validator: async (_, value) => {
        if (!editingID.value && !value) throw new Error('请输入登录密码')
        if (value && value.length < 6) throw new Error('登录密码不能少于6位')
      },
    },
  ],
}

function openCreate() {
  editingID.value = undefined
  roleCode.value = ''
  Object.assign(form, {
    username: '',
    nickname: '',
    avatar: '',
    status: 1,
    email: '',
    mobile: '',
    departmentId: undefined,
    password: '',
  })
  visible.value = true
  nextTick(() => formRef.value?.clearValidate())
}

async function openEdit(id) {
  try {
    const record = await getUserDetail(id)
    editingID.value = record.id
    roleCode.value = record.role_code || ''
    Object.assign(form, {
      username: record.username,
      nickname: record.nickname,
      avatar: record.avatar,
      status: record.status,
      email: record.email,
      mobile: record.mobile,
      departmentId: record.department_id || undefined,
      password: '',
    })
    visible.value = true
    nextTick(() => formRef.value?.clearValidate())
  } catch (error) {
    message.error(error.message)
  }
}

async function submit() {
  try {
    await formRef.value.validate()
    submitting.value = true
    const payload = {
      username: form.username,
      nickname: form.nickname,
      avatar: form.avatar,
      email: form.email,
      mobile: form.mobile,
      department_id: form.departmentId || '0',
      status: form.status,
      password: form.password,
    }
    if (editingID.value) {
      await updateUser(editingID.value, payload)
      message.success('用户修改成功')
    } else {
      await createUser(payload)
      message.success('用户新增成功')
    }
    visible.value = false
    emit('success')
  } catch (error) {
    if (error?.errorFields) return
    message.error(error.message)
  } finally {
    submitting.value = false
  }
}

defineExpose({ openCreate, openEdit })
</script>

<style scoped lang="less">
.user-form {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  column-gap: 24px;
  padding-top: 8px;
}

.user-form__password-tip {
  display: flex;
  align-items: center;
  height: 32px;
  margin-top: 30px;
  color: #999;
  font-size: 12px;
  line-height: 1.6;
}

@media (max-width: 560px) {
  .user-form {
    grid-template-columns: 1fr;
  }

  .user-form__password-tip {
    margin-top: 0;
  }
}
</style>
