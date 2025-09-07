<template>
  <div v-if="visible" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white rounded-lg p-6 w-full max-w-md mx-4">
      <div class="flex items-center justify-between mb-6">
        <h2 class="text-lg font-medium text-gray-900">相册设置</h2>
        <button @click="handleClose" class="text-gray-400 hover:text-gray-600 transition-colors">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- 封面与标题 -->
      <div class="mb-6">
        <div class="flex items-start space-x-4">
          <div @click="handleCoverClick" class="w-24 h-24 bg-gray-100 rounded-lg border-2 border-dashed border-gray-300 flex items-center justify-center cursor-pointer hover:border-gray-400 transition-colors">
            <div v-if="!coverImage" class="text-gray-400">
              <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
              </svg>
            </div>
            <img v-else :src="coverImage" alt="封面图片" class="w-full h-full object-cover rounded-lg" />
          </div>
          <input ref="fileInput" type="file" accept="image/*" class="hidden" @change="handleFileChange" />

          <div class="flex-1">
            <label class="block text-sm font-medium text-gray-700 mb-2">标题</label>
            <input v-model="albumTitle" type="text" placeholder="请输入相册标题"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent" />
          </div>
        </div>
      </div>

      <!-- 管理员列表 -->
      <div class="mb-6">
        <label class="block text-sm font-medium text-gray-700 mb-3">管理员</label>

        <div class="flex space-x-2 mb-3">
          <input v-model="adminSearchKeyword" type="text" placeholder="搜索账号添加"
            class="flex-1 px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            @input="searchAdmins" />
          <button @click="showAddAdminDialog = true" class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors">搜索添加</button>
        </div>
        <div class="flex flex-wrap gap-2">
          <div v-if="adminUsers.length === 0" class="text-sm text-gray-400 py-2">
            暂未添加管理员
          </div>
          <div v-else v-for="admin in adminUsers" :key="admin.id" class="flex items-center space-x-2 px-3 py-1 bg-gray-100 rounded-lg">
            <div class="flex flex-col">
              <span class="text-sm font-medium text-gray-700">{{ admin.name || `用户${admin.id}` }}</span>
              <span class="text-xs text-gray-500">ID: {{ admin.id }}</span>
            </div>
            <button @click="removeAdmin(admin.id)" class="text-gray-400 hover:text-gray-600 transition-colors">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>
      </div>

      <div class="flex justify-end space-x-3">
        <button @click="handleClose" class="px-4 py-2 text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors">取消</button>
        <button @click="handleConfirm" class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors">保存</button>
      </div>
    </div>

    <!-- 添加管理员搜索对话框 -->
    <div v-if="showAddAdminDialog" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-60">
      <div class="bg-white rounded-lg p-6 w-full max-w-md mx-4">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-medium text-gray-900">搜索并添加管理员</h3>
          <button 
            @click="showAddAdminDialog = false" 
            class="text-gray-400 hover:text-gray-600 transition-colors"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>
        
        <div class="mb-4">
          <input 
            v-model="adminSearchKeyword" 
            type="text" 
            placeholder="搜索账号添加"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            @input="searchAdmins"
          />
        </div>
        
        <div class="max-h-60 overflow-y-auto">
          <div v-if="searchLoading" class="p-4 text-center text-gray-500">
            搜索中...
          </div>
          <div v-else-if="adminSearchKeyword.trim() && filteredAdmins.length === 0" class="p-4 text-center text-gray-500">
            未找到匹配的用户
          </div>
          <div v-else-if="!adminSearchKeyword.trim()" class="p-4 text-center text-gray-500">
            请输入关键词搜索用户
          </div>
          <div
            v-else
            v-for="admin in filteredAdmins"
            :key="admin.id"
            class="flex items-center justify-between p-2 hover:bg-gray-50 rounded cursor-pointer"
            @click="addAdmin(admin)"
          >
            <div class="flex flex-col">
              <span class="text-sm font-medium">{{ admin.name || `用户${admin.id}` }}</span>
              <span class="text-xs text-gray-500">ID: {{ admin.id }}</span>
            </div>
            <span class="text-xs text-gray-500">{{ admin.authority?.authorityName || '未知权限' }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, defineProps, defineEmits, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getUserList } from '@/api/user'

