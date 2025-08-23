<template>
  <div class="flex h-screen bg-[#F6EEEE] p-[15px]">
    <!-- 侧边栏 -->
    <AppSidebar :active-menu="activeMenu" @menu-change="handleMenuChange" />
    
    <!-- 主内容区域 -->
    <div class="flex-1 flex flex-col ml-[15px] bg-white rounded-[10px]">
      <!-- 主要内容 -->
      <div class="flex-1 py-[30px] px-[48px]">
        <!-- 页面标题 -->
        <div class="flex items-center justify-between mb-8">
          <h1 class="text-2xl font-semibold text-gray-900">我的图纸</h1>

          <!-- 筛选和搜索 -->
          <div class="flex items-center space-x-4">
            <!-- 筛选下拉框 -->
            <div class="relative">
              <select v-model="filterStatus"
                class="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">
                <option value="all">全部图纸</option>
                <option value="downloaded">已下载</option>
                <option value="notDownloaded">未下载</option>
              </select>
            </div>

            <!-- 搜索框 -->
            <div class="relative">
              <input v-model="searchKeyword" 
                type="text" 
                placeholder="搜索图纸名称或序号 (Ctrl+F)" 
                @keydown.ctrl.f.prevent="focusSearch"
                class="pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 w-64" />
              <svg class="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400" fill="none"
                stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>

            <!-- 批量下载按钮 -->
            <button @click="batchDownload" :disabled="selectedDrawings.length === 0"
              class="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors disabled:bg-gray-300 disabled:cursor-not-allowed">
              下载选中图纸
            </button>
            
            <!-- 刷新按钮 -->
            <button @click="refreshData" :disabled="isLoading"
              class="px-4 py-2 bg-gray-500 text-white rounded-lg hover:bg-gray-600 transition-colors disabled:bg-gray-300 disabled:cursor-not-allowed">
              <svg v-if="!isLoading" class="w-4 h-4 inline mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
              </svg>
              <svg v-else class="animate-spin w-4 h-4 inline mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
              </svg>
              刷新
            </button>
          </div>
        </div>

        <!-- 图纸表格 -->
        <div class="bg-white rounded-lg shadow overflow-hidden">
          <!-- 统计信息 -->
          <div class="p-4 bg-gray-50 border-b text-sm text-gray-600">
            <div class="flex justify-between items-center">
              <div>
                图纸总数: {{ drawings.length }}, 过滤后: {{ filteredDrawings.length }}, 已选择: {{ selectedDrawings.length }}
              </div>
              <div v-if="isLoading" class="flex items-center text-blue-600">
                <svg class="animate-spin -ml-1 mr-3 h-4 w-4 text-blue-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                加载中...
              </div>
              <div v-else-if="hasError" class="flex items-center text-red-600">
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                </svg>
                {{ errorMessage }}
              </div>
            </div>
          </div>
          
          <!-- 没有数据时的提示 -->
          <div v-if="!isLoading && filteredDrawings.length === 0" class="p-8 text-center text-gray-500">
            <div class="text-lg mb-2">
              <span v-if="hasError">加载失败</span>
              <span v-else-if="filterStatus === 'all'">暂无图纸数据</span>
              <span v-else-if="filterStatus === 'downloaded'">暂无已下载图纸</span>
              <span v-else-if="filterStatus === 'notDownloaded'">暂无未下载图纸</span>
            </div>
            <div class="text-sm mb-4">
              <span v-if="hasError">{{ errorMessage }}</span>
              <span v-else-if="filterStatus === 'all'">您还没有获得任何图纸的下载权限，请联系管理员或相册创建者</span>
              <span v-else-if="filterStatus === 'downloaded'">您还没有下载过任何图纸</span>
              <span v-else-if="filterStatus === 'notDownloaded'">您已经下载了所有可下载的图纸</span>
            </div>
            <div v-if="hasError" class="flex justify-center space-x-3">
              <button @click="fetchMyDrawings" class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors">
                重试
              </button>
              <button @click="refreshData" class="px-4 py-2 bg-gray-500 text-white rounded-lg hover:bg-gray-600 transition-colors">
                刷新
              </button>
            </div>
          </div>
          
          <!-- 有数据时显示表格 -->
          <table v-else-if="!isLoading" class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  <div class="flex items-center space-x-2">
                    <input type="checkbox" 
                      :checked="isAllSelected" 
                      @change="toggleSelectAll"
                      class="w-4 h-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500 focus:ring-2 cursor-pointer appearance-auto" 
                      style="-webkit-appearance: auto; -moz-appearance: auto; appearance: auto;" />
                    <span class="text-xs text-gray-500">全选</span>
                  </div>
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">序号</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">图纸名称</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">豆量</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">海报图</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">所属相册</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">下载</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="drawing in filteredDrawings" :key="drawing.id" class="hover:bg-gray-50">
                <td class="px-6 py-4 whitespace-nowrap">
                  <input type="checkbox" 
                    :checked="selectedDrawings.includes(drawing.id)"
                    @change="toggleDrawingSelection(drawing.id)"
                    class="w-4 h-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500 focus:ring-2 cursor-pointer appearance-auto" 
                    style="-webkit-appearance: auto; -moz-appearance: auto; appearance: auto;" />
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ drawing.serialNumber }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ drawing.name }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ drawing.beanQuantity || '-' }}</td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <img :src="drawing.posterImage" :alt="drawing.name"
                    class="w-16 h-16 object-cover rounded cursor-pointer hover:opacity-80"
                    @click="previewImage(drawing.posterImage)" />
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ drawing.albumName || '-' }}</td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <button @click="downloadDrawing(drawing)" :disabled="!drawing.canDownload"
                    class="px-3 py-1 text-sm rounded transition-colors" :class="drawing.downloaded
                      ? 'bg-gray-300 text-gray-600 cursor-not-allowed'
                      : 'bg-red-500 text-white hover:bg-red-600'">
                    {{ drawing.downloaded ? '重新下载' : '下载图纸' }}
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
          
          <!-- 加载骨架屏 -->
          <div v-else-if="isLoading" class="p-8">
            <div class="animate-pulse">
              <div class="h-4 bg-gray-200 rounded w-1/4 mb-4"></div>
              <div class="space-y-3">
                <div v-for="i in 5" :key="i" class="flex space-x-4">
                  <div class="h-4 bg-gray-200 rounded w-4"></div>
                  <div class="h-4 bg-gray-200 rounded w-24"></div>
                  <div class="h-4 bg-gray-200 rounded w-32"></div>
                  <div class="h-4 bg-gray-200 rounded w-16"></div>
                  <div class="h-16 bg-gray-200 rounded w-16"></div>
                  <div class="h-4 bg-gray-200 rounded w-20"></div>
                  <div class="h-4 bg-gray-200 rounded w-20"></div>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 分页组件 -->
        <div v-if="!isLoading && filteredDrawings.length > 0" class="px-6 py-4 bg-white border-t border-gray-200">
          <div class="flex items-center justify-between">
            <div class="text-sm text-gray-700">
              显示第 {{ (currentPage - 1) * pageSize + 1 }} 到 {{ Math.min(currentPage * pageSize, filteredDrawings.length) }} 条，
              共 {{ filteredDrawings.length }} 条记录
            </div>
            <div class="flex items-center space-x-2">
              <button 
                @click="changePage(currentPage - 1)" 
                :disabled="currentPage === 1"
                class="px-3 py-1 text-sm border border-gray-300 rounded-md hover:bg-gray-50 disabled:bg-gray-100 disabled:text-gray-400 disabled:cursor-not-allowed">
                上一页
              </button>
              <span class="px-3 py-1 text-sm text-gray-700">
                第 {{ currentPage }} 页，共 {{ totalPages }} 页
              </span>
              <button 
                @click="changePage(currentPage + 1)" 
                :disabled="currentPage === totalPages"
                class="px-3 py-1 text-sm border border-gray-300 rounded-md hover:bg-gray-50 disabled:bg-gray-100 disabled:text-gray-400 disabled:cursor-not-allowed">
                下一页
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 图片预览对话框 -->
    <!-- <ImagePreviewDialog :visible="showImagePreview" :image-url="previewImageUrl" @close="showImagePreview = false" /> -->
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useUserStore } from '@/pinia/modules/user'
import { ElMessage } from 'element-plus'
import AppSidebar from '@/components/AppSidebar.vue'
// import ImagePreviewDialog from '@/components/ImagePreviewDialog.vue'
import { getMyDrawings, downloadDrawing as downloadDrawingApi, batchDownloadDrawings } from '@/api/album'
import { getBaseUrl } from '@/utils/format'

