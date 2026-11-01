<template>
  <a-modal
    v-model:open="visible"
    :title="`${activityName} · 签到模板`"
    width="1040px"
    ok-text="确定"
    cancel-text="取消"
    destroy-on-close
    @ok="confirm"
  >
    <TemplateSelector v-model="selectedTemplate" />
  </a-modal>
</template>

<script setup>
import { ref } from 'vue'
import TemplateSelector from './template-selector.vue'

const emit = defineEmits(['confirm'])
const visible = ref(false)
const activityId = ref()
const activityName = ref('')
const selectedTemplate = ref('invitation')

function open(activity, template) {
  activityId.value = activity.id
  activityName.value = activity.name
  selectedTemplate.value = template || 'invitation'
  visible.value = true
}

function confirm() {
  emit('confirm', { activityId: activityId.value, template: selectedTemplate.value })
  visible.value = false
}

defineExpose({ open })
</script>
