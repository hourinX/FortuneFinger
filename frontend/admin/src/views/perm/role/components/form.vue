<template>
  <a-modal
    v-model:open="visible"
    :title="editingID ? '编辑角色' : '新增角色'"
    width="560px"
    :confirm-loading="submitting"
    ok-text="确定"
    cancel-text="取消"
    destroy-on-close
    @ok="submit"
  >
    <a-form ref="formRef" class="role-form" :model="form" :rules="rules" layout="vertical">
      <a-form-item label="所属父级" name="parentId">
        <a-select
          v-model:value="form.parentId"
          :options="parentOptions"
          :loading="optionsLoading"
          :disabled="editingTopLevel"
          placeholder="请选择父角色"
        />
      </a-form-item>
      <a-form-item label="状态" name="status">
        <a-radio-group v-model:value="form.status">
          <a-radio :value="1">启用</a-radio>
          <a-radio :value="2" :disabled="editingTopLevel">禁用</a-radio>
        </a-radio-group>
      </a-form-item>
      <a-form-item label="角色名称" name="name">
        <a-input v-model:value="form.name" :maxlength="100" placeholder="请输入角色名称" />
      </a-form-item>
      <a-form-item label="角色编码" name="code">
        <a-input
          v-model:value="form.code"
          :maxlength="100"
          :disabled="Boolean(editingID)"
          placeholder="请输入角色编码"
        />
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import { computed, nextTick, reactive, ref } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { createRole, getRoleDetail, getRoleOptions, updateRole } from '@/api/role'

const emit = defineEmits(['success'])
const visible = ref(false)
const submitting = ref(false)
const editingID = ref()
const originalParentID = ref('0')
const formRef = ref()
const form = reactive({ parentId: '0', name: '', code: '', status: 1 })
const optionsLoading = ref(false)
const parentOptions = ref([{ label: '无（顶级角色）', value: '0' }])
const editingTopLevel = computed(() => Boolean(editingID.value) && originalParentID.value === '0')

const rules = {
  parentId: [{ required: true, message: '请选择父角色' }],
  name: [{ required: true, whitespace: true, message: '请输入角色名称' }],
  code: [{ required: true, whitespace: true, message: '请输入角色编码' }],
  status: [{ required: true, message: '请选择状态' }],
}

async function loadParentOptions() {
  optionsLoading.value = true
  try {
    const data = await getRoleOptions()
    const mapOptions = (items, depth = 0) =>
      items.flatMap((item) => {
        if (item.key === editingID.value) return []
        const prefix = depth === 0 ? '' : `${'　'.repeat(depth - 1)}└ `
        return [
          { label: `${prefix}${item.label}`, value: item.key },
          ...mapOptions(item.children || [], depth + 1),
        ]
      })
    parentOptions.value = [{ label: '无（顶级角色）', value: '0' }, ...mapOptions(data || [])]
  } catch (error) {
    message.error(error.message)
  } finally {
    optionsLoading.value = false
  }
}

async function openCreate(parent) {
  editingID.value = undefined
  originalParentID.value = '0'
  Object.assign(form, { parentId: parent?.id || '0', name: '', code: '', status: 1 })
  await loadParentOptions()
  visible.value = true
  nextTick(() => formRef.value?.clearValidate())
}

async function openEdit(id) {
  try {
    const record = await getRoleDetail(id)
    editingID.value = record.id
    originalParentID.value = record.parent_id
    await loadParentOptions()
    Object.assign(form, {
      parentId: record.parent_id,
      name: record.name,
      code: record.code,
      status: record.status,
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
      parent_id: form.parentId,
      name: form.name,
      status: form.status,
    }
    if (editingID.value) {
      if (form.parentId !== originalParentID.value && !(await confirmParentChange())) {
        return
      }
      await updateRole(editingID.value, payload)
      message.success('角色修改成功')
    } else {
      await createRole({ ...payload, code: form.code })
      message.success('角色新增成功')
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

function confirmParentChange() {
  return new Promise((resolve) => {
    Modal.confirm({
      title: '确认更改所属父级？',
      content:
        '更改上级角色后，当前角色及其全部下级角色的权限将被清空，需要重新分配。如果其中任意角色已绑定用户，系统将拒绝本次修改。',
      okText: '确认更改',
      cancelText: '取消',
      onOk: () => resolve(true),
      onCancel: () => resolve(false),
    })
  })
}

defineExpose({ openCreate, openEdit })
</script>

<style scoped lang="less">
.role-form {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  column-gap: 24px;
  padding-top: 8px;
}

@media (max-width: 520px) {
  .role-form {
    grid-template-columns: 1fr;
  }
}
</style>
