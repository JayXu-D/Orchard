<template>
  <div v-if="visible" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white rounded-lg p-6 w-full max-w-md mx-4">
      <!-- 对话框标题 -->
      <div class="flex items-center justify-between mb-6">
        <h2 class="text-lg font-medium text-gray-900">相册设置</h2>
        <button 
          @click="handleClose" 
          class="text-gray-400 hover:text-gray-600 transition-colors"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
          </svg>
        </button>
      </div>

      <!-- 相册封面 -->
      <div class="mb-6">
        <div class="flex items-start space-x-4">
          <!-- 封面图片上传区域 -->
          <div 
            @click="handleCoverClick"
            class="w-24 h-24 bg-gray-100 rounded-lg border-2 border-dashed border-gray-300 flex items-center justify-center cursor-pointer hover:border-gray-400 transition-colors"
          >
            <div v-if="!coverImage" class="text-gray-400">
              <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
              </svg>
            </div>
            <img 
              v-else 
              :src="coverImage" 
              alt="封面图片" 
              class="w-full h-full object-cover rounded-lg"
            />
          </div>
          <input 
            ref="fileInput" 
            type="file" 
            accept="image/*" 
            class="hidden" 
            @change="handleFileChange"
          />

          <!-- 标题输入 -->
          <div class="flex-1">
            <label class="block text-sm font-medium text-gray-700 mb-2">标题</label>
            <input 
              v-model="albumTitle" 
              type="text" 
              placeholder="请输入相册标题"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>
        </div>
      </div>

      <!-- 管理员管理 -->
      <div class="mb-6">
        <label class="block text-sm font-medium text-gray-700 mb-3">添加管理员</label>
        
        <!-- 添加管理员输入 -->
        <div class="flex space-x-2 mb-3">
          <input 
            v-model="newAdminId" 
            type="text" 
            placeholder="输入管理员ID"
            class="flex-1 px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            @keyup.enter="addAdmin"
          />
          <button 
            @click="addAdmin"
            class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
          >
            添加
          </button>
        </div>

        <!-- 管理员标签列表 -->
        <div class="flex flex-wrap gap-2">
          <div 
            v-for="adminId in adminIds" 
            :key="adminId"
            class="flex items-center space-x-2 px-3 py-1 bg-gray-100 rounded-lg"
          >
            <span class="text-sm text-gray-700">{{ adminId }}</span>
            <button 
              @click="removeAdmin(adminId)"
              class="text-gray-400 hover:text-gray-600 transition-colors"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
              </svg>
            </button>
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="flex justify-end space-x-3">
        <button 
          @click="handleClose"
          class="px-4 py-2 text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors"
        >
          取消
        </button>
        <button 
          @click="handleConfirm"
          class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
        >
          确定
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, defineProps, defineEmits } from 'vue'
import { ElMessage } from 'element-plus'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close', 'confirm'])

// 响应式数据
const albumTitle = ref('')
const coverImage = ref('')
const coverFile = ref(null)
const adminIds = ref([])
const newAdminId = ref('')
const fileInput = ref(null)

// 处理封面点击
const handleCoverClick = () => {
  fileInput.value?.click()
}

// 处理文件选择
const handleFileChange = (event) => {
  const file = event.target.files[0]
  if (file) {
    coverFile.value = file
    const reader = new FileReader()
    reader.onload = (e) => {
      coverImage.value = e.target.result
    }
    reader.readAsDataURL(file)
  }
}

// 添加管理员
const addAdmin = () => {
  if (newAdminId.value.trim() && !adminIds.value.includes(newAdminId.value.trim())) {
    adminIds.value.push(newAdminId.value.trim())
    newAdminId.value = ''
  }
}

// 移除管理员
const removeAdmin = (adminId) => {
  const index = adminIds.value.indexOf(adminId)
  if (index > -1) {
    adminIds.value.splice(index, 1)
  }
}

// 处理关闭
const handleClose = () => {
  emit('close')
  resetForm()
}

// 处理确认
const handleConfirm = () => {
  if (!albumTitle.value.trim()) {
    ElMessage.warning('请输入相册名称')
    return
  }
  if (!coverFile.value) {
    ElMessage.warning('请上传相册封面')
    return
  }

  emit('confirm', {
    title: albumTitle.value.trim(),
    coverImageURL: coverImage.value,
    coverFile: coverFile.value,
    adminUserIDs: adminIds.value
  })
  
  resetForm()
}

// 重置表单
const resetForm = () => {
  albumTitle.value = ''
  coverImage.value = ''
  coverFile.value = null
  adminIds.value = []
  newAdminId.value = ''
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}
</script> 