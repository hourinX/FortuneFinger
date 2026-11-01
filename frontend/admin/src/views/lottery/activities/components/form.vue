<template>
  <a-modal
    v-model:open="visible"
    :title="editingId ? '编辑活动' : '新增活动'"
    width="760px"
    :confirm-loading="submitting"
    ok-text="确定"
    cancel-text="取消"
    destroy-on-close
    @ok="submit"
  >
    <a-form ref="formRef" class="activity-form" layout="vertical" :model="form" :rules="rules">
      <div class="activity-form__grid">
        <a-form-item class="activity-form__wide" label="活动名称" name="name">
          <a-input
            v-model:value="form.name"
            :maxlength="100"
            placeholder="例如：2026 FortuneFinger 年度盛典"
          />
        </a-form-item>
        <a-form-item label="归属企业" name="orgId">
          <a-select v-model:value="form.orgId" :options="companies" placeholder="请选择归属企业" />
        </a-form-item>
        <a-form-item label="活动场地" name="venue">
          <a-input v-model:value="form.venue" :maxlength="100" placeholder="请输入活动场地" />
        </a-form-item>
        <a-form-item label="活动开始时间" name="startTime">
          <a-date-picker v-model:value="form.startTime" show-time format="YYYY-MM-DD HH:mm" />
        </a-form-item>
        <a-form-item label="活动结束时间" name="endTime">
          <a-date-picker v-model:value="form.endTime" show-time format="YYYY-MM-DD HH:mm" />
        </a-form-item>
        <a-form-item label="签到开始时间" name="checkinStartTime">
          <a-date-picker
            v-model:value="form.checkinStartTime"
            show-time
            format="YYYY-MM-DD HH:mm"
          />
        </a-form-item>
        <a-form-item label="签到结束时间" name="checkinEndTime">
          <a-date-picker v-model:value="form.checkinEndTime" show-time format="YYYY-MM-DD HH:mm" />
        </a-form-item>
        <a-form-item label="状态" name="status">
          <a-radio-group v-model:value="form.status">
            <a-radio :value="1">启用</a-radio>
            <a-radio :value="2">禁用</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item class="activity-form__wide" label="活动说明" name="remark">
          <a-textarea
            v-model:value="form.remark"
            :maxlength="255"
            :rows="3"
            show-count
            placeholder="请输入活动说明"
          />
        </a-form-item>
      </div>
    </a-form>
  </a-modal>
</template>

<script setup>
import { nextTick, reactive, ref } from 'vue'
import dayjs from 'dayjs'
import { message } from 'ant-design-vue'
import { createActivity, getActivityDetail, updateActivity } from '@/api/lottery'
import { DATE_TIME_FORMAT } from '@/utils/datetime'

defineProps({ companies: { type: Array, default: () => [] } })
const emit = defineEmits(['success'])
const visible = ref(false)
const submitting = ref(false)
const editingId = ref()
const formRef = ref()
const form = reactive(initialForm())
const rules = {
  name: [{ required: true, whitespace: true, message: '请输入活动名称' }],
  orgId: [{ required: true, message: '请选择归属企业' }],
  venue: [{ required: true, whitespace: true, message: '请输入活动场地' }],
  startTime: [{ required: true, message: '请选择活动开始时间' }],
  endTime: [{ required: true, message: '请选择活动结束时间' }],
  checkinStartTime: [{ required: true, message: '请选择签到开始时间' }],
  checkinEndTime: [{ required: true, message: '请选择签到结束时间' }],
  status: [{ required: true, message: '请选择状态' }],
}

function initialForm() {
  return {
    name: '',
    orgId: undefined,
    venue: '',
    startTime: undefined,
    endTime: undefined,
    checkinStartTime: undefined,
    checkinEndTime: undefined,
    checkinTemplate: 'invitation',
    status: 1,
    remark: '',
  }
}

function openCreate() {
  editingId.value = undefined
  Object.assign(form, initialForm())
  visible.value = true
  nextTick(() => formRef.value?.clearValidate())
}

async function openEdit(id) {
  try {
    const resp = await getActivityDetail(id)
    editingId.value = resp.id
    Object.assign(form, initialForm(), {
      name: resp.name,
      orgId: resp.org_id,
      venue: resp.venue,
      startTime: dayjs(resp.activity_start_time),
      endTime: dayjs(resp.activity_end_time),
      checkinStartTime: dayjs(resp.checkin_start_time),
      checkinEndTime: dayjs(resp.checkin_end_time),
      checkinTemplate: resp.checkin_template,
      status: resp.status,
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
    if (!form.startTime.isBefore(form.endTime)) {
      message.warning('活动结束时间必须晚于开始时间')
      return
    }
    if (!form.checkinStartTime.isBefore(form.checkinEndTime)) {
      message.warning('签到结束时间必须晚于开始时间')
      return
    }
    if (form.checkinEndTime.isAfter(form.endTime)) {
      message.warning('签到结束时间不能晚于活动结束时间')
      return
    }
    if (form.checkinEndTime.diff(form.checkinStartTime, 'day', true) > 7) {
      message.warning('签到时间范围不能超过七天')
      return
    }
    submitting.value = true
    const payload = {
      name: form.name,
      org_id: form.orgId,
      venue: form.venue,
      activity_start_time: form.startTime.second(0).format(DATE_TIME_FORMAT),
      activity_end_time: form.endTime.second(0).format(DATE_TIME_FORMAT),
      checkin_start_time: form.checkinStartTime.second(0).format(DATE_TIME_FORMAT),
      checkin_end_time: form.checkinEndTime.second(0).format(DATE_TIME_FORMAT),
      checkin_template: form.checkinTemplate,
      status: form.status,
      remark: form.remark,
    }
    if (editingId.value) {
      await updateActivity(editingId.value, payload)
      message.success('活动修改成功')
    } else {
      await createActivity(payload)
      message.success('活动新增成功')
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
.activity-form {
  padding-top: 8px;
}
.activity-form__grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0 20px;
}
.activity-form__wide {
  grid-column: 1 / -1;
}
.activity-form :deep(.ant-picker),
.activity-form :deep(.ant-select) {
  width: 100%;
}
</style>
