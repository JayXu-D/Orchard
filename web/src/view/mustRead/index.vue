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

defineOptions({
  name: 'MustReadPage'
})

const activeMenu = ref('must-read')
const isEditMode = ref(false)
const saving = ref(false)
const editorRef = ref()
const userStore = useUserStore()
const editorConfig = ref({
  placeholder: '请输入内容...',
  MENU_CONF: {}
})

const toolbarConfig = ref({
  excludeKeys: []
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
