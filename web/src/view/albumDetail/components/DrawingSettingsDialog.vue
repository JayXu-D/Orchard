<template>
  <el-dialog
    :model-value="visible"
    @update:model-value="handleVisibleChange"
    title="图纸设置"
    width="800px"
    :before-close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      label-width="120px"
      class="drawing-settings-form"
    >
      <!-- 序号 -->
      <el-form-item label="序号" prop="serialNumber">
        <el-input
          v-model="formData.serialNumber"
          placeholder="XX-88-88"
          maxlength="8"
          @input="formatSerialNumber"
        />
        <div class="text-xs text-gray-500 mt-1">
          格式：2个大写字母-2位数字-2位数字
        </div>
      </el-form-item>

      <!-- 图纸名称 -->
      <el-form-item label="图纸名称" prop="name">
        <el-input
          v-model="formData.name"
          placeholder="请输入图纸名称"
        />
      </el-form-item>

      <!-- 豆量 -->
      <el-form-item label="豆量" prop="beanQuantity">
        <el-input
          v-model.number="formData.beanQuantity"
          type="number"
          placeholder="请输入豆量（可选）"
        />
      </el-form-item>

      <!-- 海报图 -->
      <el-form-item label="海报图" prop="posterImage">
        <div class="poster-upload">
          <input
            ref="posterInput"
            type="file"
            accept="image/*"
            @change="handlePosterChange"
            class="hidden"
          />
          <div
            v-if="!formData.posterImage"
            @click="triggerPosterUpload"
            class="poster-upload-placeholder"
          >
            <svg class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
            <span class="text-sm text-gray-500">点击上传海报图</span>
          </div>
          <div v-else class="poster-preview">
            <img :src="formData.posterImage" alt="海报图" class="w-32 h-32 object-cover rounded" />
            <button
              @click="removePoster"
              class="absolute top-2 right-2 bg-red-500 text-white rounded-full w-6 h-6 flex items-center justify-center hover:bg-red-600"
            >
              ×
            </button>
          </div>
        </div>
      </el-form-item>

      <!-- 图纸 -->
      <el-form-item label="图纸" prop="drawings">
        <div class="drawings-upload">
          <input
            ref="drawingsInput"
            type="file"
            accept="image/*"
            multiple
            @change="handleDrawingsChange"
            class="hidden"
          />
          <div
            @click="triggerDrawingsUpload"
            class="drawings-upload-placeholder"
          >
            <svg class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
            <span class="text-sm text-gray-500">点击上传图纸（可多选）</span>
          </div>
          <div v-if="formData.drawings.length > 0" class="drawings-preview mt-4">
            <div class="grid grid-cols-4 gap-2">
              <div
                v-for="(drawing, index) in formData.drawings"
                :key="index"
                class="relative"
              >
                <img :src="drawing" alt="图纸" class="w-24 h-24 object-cover rounded" />
                <button
                  @click="removeDrawing(index)"
                  class="absolute top-1 right-1 bg-red-500 text-white rounded-full w-5 h-5 flex items-center justify-center hover:bg-red-600 text-xs"
                >
                  ×
                </button>
              </div>
            </div>
          </div>
        </div>
      </el-form-item>

      <!-- 允许下载成员 -->
      <el-form-item label="允许下载成员">
        <div class="allowed-members">
          <div class="flex items-center justify-between mb-2">
            <span class="text-sm text-gray-600">已授权成员</span>
            <button
              @click="showAddMemberDialog = true"
              class="px-3 py-1 bg-blue-500 text-white rounded text-sm hover:bg-blue-600"
            >
              添加成员
            </button>
          </div>
          <div v-if="formData.allowedMembers.length > 0" class="members-list">
            <div
              v-for="member in formData.allowedMembers"
              :key="member.id"
              class="flex items-center justify-between p-2 bg-gray-50 rounded mb-2"
            >
              <span class="text-sm">{{ member.name }} (ID: {{ member.id }})</span>
              <button
                @click="removeMember(member.id)"
                class="text-red-500 hover:text-red-700 text-sm"
              >
                移除
              </button>
            </div>
          </div>
          <div v-else class="text-sm text-gray-400">
            暂未添加下载权限成员
          </div>
        </div>
      </el-form-item>
    </el-form>

    <template #footer>
      <div class="flex justify-end space-x-3">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" @click="handleConfirm" :loading="loading">
          确认修改
        </el-button>
      </div>
    </template>

    <!-- 添加成员对话框 -->
    <el-dialog
      v-model="showAddMemberDialog"
      title="添加下载权限成员"
      width="500px"
      append-to-body
    >
      <div class="add-member-dialog">
        <el-input
          v-model="memberSearchKeyword"
          placeholder="搜索成员姓名或ID"
          @input="searchMembers"
          class="mb-4"
        />
        <div class="members-search-results">
          <div
            v-for="member in filteredMembers"
            :key="member.id"
            class="flex items-center justify-between p-2 hover:bg-gray-50 rounded cursor-pointer"
            @click="addMember(member)"
          >
            <span>{{ member.name }} (ID: {{ member.id }})</span>
            <span class="text-xs text-gray-500">{{ member.authority.authorityName }}</span>
          </div>
        </div>
      </div>
    </el-dialog>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, watch } from 'vue'
import { ElMessage } from 'element-plus'

defineOptions({
  name: 'DrawingSettingsDialog'
})

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  drawing: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['close', 'confirm'])

// 表单引用
const formRef = ref()
const posterInput = ref()
const drawingsInput = ref()

// 表单数据
const formData = reactive({
  id: '',
  serialNumber: '',
  name: '',
  beanQuantity: null,
  posterImage: '',
  drawings: [],
  allowedMembers: []
})

