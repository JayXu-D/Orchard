<template>
  <div class="w-20 bg-[#CA898F] flex flex-col items-center py-6 rounded-[10px]">
    <!-- Logo -->
    <div
      class="w-12 h-12 flex items-center justify-center mb-8 cursor-pointer transition-colors hover:bg-[#B87A80] rounded-lg"
      @click="handleMenuClick('home')">
      <img src="@/assets/orchard_logo.png" alt="Orchard Logo" class="w-9 h-9" />
    </div>

    <!-- 导航菜单 -->
    <nav class="flex-1 flex flex-col items-center space-y-4">
      <div v-for="item in menuItems" :key="item.key" class="flex flex-col items-center cursor-pointer transition-colors"
        @click="handleMenuClick(item.key)">
        <div class="w-12 h-12 rounded-lg flex items-center justify-center" :class="[
          activeMenu === item.key
            ? 'bg-[#B87A80]'
            : 'bg-transparent hover:bg-[#B87A80]'
        ]">
          <img :src="item.icon" :alt="item.label" class="w-9 h-9" />
        </div>
        <span class="text-white text-xs mt-1">{{ item.label }}</span>
      </div>
    </nav>

    <!-- 底部用户图标 -->
    <div class="relative user-menu-container">
      <div
        class="w-12 h-12 bg-transparent rounded-lg flex items-center justify-center cursor-pointer transition-colors hover:bg-[#B87A80]"
        @click="toggleUserMenu">
        <CustomPic :picSrc="userStore.userInfo.headerImg" :picType="'avatar'" :size="30" />
      </div>

      <!-- 用户下拉菜单 -->
      <div v-if="showUserMenu"
        class="absolute bottom-0 left-full ml-2 w-32 bg-white rounded-lg shadow-lg border border-gray-200 py-1 z-50">
        <div v-if="userStore.userInfo.authorityId === 888" class="px-3 py-2 text-sm text-gray-700 hover:bg-gray-100 cursor-pointer"
          @click="handleUserMenuClick('member')">
          成员管理
        </div>
        <div class="px-3 py-2 text-sm text-gray-700 hover:bg-gray-100 cursor-pointer"
          @click="handleUserMenuClick('password')">
          修改密码
        </div>
        <div class="px-3 py-2 text-sm text-gray-700 hover:bg-gray-100 cursor-pointer"
          @click="handleUserMenuClick('logout')">
          退出登录
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/pinia/modules/user'
import CustomPic from '@/components/customPic/index.vue'
import { jsonInBlacklist } from '@/api/jwt'
import { useCookies } from '@vueuse/integrations/useCookies'
import { useStorage } from '@vueuse/core'

const userStore = useUserStore()
const token = useStorage('token', '')
const xToken = useCookies('x-token')

defineOptions({
  name: 'AppSidebar'
})

const props = defineProps({
  activeMenu: {
    type: String,
    default: 'home'
  }
})

const emit = defineEmits(['menu-change'])

const router = useRouter()
const activeMenu = computed(() => props.activeMenu)
const showUserMenu = ref(false)

// 菜单项配置
const menuItems = [
  {
    key: 'home',
    icon: new URL('@/assets/orchard_home.png', import.meta.url).href,
    label: '首页'
  },
  {
    key: 'my',
    icon: new URL('@/assets/orchard_my.png', import.meta.url).href,
    label: '我的'
  },
  {
    key: 'must-read',
    icon: new URL('@/assets/orchard_must_read.png', import.meta.url).href,
    label: '必读'
  }
]



const handleMenuClick = (menuKey) => {
  emit('menu-change', menuKey)

  // 根据菜单项进行路由跳转
  switch (menuKey) {
    case 'home':
      router.push('/home')
      break
    case 'my':
      router.push('/my')
      break
    case 'must-read':
      router.push('/mustRead')
      break
  }
}

// 切换用户菜单显示状态
const toggleUserMenu = () => {
  showUserMenu.value = !showUserMenu.value
}

// 处理用户菜单点击
const handleUserMenuClick = (action) => {
  showUserMenu.value = false

  switch (action) {
    case 'member':
      // 跳转到成员管理页面
      router.push('/memberManagement')
      break
    case 'password':
      // 跳转到修改密码页面
      router.push('/changePassword')
      break
    case 'logout':
      // 退出登录
      LoginOut()
      break
  }
}

/* 登出*/
const LoginOut = async () => {
  const res = await jsonInBlacklist()
  console.log(res)
  // 登出失败
  if (res.code !== 0) {
    return
  }

  await ClearStorage()

  // 把路由定向到登录页，无需等待直接reload
  router.push({ name: 'Login', replace: true })
  window.location.reload()
}

/* 清理数据 */
const ClearStorage = async () => {
  token.value = ''
  // 使用remove方法正确删除cookie
  xToken.remove()
  sessionStorage.clear()
  // 清理所有相关的localStorage项
  localStorage.removeItem('originSetting')
  localStorage.removeItem('token')
}


// 点击外部关闭下拉菜单
const handleClickOutside = (event) => {
  const userMenuElement = event.target.closest('.user-menu-container')
  if (!userMenuElement && showUserMenu.value) {
    showUserMenu.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>