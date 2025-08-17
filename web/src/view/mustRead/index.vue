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
            <h1 class="text-2xl font-semibold text-gray-900">必读</h1>
          </div>
          
          <!-- 占位内容 -->
          <div class="bg-white rounded-lg shadow p-8 text-center">
            <div class="text-gray-500 text-lg mb-4">
              <svg class="w-16 h-16 mx-auto mb-4 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
              </svg>
              <p>必读页面</p>
              <p class="text-sm text-gray-400 mt-2">这里将显示重要的阅读内容和指南</p>
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
  name: 'MustReadPage'
})

const activeMenu = ref('must-read')

const handleMenuChange = (menu) => {
  activeMenu.value = menu
  // 路由跳转逻辑已经在 AppSidebar 组件中处理
}
</script> 