<template>
    <div class="flex h-screen bg-[#F6EEEE] p-[15px]">
        <!-- 侧边栏 -->
        <AppSidebar :active-menu="activeMenu" @menu-change="handleMenuChange" />

        <!-- 主内容区域 -->
        <div class="flex-1 ml-3 bg-white rounded-[10px] overflow-hidden">
            <div class="h-full p-6">
            <div class="gva-search-box">
                <el-form ref="searchForm" :inline="true" :model="searchInfo">
                    <el-form-item label="用户名">
                        <el-input v-model="searchInfo.username" placeholder="用户名" />
                    </el-form-item>
                    <el-form-item label="昵称">
                        <el-input v-model="searchInfo.nickname" placeholder="昵称" />
                    </el-form-item>
                    <el-form-item label="手机号">
                        <el-input v-model="searchInfo.phone" placeholder="手机号" />
                    </el-form-item>
                    <el-form-item label="邮箱">
                        <el-input v-model="searchInfo.email" placeholder="邮箱" />
                    </el-form-item>
                    <el-form-item>
                        <el-button type="primary" icon="search" @click="onSubmit">
                            查询
                        </el-button>
                        <el-button icon="refresh" @click="onReset"> 重置 </el-button>
                    </el-form-item>
                </el-form>
            </div>
            <div class="gva-table-box">
                <div class="gva-btn-list">
                    <el-button type="primary" icon="plus" @click="addUser">新增用户</el-button>
                </div>
                <el-table :data="tableData" row-key="ID">
                    <el-table-column align="left" label="头像" min-width="75">
                        <template #default="scope">
                            <CustomPic style="margin-top: 8px" :pic-src="scope.row.headerImg" />
                        </template>
                    </el-table-column>
                    <el-table-column align="left" label="ID" min-width="50" prop="ID" />
                    <el-table-column align="left" label="用户名" min-width="150" prop="userName" />
                    <el-table-column align="left" label="昵称" min-width="150" prop="nickName" />
                    <el-table-column align="left" label="手机号" min-width="180" prop="phone" />
                    <el-table-column align="left" label="邮箱" min-width="180" prop="email" />
                    <el-table-column align="left" label="用户角色" min-width="200">
                        <template #default="scope">
                            <el-cascader
                                v-if="userStore.userInfo.ID !== scope.row.ID"
                                v-model="scope.row.authorityId"
                                :options="authOptions"
                                :show-all-levels="false"
                                :props="{
                                    multiple: false,
                                    checkStrictly: true,
                                    label: 'authorityName',
                                    value: 'authorityId',
                                    disabled: 'disabled',
                                    emitPath: false
                                }"
                                :clearable="false"
                                @change="() => changeAuthority(scope.row)"
                            />
                            <span style="font-size: 14px; color: #A6A6A6;" v-else>{{ "我（超级管理员）" }}</span>
                        </template>
                    </el-table-column>
                    <el-table-column align="left" label="启用" min-width="150">
                        <template #default="scope">
                            <el-switch v-if="userStore.userInfo.ID !== scope.row.ID" v-model="scope.row.enable" inline-prompt :active-value="1" :inactive-value="2"
                                @change="
                                    () => {
                                        switchEnable(scope.row)
                                    }
                                " />
                        </template>
                    </el-table-column>

                    <el-table-column label="操作" :min-width="appStore.operateMinWith" fixed="right">
                        <template #default="scope">
                            <el-button type="primary" link icon="delete" v-if="userStore.userInfo.ID !== scope.row.ID"
                                @click="deleteUserFunc(scope.row)">删除</el-button>
                            <el-button type="primary" link icon="edit" v-if="userStore.userInfo.ID !== scope.row.ID" @click="openEdit(scope.row)">编辑</el-button>
                            <el-button type="primary" link icon="magic-stick" v-if="userStore.userInfo.ID !== scope.row.ID"
                                @click="resetPasswordFunc(scope.row)">重置密码</el-button>
                        </template>
                    </el-table-column>
                </el-table>
                <div class="gva-pagination">
                    <el-pagination :current-page="page" :page-size="pageSize" :page-sizes="[10, 30, 50, 100]"
                        :total="total" layout="total, sizes, prev, pager, next, jumper"
                        @current-change="handleCurrentChange" @size-change="handleSizeChange" />
                </div>
            </div>
            <!-- 重置密码对话框 -->
            <el-dialog v-model="resetPwdDialog" title="重置密码" width="500px" :close-on-click-modal="false"
                :close-on-press-escape="false">
                <el-form :model="resetPwdInfo" ref="resetPwdForm" label-width="100px">
                    <el-form-item label="用户账号">
                        <el-input v-model="resetPwdInfo.userName" disabled />
                    </el-form-item>
                    <el-form-item label="用户昵称">
                        <el-input v-model="resetPwdInfo.nickName" disabled />
                    </el-form-item>
                    <el-form-item label="新密码">
                        <div class="flex w-full">
                            <el-input class="flex-1" v-model="resetPwdInfo.password" placeholder="请输入新密码"
                                show-password />
                            <el-button type="primary" @click="generateRandomPassword" style="margin-left: 10px">
                                生成随机密码
                            </el-button>
                        </div>
                    </el-form-item>
                </el-form>
                <template #footer>
                    <div class="dialog-footer">
                        <el-button @click="closeResetPwdDialog">取 消</el-button>
                        <el-button type="primary" @click="confirmResetPassword">确 定</el-button>
                    </div>
                </template>
            </el-dialog>

            <el-drawer v-model="addUserDialog" :size="appStore.drawerSize" :show-close="false"
                :close-on-press-escape="false" :close-on-click-modal="false">
                <template #header>
                    <div class="flex justify-between items-center">
                        <span class="text-lg">用户</span>
                        <div>
                            <el-button @click="closeAddUserDialog">取 消</el-button>
                            <el-button type="primary" @click="enterAddUserDialog">确 定</el-button>
                        </div>
                    </div>
                </template>

                <el-form ref="userForm" :rules="rules" :model="userInfo" label-width="80px">
                    <el-form-item v-if="dialogFlag === 'add'" label="用户名" prop="userName">
                        <el-input v-model="userInfo.userName" />
                    </el-form-item>
                    <el-form-item v-if="dialogFlag === 'add'" label="密码" prop="password">
                        <el-input v-model="userInfo.password" />
                    </el-form-item>
                    <el-form-item label="昵称" prop="nickName">
                        <el-input v-model="userInfo.nickName" />
                    </el-form-item>
                    <el-form-item label="手机号" prop="phone">
                        <el-input v-model="userInfo.phone" />
                    </el-form-item>
                    <el-form-item label="邮箱" prop="email">
                        <el-input v-model="userInfo.email" />
                    </el-form-item>
                    <el-form-item label="用户角色" prop="authorityId">
                        <el-cascader
                            v-model="userInfo.authorityId"
                            style="width: 100%"
                            :options="authOptions"
                            :show-all-levels="false"
                            :props="{
                                multiple: false,
                                checkStrictly: true,
                                label: 'authorityName',
                                value: 'authorityId',
                                disabled: 'disabled',
                                emitPath: false
                            }"
                            :clearable="false"
                        />
                    </el-form-item>
                    <el-form-item label="启用" prop="disabled">
                        <el-switch v-model="userInfo.enable" inline-prompt :active-value="1" :inactive-value="2" />
                    </el-form-item>
                    <el-form-item label="头像" label-width="80px">
                        <SelectImage v-model="userInfo.headerImg" />
                    </el-form-item>
                </el-form>
            </el-drawer>
            </div>
        </div>
    </div>
