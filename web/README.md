# gin-vue-admin web

## Project setup

```
npm install
```

### Compiles and hot-reloads for development

```
npm run serve
```

### Compiles and minifies for production

```
npm run build
```

### Run your tests

```
npm run test
```

### Lints and fixes files

```
npm run lint
```

整理代码结构

```lua
web
 ├── babel.config.js
 ├── Dockerfile
 ├── favicon.ico
 ├── index.html                 -- 主页面
 ├── limit.js                   -- 助手代码
 ├── package.json               -- 包管理器代码
 ├── src                        -- 源代码
 │   ├── api                    -- api 组
 │   ├── App.vue                -- 主页面
 │   ├── assets                 -- 静态资源
 │   ├── components             -- 全局组件
 │   ├── core                   -- gva 组件包
 │   │   ├── config.js          -- gva网站配置文件
 │   │   ├── gin-vue-admin.js   -- 注册欢迎文件
 │   │   └── global.js          -- 统一导入文件
 │   ├── directive              -- v-auth 注册文件
 │   ├── main.js                -- 主文件
 │   ├── permission.js          -- 路由中间件
 │   ├── pinia                  -- pinia 状态管理器，取代vuex
 │   │   ├── index.js           -- 入口文件
 │   │   └── modules            -- modules
 │   │       ├── dictionary.js
 │   │       ├── router.js
 │   │       └── user.js
 │   ├── router                 -- 路由声明文件
 │   │   └── index.js
 │   ├── style                  -- 全局样式
 │   │   ├── base.scss
 │   │   ├── basics.scss
 │   │   ├── element_visiable.scss  -- 此处可以全局覆盖 element-plus 样式
 │   │   ├── iconfont.css           -- 顶部几个icon的样式文件
 │   │   ├── main.scss
 │   │   ├── mobile.scss
 │   │   └── newLogin.scss
 │   ├── utils                  -- 方法包库
 │   │   ├── asyncRouter.js     -- 动态路由相关
 │   │   ├── bus.js             -- 全局mitt声明文件
 │   │   ├── date.js            -- 日期相关
 │   │   ├── dictionary.js      -- 获取字典方法
 │   │   ├── downloadImg.js     -- 下载图片方法
 │   │   ├── format.js          -- 格式整理相关
 │   │   ├── image.js           -- 图片相关方法
 │   │   ├── page.js            -- 设置页面标题
 │   │   ├── request.js         -- 请求
 │   │   └── stringFun.js       -- 字符串文件
 |   ├── view -- 主要view代码
 |   |   ├── about -- 关于我们
 |   |   ├── dashboard -- 面板
 |   |   ├── error -- 错误
 |   |   ├── example --上传案例
 |   |   ├── iconList -- icon列表
 |   |   ├── init -- 初始化数据
 |   |   |   ├── index -- 新版本
 |   |   |   ├── init -- 旧版本
 |   |   ├── layout  --  layout约束页面
 |   |   |   ├── aside
 |   |   |   ├── bottomInfo     -- bottomInfo
 |   |   |   ├── screenfull     -- 全屏设置
 |   |   |   ├── setting        -- 系统设置
 |   |   |   └── index.vue      -- base 约束
 |   |   ├── login              --登录
 |   |   ├── person             --个人中心
 |   |   ├── superAdmin         -- 超级管理员操作
 |   |   ├── system             -- 系统检测页面
 |   |   ├── systemTools        -- 系统配置相关页面
 |   |   └── routerHolder.vue   -- page 入口页面
 ├── vite.config.js             -- vite 配置文件
 └── yarn.lock

```

打包web
npm run build

打包server
set GOARCH=amd64
set GOOS=linux
go build

服务器
安装mysql

# 更新包列表
sudo apt update

# 安装 MySQL
sudo apt install mysql-server

# 启动 MySQL 服务
sudo systemctl start mysql
sudo systemctl enable mysql

安装radis
sudo apt install redis-server
sudo systemctl start redis-server
sudo systemctl enable redis-server

安装nginx
sudo apt install nginx -y
sudo systemctl start nginx         # 启动Nginx
sudo systemctl enable nginx        # 设置开机自启
sudo systemctl status nginx         # 查看运行状态
sudo ufw allow 'Nginx Full'
重载ngnix
sudo nginx -t
sudo systemctl reload nginx

sudo vim /etc/ngnix/ngnix.conf
server {
            listen 80;
            server_name 121.43.127.128; # 你的域名或服务器IP
            root /var/www/orchard; # dist目录的绝对路径
            index index.html;

            location / {
                 try_files $uri $uri/ /index.html; # 支持前端路由
            }

            # 可选：配置反向代理解决跨域（如需调用后端API）
            location /api/ {
                proxy_pass http://127.0.0.1:8888/;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            }

            location /apiFileupload/ {
                alias /var/www/orchard/pipixia/uploads/file/;
                access_log off;
                expires 30d;
                add_header Cache-Control "public, max-age=2592000";
                try_files $uri =404;
                location ~* \.(php|pl|py|sh|cgi)$ {
                    deny all;
                 }
                }
        }

        启动服务
        nohup ./server > service.log 2>&1 &