<template>
  <div class="flex h-screen bg-[#F6EEEE] p-[15px]">
    <!-- 侧边栏 -->
    <AppSidebar :active-menu="activeMenu" @menu-change="handleMenuChange" />
    
    <!-- 主内容区域 -->
    <div class="flex-1 flex flex-col ml-[15px] bg-white rounded-[10px]">
      <!-- 主要内容 -->
      <div class="flex-1 py-[30px] px-[48px]">
        <!-- 相册标题 -->
        <div class="flex items-center justify-between mb-8">
          <h1 class="text-2xl font-semibold text-gray-900">{{ albumTitle }}</h1>
          
          <!-- 筛选和搜索 -->
          <div class="flex items-center space-x-4">
            <!-- 筛选下拉框 -->
            <div class="relative">
              <select v-model="filterStatus" class="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">
                <option value="all">全部图纸</option>
                <option value="downloaded">已下载</option>
                <option value="notDownloaded">未下载</option>
                <option value="hasPermission">已获得</option>
                <option value="noPermission">未获得</option>
              </select>
            </div>
            
            <!-- 搜索框 -->
            <div class="relative">
              <input 
                v-model="searchKeyword" 
                type="text" 
                placeholder="搜索图纸名称或序号"
                class="pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 w-64"
              />
              <svg class="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
            
            <!-- 批量下载按钮 -->
            <button 
              @click="batchDownload"
              :disabled="selectedDrawings.length === 0"
              class="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors disabled:bg-gray-300 disabled:cursor-not-allowed"
            >
              下载选中图纸
            </button>
          </div>
        </div>
        
        <!-- 操作按钮 -->
        <div v-if="canManage" class="flex justify-end space-x-3 mb-6">
          <button 
            @click="showAddPermissionDialog = true"
            class="px-4 py-2 bg-gray-300 text-gray-700 rounded-lg hover:bg-gray-400 transition-colors"
          >
            添加权限
          </button>
          <button 
            @click="showUploadDialog = true"
            class="px-4 py-2 bg-gray-300 text-gray-700 rounded-lg hover:bg-gray-400 transition-colors"
          >
            上传新图纸
          </button>
        </div>
        
        <!-- 图纸表格 -->
        <div class="bg-white rounded-lg shadow overflow-hidden">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  <input 
                    type="checkbox" 
                    :checked="isAllSelected"
                    @change="toggleSelectAll"
                    class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                  />
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">序号</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">图纸名称</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">豆量</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">海报图</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">下载</th>
                <th v-if="canManage" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">操作</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="drawing in filteredDrawings" :key="drawing.id" class="hover:bg-gray-50">
                <td class="px-6 py-4 whitespace-nowrap">
                  <input 
                    type="checkbox" 
                    :checked="selectedDrawings.includes(drawing.id)"
                    @change="toggleDrawingSelection(drawing.id)"
                    class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                  />
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ drawing.serialNumber }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ drawing.name }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ drawing.beanQuantity || '-' }}</td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <img 
                    :src="drawing.posterImage" 
                    :alt="drawing.name"
                    class="w-16 h-16 object-cover rounded cursor-pointer hover:opacity-80"
                    @click="previewImage(drawing.posterImage)"
                  />
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <button 
                    @click="downloadDrawing(drawing)"
                    :disabled="!drawing.canDownload"
                    class="px-3 py-1 text-sm rounded transition-colors"
                    :class="drawing.downloaded 
                      ? 'bg-gray-300 text-gray-600 cursor-not-allowed' 
                      : 'bg-red-500 text-white hover:bg-red-600'"
                  >
                    {{ drawing.downloaded ? '重新下载' : '下载图纸' }}
                  </button>
                </td>
                <td v-if="canManage" class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  <button 
                    v-if="drawing.canEdit"
                    @click="editDrawing(drawing)"
                    class="text-blue-600 hover:text-blue-800 cursor-pointer"
                  >
                    编辑
                  </button>
                  <span v-else class="text-gray-400 cursor-not-allowed">编辑</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- 上传新图纸对话框 -->
    <UploadDrawingDialog
      :visible="showUploadDialog"
      :album-id="albumId"
      @close="showUploadDialog = false"
      @confirm="handleUploadDrawing"
    />

    <!-- 添加权限对话框 -->
    <AddPermissionDialog
      :visible="showAddPermissionDialog"
      :album-id="albumId"
      @close="showAddPermissionDialog = false"
      @confirm="handleAddPermission"
    />

    <!-- 图纸设置对话框 -->
    <DrawingSettingsDialog
      :visible="showSettingsDialog"
      :drawing="editingDrawing"
      @close="showSettingsDialog = false"
      @confirm="handleUpdateDrawing"
    />

    <!-- 图片预览对话框 -->
    <ImagePreviewDialog
      :visible="showImagePreview"
      :image-url="previewImageUrl"
      @close="showImagePreview = false"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/pinia/modules/user'