</template>

<script setup>
import {
    getUserList,
    setUserAuthorities,
    register,
    deleteUser
} from '@/api/user'

import { getAuthorityList } from '@/api/authority'
import CustomPic from '@/components/customPic/index.vue'
import WarningBar from '@/components/warningBar/warningBar.vue'
import { setUserInfo, resetPassword } from '@/api/user.js'

import { nextTick, ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import SelectImage from '@/components/selectImage/selectImage.vue'
import { useAppStore } from "@/pinia";
import AppSidebar from '@/components/AppSidebar.vue'
import { useUserStore } from '@/pinia/modules/user'

const userStore = useUserStore()

defineOptions({
    name: 'User'
})

const appStore = useAppStore()
const activeMenu = ref('member')

const searchInfo = ref({
    username: '',
    nickname: '',
    phone: '',
    email: ''
})

const onSubmit = () => {
    page.value = 1
    getTableData()
}

const handleMenuChange = (menu) => {
    activeMenu.value = menu
    // 路由跳转逻辑已经在 AppSidebar 组件中处理
}
const onReset = () => {
    searchInfo.value = {
        username: '',
        nickname: '',
        phone: '',
        email: ''
    }
    getTableData()
}
// 初始化相关
const setAuthorityOptions = (AuthorityData, optionsData) => {
    AuthorityData &&
        AuthorityData.forEach((item) => {
            if (item.children && item.children.length) {
                const option = {
                    authorityId: item.authorityId,
                    authorityName: item.authorityName,
                    children: []
                }
                setAuthorityOptions(item.children, option.children)
                optionsData.push(option)
            } else {
                const option = {
                    authorityId: item.authorityId,
                    authorityName: item.authorityName
                }
                optionsData.push(option)
            }
        })
}

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
// 分页
const handleSizeChange = (val) => {
    pageSize.value = val
    getTableData()
}

const handleCurrentChange = (val) => {
    page.value = val
    getTableData()
}

// 查询
const getTableData = async () => {
    const table = await getUserList({
        page: page.value,
        pageSize: pageSize.value,
        ...searchInfo.value
    })
    if (table.code === 0) {
        tableData.value = table.data.list
        total.value = table.data.total
        page.value = table.data.page
        pageSize.value = table.data.pageSize
    }
}

watch(
    () => tableData.value,
    () => {
        setAuthorityIds()
    }
)

const initPage = async () => {
    getTableData()
    const res = await getAuthorityList()
    setOptions(res.data)
}

initPage()

// 重置密码对话框相关
const resetPwdDialog = ref(false)
const resetPwdForm = ref(null)
const resetPwdInfo = ref({
    ID: '',
    userName: '',
    nickName: '',
    password: ''
})

// 生成随机密码
const generateRandomPassword = () => {
    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*'
    let password = ''
    for (let i = 0; i < 12; i++) {
        password += chars.charAt(Math.floor(Math.random() * chars.length))
    }
    resetPwdInfo.value.password = password
    // 复制到剪贴板
    navigator.clipboard.writeText(password).then(() => {
        ElMessage({
            type: 'success',
            message: '密码已复制到剪贴板'
        })
    }).catch(() => {
        ElMessage({
            type: 'error',
            message: '复制失败，请手动复制'
        })
    })
}

// 打开重置密码对话框
const resetPasswordFunc = (row) => {
    resetPwdInfo.value.ID = row.ID
    resetPwdInfo.value.userName = row.userName
    resetPwdInfo.value.nickName = row.nickName
    resetPwdInfo.value.password = ''
    resetPwdDialog.value = true
}

// 确认重置密码
const confirmResetPassword = async () => {
    if (!resetPwdInfo.value.password) {
        ElMessage({
            type: 'warning',
            message: '请输入或生成密码'
        })
        return
    }

    const res = await resetPassword({
        ID: resetPwdInfo.value.ID,
        password: resetPwdInfo.value.password
    })

    if (res.code === 0) {
        ElMessage({
            type: 'success',
            message: res.msg || '密码重置成功'
        })
        resetPwdDialog.value = false
    } else {
        ElMessage({
            type: 'error',
            message: res.msg || '密码重置失败'
        })
    }
}

// 关闭重置密码对话框
const closeResetPwdDialog = () => {
    resetPwdInfo.value.password = ''
    resetPwdDialog.value = false
}
const setAuthorityIds = () => {
    tableData.value &&
        tableData.value.forEach((user) => {
            // 取该用户的第一个角色作为当前选中角色
            if (user.authorities && user.authorities.length) {
                user.authorityId = user.authorities[0].authorityId
            } else if (user.authorityIds && user.authorityIds.length) {
                user.authorityId = user.authorityIds[0]
            } else {
                user.authorityId = undefined
            }
        })
}

const authOptions = ref([])
const setOptions = (authData) => {
    authOptions.value = []
    setAuthorityOptions(authData, authOptions.value)
    // 禁用超级管理员（ID为 888）
    const disableSuperAdmin = (list) => {
        list.forEach((opt) => {
            if (opt.authorityId === 888) {
                opt.disabled = true
            }
            if (opt.children && opt.children.length) disableSuperAdmin(opt.children)
        })
    }
    disableSuperAdmin(authOptions.value)
}

const deleteUserFunc = async (row) => {
    ElMessageBox.confirm('确定要删除吗?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(async () => {
        const res = await deleteUser({ id: row.ID })
        if (res.code === 0) {
            ElMessage.success('删除成功')
            await getTableData()
        }
    })
}

