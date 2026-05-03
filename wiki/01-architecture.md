# 01 整体架构

## 系统概览

```
浏览器
  │
  ├── GET /、/blog.html、/admin/...  →  Go 服务器直接返回 HTML 文件
  │
  └── fetch /api/...               →  Go REST API 处理业务逻辑
                                           │
                                      SQLite 数据库
```

整个系统只有**一个 Go 进程**，同时承担两件事：

1. **静态文件服务**：把 `frontend/` 目录里的 HTML/CSS/JS 文件直接返回给浏览器
2. **REST API**：处理所有 `/api/` 开头的接口请求，读写 SQLite 数据库

这样部署最简单，只需启动一个程序，访问一个端口。

## 技术选型理由

### 为什么用原生 HTML/CSS/JS 而不是 React/Vue？

- 个人网站内容相对静态，没有复杂的组件状态管理需求
- 无需构建步骤，改完 HTML 直接刷新浏览器看效果
- 学习成本低，代码直观易读
- 缺点：缺乏组件化，重复代码较多（通过 `components.js` 部分缓解）

### 为什么用 Go 而不是 Java/Node？

- 编译成单个二进制文件，Docker 镜像只有 40MB
- 启动极快（毫秒级），资源占用低
- 原生支持并发，适合处理统计埋点等轻量异步任务

### 为什么用 SQLite 而不是 MySQL/PostgreSQL？

- 个人博客并发量低，SQLite 完全够用
- 零配置，数据库就是一个文件 `data/blog.db`
- 备份简单：复制文件即可
- 迁移方便：换服务器时把文件复制过去

## 目录职责

```
frontend/
├── css/
│   ├── common.css       # 所有页面共用的布局和组件样式
│   ├── theme-dark.css   # 暗黑主题的 CSS 变量定义
│   └── theme-light.css  # 明亮主题的 CSS 变量定义
├── js/
│   ├── api.js           # 所有后端接口的封装（fetch 请求）
│   ├── mock.js          # 模拟数据，开发时不需要启动后端
│   ├── theme.js         # 主题切换逻辑
│   ├── components.js    # 共用组件（导航栏、标签渲染等）
│   └── nav.js           # 导航栏当前页高亮
├── *.html               # 6 个前台页面
└── admin/               # 后台管理页面

backend/
├── cmd/server/main.go   # 程序入口，组装所有依赖
└── internal/
    ├── config/          # 读取环境变量
    ├── db/              # 连接数据库，自动建表
    ├── model/           # 数据库表结构定义（GORM 模型）
    ├── service/         # 业务逻辑（增删改查、字数统计等）
    ├── handler/         # HTTP 请求处理（参数解析、调 service、返回 JSON）
    ├── middleware/       # JWT 鉴权、CORS 跨域
    └── router/          # 注册所有路由
```

## 请求流程示例

以访问博客列表页为例：

```
用户访问 http://localhost:8080/blog.html
  │
  ├── 1. Nginx/Go 返回 blog.html 文件
  │
  ├── 2. 浏览器执行 JS，调用 api.getTags() 和 api.getPosts()
  │         发送 GET /api/tags 和 GET /api/posts?page=1
  │
  ├── 3. Go router 匹配路由，调用 PostHandler.List()
  │
  ├── 4. PostHandler 调用 PostService.List()
  │
  ├── 5. PostService 查询 SQLite，返回文章列表
  │
  └── 6. PostHandler 将结果序列化为 JSON 返回
            浏览器渲染文章卡片
```

## Phase 2（未实现）

当前是 Phase 1，以下功能计划在 Phase 2 实现：
- 可视化页面构建器（拖拽编排页面布局）
- 评论系统
- 全文搜索
