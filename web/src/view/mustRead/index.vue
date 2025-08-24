<template>
  <div class="flex h-screen bg-[#F6EEEE] p-[15px]">
    <!-- 侧边栏 -->
    <AppSidebar :active-menu="activeMenu" @menu-change="handleMenuChange" />
    
    <!-- 主内容区域 -->
    <div class="flex-1 flex flex-col ml-[15px] bg-white rounded-[10px]">
      <!-- 主要内容 -->
      <div class="flex-1 py-[30px] px-[48px]">
        <div class="flex-1">
          <!-- 标题和编辑按钮 -->
          <div class="flex justify-between items-center mb-8">
            <h1 class="text-2xl font-semibold text-gray-900">必读</h1>
            <!-- 只有超级管理员可见编辑按钮 -->
            <el-button 
              v-if="isSuperAdmin" 
              type="primary" 
              @click="toggleEditMode"
              :icon="isEditMode ? 'Close' : 'Edit'"
            >
              {{ isEditMode ? '取消' : '编辑' }}
            </el-button>
          </div>
          
          <!-- 编辑模式 -->
          <div v-if="isEditMode && isSuperAdmin" class="bg-white rounded-lg shadow p-8">
            <!-- 标题输入 -->
            <div class="mb-6">
              <label class="block text-sm font-medium text-gray-700 mb-2">标题</label>
              <el-input 
                v-model="editData.title" 
                placeholder="请输入标题"
                size="large"
                class="w-full"
              />
            </div>
            
            <!-- 富文本编辑器 -->
            <div class="mb-6">
              <label class="block text-sm font-medium text-gray-700 mb-2">内容</label>
              <div class="border border-gray-300 rounded-md">
                <Toolbar
                  :editor="editorRef"
                  :defaultConfig="toolbarConfig"
                  mode="default"
                  style="border-bottom: 1px solid #ccc"
                />
                <Editor
                  :defaultConfig="editorConfig"
                  mode="default"
                  v-model="editData.content"
                  style="height: 400px; overflow-y: hidden;"
                  @onCreated="handleCreated"
                />
              </div>
              <!-- 自定义图片上传按钮 -->
              <div class="mt-2 flex items-center space-x-2">
                <input
                  ref="imageInput"
                  type="file"
                  accept="image/*"
                  style="display: none"
                  @change="handleImageUpload"
                />
                <el-button 
                  size="small" 
                  @click="imageInput.click()"
                  icon="Picture"
                >
                  插入图片
                </el-button>
                <span class="text-sm text-gray-500">支持 JPG、PNG、GIF 格式，最大 10MB</span>
              </div>
            </div>
            
            <!-- 操作按钮 -->
            <div class="flex justify-end space-x-4">
              <el-button @click="cancelEdit">取消</el-button>
              <el-button type="primary" @click="saveContent" :loading="saving">保存</el-button>
            </div>
          </div>
          
          <!-- 展示模式 -->
          <div v-else class="bg-white rounded-lg shadow p-8">
            <!-- 标题显示 -->
            <div class="mb-6">
              <h2 class="text-xl font-semibold text-gray-900">{{ contentData.title || '标题标题' }}</h2>
            </div>
            
            <!-- 内容显示 -->
            <div class="prose max-w-none">
              <div v-html="contentData.content || defaultContent" class="text-gray-700 leading-relaxed"></div>
            </div>
            
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import AppSidebar from '@/components/AppSidebar.vue'
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { ElMessage } from 'element-plus'
import '@wangeditor/editor/dist/css/style.css'
import { Editor, Toolbar } from '@wangeditor/editor-for-vue'
import { useUserStore } from '@/pinia/modules/user'
import { createMustRead, updateMustRead, getLatestMustRead } from '@/api/mustRead'
import { uploadFile } from '@/api/fileUploadAndDownload'
import { getBaseUrl } from '@/utils/format'

defineOptions({
  name: 'MustReadPage'
})

const activeMenu = ref('must-read')
const isEditMode = ref(false)
const saving = ref(false)
const editorRef = ref()
const imageInput = ref()
const userStore = useUserStore()
const editorConfig = ref({
  placeholder: '请输入内容...',
  MENU_CONF: {
    uploadImage: {
      server: `${getBaseUrl()}/fileUploadAndDownload/upload`,
      fieldName: 'file',
      maxFileSize: 10 * 1024 * 1024, // 10MB
      maxNumberOfFiles: 10,
      allowedFileTypes: ['image/*'],
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token') || ''}`
      },
      customInsert(res, insertFn) {
        // 自定义插入图片的逻辑
        if (res.code === 0 && res.data && res.data.file) {
          const url = res.data.file.url
          // 如果返回的是相对路径，需要拼接完整的URL
          const fullUrl = url.startsWith('http') ? url : `${getBaseUrl()}${url}`
          insertFn(fullUrl, res.data.file.name || '图片', fullUrl)
        } else {
          ElMessage.error('图片上传失败')
        }
      },
      onError(file, err, res) {
        console.error('图片上传失败:', err, res)
        ElMessage.error('图片上传失败')
      },
      onProgress(file, progress) {
        console.log('上传进度:', progress)
      },
      // 添加更多配置选项
      timeout: 30000, // 30秒超时
      onBeforeUpload(file) {
        console.log('准备上传图片:', file.name)
        return file
      },
      onSuccess(file, res) {
        console.log('图片上传成功:', res)
      }
    }
  }
})

const toolbarConfig = ref({
  excludeKeys: [
    // 排除一些不需要的工具栏按钮，但保留图片上传
    'group-video',
    'insertTable',
    'codeBlock',
    'fullScreen'
  ]
})

