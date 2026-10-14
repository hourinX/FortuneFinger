<template>
  <a-modal
    v-model:open="visible"
    :title="editingID ? '编辑字典数据' : '新增字典数据'"
    width="680px"
    :confirm-loading="submitting"
    ok-text="确定"
    cancel-text="取消"
    destroy-on-close
    @ok="submit"
  >
    <a-form ref="formRef" class="dict-data-form" :model="form" :rules="rules" layout="vertical">
      <a-form-item label="字典类型">
        <a-input :value="`${dictType.name}（${dictType.code}）`" disabled />
      </a-form-item>
      <a-form-item label="字典标签" name="label">
        <a-input v-model:value="form.label" :maxlength="100" placeholder="请输入字典标签" />
      </a-form-item>
      <a-form-item label="字典值" name="value">
        <a-input v-model:value="form.value" :maxlength="255" placeholder="请输入字典值" />
      </a-form-item>
      <a-form-item label="展示颜色" name="color">
        <a-auto-complete
          v-model:value="form.color"
          :options="colorOptions"
          placeholder="选择预设颜色或输入 #RRGGBB"
        >
          <template #option="option">
            <a-tag :color="tagColor(option.value)">{{ option.label }}</a-tag>
          </template>
        </a-auto-complete>
      </a-form-item>
      <a-form-item label="排序" name="sort">
        <a-input-number v-model:value="form.sort" :min="0" :max="9999" />
      </a-form-item>
      <a-form-item label="默认选项" name="isDefault">
        <a-radio-group v-model:value="form.isDefault">
          <a-radio :value="1">是</a-radio>
          <a-radio :value="2">否</a-radio>
        </a-radio-group>
      </a-form-item>
      <a-form-item label="状态" name="status">
        <a-radio-group v-model:value="form.status">
          <a-radio :value="1">启用</a-radio>
          <a-radio :value="2">禁用</a-radio>
        </a-radio-group>
      </a-form-item>
      <a-form-item class="dict-data-form__remark" label="备注" name="remark">
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
import { createDictData, getDictDataDetail, updateDictData } from '@/api/dict'

const emit = defineEmits(['success'])
const visible = ref(false)
const submitting = ref(false)
const editingID = ref()
const formRef = ref()
const dictType = reactive({ name: '', code: '' })
const form = reactive({
  label: '',
  value: '',
  color: '',
  sort: 0,
  isDefault: 2,
  status: 1,
  remark: '',
})
const colorOptions = [
  { label: '默认', value: 'default' },
  { label: '蓝色', value: 'blue' },
  { label: '绿色', value: 'green' },
  { label: '橙色', value: 'orange' },
  { label: '红色', value: 'red' },
  { label: '紫色', value: 'purple' },
]
const rules = {
  label: [{ required: true, whitespace: true, message: '请输入字典标签' }],
  value: [{ required: true, whitespace: true, message: '请输入字典值' }],
  color: [
    { required: true, message: '请选择或输入展示颜色' },
    {
      pattern: /^(default|blue|green|orange|red|purple|#[0-9a-fA-F]{6})$/,
      message: '请输入有效的十六进制颜色码，例如 #E43D36',
    },
  ],
  sort: [{ required: true, message: '请输入排序' }],
  isDefault: [{ required: true, message: '请选择是否为默认选项' }],
  status: [{ required: true, message: '请选择状态' }],
}

function tagColor(color) {
  return color === 'default' ? undefined : color
}

function resetForm(record = {}) {
  Object.assign(form, {
    label: record.label || '',
    value: record.value || '',
    color: record.color || '',
    sort: record.sort || 0,
    isDefault: record.is_default || 2,
    status: record.status || 1,
    remark: record.remark || '',
  })
}

function openCreate(type) {
  editingID.value = undefined
  Object.assign(dictType, type)
  resetForm()
  visible.value = true
  nextTick(() => formRef.value?.clearValidate())
}

async function openEdit(type, id) {
  try {
    const record = await getDictDataDetail(id)
    editingID.value = record.id
    Object.assign(dictType, type)
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
    const payload = {
      label: form.label,
      value: form.value,
      color: form.color,
      is_default: form.isDefault,
      sort: form.sort,
      status: form.status,
      remark: form.remark,
    }
    if (editingID.value) {
      await updateDictData(editingID.value, payload)
      message.success('字典数据修改成功')
    } else {
      await createDictData({ ...payload, dict_code: dictType.code })
      message.success('字典数据新增成功')
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
.dict-data-form {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  column-gap: 24px;
  padding-top: 8px;

  :deep(.ant-input-number) {
    width: 100%;
  }

  .dict-data-form__remark {
    grid-column: 1 / -1;
  }
}

@media (max-width: 680px) {
  .dict-data-form {
    grid-template-columns: 1fr;

    .dict-data-form__remark {
      grid-column: auto;
    }
  }
}
</style>