// 防抖函数
const debounce = (func, wait) => {
  let timeout
  return function executedFunction(...args) {
    const later = () => {
      clearTimeout(timeout)
      func(...args)
    }
    clearTimeout(timeout)
    timeout = setTimeout(later, wait)
  }
}

defineOptions({
  name: 'MyDrawings'
})

const userStore = useUserStore()

// 响应式数据
const drawings = ref([])
const selectedDrawings = ref([])
const filterStatus = ref('all')
const searchKeyword = ref('')
const showImagePreview = ref(false)
const previewImageUrl = ref('')
const isLoading = ref(false) // 新增加载状态
const hasError = ref(false) // 新增错误状态
const errorMessage = ref('') // 新增错误信息

// 分页相关
const pageSize = ref(10) // 每页显示数量
const currentPage = ref(1) // 当前页码
const totalPages = computed(() => {
  return Math.ceil(filteredDrawings.value.length / pageSize.value)
})

// 计算属性
const activeMenu = ref('my')

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
  }

  // 关键词搜索
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    filtered = filtered.filter(d =>
      d.serialNumber.toLowerCase().includes(keyword) ||
      d.name.toLowerCase().includes(keyword) ||
      (d.albumName && d.albumName.toLowerCase().includes(keyword))
    )
  }

  // 分页
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filtered.slice(start, end)
})

