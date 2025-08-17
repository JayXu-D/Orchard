<template>
  <div class="flex h-screen bg-[#F6EEEE] p-[15px]">
    <!-- 侧边栏 -->
    <AppSidebar :active-menu="activeMenu" @menu-change="handleMenuChange" />
    
    <!-- 主内容区域 -->
    <div class="flex-1 flex flex-col ml-[15px] bg-white rounded-[10px]">
    
      <!-- 主要内容 -->
      <div class="flex-1 py-[30px] px-[48px]">
        <div class="flex-1">
          <!-- 标题和按钮 -->
          <div class="flex items-center justify-between mb-8">
            <h1 class="text-2xl font-semibold text-gray-900">所有相册</h1>
            <button 
              @click="showCreateDialog = true"
              class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
            >
              新建相册
            </button>
          </div>
          
          <!-- 相册网格 -->
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-16">
            <AlbumCard 
              v-for="album in albums" 
              :key="album.id" 
              :album="album" 
              :currentUserUUID="userStore.userInfo.uuid"
              @settings="openSettings"
              @delete="confirmDelete"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- 相册创建对话框 -->
    <AlbumCreateDialog 
      :visible="showCreateDialog"
      @close="showCreateDialog = false"
      @confirm="handleCreateAlbum"
    />
    <AlbumSettingsDialog
      :visible="showSettingsDialog"
      :album="editingAlbum"
      @close="closeSettings"
      @confirm="handleUpdateAlbum"
    />
  </div>
</template>

<script setup>
import AppSidebar from '@/components/AppSidebar.vue'
import AlbumCard from './components/AlbumCard.vue'
import AlbumCreateDialog from './components/AlbumCreateDialog.vue'
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/pinia/modules/user'
import { createAlbum, getAlbumList, updateAlbum, deleteAlbumApi } from '@/api/album'
import { uploadFile } from '@/api/fileUploadAndDownload'
import { getBaseUrl } from '@/utils/format'
import AlbumSettingsDialog from './components/AlbumSettingsDialog.vue'

defineOptions({
  name: 'Home'
})

const activeMenu = ref('home')
const showCreateDialog = ref(false)
const showSettingsDialog = ref(false)
const editingAlbum = ref(null)

const handleMenuChange = (menu) => {
  activeMenu.value = menu
  // 路由跳转逻辑已经在 AppSidebar 组件中处理
}

// 相册数据
const albums = ref([])
const userStore = useUserStore()

// 处理创建相册
const handleCreateAlbum = async (albumData) => {
  try {
    // 1. 若有选择封面文件，先上传到 /fileUploadAndDownload/upload 拿到 url
    let coverUrl = albumData.coverImageURL || ''
    if (albumData.coverFile) {
      const form = new FormData()
      form.append('file', albumData.coverFile)
      const uploadRes = await uploadFile(form)
      if (uploadRes.code === 0 && uploadRes.data && uploadRes.data.file) {
        coverUrl = uploadRes.data.file.url || ''
      }
    }

    // 2. 调用创建相册接口
    const payload = {
      creatorUUID: userStore.userInfo.uuid,
      title: albumData.title,
      coverImageURL: coverUrl,
      description: '',
      adminUserIDs: (albumData.adminUserIDs || []).map(id => Number(id)).filter(Boolean)
    }
    const res = await createAlbum(payload)
    if (res.code === 0) {
      await fetchAlbums()
    }
  } finally {
    showCreateDialog.value = false
  }
}

const openSettings = (album) => {
  editingAlbum.value = album
  showSettingsDialog.value = true
}
const closeSettings = () => {
  showSettingsDialog.value = false
  editingAlbum.value = null
}

const handleUpdateAlbum = async (data) => {
  try {
    let coverUrl = data.coverImageURL || ''
    if (data.coverFile) {
      const form = new FormData()
      form.append('file', data.coverFile)
      const uploadRes = await uploadFile(form)
      if (uploadRes.code === 0 && uploadRes.data && uploadRes.data.file) {
        coverUrl = uploadRes.data.file.url || ''
      }
    }
    const payload = {
      id: data.id,
      title: data.title,
      coverImageURL: coverUrl,
      description: '',
      adminUserIDs: (data.adminUserIDs || []).map(n => Number(n)).filter(Boolean)
    }
    const res = await updateAlbum(payload)
    if (res.code === 0) {
      await fetchAlbums()
    }
  } finally {
    closeSettings()
  }
}

const confirmDelete = async (album) => {
  const res = await deleteAlbumApi({ id: album.id })
  if (res.code === 0) {
    await fetchAlbums()
  }
}

// 获取相册列表（分页，按当前用户筛选）
const fetchAlbums = async () => {
  const res = await getAlbumList({
    page: 1,
    pageSize: 100,
    creatorUUID: userStore.userInfo.uuid
  })
  if (res.code === 0 && res.data && (res.data.albums || res.data.list)) {
    // 适配 AlbumCard 所需字段
    const arr = res.data.albums || res.data.list
    albums.value = arr.map(a => ({
      id: a.id || a.ID,
      name: a.title,
      progress: 0,
      total: 0,
      cover: getBaseUrl() + a.coverImageURL,
      creatorUUID: a.creatorUUID || a.creator?.uuid
    }))

    console.log("albums", albums.value)
    console.log("userStore.userInfo.uuid", userStore.userInfo.uuid)
  }
}

onMounted(() => {
  fetchAlbums()
})
</script>