import AppSidebar from '@/components/AppSidebar.vue'
import UploadDrawingDialog from './components/UploadDrawingDialog.vue'
import AddPermissionDialog from './components/AddPermissionDialog.vue'
import DrawingSettingsDialog from './components/DrawingSettingsDialog.vue'
import ImagePreviewDialog from './components/ImagePreviewDialog.vue'
import { getAlbumDetail, getDrawingList, downloadDrawing as downloadDrawingApi } from '@/api/album'
import { getBaseUrl } from '@/utils/format'

defineOptions({
  name: 'AlbumDetail'
})

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

// 响应式数据
const albumId = ref('')
const albumTitle = ref('')
const drawings = ref([])
const selectedDrawings = ref([])
const filterStatus = ref('all')
const searchKeyword = ref('')
const showUploadDialog = ref(false)
const showAddPermissionDialog = ref(false)
const showSettingsDialog = ref(false)
const showImagePreview = ref(false)
const previewImageUrl = ref('')
const editingDrawing = ref(null)

// 计算属性
const activeMenu = ref('home')

// 添加调试信息
console.log('相册详情页面初始化，相册ID:', albumId.value, '路由参数:', route.params, '完整路由:', route, '当前路径:', window.location.href, 'hash:', window.location.hash, 'search:', window.location.search, 'pathname:', window.location.pathname, 'origin:', window.location.origin, 'protocol:', window.location.protocol, 'host:', window.location.host, 'port:', window.location.port, 'href:', window.location.href, 'assign:', window.location.assign, 'reload:', window.location.reload)

const canManage = computed(() => {
  // 相册创建者或管理员可以管理
  console.log('计算管理权限，相册ID:', albumId.value)
  return true // TODO: 根据实际权限判断
})

const filteredDrawings = computed(() => {
  let filtered = drawings.value

  // 状态筛选
  switch (filterStatus.value) {
    case 'downloaded':
      filtered = filtered.filter(d => d.downloaded)
      break
    case 'notDownloaded':
      filtered = filtered.filter(d => !d.downloaded)
      break
    case 'hasPermission':
      filtered = filtered.filter(d => d.canDownload)
      break
    case 'noPermission':
      filtered = filtered.filter(d => !d.canDownload)
      break
  }

  // 关键词搜索
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    filtered = filtered.filter(d => 
      d.serialNumber.toLowerCase().includes(keyword) ||
      d.name.toLowerCase().includes(keyword)
    )
  }

  console.log('过滤图纸，相册ID:', albumId.value, '结果数量:', filtered.length)
  return filtered
})

const isAllSelected = computed(() => {
  const result = filteredDrawings.value.length > 0 && 
         filteredDrawings.value.every(d => selectedDrawings.value.includes(d.id))
  console.log('计算全选状态，相册ID:', albumId.value, '结果:', result)
  return result
})

// 方法
const handleMenuChange = (menu) => {
  activeMenu.value = menu
  console.log('菜单切换:', menu)
  console.log('相册ID:', albumId.value)
}

const toggleSelectAll = () => {
  if (isAllSelected.value) {
    selectedDrawings.value = []
  } else {
    selectedDrawings.value = filteredDrawings.value.map(d => d.id)
  }
  console.log('全选/取消全选:', selectedDrawings.value.length)
  console.log('相册ID:', albumId.value)
}

const toggleDrawingSelection = (drawingId) => {
  const index = selectedDrawings.value.indexOf(drawingId)
  if (index > -1) {
    selectedDrawings.value.splice(index, 1)
  } else {
    selectedDrawings.value.push(drawingId)
  }
  console.log('选择图纸:', drawingId)
  console.log('相册ID:', albumId.value)
}

