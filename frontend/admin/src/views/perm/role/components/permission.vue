<template>
  <a-drawer
    v-model:open="visible"
    class="permission-drawer"
    width="860"
    placement="right"
    :closable="true"
    destroy-on-close
  >
    <template #title>
      <div class="drawer-title">
        <span>配置权限</span>
        <span v-if="role" class="role-caption">{{ role.name }} · {{ role.code }}</span>
      </div>
    </template>

    <div class="permission-toolbar">
      <span class="selected-count">已选择 {{ checkedKeys.length }} 项</span>
      <a-button type="link" :disabled="loading || readonly" @click="checkAll">全部选择</a-button>
      <a-button type="link" :disabled="loading || readonly" @click="checkedKeys = []"
        >清空</a-button
      >
    </div>

    <div class="permission-content" :class="{ 'is-loading': loading }">
      <div class="matrix-header">
        <span>目录/菜单</span>
        <span>访问权限</span>
        <span>操作权限</span>
      </div>
      <a-spin :spinning="loading">
        <a-empty v-if="!loading && !permissionGroups.length" description="暂无可配置权限" />
        <div v-for="group in permissionGroups" :key="group.key" class="permission-group">
          <div v-if="group.title" class="group-title">
            <a-checkbox
              :checked="isGroupChecked(group)"
              :indeterminate="isGroupIndeterminate(group)"
              :disabled="readonly"
              @change="toggleGroup(group, $event.target.checked)"
            />
            <span>{{ group.title }}</span>
          </div>
          <div v-for="menu in group.menus" :key="menu.id" class="permission-row">
            <div class="menu-cell">
              <span>{{ menu.name }}</span>
              <code>{{ menu.code }}</code>
            </div>
            <div class="access-cell">
              <a-checkbox
                :checked="checkedKeys.includes(menu.id)"
                :disabled="readonly"
                @change="togglePermission(menu.id, $event.target.checked)"
                >访问</a-checkbox
              >
            </div>
            <div class="action-cell">
              <a-checkbox
                v-for="action in menu.actions"
                :key="action.id"
                :checked="checkedKeys.includes(action.id)"
                :disabled="readonly"
                @change="toggleAction(menu.id, action.id, $event.target.checked)"
                >{{ action.name }}</a-checkbox
              >
              <span v-if="!menu.actions.length" class="empty-action">—</span>
            </div>
          </div>
        </div>
      </a-spin>
    </div>

    <template #footer>
      <div class="drawer-footer">
        <a-button @click="visible = false">{{ readonly ? '关闭' : '取消' }}</a-button>
        <a-button v-if="!readonly" type="primary" :loading="submitting" @click="save">
          保存权限
        </a-button>
      </div>
    </template>
  </a-drawer>
</template>

<script setup>
import { computed, ref } from 'vue'
import { message } from 'ant-design-vue'
import { getPermissionTree } from '@/api/menu'
import { updateRolePermission } from '@/api/role'

const visible = ref(false)
const role = ref()
const checkedKeys = ref([])
const loading = ref(false)
const submitting = ref(false)
const permissionGroups = ref([])
const readonly = computed(() => Number(role.value?.parent_id) === 0)

function menuCodes(menu) {
  return [menu.id, ...menu.actions.map((item) => item.id)]
}

function groupCodes(group) {
  return group.menus.flatMap(menuCodes)
}

function allCodes() {
  return permissionGroups.value.flatMap(groupCodes)
}

function checkAll() {
  checkedKeys.value = allCodes()
}

function togglePermission(code, checked) {
  const next = new Set(checkedKeys.value)
  if (checked) next.add(code)
  else next.delete(code)
  checkedKeys.value = [...next]
}

function toggleAction(menuId, actionId, checked) {
  const next = new Set(checkedKeys.value)
  if (checked) {
    next.add(menuId)
    next.add(actionId)
  } else {
    next.delete(actionId)
  }
  checkedKeys.value = [...next]
}