// 内容数据
const contentData = ref({
  id: null,
  title: '标题标题',
  content: ''
})

// 编辑时的临时数据
const editData = ref({
  title: '',
  content: ''
})

// 默认内容（参考图片中的内容）
const defaultContent = `
<div style="margin-bottom: 20px;">
  <p>必读内容</p>
  <ol style="margin-left: 20px; margin-top: 10px;">
    <li>暂无必读内容</li>
  </ol>
</div>
`

// 检查是否为超级管理员
const isSuperAdmin = ref(false)

// 检查用户权限
const checkUserPermission = () => {
  // 使用 useUserStore 获取用户信息
  try {
    const userInfo = userStore.userInfo
    // 根据用户权限判断是否为超级管理员
    // 这里可以根据实际的权限字段进行调整
    isSuperAdmin.value = userInfo.authorityId === 888
    console.log("userInfo:", userInfo)
    console.log("isSuperAdmin:", isSuperAdmin.value)
  } catch (error) {
    console.error('获取用户权限失败:', error)
    isSuperAdmin.value = false
  }
}

// 切换编辑模式
const toggleEditMode = () => {
  if (!isSuperAdmin.value) return
  
  if (!isEditMode.value) {
    // 进入编辑模式，复制当前数据到编辑区域
    editData.value = {
      title: contentData.value.title || '标题标题',
      content: contentData.value.content || defaultContent
    }
  }
  isEditMode.value = !isEditMode.value
}

// 取消编辑
const cancelEdit = () => {
  isEditMode.value = false
  editData.value = {
    title: '',
    content: ''
  }
}

// 保存内容
const saveContent = async () => {
  if (!editData.value.title.trim()) {
    ElMessage.warning('请输入标题')
    return
  }
  
  if (!editData.value.content.trim()) {
    ElMessage.warning('请输入内容')
    return
  }
  
  saving.value = true
  
  try {
    // 调用API保存数据
    const payload = {
      creatorUUID: userStore.userInfo.uuid,
      title: editData.value.title,
      content: editData.value.content
    }
    
    // 如果有ID，则更新；否则创建
    if (contentData.value.id) {
      await updateMustRead({
        id: contentData.value.id,
        title: editData.value.title,
        content: editData.value.content
      })
    } else {
      const res = await createMustRead(payload)
      if (res.code === 0 && res.data) {
        contentData.value.id = res.data.id
      }
    }
    
    // 更新显示数据
    contentData.value = {
      ...contentData.value,
      title: editData.value.title,
      content: editData.value.content
    }
    
    // 退出编辑模式
    isEditMode.value = false
    
    ElMessage.success('保存成功')
  } catch (error) {
    console.error('保存失败:', error)
    ElMessage.error('保存失败，请重试')
  } finally {
    saving.value = false
  }
}

// 编辑器创建完成
const handleCreated = (editor) => {
  editorRef.value = editor
  console.log('编辑器创建完成:', editor)
  
  // 检查编辑器配置
  console.log('编辑器配置:', editor.getConfig())
  
  // 检查是否有图片上传菜单
  const menus = editor.getMenuConfig()
  console.log('菜单配置:', menus)
}

// 处理图片上传
const handleImageUpload = async (event) => {
  const file = event.target.files[0]
  if (!file) return
  
  // 检查文件大小
  if (file.size > 10 * 1024 * 1024) {
    ElMessage.error('图片大小不能超过 10MB')
    return
  }
  
  // 检查文件类型
  if (!file.type.startsWith('image/')) {
    ElMessage.error('请选择图片文件')
    return
  }
  
  try {
    const formData = new FormData()
    formData.append('file', file)
    
    const res = await uploadFile(formData)
    
    if (res.code === 0 && res.data && res.data.file) {
      const url = res.data.file.url
      const fullUrl = url.startsWith('http') ? url : `${getBaseUrl()}${url}`
      
      // 在编辑器中插入图片
      if (editorRef.value) {
        editorRef.value.insertImage({
          url: fullUrl,
          alt: file.name,
          href: fullUrl
        })
        ElMessage.success('图片插入成功')
      }
    } else {
      ElMessage.error('图片上传失败')
    }
  } catch (error) {
    console.error('图片上传失败:', error)
    ElMessage.error('图片上传失败')
  }
  
  // 清空文件输入框
  event.target.value = ''
}

// 菜单切换处理
const handleMenuChange = (menu) => {
  activeMenu.value = menu
}

// 组件挂载时检查权限
onMounted(async () => {
  checkUserPermission()
  
  // 从API加载最新的必读内容
  try {
    const res = await getLatestMustRead()
    if (res.code === 0 && res.data) {
      contentData.value = {
        id: res.data.id,
        title: res.data.title,
        content: res.data.content
      }
    } else {
      // 如果没有数据，使用默认内容
      contentData.value.content = defaultContent
    }
  } catch (error) {
    console.error('加载必读内容失败:', error)
    // 使用默认内容
    contentData.value.content = defaultContent
  }
})

// 组件卸载时销毁编辑器
onBeforeUnmount(() => {
  if (editorRef.value) {
    editorRef.value.destroy()
  }
})
</script>

<style scoped>
.prose {
  line-height: 1.6;
}

.prose h1, .prose h2, .prose h3, .prose h4, .prose h5, .prose h6 {
  margin-top: 1.5em;
  margin-bottom: 0.5em;
  font-weight: 600;
}

.prose p {
  margin-bottom: 1em;
}

.prose ul, .prose ol {
  margin-bottom: 1em;
  padding-left: 1.5em;
}

.prose li {
  margin-bottom: 0.5em;
}

.prose strong {
  font-weight: 600;
}
</style>
