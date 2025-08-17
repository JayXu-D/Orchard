<template>
  <el-dialog
    :model-value="visible"
    @update:model-value="handleVisibleChange"
    title="添加下载权限"
    width="600px"
    :before-close="handleClose"
  >
    <div class="add-permission-dialog">
      <!-- 搜索成员 -->
      <div class="mb-4">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索成员姓名或ID"
          @input="searchMembers"
          clearable
        >
          <template #prefix>
            <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </template>
        </el-input>
      </div>

      <!-- 搜索结果 -->
      <div v-if="filteredMembers.length > 0" class="search-results mb-4">
        <div class="text-sm text-gray-600 mb-2">搜索结果：</div>
        <div class="space-y-2">
          <div
            v-for="member in filteredMembers"
            :key="member.id"
            class="flex items-center justify-between p-3 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
          >
            <div class="flex items-center space-x-3">
              <div class="w-8 h-8 bg-blue-500 rounded-full flex items-center justify-center text-white text-sm font-medium">
                {{ member.name.charAt(0) }}
              </div>
              <div>
                <div class="font-medium text-gray-900">{{ member.name }}</div>
                <div class="text-sm text-gray-500">ID: {{ member.id }}</div>
                <div class="text-xs text-gray-400">{{ member.authority?.authorityName || '未知角色' }}</div>
              </div>
            </div>
            <button
              @click="addMember(member)"
              class="px-3 py-1 bg-blue-500 text-white rounded text-sm hover:bg-blue-600 transition-colors"
            >
              添加权限
            </button>
          </div>
        </div>
      </div>

      <!-- 已添加成员列表 -->
      <div v-if="selectedMembers.length > 0" class="selected-members">
        <div class="text-sm text-gray-600 mb-2">已添加成员：</div>
        <div class="space-y-2">
          <div
            v-for="member in selectedMembers"
            :key="member.id"
            class="flex items-center justify-between p-3 bg-green-50 border border-green-200 rounded-lg"
          >
            <div class="flex items-center space-x-3">
              <div class="w-8 h-8 bg-green-500 rounded-full flex items-center justify-center text-white text-sm font-medium">
                {{ member.name.charAt(0) }}
              </div>
              <div>
                <div class="font-medium text-gray-900">{{ member.name }}</div>
                <div class="text-sm text-gray-500">ID: {{ member.id }}</div>
                <div class="text-xs text-green-600">已添加</div>
              </div>
            </div>
            <button
              @click="removeMember(member.id)"
              class="px-3 py-1 bg-red-500 text-white rounded text-sm hover:bg-red-600 transition-colors"
            >
              移除
            </button>
          </div>
        </div>
      </div>

      <!-- 无搜索结果提示 -->
      <div v-else-if="searchKeyword && filteredMembers.length === 0" class="text-center py-8 text-gray-500">
        未找到匹配的成员
      </div>

      <!-- 初始提示 -->
      <div v-else class="text-center py-8 text-gray-500">
        在搜索框中输入成员姓名或ID来搜索
      </div>
    </div>

    <template #footer>
      <div class="flex justify-end space-x-3">
        <el-button @click="handleClose">取消</el-button>
        <el-button 
          type="primary" 
          @click="handleConfirm"
          :disabled="selectedMembers.length === 0"
        >
          确认添加 ({{ selectedMembers.length }})
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'
import { ElMessage } from 'element-plus'

defineOptions({
  name: 'AddPermissionDialog'
})

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  albumId: {
    type: [String, Number],
    required: true
  }
})

const emit = defineEmits(['close', 'confirm'])

// 响应式数据
const searchKeyword = ref('')
const filteredMembers = ref([])
const selectedMembers = ref([])

// 搜索成员
const searchMembers = async () => {
  if (!searchKeyword.value.trim()) {
    filteredMembers.value = []
    return
  }

  try {
    // TODO: 调用后端API搜索成员
    // 这里使用模拟数据
    const mockMembers = [
      { id: 1, name: '张三', authority: { authorityName: '普通用户' } },
      { id: 2, name: '李四', authority: { authorityName: '管理员' } },
      { id: 3, name: '王五', authority: { authorityName: '普通用户' } },
      { id: 4, name: '赵六', authority: { authorityName: '普通用户' } }
    ]

    // 过滤搜索结果
    const keyword = searchKeyword.value.toLowerCase()
    filteredMembers.value = mockMembers.filter(member => 
      member.name.toLowerCase().includes(keyword) ||
      member.id.toString().includes(keyword)
    )
  } catch (error) {
    console.error('搜索成员失败:', error)
    ElMessage.error('搜索成员失败')
  }
}

// 添加成员
const addMember = (member) => {
  if (!selectedMembers.value.find(m => m.id === member.id)) {
    selectedMembers.value.push({
      ...member,
      addedTime: new Date().toISOString()
    })
    ElMessage.success(`已添加 ${member.name} 的下载权限`)
  }
}

// 移除成员
const removeMember = (memberId) => {
  const index = selectedMembers.value.findIndex(m => m.id === memberId)
  if (index > -1) {
    const memberName = selectedMembers.value[index].name
    selectedMembers.value.splice(index, 1)
    ElMessage.info(`已移除 ${memberName} 的下载权限`)
  }
}

// 关闭对话框
const handleClose = () => {
  emit('close')
}

// 处理visible变化
const handleVisibleChange = (value) => {
  emit('close')
}

// 确认添加
const handleConfirm = () => {
  if (selectedMembers.value.length === 0) {
    ElMessage.warning('请先选择要添加权限的成员')
    return
  }

  const permissionData = {
    albumId: props.albumId,
    members: selectedMembers.value
  }

  emit('confirm', permissionData)
}

// 监听visible变化，重置数据
watch(() => props.visible, (newVal) => {
  if (newVal) {
    searchKeyword.value = ''
    filteredMembers.value = []
    selectedMembers.value = []
  }
})
</script>

<style scoped>
.add-permission-dialog {
  min-height: 400px;
}

.search-results {
  max-height: 300px;
  overflow-y: auto;
}

.selected-members {
  max-height: 300px;
  overflow-y: auto;
}
</style>
