<template>
  <a-modal
    v-model:open="visible"
    title="分配角色"
    width="460px"
    :confirm-loading="submitting"
    ok-text="确定"
    cancel-text="取消"
    destroy-on-close
    @ok="submit"
  >
    <a-form ref="formRef" :model="form" :rules="rules" layout="vertical">
      <a-form-item label="用户">
        <a-input :value="userLabel" disabled />
      </a-form-item>
      <a-form-item label="所属角色" name="roleId">
        <a-tree-select
          v-model:value="form.roleId"
          :tree-data="roleOptions"
          :loading="roleLoading"
          tree-default-expand-all
          :placeholder="rolePlaceholder"
        />
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import { computed, nextTick, reactive, ref } from 'vue'
import { message } from 'ant-design-vue'
import { getRoleOptions } from '@/api/role'
import { assignUserRole } from '@/api/user'

const emit = defineEmits(['success'])
const visible = ref(false)
const submitting = ref(false)
const roleLoading = ref(false)
const roleOptions = ref([])
const currentUser = ref()
const formRef = ref()
const form = reactive({ roleId: undefined })
const userLabel = computed(() => {
  if (!currentUser.value) return ''
  return `${currentUser.value.nickname}（${currentUser.value.username}）`
})
const rolePlaceholder = computed(() =>
  currentUser.value?.role_code &&
  (!currentUser.value?.role_id || currentUser.value.role_id === '0')
    ? '角色已禁用/删除'
    : '请选择角色',
)
const rules = {
  roleId: [{ required: true, message: '请选择角色' }],
}

async function loadRoleOptions() {
  roleLoading.value = true
  try {
    const data = await getRoleOptions()
    const mapOptions = (items) =>
      (items || []).map((item) => ({
        title: item.label,
        value: item.key,
        key: item.key,
        disabled: item.status !== 1,
        children: mapOptions(item.children),
      }))
    roleOptions.value = mapOptions(data)
  } catch (error) {
    message.error(error.message)
  } finally {
    roleLoading.value = false
  }
}

async function open(record) {
  currentUser.value = record
  form.roleId = record.role_id && record.role_id !== '0' ? record.role_id : undefined
  await loadRoleOptions()
  visible.value = true
  nextTick(() => formRef.value?.clearValidate())
}

async function submit() {
  try {
    await formRef.value.validate()
    submitting.value = true
    await assignUserRole(currentUser.value.id, { role_id: form.roleId })
    message.success('角色分配成功')
    visible.value = false
    emit('success')
  } catch (error) {
    if (error?.errorFields) return
    message.error(error.message)
  } finally {
    submitting.value = false
  }
}

defineExpose({ open })
</script>
