<template>
  <div 
    class="album-card-root bg-gray-100 rounded-lg border-2 border-[#CA898F] overflow-hidden cursor-pointer hover:shadow-md transition-shadow"
    @click="handleCardClick"
  >
    <!-- 相册封面 -->
    <div class="aspect-square bg-gray-200 flex items-center justify-center">
      <div v-if="!album.cover" class="text-gray-400">
        <svg class="w-12 h-12" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
        </svg>
      </div>
      <img v-else :src="album.cover" :alt="album.name" class="w-full h-full object-cover" />
    </div>
    
    <!-- 相册信息 -->
    <div class="p-3 bg-white">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2">
          <h3 class="text-sm font-medium text-gray-900 truncate">{{ album.name }}</h3>
          <!-- 拥有者菜单按钮 -->
          <div v-if="album.creatorUUID && currentUserUUID && album.creatorUUID === currentUserUUID" class="relative">
            <button @click.stop="toggleMenu" class="p-1 rounded hover:bg-gray-100">
              <svg class="w-5 h-5 text-gray-600" viewBox="0 0 20 20" fill="currentColor">
                <path d="M6 10a2 2 0 11-4 0 2 2 0 014 0zm6 0a2 2 0 11-4 0 2 2 0 014 0zm6 0a2 2 0 11-4 0 2 2 0 014 0z" />
              </svg>
            </button>
            <!-- 弹框菜单：显示在三点右侧，底部对齐 -->
            <div v-if="showMenu" class="absolute left-full bottom-0 ml-2 w-36 bg-white rounded-[10px] border border-[#CA898F] shadow-lg z-50">
              <div @click.stop="emitSettings" class="px-4 py-2 text-[14px] text-gray-800 hover:bg-gray-100 rounded-t-[10px]">相册设置</div>
              <div @click.stop="emitDelete" class="px-4 py-2 text-[14px] text-gray-800 hover:bg-gray-100 rounded-b-[10px]">删除相册</div>
            </div>
          </div>
        </div>
        <div class="flex items-center text-red-500 text-xs">
          <svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"></path>
          </svg>
          {{ album.progress }}/{{ album.total }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'

defineOptions({
  name: 'AlbumCard'
})

const router = useRouter()

const props = defineProps({
  album: {
    type: Object,
    required: true
  },
  currentUserUUID: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['settings', 'delete'])

const showMenu = ref(false)
const toggleMenu = () => {
  showMenu.value = !showMenu.value
}
const handleClickOutside = (e) => {
  // 关闭菜单：若点击不在当前卡片内
  const card = e.target.closest('.album-card-root')
  if (!card) {
    showMenu.value = false
  }
}
onMounted(() => document.addEventListener('click', handleClickOutside))
onUnmounted(() => document.removeEventListener('click', handleClickOutside))

const emitSettings = () => {
  showMenu.value = false
  emit('settings', props.album)
}
const emitDelete = () => {
  showMenu.value = false
  emit('delete', props.album)
}

const handleCardClick = () => {
  // 跳转到相册详情页
  router.push(`/albumDetail/${props.album.id}`)
}
</script> 