// 弹窗相关
const userInfo = ref({
    userName: '',
    password: '',
    nickName: '',
    headerImg: '',
    authorityId: '',
    enable: 1
})

const rules = ref({
    userName: [
        { required: true, message: '请输入用户名', trigger: 'blur' },
        { min: 5, message: '最低5位字符', trigger: 'blur' }
    ],
    password: [
        { required: true, message: '请输入用户密码', trigger: 'blur' },
        { min: 6, message: '最低6位字符', trigger: 'blur' }
    ],
    nickName: [{ required: true, message: '请输入用户昵称', trigger: 'blur' }],
    phone: [
        {
            pattern: /^1([38][0-9]|4[014-9]|[59][0-35-9]|6[2567]|7[0-8])\d{8}$/,
            message: '请输入合法手机号',
            trigger: 'blur'
        }
    ],
    email: [
        {
            pattern: /^([0-9A-Za-z\-_.]+)@([0-9a-z]+\.[a-z]{2,3}(\.[a-z]{2})?)$/g,
            message: '请输入正确的邮箱',
            trigger: 'blur'
        }
    ],
    authorityId: [
        { required: true, message: '请选择用户角色', trigger: 'blur' }
    ]
})
const userForm = ref(null)
const enterAddUserDialog = async () => {
    userForm.value.validate(async (valid) => {
        if (valid) {
            const req = {
                ...userInfo.value,
                authorityIds: userInfo.value.authorityId ? [userInfo.value.authorityId] : []
            }
            if (dialogFlag.value === 'add') {
                const res = await register(req)
                if (res.code === 0) {
                    ElMessage({ type: 'success', message: '创建成功' })
                    await getTableData()
                    closeAddUserDialog()
                }
            }
            if (dialogFlag.value === 'edit') {
                const res = await setUserInfo(req)
                if (res.code === 0) {
                    ElMessage({ type: 'success', message: '编辑成功' })
                    await getTableData()
                    closeAddUserDialog()
                }
            }
        }
    })
}

