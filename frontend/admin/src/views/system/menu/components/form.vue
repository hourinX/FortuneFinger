<template>
  <a-modal
    v-model:open="visible"
    :title="modalTitle"
    width="720px"
    :confirm-loading="submitting"
    ok-text="确定"
    cancel-text="取消"
    destroy-on-close
    @ok="submit"
  >
    <a-form ref="formRef" class="menu-form" :model="form" :rules="rules" layout="vertical">
      <div class="form-grid">
        <a-form-item label="所属父级" name="pid">
          <a-tree-select
            v-model:value="form.pid"
            :tree-data="parentOptions"
            tree-default-expand-all
            popup-class-name="menu-parent-dropdown"
            placeholder="请选择父级菜单"
            @change="changeParent"
          />
        </a-form-item>

        <a-form-item label="节点类型" name="types">
          <a-radio-group v-model:value="form.types" button-style="solid" @change="changeType">
            <a-radio-button :value="1" :disabled="directoryOrMenuDisabled">目录</a-radio-button>
            <a-radio-button :value="2" :disabled="directoryOrMenuDisabled">菜单</a-radio-button>
            <a-radio-button :value="3" :disabled="permissionDisabled">权限</a-radio-button>
          </a-radio-group>
        </a-form-item>

        <a-form-item label="节点名称" name="title">
          <a-input v-model:value="form.title" :maxlength="100" placeholder="请输入节点名称" />
        </a-form-item>

        <a-form-item label="节点" name="node">
          <a-input v-model:value="form.node" :maxlength="255" placeholder="例如 system/menu" />
        </a-form-item>

        <a-form-item v-if="form.types === 2" label="路由" name="route">
          <a-input v-model:value="form.route" :maxlength="255" placeholder="例如 /system/menu" />
        </a-form-item>

        <a-form-item v-if="form.types === 3" label="权限编码" name="code">
          <a-select
            v-if="permissionBase && permissionAction !== 'custom'"
            v-model:value="permissionAction"
            :options="permissionOptions"
            placeholder="请选择权限编码"
            @change="changePermissionAction"
          >
            <template #option="option">
              {{ option.value === 'custom' ? '自定义' : option.label }}
            </template>
          </a-select>
          <div v-else-if="permissionBase" class="permission-code-custom">
            <a-select
              v-model:value="permissionAction"
              :options="permissionOptions"
              @change="changePermissionAction"
            >
              <template #option="option">
                {{ option.value === 'custom' ? '自定义' : option.label }}
              </template>
            </a-select>
            <span>:</span>
            <a-input
              ref="customPermissionInput"
              v-model:value="customPermissionCode"
              :maxlength="customPermissionMaxLength"
              allow-clear
              placeholder="自定义编码"
              @change="syncCustomPermissionCode"
            />
          </div>
          <a-input
            v-else
            v-model:value="form.code"
            :maxlength="100"
            placeholder="请输入完整权限编码"
          />
        </a-form-item>

        <a-form-item v-if="form.types !== 3" label="菜单图标" name="icon">
          <IconPicker v-model:value="form.icon" />
        </a-form-item>

        <a-form-item label="同级排序" name="sort">
          <a-input-number v-model:value="form.sort" :min="0" :max="9999" />
        </a-form-item>

        <a-form-item v-if="form.types === 2" label="页面缓存" name="keepAlive">
          <a-radio-group v-model:value="form.keepAlive">
            <a-radio :value="1">缓存</a-radio>
            <a-radio :value="2">不缓存</a-radio>
          </a-radio-group>
        </a-form-item>

        <a-form-item label="状态" name="status">
          <a-radio-group v-model:value="form.status">
            <a-radio :value="1">启用</a-radio>
            <a-radio :value="2">禁用</a-radio>
          </a-radio-group>
        </a-form-item>
      </div>

      <a-form-item label="备注" name="remark">
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
import IconPicker from '@/components/IconPicker.vue'
import { createMenu, getMenuDetail, updateMenu } from '@/api/menu'

