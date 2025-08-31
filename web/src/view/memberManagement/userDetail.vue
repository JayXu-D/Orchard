<template>
    <div class="flex h-screen bg-[#F6EEEE] p-[15px]">
        <!-- 侧边栏 -->
        <AppSidebar :active-menu="activeMenu" @menu-change="handleMenuChange" />

        <!-- 主内容区域 -->
        <div class="flex-1 ml-3 bg-white rounded-[10px] overflow-hidden">
            <div class="h-full p-6">
                <!-- 返回按钮和用户信息 -->
                <div class="mb-6">
                    <el-button 
                        type="text" 
                        icon="ArrowLeft" 
                        @click="goBack"
                        class="mb-4 text-gray-600 hover:text-blue-600"
                    >
                        返回成员管理
                    </el-button>
                    
                    <div class="text-center mb-6">
                        <div class="text-4xl font-bold text-gray-800 mb-2">{{ userInfo.userName || '用户ID' }}</div>
                        <div class="text-lg text-gray-600">注册时间: {{ formatDate(userInfo.createdAt) }}</div>
                    </div>
                </div>

                <!-- 已获得的图纸 -->
                <div class="bg-white rounded-lg">
                    <h2 class="text-2xl font-bold text-gray-800 mb-6 text-center">已获得的图纸</h2>
                    
                    <el-table :data="drawingsList" style="width: 100%" v-loading="loading">
                        <el-table-column prop="serialNumber" label="序号" min-width="120" align="center">
                            <template #default="scope">
                                <span class="font-medium">{{ scope.row.serialNumber }}</span>
                            </template>
                        </el-table-column>
                        
                        <el-table-column prop="drawingName" label="图纸名称" min-width="300" align="center">
                            <template #default="scope">
                                <span class="text-gray-700">{{ scope.row.drawingName }}</span>
                            </template>
                        </el-table-column>
                        
                        <el-table-column prop="acquisitionTime" label="获得图纸时间" min-width="180" align="center">
                            <template #default="scope">
                                <span class="text-gray-600">{{ formatDate(scope.row.acquisitionTime) }}</span>
                            </template>
                        </el-table-column>
                        
                        <el-table-column prop="lastDownloadTime" label="最后下载时间" min-width="180" align="center">
                            <template #default="scope">
                                <span v-if="scope.row.lastDownloadTime" class="text-gray-600">
                                    {{ formatDate(scope.row.lastDownloadTime) }}
                                </span>
                                <span v-else class="text-gray-400">未下载</span>
                            </template>
                        </el-table-column>
                    </el-table>
                    
                    <!-- 空状态 -->
                    <div v-if="!loading && drawingsList.length === 0" class="text-center py-12">
                        <el-empty description="该用户暂无获得的图纸" />
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import AppSidebar from '@/components/AppSidebar.vue'
import { getUserDetail, getUserDrawings } from '@/api/user'

defineOptions({
    name: 'UserDetail'
})

const route = useRoute()
const router = useRouter()
const activeMenu = ref('member')

// 用户信息
const userInfo = ref({})
// 图纸列表
const drawingsList = ref([])
// 加载状态
const loading = ref(false)

// 获取用户详情
const fetchUserDetail = async () => {
    const userId = route.params.id
    if (!userId) {
        ElMessage.error('用户ID不能为空')
        return
    }
    
    try {
        loading.value = true
        const res = await getUserDetail(userId)
        if (res.code === 0) {
            userInfo.value = res.data
        } else {
            ElMessage.error(res.msg || '获取用户详情失败')
        }
        
        // 获取用户图纸列表
        await fetchUserDrawings(userId)
    } catch (error) {
        console.error('获取用户详情失败:', error)
        ElMessage.error('获取用户详情失败')
    } finally {
        loading.value = false
    }
}

// 获取用户图纸列表
const fetchUserDrawings = async (userId) => {
    try {
        const res = await getUserDrawings(userId)
        if (res.code === 0) {
            drawingsList.value = res.data
        } else {
            ElMessage.error(res.msg || '获取用户图纸列表失败')
        }
    } catch (error) {
        console.error('获取用户图纸列表失败:', error)
        ElMessage.error('获取用户图纸列表失败')
    }
}

// 格式化日期
const formatDate = (dateString) => {
    if (!dateString) return ''
    const date = new Date(dateString)
    return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: 'numeric',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    }).replace(/\//g, '-')
}

// 返回上一页
const goBack = () => {
    router.go(-1)
}

// 菜单切换
const handleMenuChange = (menu) => {
    activeMenu.value = menu
    // 路由跳转逻辑已经在 AppSidebar 组件中处理
}

onMounted(() => {
    fetchUserDetail()
})
</script>

<style lang="scss" scoped>
.el-table {
    border-radius: 8px;
    overflow: hidden;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.el-table th {
    background-color: #f5f7fa;
    color: #606266;
    font-weight: 600;
}

.el-table td {
    border-bottom: 1px solid #ebeef5;
}

.el-button--text {
    font-size: 16px;
    padding: 0;
}

.el-button--text:hover {
    background-color: transparent;
}
</style>