const addUserDialog = ref(false)
const closeAddUserDialog = () => {
    userForm.value.resetFields()
    userInfo.value.headerImg = ''
    userInfo.value.authorityIds = []
    addUserDialog.value = false
}

const dialogFlag = ref('add')

const addUser = () => {
    dialogFlag.value = 'add'
    addUserDialog.value = true
}

const tempAuth = {}
const changeAuthority = async (row) => {
    await nextTick()
    const res = await setUserAuthorities({
        ID: row.ID,
        authorityIds: [row.authorityId]
    })
    if (res.code === 0) {
        ElMessage({ type: 'success', message: '角色设置成功' })
    }
}

const openEdit = (row) => {
    dialogFlag.value = 'edit'
    userInfo.value = JSON.parse(JSON.stringify(row))
    addUserDialog.value = true
}

const switchEnable = async (row) => {
    userInfo.value = JSON.parse(JSON.stringify(row))
    await nextTick()
    const req = {
        ...userInfo.value
    }
    const res = await setUserInfo(req)
    if (res.code === 0) {
        ElMessage({
            type: 'success',
            message: `${req.enable === 2 ? '禁用' : '启用'}成功`
        })
        await getTableData()
        userInfo.value.headerImg = ''
        userInfo.value.authorityIds = []
    }
}
</script>

<style lang="scss">
.header-img-box {
    @apply w-52 h-52 border border-solid border-gray-300 rounded-xl flex justify-center items-center cursor-pointer;
}
</style>