const isAllSelected = computed(() => {
  return filteredDrawings.value.length > 0 &&
    filteredDrawings.value.every(d => selectedDrawings.value.includes(d.id))
})

// 方法
const handleMenuChange = (menu) => {
  activeMenu.value = menu
}

const toggleSelectAll = () => {
  if (isAllSelected.value) {
    selectedDrawings.value = []
  } else {
    selectedDrawings.value = filteredDrawings.value.map(d => d.id)
  }
}

const toggleDrawingSelection = (drawingId) => {
  const index = selectedDrawings.value.indexOf(drawingId)
  if (index > -1) {
    selectedDrawings.value.splice(index, 1)
  } else {
    selectedDrawings.value.push(drawingId)
  }
}

const batchDownload = async () => {
  if (selectedDrawings.value.length === 0) return

  try {
    // 调用批量下载接口
    const result = await batchDownloadDrawings({
      drawingIds: selectedDrawings.value,
      addWatermark: true,
      watermarkText: '批量下载图纸'
    })
    
    if (result.code === 0) {
      ElMessage.success(`成功下载 ${selectedDrawings.value.length} 个图纸`)
      
      // 标记所有选中的图纸为已下载
      selectedDrawings.value.forEach(drawingId => {
        const drawing = drawings.value.find(d => d.id === drawingId)
        if (drawing) {
          drawing.downloaded = true
        }
      })
      
      // 清空选择
      selectedDrawings.value = []
      
      // 如果返回了文件路径列表，触发浏览器下载
      if (result.data && result.data.filePaths && result.data.filePaths.length > 0) {
        result.data.filePaths.forEach((filePath, index) => {
          const downloadUrl = filePath.startsWith('/') ? filePath : `/api/v1/drawing/${filePath}`
          
          const link = document.createElement('a')
          link.href = downloadUrl
          link.download = filePath.split('/').pop()
          link.style.display = 'none'
          document.body.appendChild(link)
          
          setTimeout(() => {
            link.click()
            document.body.removeChild(link)
          }, index * 100)
        })
        
        ElMessage.success(`批量下载成功，共 ${result.data.filePaths.length} 个文件`)
      }
    } else {
      ElMessage.error('批量下载失败: ' + (result.msg || '未知错误'))
    }
  } catch (error) {
    console.error('批量下载失败:', error)
    ElMessage.error('批量下载失败: ' + error.message)
  }
}

const downloadDrawing = async (drawing) => {
  if (!drawing.canDownload) return
  
  try {
    // 调用下载接口
    const result = await downloadDrawingApi({
      drawingId: drawing.id,
      albumId: drawing.albumId,
      addWatermark: true,
      watermarkText: `创建者: ${drawing.creator?.username || '未知'}`
    })
    
    if (result.code === 0) {
      // 标记为已下载
      drawing.downloaded = true
      ElMessage.success('图纸下载成功')
      
      // 如果返回了文件路径列表，触发浏览器下载
      if (result.data && result.data.filePaths && result.data.filePaths.length > 0) {
        result.data.filePaths.forEach((filePath, index) => {
          const downloadUrl = filePath.startsWith('/') ? filePath : `/api/v1/drawing/${filePath}`
          
          const link = document.createElement('a')
          link.href = downloadUrl
          link.download = filePath.split('/').pop()
          link.style.display = 'none'
          document.body.appendChild(link)
          
          setTimeout(() => {
            link.click()
            document.body.removeChild(link)
          }, index * 100)
        })
        
        ElMessage.success(`图纸下载成功，共 ${result.data.filePaths.length} 个文件`)
      }
    } else {
      ElMessage.error('图纸下载失败: ' + (result.msg || '未知错误'))
    }
  } catch (error) {
    console.error('下载失败:', error)
    ElMessage.error('下载失败: ' + error.message)
  }
}

