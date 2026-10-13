<template>
  <a-modal
    v-model:open="visible"
    :title="editingID ? '编辑组织' : '新增组织'"
    width="640px"
    :confirm-loading="submitting"
    ok-text="确定"
    cancel-text="取消"
    destroy-on-close
    @ok="submit"
  >
    <a-form ref="formRef" class="org-form" :model="form" :rules="rules" layout="vertical">
      <a-form-item label="所属父级" name="parentId">
        <a-select v-model:value="form.parentId" :options="parentOptions" />
      </a-form-item>
      <a-form-item label="组织类型" name="types">
        <a-select v-model:value="form.types" :options="typeOptions" />
      </a-form-item>
      <a-form-item label="组织名称" name="name">
        <a-input v-model:value="form.name" :maxlength="100" placeholder="请输入组织名称" />
      </a-form-item>
      <a-form-item label="组织编码" name="code">
        <a-input v-model:value="form.code" :maxlength="100" placeholder="请输入组织编码" />
      </a-form-item>
      <a-form-item label="排序" name="sort">
        <a-input-number v-model:value="form.sort" :min="0" :max="9999" />
      </a-form-item>
      <a-form-item label="状态" name="status">
        <a-radio-group v-model:value="form.status">
          <a-radio :value="1">启用</a-radio>
          <a-radio :value="2">禁用</a-radio>
        </a-radio-group>
      </a-form-item>
      <a-form-item class="org-form__remark" label="备注" name="remark">
        <a-textarea
          v-model:value="form.remark"
          :maxlength="255"
          :rows="3"
          show-count
          placeholder="请输入备注"
        />
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import { computed, nextTick, reactive, ref } from 'vue'
import { message } from 'ant-design-vue'
import { createOrg, getOrgDetail, updateOrg } from '@/api/org'

const props = defineProps({ organizations: { type: Array, default: () => [] } })
const emit = defineEmits(['success'])

const visible = ref(false)
const submitting = ref(false)
const editingID = ref()
const formRef = ref()
const form = reactive({
  parentId: '0',
  name: '',
  code: '',
  types: 1,
  sort: 0,
  status: 1,
  remark: '',
})
const typeOptions = [
  { label: '集团', value: 1 },
  { label: '公司', value: 2 },
  { label: '部门', value: 3 },
]
const parentOptions = computed(() => [
  { label: '无（顶级组织）', value: '0' },
  ...props.organizations
    .filter((item) => item.id !== editingID.value && item.types !== 3)
    .map((item) => ({ label: item.name, value: item.id })),
])
const rules = {
  parentId: [{ required: true, message: '请选择所属父级' }],
  types: [{ required: true, message: '请选择组织类型' }],
  name: [{ required: true, whitespace: true, message: '请输入组织名称' }],
  code: [{ required: true, whitespace: true, message: '请输入组织编码' }],
  sort: [{ required: true, message: '请输入排序' }],
  status: [{ required: true, message: '请选择状态' }],
}

function resetForm(record = {}) {
  Object.assign(form, {
    parentId: record.parentId ?? '0',
    name: record.name ?? '',
    code: record.code ?? '',
    types: record.types ?? 1,
    sort: record.sort ?? 0,
    status: record.status ?? 1,
    remark: record.remark ?? '',
  })
}

function openCreate(parent) {
  editingID.value = undefined
  const types = parent?.types === 1 ? 2 : parent?.types === 2 ? 3 : 1
  resetForm({ parentId: parent?.id || '0', types })
  visible.value = true
  nextTick(() => formRef.value?.clearValidate())
}

async function openEdit(record) {
  try {
    const detail = await getOrgDetail(record.id)
    editingID.value = detail.id
    resetForm({ ...detail, parentId: detail.parent_id })
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
      code: form.code,
      types: form.types,
      sort: form.sort,
      status: form.status,
      remark: form.remark,
    }
    if (editingID.value) {
      await updateOrg(editingID.value, payload)
      message.success('组织修改成功')
    } else {
      await createOrg(payload)
      message.success('组织新增成功')
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
.org-form {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  column-gap: 24px;
  padding-top: 8px;

  :deep(.ant-input-number) {
    width: 100%;
  }
}

.org-form__remark {
  grid-column: 1 / -1;
}

@media (max-width: 620px) {
  .org-form {
    grid-template-columns: 1fr;
  }

  .org-form__remark {
    grid-column: auto;
  }
}
</style>