const props = defineProps({
  visible: { type: Boolean, default: false },
  album: { type: Object, default: () => ({}) }
})
const emit = defineEmits(['close', 'confirm'])

const albumTitle = ref('')
const coverImage = ref('')
const coverFile = ref(null)
const adminUsers = ref([])
const adminSearchKeyword = ref('')
const filteredAdmins = ref([])
const searchLoading = ref(false)
const searchTimeout = ref(null)
const showAddAdminDialog = ref(false)
const fileInput = ref(null)

watch(
  () => props.visible,
  (v) => {
    if (v && props.album) {
      console.log('相册数据:', props.album)
      albumTitle.value = props.album.name || props.album.title || ''
      coverImage.value = props.album.cover || props.album.coverImageURL || ''
      
      // 处理管理员数据
   if (Array.isArray(props.album.adminUsers)) {
        // 如果有完整的用户对象数组
        adminUsers.value = props.album.adminUsers.map(user => ({
          id: user.ID || user.id,
          uuid: user.uuid || '',
          name: user.nickName || user.username || `用户${user.ID || user.id}`,
          addedTime: new Date().toISOString()
        }))
      } else {
        adminUsers.value = []
      }
      
      coverFile.value = null
    }
  },
  { immediate: true }
)

const handleCoverClick = () => fileInput.value?.click()
const handleFileChange = (e) => {
  const f = e.target.files[0]
  if (f) {
    coverFile.value = f
    const reader = new FileReader()
    reader.onload = (ev) => (coverImage.value = ev.target.result)
    reader.readAsDataURL(f)
  }
}

// 搜索管理员
const searchAdmins = async () => {
  // 清除之前的定时器
  if (searchTimeout.value) {
    clearTimeout(searchTimeout.value)
  }
  
  // 设置防抖延迟
  searchTimeout.value = setTimeout(async () => {
    if (adminSearchKeyword.value.trim()) {
      try {
        searchLoading.value = true
        const response = await getUserList({
          page: 1,
          pageSize: 50,
          username: adminSearchKeyword.value.trim()
        })
        
        if (response.code === 0 && response.data && response.data.list) {
          // 过滤掉已经在管理员列表中的用户
          filteredAdmins.value = response.data.list.filter(user => 
            !adminUsers.value.find(admin => admin.id === user.ID)
          ).map(user => ({
            id: user.ID,
            uuid: user.uuid,
            name: user.nickName || user.username,
            authority: { authorityName: user.authority?.authorityName || '未知权限' }
          }))
        } else {
          filteredAdmins.value = []
        }
      } catch (error) {
        console.error('搜索用户失败:', error)
        ElMessage.error('搜索用户失败，请稍后重试')
        filteredAdmins.value = []
      } finally {
        searchLoading.value = false
      }
    } else {
      filteredAdmins.value = []
      searchLoading.value = false
    }
  }, 300) // 300ms 防抖延迟
}

// 添加管理员
const addAdmin = (admin) => {
  if (!adminUsers.value.find(a => a.id === admin.id)) {
    adminUsers.value.push({
      id: admin.id,
      uuid: admin.uuid,
      name: admin.name,
      addedTime: new Date().toISOString()
    })
  }
  showAddAdminDialog.value = false
  adminSearchKeyword.value = ''
  filteredAdmins.value = []
}
const removeAdmin = (id) => {
  const i = adminUsers.value.findIndex(admin => admin.id === id)
  if (i > -1) adminUsers.value.splice(i, 1)
}

const handleClose = () => emit('close')
const handleConfirm = () => {
  if (!albumTitle.value.trim()) {
    ElMessage.warning('请输入相册标题')
    return
  }
  emit('confirm', {
    id: props.album.id || props.album.ID,
    title: albumTitle.value.trim(),
    coverImageURL: coverImage.value,
    coverFile: coverFile.value,
    adminUserIDs: adminUsers.value.map(admin => admin.id)
  })
}

// 清理定时器
onUnmounted(() => {
  if (searchTimeout.value) {
    clearTimeout(searchTimeout.value)
  }
})
</script>


