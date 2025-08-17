<template>
  <div class="flex h-screen bg-[#F6EEEE] p-[15px]">
    <!-- 侧边栏 -->
    <AppSidebar :active-menu="activeMenu" @menu-change="handleMenuChange" />

    <!-- 主内容区域 -->
    <div class="flex-1 flex flex-col ml-[15px] bg-white rounded-[10px]">
      <!-- 主要内容 -->
      <div class="flex-1 py-[30px] px-[48px]">
        <div class="flex-1">
          <!-- 标题 -->
          <div class="mb-8">
            <h1 class="text-[24px] font-semibold text-[#CA898F]">修改密码</h1>
          </div>

          <!-- 修改密码表单 -->
          <div>
            <form @submit.prevent="handleSubmit">
              <!-- 原密码字段 -->
              <div class="mb-[41px]">
                <label for="oldPassword" class="block text-sm font-medium text-gray-700 mb-2">
                  原密码
                </label>
                <input id="oldPassword" v-model="formData.oldPassword" type="password" placeholder="请输入原密码"
                  class="w-[394px] h-[46px] border border-gray-300 rounded-[5px] focus:outline-none focus:ring-2 focus:ring-[#CA898F] focus:border-transparent text-xs"
                  style="padding: 12px; font-size: 16px;" :class="{ 'border-red-500': errors.oldPassword }" />
                <p v-if="errors.oldPassword" class="mt-1 text-sm text-red-600">
                  {{ errors.oldPassword }}
                </p>
              </div>

              <!-- 新密码字段 -->
              <div class="mb-[41px]">
                <label for="newPassword" class="block text-sm font-medium text-gray-700 mb-2">
                  新密码
                </label>
                <div class="relative">
                  <input id="newPassword" v-model="formData.newPassword" type="password" placeholder="请输入密码, 6-16个字符"
                    class="w-[394px] h-[46px] border border-gray-300 rounded-[5px] focus:outline-none focus:ring-2 focus:ring-[#CA898F] focus:border-transparent text-xs"
                    style="padding: 12px; font-size: 16px;" :class="{ 'border-red-500': errors.newPassword }" />
                  <span class="absolute left-[420px] top-1/2 -translate-y-1/2 text-[16px]"
                    style="color: #383838;">6-16个字符</span>
                </div>
                <p v-if="errors.newPassword" class="mt-1 text-sm text-red-600">
                  {{ errors.newPassword }}
                </p>
              </div>

              <!-- 确认密码字段 -->
              <div class="mb-[41px]">
                <label for="confirmPassword" class="block text-sm font-medium text-gray-700 mb-2">
                  确认密码
                </label>
                <input id="confirmPassword" v-model="formData.confirmPassword" type="password" placeholder="请再次输入新密码"
                  class="w-[394px] h-[46px] border border-gray-300 rounded-[5px] focus:outline-none focus:ring-2 focus:ring-[#CA898F] focus:border-transparent text-xs"
                  style="padding: 12px; font-size: 16px;" :class="{ 'border-red-500': errors.confirmPassword }" />
                <p v-if="errors.confirmPassword" class="mt-1 text-sm text-red-600">
                  {{ errors.confirmPassword }}
                </p>
              </div>

              <!-- 确认按钮 -->
              <div>
                <button type="submit" :disabled="isSubmitting"
                  class="w-[126px] h-[46px] bg-[#CA898F] text-white rounded-[5px] hover:bg-[#B87A80] transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                  style="padding-right: 31px; padding-left: 31px;">
                  {{ isSubmitting ? '修改中...' : '确认修改' }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import AppSidebar from '@/components/AppSidebar.vue'
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { changePassword } from '@/api/user'

defineOptions({
  name: 'ChangePassword'
})

const router = useRouter()
const activeMenu = ref('my')
const isSubmitting = ref(false)

// 表单数据
const formData = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 错误信息
const errors = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const handleMenuChange = (menu) => {
  activeMenu.value = menu
  // 路由跳转逻辑已经在 AppSidebar 组件中处理
}

// 验证表单
const validateForm = () => {
  let isValid = true

  // 重置错误信息
  errors.oldPassword = ''
  errors.newPassword = ''
  errors.confirmPassword = ''

  // 验证原密码
  if (!formData.oldPassword) {
    errors.oldPassword = '请输入原密码'
    isValid = false
  }

  // 验证新密码
  if (!formData.newPassword) {
    errors.newPassword = '请输入新密码'
    isValid = false
  } else if (formData.newPassword.length < 6 || formData.newPassword.length > 16) {
    errors.newPassword = '密码长度必须在6-16个字符之间'
    isValid = false
  }

  // 验证确认密码
  if (!formData.confirmPassword) {
    errors.confirmPassword = '请确认密码'
    isValid = false
  } else if (formData.newPassword !== formData.confirmPassword) {
    errors.confirmPassword = '两次输入的密码不一致'
    isValid = false
  }

  return isValid
}

// 提交表单
const handleSubmit = async () => {
  if (!validateForm()) {
    return
  }

  isSubmitting.value = true

  try {
    // TODO: 调用后端API修改密码
    // 示例：await changePassword(formData.newPassword)

    changePassword({
      password: formData.oldPassword,
      newPassword: formData.newPassword
    }).then((res) => {
      if (res.code === 0) {
        ElMessage.success('修改密码成功！')

        // 清空表单
        formData.oldPassword = ''
        formData.newPassword = ''
        formData.confirmPassword = ''
      }
      isSubmitting.value = false
    })

  } catch (error) {
    ElMessage.error('密码修改失败，请重试')
    console.error('修改密码失败:', error)
  } finally {
    isSubmitting.value = false
  }
}
</script>

<style scoped>
input::placeholder {
  font-size: 16px !important;
}
</style>