<template>
  <a-modal
    v-model:open="visible"
    :title="editingID ? '编辑奖品' : '新增奖品'"
    width="700px"
    :confirm-loading="submitting"
    ok-text="确定"
    cancel-text="取消"
    destroy-on-close
    @ok="submit"
  >
    <a-form ref="formRef" class="prize-form" layout="vertical" :model="form" :rules="rules">
      <div class="prize-form__main">
        <ImageUpload
          v-model="form.storageId"
          class="prize-image-upload"
          scene="prize"
          :src="form.imageUrl"
          @removed="removeImage"
        >
          <template #default="{ preview, hasImage }">
            <div class="cover-upload">
              <img v-if="preview" :src="preview" alt="奖品图片预览" />
              <span v-else class="cover-upload__placeholder">
                <span>上传奖品图片</span>
                <small>{{ hasImage ? '已设置奖品图片' : '建议使用 1:1 图片' }}</small>
              </span>
            </div>
          </template>
        </ImageUpload>

        <div class="prize-form__identity">
          <a-form-item label="奖品名称" name="name">
            <a-input v-model:value="form.name" :maxlength="100" placeholder="请输入奖品名称" />
          </a-form-item>
          <div class="prize-form__identity-row">
            <a-form-item label="奖品类型" name="types">
              <a-select
                v-model:value="form.types"
                :options="typeOptions"
                placeholder="请选择奖品类型"
              />
            </a-form-item>
            <a-form-item label="归属企业" name="orgId">
              <a-select
                v-model:value="form.orgId"
                :options="companyOptions"
                placeholder="请选择归属企业"
              />
            </a-form-item>
          </div>
        </div>
      </div>

      <div class="prize-form__grid">
        <a-form-item label="奖项等级" name="level">
          <a-select v-model:value="form.level" :options="levelOptions" placeholder="请选择奖项等级" />
        </a-form-item>
        <a-form-item label="奖品数量" name="quantity">
          <a-input-number v-model:value="form.quantity" :min="1" :max="999999" />
        </a-form-item>
        <a-form-item label="展示排序" name="sort">
          <a-input-number v-model:value="form.sort" :min="0" :max="9999" />
        </a-form-item>
        <a-form-item label="状态" name="status">
          <a-radio-group v-model:value="form.status">
            <a-radio :value="1">启用</a-radio>
            <a-radio :value="2">禁用</a-radio>
          </a-radio-group>
        </a-form-item>
      </div>

      <a-form-item label="奖品说明" name="remark">
        <a-textarea
          v-model:value="form.remark"
          :maxlength="255"
          :rows="3"
          show-count
          placeholder="请输入奖品说明"
        />
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import { nextTick, reactive, ref } from 'vue'
import { message } from 'ant-design-vue'
import { getOrgCompanies } from '@/api/org'
import { createPrize, getPrizeDetail, updatePrize } from '@/api/prize'
import ImageUpload from '@/components/ImageUpload/index.vue'
import { useDictionaryOptions } from '@/hooks/useDictionary'

const emit = defineEmits(['success'])

const levelOptions = useDictionaryOptions('prize_level')
const typeOptions = useDictionaryOptions('prize_types', true)
const companyOptions = ref([])
const companiesLoaded = ref(false)
const visible = ref(false)
const submitting = ref(false)
const editingID = ref()
const formRef = ref()
const form = reactive(initialForm())
const rules = {
  name: [{ required: true, whitespace: true, message: '请输入奖品名称' }],
  level: [{ required: true, message: '请选择奖项等级' }],
  types: [{ required: true, message: '请选择奖品类型' }],
  orgId: [{ required: true, message: '请选择归属企业' }],
  quantity: [{ required: true, type: 'number', min: 1, message: '奖品数量至少为1' }],
  status: [{ required: true, message: '请选择状态' }],
  sort: [{ required: true, message: '请输入展示排序' }],
}

function initialForm() {
  return {
    name: '',
    storageId: '0',
    imageUrl: '',
    orgId: undefined,
    level: undefined,
    types: undefined,
    quantity: 1,
    status: 1,
    sort: 0,
    remark: '',
  }
}

async function openCreate() {
  try {
    await loadCompanyOptions()
    editingID.value = undefined
    Object.assign(form, initialForm())
    visible.value = true
    nextTick(() => formRef.value?.clearValidate())
  } catch (error) {
    message.error(error.message)
  }
}

async function openEdit(id) {
  try {
    const [resp] = await Promise.all([getPrizeDetail(id), loadCompanyOptions()])
    editingID.value = resp.id
    Object.assign(form, initialForm(), {
      name: resp.name,
      storageId: resp.storage_id,
      imageUrl: resp.image_url,
      orgId: resp.org_id,
      level: resp.level,
      types: resp.types,
      quantity: resp.quantity,
      status: resp.status,
      sort: resp.sort,
      remark: resp.remark,
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
      name: form.name,
      storage_id: form.storageId,
      org_id: form.orgId,
      level: form.level,
      types: form.types,
      quantity: form.quantity,
      status: form.status,
      sort: form.sort,
      remark: form.remark,
    }
    if (editingID.value) {
      await updatePrize(editingID.value, payload)
      message.success('奖品修改成功')
    } else {
      await createPrize(payload)
      message.success('奖品新增成功')
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

function removeImage() {
  form.imageUrl = ''
}

async function loadCompanyOptions() {
  if (companiesLoaded.value) return
  const resp = await getOrgCompanies()
  companyOptions.value = (resp.list || []).map((company) => ({
    label: company.name,
    value: company.id,
  }))
  companiesLoaded.value = true
}

defineExpose({ openCreate, openEdit })
</script>

<style scoped lang="less">
.prize-form {
  padding-top: 8px;
}

.prize-form__main {
  display: grid;
  grid-template-columns: 132px minmax(0, 1fr);
  align-items: center;
  gap: 20px;
  margin-bottom: 4px;
}

.prize-image-upload,
.cover-upload {
  width: 132px;
  height: 132px;
}

.cover-upload {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 7px;
  padding: 0;
  color: #777;
  background: #fafafa;
  border: 1px dashed #d8d8d8;
  border-radius: 8px;
  overflow: hidden;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

.cover-upload__placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;

  small {
    color: #aaa;
  }
}

.prize-form__identity {
  min-width: 0;
}

.prize-form__identity-row {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0 20px;
}

.prize-form__grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0 20px;
}

.prize-form :deep(.ant-input-number),
.prize-form :deep(.ant-select) {
  width: 100%;
}

@media (max-width: 700px) {
  .prize-form__main,
  .prize-form__grid,
  .prize-form__identity-row {
    grid-template-columns: 1fr;
  }
}
</style>