const props = defineProps({
  menus: { type: Array, default: () => [] },
})
const emit = defineEmits(['success'])

const initialForm = () => ({
  pid: '0',
  types: 1,
  title: '',
  node: '',
  route: '',
  code: '',
  icon: 'appstore',
  sort: 0,
  keepAlive: 1,
  status: 1,
  remark: '',
})

const visible = ref(false)
const submitting = ref(false)
const formRef = ref()
const editingID = ref()
const permissionAction = ref()
const customPermissionCode = ref('')
const customPermissionInput = ref()
const form = reactive(initialForm())
const modalTitle = computed(() => (editingID.value ? '编辑菜单' : '新增菜单'))
const selectedParent = computed(() => props.menus.find((item) => item.id === form.pid))
const directoryOrMenuDisabled = computed(() => selectedParent.value?.types === 2)
const permissionDisabled = computed(() => selectedParent.value?.types !== 2)

const permissionBase = computed(() => {
  const parent = selectedParent.value
  if (!parent?.code) return ''
  return parent.code.endsWith(':view') ? parent.code.slice(0, -5) : parent.code
})
const permissionOptions = computed(() => {
  const prefix = permissionBase.value
  return [
    { label: `${prefix}:list (查询列表)`, value: 'list' },
    { label: `${prefix}:info (查询详情)`, value: 'info' },
    { label: `${prefix}:add (新增数据)`, value: 'add' },
    { label: `${prefix}:edit (修改数据)`, value: 'edit' },
    { label: `${prefix}:del (删除数据)`, value: 'del' },
    { label: prefix, value: 'custom' },
  ]
})
const customPermissionMaxLength = computed(() => Math.max(0, 99 - permissionBase.value.length))

const parentOptions = computed(() => {
  const excludedIDs = new Set()
  const collectChildren = (id) => {
    excludedIDs.add(id)
    props.menus.filter((item) => item.pid === id).forEach((item) => collectChildren(item.id))
  }
  if (editingID.value) collectChildren(editingID.value)

  const buildTree = (pid) =>
    props.menus
      .filter((item) => item.pid === pid && item.types !== 3 && !excludedIDs.has(item.id))
      .map((item) => {
        const children = buildTree(item.id)
        if (form.types === 3) {
          if (item.types === 1) {
            return children.length
              ? { title: item.title, value: `directory-${item.id}`, selectable: false, children }
              : null
          }
          return { title: item.title, value: item.id }
        }
        if (item.types !== 1) return null
        return { title: item.title, value: item.id, children }
      })
      .filter(Boolean)

  const tree = buildTree('0')
  if (form.types === 3) return tree
  return [{ title: '无（顶级菜单）', value: '0', class: 'root-option' }, ...tree]
})

const rules = {
  pid: [{ required: true, message: '请选择所属父级' }],
  types: [{ required: true, message: '请选择节点类型' }],
  title: [{ required: true, whitespace: true, message: '请输入节点名称' }],
  node: [
    { required: true, whitespace: true, message: '请输入节点' },
    {
      pattern: /^[a-z][a-z0-9-]*(\/[a-z][a-z0-9-]*)*$/,
      message: '节点只能使用小写字母、数字、短横线和斜杠',
    },
  ],
  route: [
    { required: true, whitespace: true, message: '请输入路由' },
    { pattern: /^\/[a-z0-9-/]+$/, message: '路由必须以 / 开头' },
  ],
  code: [
    { required: true, whitespace: true, message: '请输入权限编码' },
    { pattern: /^[a-z][a-z0-9-]*(?::[a-z][a-z0-9-]*)+$/, message: '请输入正确的权限编码' },
  ],
  icon: [{ required: true, message: '请选择菜单图标' }],
  sort: [{ required: true, message: '请输入同级排序' }],
  keepAlive: [{ required: true, message: '请选择页面缓存方式' }],
  status: [{ required: true, message: '请选择状态' }],
}

