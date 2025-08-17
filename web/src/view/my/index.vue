<template>
  <div class="flex h-screen bg-[#F6EEEE] p-[15px]">
    <!-- 侧边栏 -->
    <AppSidebar :active-menu="activeMenu" @menu-change="handleMenuChange" />
    
    <!-- 主内容区域 -->
    <div class="flex-1 flex flex-col ml-[15px] bg-white rounded-[10px]">
      <!-- 主要内容 -->
      <div class="flex-1 p-6">
        <div class="max-w-7xl mx-auto">
          <!-- 标题 -->
          <div class="mb-8">
            <h1 class="text-2xl font-semibold text-gray-900">我的</h1>
          </div>
          
          <!-- 占位内容 -->
          <div class="bg-white rounded-lg shadow p-8 text-center">
            <div class="text-gray-500 text-lg mb-4">
              <svg class="w-16 h-16 mx-auto mb-4 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
              </svg>
              <p>我的页面</p>
              <p class="text-sm text-gray-400 mt-2">这里将显示用户个人信息和相关功能</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import AppSidebar from '@/components/AppSidebar.vue'
import { ref } from 'vue'

defineOptions({
  name: 'MyPage'
})

const activeMenu = ref('my')

const handleMenuChange = (menu) => {
  activeMenu.value = menu
  // 路由跳转逻辑已经在 AppSidebar 组件中处理
}
</script> 