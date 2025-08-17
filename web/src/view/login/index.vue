<template>
  <div id="userLayout" class="w-full h-full relative" style="background-color: #F6EEEE;">
    <div class="w-full h-full flex items-center justify-center p-4">
             <!-- 白色卡片容器 -->
       <div class="bg-white rounded-lg shadow-lg p-8" style="width: 100%; max-width: 541px;">
                 <!-- 标题区域 -->
         <div class="text-left mb-8">
           <h1 class="text-[30px] font-medium text-gray-800 mb-2">悠闲果园</h1>
           <p class="text-[14px] font-normal" style="color: #AEAEAE;">注册账号&忘记密码? 请联系管理员</p>
         </div>

        <!-- 登录表单 -->
        <el-form
          ref="loginForm"
          :model="loginFormData"
          :rules="rules"
          :validate-on-rule-change="false"
          @keyup.enter="submitForm"
          style="margin-top: 38px;"
        >
          <!-- 用户名输入框 -->
          <el-form-item prop="username" class="mb-6" >
            <el-input
              v-model="loginFormData.username"
              size="large"
              placeholder="输入账号"
              class="rounded-lg"
            />
          </el-form-item>

          <!-- 密码输入框 -->
          <el-form-item prop="password" class="mb-6"> 
            <el-input
              v-model="loginFormData.password"
              show-password
              size="large"
              type="password"
              placeholder="输入密码"
              class="rounded-lg"
            />
          </el-form-item>

          <!-- 验证码输入框 -->
          <el-form-item
            v-if="loginFormData.openCaptcha"
            prop="captcha"
            class="mb-6"
          >
            <div class="flex w-full gap-3">
              <el-input
                v-model="loginFormData.captcha"
                placeholder="请输入验证码"
                size="large"
                class="flex-1 rounded-lg"
              />
              <div class="w-24 h-11 bg-gray-100 rounded-lg overflow-hidden cursor-pointer">
                <img
                  v-if="picPath"
                  class="w-full h-full object-cover"
                  :src="picPath"
                  alt="验证码"
                  @click="loginVerify()"
                />
              </div>
            </div>
          </el-form-item>

          <!-- 登录按钮 --> 
          <el-form-item class="mb-4">
            <el-button
              class="w-full h-12 rounded-lg text-white font-medium"
              type="primary"
              size="large"
              @click="submitForm"
              style="background-color: #e5b8b8; border-color: #e5b8b8;"
            >
              登录
            </el-button>
          </el-form-item>

                     <!-- 前往初始化按钮 -->
           <!-- <el-form-item class="mb-0">
             <el-button
               class="w-full h-12 rounded-lg"
               type="default"
               size="large"
               @click="checkInit"
               style="border-color: #e5b8b8; color: #e5b8b8;"
             >
               前往初始化
             </el-button>
           </el-form-item> -->
        </el-form>
      </div>
    </div>

    <!-- 底部信息 -->
    <BottomInfo class="left-0 right-0 absolute bottom-3 mx-auto w-full z-20">
      <div class="links items-center justify-center gap-2 hidden md:flex">
        <a href="https://www.gin-vue-admin.com/" target="_blank">
          <img src="@/assets/docs.png" class="w-8 h-8" alt="文档" />
        </a>
        <a href="https://support.qq.com/product/371961" target="_blank">
          <img src="@/assets/kefu.png" class="w-8 h-8" alt="客服" />
        </a>
        <a
          href="https://github.com/flipped-aurora/gin-vue-admin"
          target="_blank"
        >
          <img src="@/assets/github.png" class="w-8 h-8" alt="github" />
        </a>
        <a href="https://space.bilibili.com/322210472" target="_blank">
          <img src="@/assets/video.png" class="w-8 h-8" alt="视频站" />
        </a>
      </div>
    </BottomInfo>
  </div>
</template>

<script setup>
  import { captcha } from '@/api/user'
  import { checkDB } from '@/api/initdb'
  import BottomInfo from '@/components/bottomInfo/bottomInfo.vue'
  import { reactive, ref } from 'vue'
  import { ElMessage } from 'element-plus'
  import { useRouter } from 'vue-router'
  import { useUserStore } from '@/pinia/modules/user'

  defineOptions({
    name: 'Login'
  })

  const router = useRouter()
  // 验证函数
  const checkUsername = (rule, value, callback) => {
    if (value.length < 5) {
      return callback(new Error('请输入正确的用户名'))
    } else {
      callback()
    }
  }
  const checkPassword = (rule, value, callback) => {
    if (value.length < 6) {
      return callback(new Error('请输入正确的密码'))
    } else {
      callback()
    }
  }

  // 获取验证码
  const loginVerify = async () => {
    const ele = await captcha()
    rules.captcha.push({
      max: ele.data.captchaLength,
      min: ele.data.captchaLength,
      message: `请输入${ele.data.captchaLength}位验证码`,
      trigger: 'blur'
    })
    picPath.value = ele.data.picPath
    loginFormData.captchaId = ele.data.captchaId
    loginFormData.openCaptcha = ele.data.openCaptcha
  }
  loginVerify()

  // 登录相关操作
  const loginForm = ref(null)
  const picPath = ref('')
  const loginFormData = reactive({
    username: 'admin',
    password: '',
    captcha: '',
    captchaId: '',
    openCaptcha: false
  })
  const rules = reactive({
    username: [{ validator: checkUsername, trigger: 'blur' }],
    password: [{ validator: checkPassword, trigger: 'blur' }],
    captcha: [
      {
        message: '验证码格式不正确',
        trigger: 'blur'
      }
    ]
  })

  const userStore = useUserStore()
  const login = async () => {
    return await userStore.LoginIn(loginFormData)
  }
  const submitForm = () => {
    loginForm.value.validate(async (v) => {
      if (!v) {
        // 未通过前端静态验证
        ElMessage({
          type: 'error',
          message: '请正确填写登录信息',
          showClose: true
        })
        await loginVerify()
        return false
      }

      // 通过验证，请求登陆
      const flag = await login()

      // 登陆失败，刷新验证码
      if (!flag) {
        await loginVerify()
        return false
      }

      // 登陆成功
      return true
    })
  }

  // 跳转初始化
  const checkInit = async () => {
    const res = await checkDB()
    if (res.code === 0) {
      if (res.data?.needInit) {
        userStore.NeedInit()
        await router.push({ name: 'Init' })
      } else {
        ElMessage({
          type: 'info',
          message: '已配置数据库信息，无法初始化'
        })
      }
    }
  }
</script>

<style scoped>
.el-button--primary {
  background-color: #CA898F !important;
  border-color: #CA898F !important; 
}

.el-button--primary:hover {
  background-color: #d4a5a5 !important;
  border-color: #d4a5a5 !important;
}

.el-input__wrapper {
  border-radius: 8px !important;
  height: 58px !important;
}

.el-input__inner {
  padding: 10px 18px !important;
  height: 58px !important;
  line-height: 58px !important;
}

.el-button {
  border-radius: 8px !important;
}
</style>
