# 个人博客网站

一个功能完整的个人网站，支持博客写作、项目展示、知识花园、统计分析，配备完整的后台 CMS 管理系统。

## 技术栈

| 层 | 技术 |
|---|---|
| 前端 | 原生 HTML5 / CSS3 / Vanilla JS（无框架） |
| 后端 | Go 1.22 + Gin + GORM |
| 数据库 | SQLite（单文件，零配置） |
| 部署 | Docker Compose + Nginx |
| 统计 | 自建埋点 + Google Analytics 4 |

## 功能模块

- **博客**：Markdown 写作，支持置顶、标签、专栏、字数/阅读时间/浏览统计
- **作品集**：项目展示、系统演示（B站/图片/iframe）、自定义板块
- **知识花园**：多主题收藏（文章/视频/书籍/网站），多标签筛选
- **主页**：个人信息、技能标签、经历时间线、社交链接（后台可配置）
- **CMS 后台**：文章/专栏/项目/花园/主页 全套管理，JWT 鉴权
- **双主题**：暗黑 / 明亮，跟随系统偏好，用户可手动切换

## 快速启动

**依赖：** Go 1.22+、Git

```bash
git clone <repo-url>
cd blog/backend
go run ./cmd/server
```

访问 http://localhost:8080

后台地址：http://localhost:8080/admin/index.html
默认账号：`admin` / `admin123`（生产环境务必通过环境变量修改）

## 目录结构

```
blog/
├── frontend/           # 前端静态文件（由后端直接 serve）
│   ├── index.html      # 主页
│   ├── blog.html       # 博客列表
│   ├── post.html       # 文章详情
│   ├── column.html     # 专栏详情
│   ├── works.html      # 作品集
│   ├── garden.html     # 知识花园
│   ├── css/            # 主题变量 + 公共样式
│   ├── js/             # api.js / mock.js / theme.js / components.js
│   └── admin/          # 后台管理页面
├── backend/
│   ├── cmd/server/     # 程序入口
│   └── internal/
│       ├── config/     # 环境变量配置
│       ├── db/         # 数据库初始化 + AutoMigrate
│       ├── model/      # GORM 数据模型
│       ├── service/    # 业务逻辑
│       ├── handler/    # HTTP 处理器
│       ├── middleware/ # JWT / CORS
│       └── router/     # 路由注册
├── docs/               # 设计文档 + 实现计划
└── data/               # SQLite 数据库文件（gitignore）
```

## 环境变量

| 变量 | 默认值 | 说明 |
|---|---|---|
| `PORT` | `8080` | 服务端口 |
| `DB_PATH` | `./data/blog.db` | SQLite 文件路径 |
| `JWT_SECRET` | `dev-secret-...` | JWT 签名密钥，生产必改 |
| `ADMIN_USER` | `admin` | 后台用户名 |
| `ADMIN_PASS` | `admin123` | 后台密码，生产必改 |
| `FRONTEND_DIR` | `../frontend` | 前端静态文件目录 |

## Docker 部署

```bash
cd backend
docker build -t blog-api:latest .
```

完整 Docker Compose 部署参见 [wiki/05-deployment.md](wiki/05-deployment.md)。

## 文档

详细文档见 [`wiki/`](wiki/) 目录，适合从零开始理解本项目。

## License

MIT
