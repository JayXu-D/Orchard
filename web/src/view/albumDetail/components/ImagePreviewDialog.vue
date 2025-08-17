<template>
  <el-dialog
    :model-value="visible"
    @update:model-value="handleVisibleChange"
    title="图片预览"
    width="800px"
    :before-close="handleClose"
    class="image-preview-dialog"
  >
    <div class="flex justify-center items-center">
      <img 
        :src="imageUrl" 
        :alt="'图片预览'"
        class="max-w-full max-h-[500px] object-contain rounded-lg shadow-lg"
        @error="handleImageError"
      />
    </div>
    
    <template #footer>
      <div class="flex justify-end space-x-3">
        <el-button @click="handleClose">关闭</el-button>
        <el-button 
          type="primary" 
          @click="downloadImage"
          :loading="downloading"
        >
          下载图片
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'

defineOptions({
  name: 'ImagePreviewDialog'
})

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  imageUrl: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['close'])

const downloading = ref(false)

const handleClose = () => {
  emit('close')
}

const handleVisibleChange = (value) => {
  emit('close')
}

const handleImageError = () => {
  ElMessage.error('图片加载失败')
}

const downloadImage = async () => {
  if (!props.imageUrl) {
    ElMessage.warning('没有可下载的图片')
    return
  }

  try {
    downloading.value = true
    
    // 创建一个临时的 a 标签来下载图片
    const link = document.createElement('a')
    link.href = props.imageUrl
    link.download = `image_${Date.now()}.jpg`
    link.target = '_blank'
    
    // 添加到 DOM 中并触发点击
    document.body.appendChild(link)
    link.click()
    
    // 清理
    document.body.removeChild(link)
    
    ElMessage.success('图片下载成功')
  } catch (error) {
    console.error('下载图片失败:', error)
    ElMessage.error('图片下载失败')
  } finally {
    downloading.value = false
  }
}

// 监听 visible 变化，重置状态
watch(() => props.visible, (newVal) => {
  if (!newVal) {
    downloading.value = false
  }
})
</script>

<style scoped>
.image-preview-dialog :deep(.el-dialog__body) {
  padding: 20px;
}

.image-preview-dialog :deep(.el-dialog__footer) {
  padding: 20px;
  border-top: 1px solid #e4e7ed;
}
</style>