// 验证规则
const rules = {
  serialNumber: [
    { required: true, message: '请输入序号', trigger: 'blur' },
    { pattern: /^[A-Z]{2}-\d{2}-\d{2}$/, message: '序号格式为XX-88-88', trigger: 'blur' }
  ],
  name: [
    { required: true, message: '请输入图纸名称', trigger: 'blur' }
  ],
  posterImage: [
    { required: true, message: '请上传海报图', trigger: 'change' }
  ],
  drawings: [
    { required: true, message: '请上传至少一张图纸', trigger: 'change' },
    { type: 'array', min: 1, message: '请上传至少一张图纸', trigger: 'change' }
  ]
}

// 状态
const loading = ref(false)
const showAddMemberDialog = ref(false)
const memberSearchKeyword = ref('')
const filteredMembers = ref([])

// 格式化序号
const formatSerialNumber = (value) => {
  let formatted = value.replace(/[^A-Z0-9]/g, '').toUpperCase()
  if (formatted.length >= 2) {
    formatted = formatted.slice(0, 2) + '-' + formatted.slice(2)
  }
  if (formatted.length >= 5) {
    formatted = formatted.slice(0, 5) + '-' + formatted.slice(5)
  }
  formData.serialNumber = formatted.slice(0, 8)
}

// 触发海报图上传
const triggerPosterUpload = () => {
  posterInput.value.click()
}

// 处理海报图变化
const handlePosterChange = (event) => {
  const file = event.target.files[0]
  if (file) {
    if (file.size > 10 * 1024 * 1024) {
      ElMessage.error('海报图大小不能超过10MB')
      return
    }
    const reader = new FileReader()
    reader.onload = (e) => {
      formData.posterImage = e.target.result
    }
    reader.readAsDataURL(file)
  }
}

// 移除海报图
const removePoster = () => {
  formData.posterImage = ''
  posterInput.value.value = ''
}

// 触发图纸上传
const triggerDrawingsUpload = () => {
  drawingsInput.value.click()
}

// 处理图纸变化
const handleDrawingsChange = (event) => {
  const files = Array.from(event.target.files)
  if (files.length > 0) {
    files.forEach(file => {
      if (file.size > 20 * 1024 * 1024) {
        ElMessage.error(`图纸 ${file.name} 大小不能超过20MB`)
        return
      }
      const reader = new FileReader()
      reader.onload = (e) => {
        formData.drawings.push(e.target.result)
      }
      reader.readAsDataURL(file)
    })
  }
}

// 移除图纸
const removeDrawing = (index) => {
  formData.drawings.splice(index, 1)
}

// 搜索成员
const searchMembers = async () => {
  if (memberSearchKeyword.value.trim()) {
    // 模拟搜索结果
    filteredMembers.value = [
      { id: 1, name: '用户1', authority: { authorityName: '普通用户' } },
      { id: 2, name: '用户2', authority: { authorityName: '管理员' } }
    ]
  } else {
    filteredMembers.value = []
  }
}

// 添加成员
const addMember = (member) => {
  if (!formData.allowedMembers.find(m => m.id === member.id)) {
    formData.allowedMembers.push({
      id: member.id,
      name: member.name,
      addedTime: new Date().toISOString()
    })
  }
  showAddMemberDialog.value = false
  memberSearchKeyword.value = ''
}

// 移除成员
const removeMember = (memberId) => {
  const index = formData.allowedMembers.findIndex(m => m.id === memberId)
  if (index > -1) {
    formData.allowedMembers.splice(index, 1)
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

// 确认修改
const handleConfirm = async () => {
  try {
    await formRef.value.validate()
    
    loading.value = true
    
    const drawingData = {
      id: formData.id,
      serialNumber: formData.serialNumber,
      name: formData.name,
      beanQuantity: formData.beanQuantity,
      posterImage: formData.posterImage,
      drawings: formData.drawings,
      allowedMembers: formData.allowedMembers
    }
    
    emit('confirm', drawingData)
    
  } catch (error) {
    console.error('表单验证失败:', error)
  } finally {
    loading.value = false
  }
}

// 监听visible和drawing变化，初始化表单
watch([() => props.visible, () => props.drawing], ([visible, drawing]) => {
  if (visible && drawing) {
    // 填充表单数据
    formData.id = drawing.id || ''
    formData.serialNumber = drawing.serialNumber || ''
    formData.name = drawing.name || ''
    formData.beanQuantity = drawing.beanQuantity || null
    formData.posterImage = drawing.posterImage || ''
    formData.drawings = drawing.drawings || []
    formData.allowedMembers = drawing.allowedMembers || []
  } else if (visible) {
    // 重置表单
    formData.id = ''
    formData.serialNumber = ''
    formData.name = ''
    formData.beanQuantity = null
    formData.posterImage = ''
    formData.drawings = []
    formData.allowedMembers = []
  }
})
</script>

<style scoped>
.drawing-settings-form {
  max-height: 60vh;
  overflow-y: auto;
}

.poster-upload {
  position: relative;
}

.poster-upload-placeholder {
  width: 128px;
  height: 128px;
  border: 2px dashed #d1d5db;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: border-color 0.2s;
}

.poster-upload-placeholder:hover {
  border-color: #3b82f6;
}

.poster-preview {
  position: relative;
  display: inline-block;
}

.drawings-upload-placeholder {
  width: 200px;
  height: 100px;
  border: 2px dashed #d1d5db;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: border-color 0.2s;
}

.drawings-upload-placeholder:hover {
  border-color: #3b82f6;
}

.drawings-preview {
  margin-top: 16px;
}

.members-list {
  max-height: 200px;
  overflow-y: auto;
}

.add-member-dialog {
  max-height: 400px;
}

.members-search-results {
  max-height: 300px;
  overflow-y: auto;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
}
</style>