function changeType() {
  form.route = ''
  form.code = ''
  const parent = props.menus.find((item) => item.id === form.pid)
  if (form.types === 3) {
    if (parent?.types !== 2) form.pid = undefined
  } else if (form.pid !== '0' && parent?.types !== 1) {
    form.pid = '0'
  }
  resetPermissionCode()
  formRef.value?.clearValidate(['route', 'code', 'icon'])
}

function changeParent() {
  if (form.types === 3) resetPermissionCode()
}

function changePermissionAction(action) {
  if (action === 'custom') {
    customPermissionCode.value = ''
    form.code = ''
    nextTick(() => customPermissionInput.value?.focus())
    return
  }
  form.code = `${permissionBase.value}:${action}`
}

function syncCustomPermissionCode() {
  form.code = customPermissionCode.value
    ? `${permissionBase.value}:${customPermissionCode.value}`
    : ''
}

function resetPermissionCode() {
  permissionAction.value = undefined
  customPermissionCode.value = ''
  form.code = ''
}

function setPermissionCode(code) {
  const prefix = permissionBase.value
  const value = code || ''
  if (!prefix) {
    form.code = value
    return
  }
  const permissionPrefix = `${prefix}:`
  const action = value.startsWith(permissionPrefix) ? value.slice(permissionPrefix.length) : ''
  const preset = ['list', 'info', 'add', 'edit', 'del']
  if (preset.includes(action)) {
    permissionAction.value = action
    form.code = value
    return
  }
  permissionAction.value = 'custom'
  customPermissionCode.value = action
  form.code = value
}

function openCreate(parent) {
  editingID.value = undefined
  const childDefaults = parent
    ? { pid: parent.id, types: parent.types === 1 ? 2 : 3 }
    : initialForm()
  Object.assign(form, initialForm(), childDefaults)
  resetPermissionCode()
  visible.value = true
  nextTick(() => formRef.value?.clearValidate())
}

async function openEdit(id) {
  try {
    const record = await getMenuDetail(id)
    editingID.value = record.id
    Object.assign(form, initialForm(), {
      pid: record.pid,
      types: record.types,
      title: record.title,
      node: record.node,
      route: record.route || '',
      code: record.code || '',
      icon: record.icon || 'appstore',
      sort: record.sort,
      keepAlive: record.keep_alive,
      status: record.status,
      remark: record.remark || '',
    })
    setPermissionCode(record.code)
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
      pid: form.pid,
      types: form.types,
      title: form.title,
      node: form.node,
      route: form.types === 2 ? form.route : '',
      code: form.types === 3 ? form.code : '',
      icon: form.types === 3 ? '' : form.icon,
      sort: form.sort,
      keep_alive: form.types === 2 ? form.keepAlive : 2,
      status: form.status,
      remark: form.remark,
    }
    if (editingID.value) {
      await updateMenu(editingID.value, payload)
      message.success('菜单修改成功')
    } else {
      await createMenu(payload)
      message.success('菜单新增成功')
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
.menu-form {
  padding-top: 8px;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  column-gap: 24px;
}

.menu-form :deep(.ant-input-number),
.menu-form :deep(.picker-trigger) {
  width: 100%;
}

.menu-form :deep(.ant-radio-button-wrapper) {
  min-width: 72px;
  text-align: center;
}

.permission-code-custom {
  display: grid;
  grid-template-columns: minmax(0, 1.8fr) auto minmax(0, 1fr);
  gap: 8px;
  align-items: center;
}

.permission-code-custom > span {
  color: #8c8c8c;
}

@media (max-width: 720px) {
  .form-grid {
    grid-template-columns: 1fr;
  }
}
</style>

<style lang="less">
.menu-parent-dropdown .root-option {
  margin-bottom: 4px;
  padding-bottom: 4px;
  border-bottom: 1px solid #eeeeee;
}

.menu-parent-dropdown .root-option .ant-select-tree-node-content-wrapper {
  color: #6f7784;
  background: #f7f8fa;
  border-radius: 4px;
}

.menu-parent-dropdown .root-option .ant-select-tree-node-selected {
  color: var(--color-primary);
  background: var(--color-primary-light);
}
</style>