const previewImage = (imageUrl) => {
  previewImageUrl.value = imageUrl
  showImagePreview.value = true
}

const refreshData = async () => {
  isLoading.value = true
  try {
    const res = await getMyDrawings({
      page: 1,
      pageSize: 100
    })
    
    if (res.code === 0 && res.data) {
      drawings.value = res.data.drawings.map(d => ({
        ...d,
        id: d.id || d.ID,
        serialNumber: d.serialNumber || d.serial_number || '',
        name: d.name || d.title || '',
        beanQuantity: d.beanQuantity || d.bean_quantity || null,
        albumId: d.albumId || d.album_id || d.AlbumID,
        albumName: d.albumName || d.album_name || d.AlbumName || '',
        posterImage: getBaseUrl() + (d.posterImageURL || ''),
        canDownload: true,
        downloaded: d.downloaded || false,
        creator: d.creator || {}
      }))
    } else {
      drawings.value = []
    }
  } catch (error) {
    console.error('刷新图纸列表失败:', error)
    drawings.value = []
  } finally {
    isLoading.value = false
  }
}

const fetchMyDrawings = async () => {
  isLoading.value = true
  hasError.value = false
  errorMessage.value = ''
  
  try {
    const res = await getMyDrawings({
      page: 1,
      pageSize: 100
    })
    
    if (res.code === 0 && res.data) {
      drawings.value = res.data.drawings.map(d => ({
        ...d,
        id: d.id || d.ID,
        serialNumber: d.serialNumber || d.serial_number || '',
        name: d.name || d.title || '',
        beanQuantity: d.beanQuantity || d.bean_quantity || null,
        albumId: d.albumId || d.album_id || d.AlbumID,
        albumName: d.albumName || d.album_name || d.AlbumName || '',
        posterImage: getBaseUrl() + (d.posterImageURL || ''),
        canDownload: true,
        downloaded: d.downloaded || false,
        creator: d.creator || {}
      }))
      console.log('成功获取我的图纸列表:', drawings.value.length)
    } else {
      drawings.value = []
      hasError.value = true
      errorMessage.value = res.msg || '获取图纸列表失败'
      console.warn('API调用失败:', res.msg)
    }
  } catch (error) {
    console.error('获取我的图纸列表失败:', error)
    drawings.value = []
    hasError.value = true
    errorMessage.value = error.message || '网络错误，请稍后重试'
  } finally {
    isLoading.value = false
  }
}

const changePage = (newPage) => {
  if (newPage >= 1 && newPage <= totalPages.value) {
    currentPage.value = newPage
  }
}

const focusSearch = () => {
  // 聚焦到搜索框
  const searchInput = document.querySelector('input[placeholder*="搜索"]')
  if (searchInput) {
    searchInput.focus()
    searchInput.select()
  }
}

// 监听筛选和搜索变化
watch([filterStatus, searchKeyword], () => {
  currentPage.value = 1 // 重置页码
  selectedDrawings.value = [] // 重置选择
})

onMounted(() => {
  fetchMyDrawings()
})
</script>

<style scoped>
/* 确保checkbox在所有浏览器中可见 */
input[type="checkbox"] {
  -webkit-appearance: auto !important;
  -moz-appearance: auto !important;
  appearance: auto !important;
  display: inline-block !important;
  visibility: visible !important;
  opacity: 1 !important;
}

/* 自定义checkbox样式 */
input[type="checkbox"] {
  width: 16px !important;
  height: 16px !important;
  border: 2px solid #d1d5db !important;
  border-radius: 4px !important;
  background-color: white !important;
  cursor: pointer !important;
  position: relative !important;
}

input[type="checkbox"]:checked {
  background-color: #3b82f6 !important;
  border-color: #3b82f6 !important;
}

input[type="checkbox"]:checked::after {
  content: '✓' !important;
  position: absolute !important;
  top: 50% !important;
  left: 50% !important;
  transform: translate(-50%, -50%) !important;
  color: white !important;
  font-size: 12px !important;
  font-weight: bold !important;
}

input[type="checkbox"]:focus {
  outline: 2px solid #3b82f6 !important;
  outline-offset: 2px !important;
}
</style> 