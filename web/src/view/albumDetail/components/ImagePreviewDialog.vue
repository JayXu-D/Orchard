<template>
  <el-dialog
    :model-value="visible"
    @update:model-value="handleVisibleChange"
    title="海报图"
    width="800px"
    :before-close="handleClose"
    class="image-preview-dialog"
  >
    <div class="flex justify-center items-center">
      <img 
        :src="imageUrl" 
        :alt="'图片预览'"
        class="max-w-full max-h-[650px] object-contain rounded-lg shadow-lg"
        @error="handleImageError"
      />
    </div>
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