function isGroupChecked(group) {
  const codes = groupCodes(group)
  return codes.every((code) => checkedKeys.value.includes(code))
}

function isGroupIndeterminate(group) {
  const codes = groupCodes(group)
  const count = codes.filter((code) => checkedKeys.value.includes(code)).length
  return count > 0 && count < codes.length
}

function toggleGroup(group, checked) {
  const next = new Set(checkedKeys.value)
  groupCodes(group).forEach((code) => (checked ? next.add(code) : next.delete(code)))
  checkedKeys.value = [...next]
}

function buildPermissionGroups(tree) {
  const groups = []
  const visit = (nodes, path = []) => {
    const menus = nodes
      .filter((node) => node.code)
      .map((node) => ({
        id: node.id,
        name: node.name,
        code: node.code,
        actions: (node.children || []).filter((child) => child.code),
      }))
    if (menus.length) {
      groups.push({
        key: path.length ? `directory-${path.join('-')}` : 'root',
        title: path.join(' / '),
        menus,
      })
    }
    nodes
      .filter((node) => !node.code)
      .forEach((directory) => visit(directory.children || [], [...path, directory.name]))
  }
  visit(tree || [])
  return groups
}

async function open(record) {
  role.value = record
  checkedKeys.value = []
  permissionGroups.value = []
  visible.value = true
  loading.value = true
  try {
    const resp = await getPermissionTree(record.id)
    if (role.value?.id !== record.id) return
    permissionGroups.value = buildPermissionGroups(resp.list)
    checkedKeys.value = resp.selected_ids || []
  } catch (error) {
    message.error(error.message || '获取角色权限失败')
  } finally {
    if (role.value?.id === record.id) loading.value = false
  }
}

async function save() {
  if (!role.value || readonly.value) return
  submitting.value = true
  try {
    await updateRolePermission(role.value.id, { menu_ids: checkedKeys.value })
    message.success('角色权限配置成功')
    visible.value = false
  } catch (error) {
    message.error(error.message || '保存角色权限失败')
  } finally {
    submitting.value = false
  }
}

defineExpose({ open })
</script>

<style scoped lang="less">
.drawer-title {
  display: flex;
  align-items: baseline;
  gap: 12px;
}

.role-caption {
  color: #8c8c8c;
  font-size: 13px;
  font-weight: 400;
}

.permission-toolbar {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
}

.selected-count {
  margin-left: auto;
  color: #8c8c8c;
  font-size: 13px;
}

.permission-content {
  overflow: hidden;
  border: 1px solid #eeeeee;
  border-radius: 7px;
}

.matrix-header,
.permission-row {
  display: grid;
  grid-template-columns: minmax(220px, 1.1fr) 120px minmax(360px, 2fr);
  align-items: center;
}

.matrix-header {
  height: 46px;
  padding: 0 16px;
  color: #555;
  font-size: 13px;
  font-weight: 500;
  background: #fafafa;
  border-bottom: 1px solid #eeeeee;
}

.permission-group + .permission-group {
  border-top: 1px solid #eeeeee;
}

.group-title {
  display: flex;
  height: 40px;
  align-items: center;
  gap: 9px;
  padding: 0 16px;
  color: #3f434a;
  font-weight: 600;
  background: #fcfcfc;
  border-bottom: 1px solid #f1f1f1;
}

.permission-row {
  min-height: 48px;
  padding: 7px 16px;
  border-bottom: 1px solid #f3f3f3;
}

.permission-row:last-child {
  border-bottom: 0;
}

.menu-cell {
  display: flex;
  min-width: 0;
  align-items: center;
  gap: 10px;
  color: #333;
  font-weight: 500;
}

.menu-cell code {
  overflow: hidden;
  color: #999;
  font-size: 11px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.action-cell {
  display: flex;
  flex-wrap: wrap;
  gap: 6px 18px;
}

.empty-action {
  color: #bfbfbf;
}

.drawer-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}
</style>