const batchDownload = async () => {
  if (selectedDrawings.value.length === 0) return
  
  try {
    // TODO: 实现批量下载逻辑
    console.log('批量下载:', selectedDrawings.value)
    console.log('相册ID:', albumId.value)
  } catch (error) {
    console.error('批量下载失败:', error)
  }
}

const downloadDrawing = async (drawing) => {
  if (!drawing.canDownload) return
  
  try {
    // TODO: 实现单个图纸下载逻辑
    console.log('下载图纸:', drawing.id)
    console.log('相册ID:', albumId.value)
    const result = await downloadDrawingApi(drawing.id)
    if (result.code === 0) {
      // 标记为已下载
      drawing.downloaded = true
    }
  } catch (error) {
    console.error('下载失败:', error)
  }
}

const editDrawing = (drawing) => {
  if (!drawing.canEdit) return
  
  editingDrawing.value = drawing
  showSettingsDialog.value = true
  console.log('编辑图纸:', drawing.id)
  console.log('相册ID:', albumId.value)
}

const previewImage = (imageUrl) => {
  previewImageUrl.value = imageUrl
  showImagePreview.value = true
  console.log('预览图片:', imageUrl)
  console.log('相册ID:', albumId.value)
}

const handleUploadDrawing = async (drawingData) => {
  try {
    // TODO: 实现上传图纸逻辑
    console.log('上传图纸:', drawingData)
    console.log('相册ID:', albumId.value)
    await fetchDrawings()
  } catch (error) {
    console.error('上传失败:', error)
  }
}

const handleAddPermission = async (permissionData) => {
  try {
    // TODO: 实现添加权限逻辑
    console.log('添加权限:', permissionData)
    console.log('相册ID:', albumId.value)
  } catch (error) {
    console.error('添加权限失败:', error)
  }
}

const handleUpdateDrawing = async (drawingData) => {
  try {
    // TODO: 实现更新图纸逻辑
    console.log('更新图纸:', drawingData)
    console.log('相册ID:', albumId.value)
    await fetchDrawings()
  } catch (error) {
    console.error('更新失败:', error)
  }
}

const fetchAlbumDetail = async () => {
  if (!albumId.value) {
    console.warn('相册ID未设置')
    return
  }
  
  try {
    const res = await getAlbumDetail(albumId.value)
    if (res.code === 0 && res.data) {
      albumTitle.value = res.data.title || res.data.name || '相册详情'
    } else {
      // 如果API调用失败，设置默认标题
      albumTitle.value = '相册详情'
    }
  } catch (error) {
    console.error('获取相册详情失败:', error)
    albumTitle.value = '相册详情'
  }
}

const fetchDrawings = async () => {
  if (!albumId.value) {
    console.warn('相册ID未设置')
    return
  }
  
  try {
    const res = await getDrawingList({ albumId: albumId.value })
    if (res.code === 0 && res.data) {
      drawings.value = res.data.map(d => ({
        ...d,
        id: d.id || d.ID,
        serialNumber: d.serialNumber || d.serial_number || '',
        name: d.name || d.title || '',
        beanQuantity: d.beanQuantity || d.bean_quantity || d.beanQuantity || null,
        posterImage: d.posterImage || d.poster_image || d.coverImageURL || '',
        canDownload: true, // TODO: 根据实际权限判断
        canEdit: d.creatorId === userStore.userInfo.ID || d.creatorUUID === userStore.userInfo.uuid, // 只有上传者可以编辑
        downloaded: false // TODO: 根据实际下载状态判断
      }))
    } else {
      // 如果API调用失败，设置空数组
      drawings.value = []
    }
  } catch (error) {
    console.error('获取图纸列表失败:', error)
    drawings.value = []
  }
}

// 监听筛选和搜索变化
watch([filterStatus, searchKeyword], () => {
  selectedDrawings.value = [] // 重置选择
  console.log('筛选或搜索变化')
  console.log('相册ID:', albumId.value)
})

// 监听路由参数变化
watch(() => route.params.id, (newId) => {
  if (newId) {
    albumId.value = newId
    console.log('路由参数变化，相册ID:', newId)
    fetchAlbumDetail()
    fetchDrawings()
  }
}, { immediate: true })

onMounted(() => {
  // 初始化相册ID
  if (route.params.id) {
    albumId.value = route.params.id
    console.log('组件挂载，相册ID:', route.params.id)
    fetchAlbumDetail()
    fetchDrawings()
  } else {
    console.warn('组件挂载时没有相册ID')
  }
})
</script>
