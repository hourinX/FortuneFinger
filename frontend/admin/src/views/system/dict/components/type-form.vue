<template>
  <a-modal
    v-model:open="visible"
    :title="editingID ? '编辑字典类型' : '新增字典类型'"
    width="620px"
    :confirm-loading="submitting"
    ok-text="确定"
    cancel-text="取消"
    destroy-on-close
    @ok="submit"
  >
    <a-form ref="formRef" class="dict-type-form" :model="form" :rules="rules" layout="vertical">
      <a-form-item label="字典名称" name="name">
        <a-input v-model:value="form.name" :maxlength="100" placeholder="请输入字典名称" />
      </a-form-item>
      <a-form-item label="字典编码" name="code">
        <a-input
          v-model:value="form.code"
          :disabled="Boolean(editingID)"
          :maxlength="100"
          placeholder="例如 sys_status"
        />
      </a-form-item>
      <a-form-item label="状态" name="status">
        <a-radio-group v-model:value="form.status">
          <a-radio :value="1">启用</a-radio>
          <a-radio :value="2">禁用</a-radio>
        </a-radio-group>
      </a-form-item>
      <a-form-item class="dict-type-form__remark" label="备注" name="remark">
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
import { nextTick, reactive, ref } from 'vue'
import { message } from 'ant-design-vue'
import { createDictType, getDictTypeDetail, updateDictType } from '@/api/dict'

const emit = defineEmits(['success'])
const visible = ref(false)
const submitting = ref(false)
const editingID = ref()
const formRef = ref()
const form = reactive({ name: '', code: '', status: 1, remark: '' })
const rules = {
  name: [{ required: true, whitespace: true, message: '请输入字典名称' }],
  code: [
    { required: true, whitespace: true, message: '请输入字典编码' },
    { pattern: /^[a-z][a-z0-9_]*$/, message: '字典编码只能使用小写字母、数字和下划线' },
  ],
  status: [{ required: true, message: '请选择状态' }],
}

function resetForm(record = {}) {
  Object.assign(form, {
    name: record.name || '',
    code: record.code || '',
    status: record.status || 1,
    remark: record.remark || '',
  })
}

function openCreate() {
  editingID.value = undefined
  resetForm()
  visible.value = true
  nextTick(() => formRef.value?.clearValidate())
}

async function openEdit(id) {
  try {
    const record = await getDictTypeDetail(id)
    editingID.value = record.id
    resetForm(record)
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
    const payload = { name: form.name, status: form.status, remark: form.remark }
    if (editingID.value) {
      await updateDictType(editingID.value, payload)
      message.success('字典类型修改成功')
    } else {
      await createDictType({ ...payload, code: form.code })
      message.success('字典类型新增成功')
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
.dict-type-form {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  column-gap: 24px;
  padding-top: 8px;

  .dict-type-form__remark {
    grid-column: 1 / -1;
  }
}

@media (max-width: 620px) {
  .dict-type-form {
    grid-template-columns: 1fr;

    .dict-type-form__remark {
      grid-column: auto;
    }
  }
}
</style>
