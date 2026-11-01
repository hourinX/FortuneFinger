<template>
  <a-modal
    v-model:open="visible"
    width="680px"
    :confirm-loading="submitting"
    ok-text="确定"
    cancel-text="取消"
    destroy-on-close
    @ok="submit"
  >
    <template #title>
      <div class="round-modal-title">
        <span>{{ editingId ? '编辑抽奖轮次' : '新增抽奖轮次' }}</span>
        <a-tooltip placement="right">
          <template #title>
            抽奖轮次进入执行阶段后将禁止修改，请在保存前仔细核对所属活动、奖品、开奖时间和中奖人数。
          </template>
          <exclamation-circle-outlined class="round-modal-title__notice" />
        </a-tooltip>
      </div>
    </template>
    <a-form ref="formRef" class="round-form" layout="vertical" :model="form" :rules="rules">
      <div class="round-form__grid">
        <a-form-item label="所属活动" name="activityId">
          <a-select
            v-model:value="form.activityId"
            :options="activities"
            placeholder="请选择所属活动"
            @change="changeActivity"
          />
          <span v-if="selectedActivity" class="field-description">
            活动时间：{{ selectedActivity.startDate }} 至 {{ selectedActivity.endDate }}
          </span>
        </a-form-item>
        <a-form-item label="轮次名称" name="name">
          <a-input v-model:value="form.name" :maxlength="100" placeholder="例如：第一轮抽奖" />
        </a-form-item>
        <a-form-item label="奖品" name="prizeId">
          <a-select
            v-model:value="form.prizeId"
            :options="availablePrizeOptions"
            :disabled="!form.activityId"
            placeholder="请选择奖品"
            @change="changePrize"
          />
        </a-form-item>
        <a-form-item label="奖项等级" name="level">
          <a-select
            v-model:value="form.level"
            :options="prizeLevelOptions"
            disabled
            placeholder="选择奖品后自动填写"
          />
        </a-form-item>
        <a-form-item label="开奖时间" name="drawAt">
          <a-date-picker
            v-model:value="form.drawAt"
            :disabled-date="disabledRoundDate"
            :show-time="{ format: 'HH:mm' }"
            format="YYYY-MM-DD HH:mm"
            placeholder="请选择开奖时间"
          />
        </a-form-item>
        <a-form-item label="本轮中奖人数" name="winnerCount">
          <a-input-number v-model:value="form.winnerCount" :min="1" :max="999" />
        </a-form-item>
      </div>
    </a-form>
  </a-modal>
</template>

<script setup>
import { computed, nextTick, reactive, ref } from 'vue'
import dayjs from 'dayjs'
import { ExclamationCircleOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { createRound, getRoundDetail, updateRound } from '@/api/lottery'
import { getPrizeOptions } from '@/api/prize'
import { useDictionaryOptions } from '@/hooks/useDictionary'

const props = defineProps({ activities: { type: Array, default: () => [] } })
const emit = defineEmits(['success'])
const prizeLevelOptions = useDictionaryOptions('prize_level')
const visible = ref(false)
const submitting = ref(false)
const editingId = ref()
const formRef = ref()
const prizeOptions = ref([])
const form = reactive(initialForm())
const rules = {
  activityId: [{ required: true, message: '请选择所属活动' }],
  name: [{ required: true, whitespace: true, message: '请输入轮次名称' }],
  prizeId: [{ required: true, message: '请选择奖品' }],
  level: [{ required: true, message: '奖项等级不能为空' }],
  drawAt: [{ required: true, message: '请选择开奖时间' }],
  winnerCount: [{ required: true, message: '请输入本轮中奖人数' }],
}
const selectedActivity = computed(() =>
  props.activities.find((item) => item.value === form.activityId),
)
const availablePrizeOptions = computed(() =>
  prizeOptions.value.map((item) => ({ label: item.name, value: item.id })),
)

function initialForm() {
  return {
    activityId: undefined,
    name: '',
    drawAt: undefined,
    level: undefined,
    prizeId: undefined,
    winnerCount: 1,
  }
}

function openCreate() {
  editingId.value = undefined
  prizeOptions.value = []
  Object.assign(form, initialForm())
  visible.value = true
  nextTick(() => formRef.value?.clearValidate())
}

async function openEdit(id) {
  try {
    const detail = await getRoundDetail(id)
    const activity = props.activities.find((item) => item.value === detail.activity_id)
    await loadPrizeOptions(activity?.orgId)
    const prize = prizeOptions.value.find((item) => item.id === detail.prize_id)
    editingId.value = detail.id
    Object.assign(form, initialForm(), {
      activityId: detail.activity_id,
      name: detail.name,
      drawAt: dayjs(`${detail.draw_date} ${detail.draw_time}`),
      level: prize?.level,
      prizeId: detail.prize_id,
      winnerCount: detail.winner_count,
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
    if (disabledRoundDate(form.drawAt)) {
      message.warning('开奖日期必须在所属活动的日期范围内')
      return
    }
    submitting.value = true
    const payload = {
      activity_id: form.activityId,
      prize_id: form.prizeId,
      name: form.name,
      draw_date: form.drawAt.format('YYYY-MM-DD'),
      draw_time: form.drawAt.second(0).format('HH:mm:ss'),
      winner_count: form.winnerCount,
    }
    if (editingId.value) {
      await updateRound(editingId.value, payload)
      message.success('轮次修改成功')
    } else {
      await createRound(payload)
      message.success('轮次新增成功')
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

async function changeActivity() {
  if (form.drawAt && disabledRoundDate(form.drawAt)) form.drawAt = undefined
  form.prizeId = undefined
  form.level = undefined
  await loadPrizeOptions(selectedActivity.value?.orgId)
}

function changePrize(id) {
  const prize = prizeOptions.value.find((item) => item.id === id)
  form.level = prize?.level
}

function disabledRoundDate(current) {
  const activity = selectedActivity.value
  if (!activity || !current) return false
  return (
    current.isBefore(dayjs(activity.startDate), 'day') ||
    current.isAfter(dayjs(activity.endDate), 'day')
  )
}

async function loadPrizeOptions(orgId) {
  prizeOptions.value = []
  if (!orgId) return
  const resp = await getPrizeOptions({ org_id: orgId })
  prizeOptions.value = resp.list || []
}

defineExpose({ openCreate, openEdit })
</script>

<style scoped lang="less">
.round-modal-title {
  display: flex;
  align-items: center;
  gap: 7px;
}

.round-modal-title__notice {
  color: #999;
  font-size: 14px;
  cursor: help;
  transition: color 0.2s ease;

  &:hover {
    color: var(--color-primary);
  }
}

.round-form {
  padding-top: 8px;
}

.round-form__grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0 20px;
}

.round-form :deep(.ant-input-number),
.round-form :deep(.ant-select),
.round-form :deep(.ant-picker) {
  width: 100%;
}

.field-description {
  display: block;
  margin-top: 5px;
  color: #999;
  font-size: 11px;
  line-height: 17px;
}
</style>
