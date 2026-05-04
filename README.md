# 个人博客网站

一个功能完整的个人网站，支持博客写作、项目展示、知识花园、统计分析，配备完整的后台 CMS 管理系统。

## 功能模块

- **博客**：Markdown 写作，支持置顶、标签、专栏、字数/阅读时间/浏览统计
- **作品集**：项目展示、系统演示（B站/图片/iframe）、自定义板块
- **知识花园**：多主题收藏（文章/视频/书籍/网站），多标签筛选
- **主页**：个人信息、技能标签、经历时间线、社交链接（后台可配置）
- **CMS 后台**：文章/专栏/项目/花园/主页 全套管理，JWT 鉴权
- **双主题**：暗黑 / 明亮，跟随系统偏好，用户可手动切换

## 技术栈

| 层 | 技术 |
|---|---|
| 前端 | 原生 HTML5 / CSS3 / Vanilla JS（无框架） |
| 后端 | Go 1.22 + Gin + GORM |
| 数据库 | SQLite（单文件，零配置） |
| 部署 | Docker Compose |
| 统计 | 自建埋点 + Google Analytics 4 |

---

## 部署

### 方式一：Docker（推荐）

**依赖：** [Docker](https://docs.docker.com/get-docker/)（macOS：`brew install --cask docker`）

```bash
git clone <repo-url>
cd blog
./deploy.sh
```

脚本全程交互引导，**无需手动编辑任何文件**：
- 自动检测 Docker 是否就绪
- 询问端口、用户名、密码（全有默认值，直接回车跳过）
- 密码留空时自动生成强密码并显示
- JWT 密钥完全自动生成

```bash
./deploy.sh              # 启动 / 重新部署
./deploy.sh stop         # 停止
./deploy.sh restart      # 重启
./deploy.sh logs         # 查看日志
./deploy.sh update       # 更新到最新版本
```

### 方式二：直接运行（开发者）

**依赖：** [Go 1.22+](https://go.dev/dl/)（macOS：`brew install go`）

```bash
git clone <repo-url>
cd blog
./start.sh
```

```bash
./start.sh stop|restart|logs|status
```

### 访问

| 地址 | 说明 |
|---|---|
| `http://localhost:8080` | 前台网站 |
| `http://localhost:8080/admin` | 后台管理 |

服务器部署、内网穿透、数据备份等详见 [wiki/05-deployment.md](wiki/05-deployment.md)。

---

## 目录结构

```
blog/
├── frontend/           # 前端静态文件
│   ├── *.html          # 前台页面
│   ├── admin/          # 后台管理页面
│   ├── css/            # 主题样式
│   └── js/             # api / mock / theme / components
├── backend/
│   ├── cmd/server/     # 程序入口
│   └── internal/       # config / db / model / service / handler / router
├── wiki/               # 详细文档
├── .env.example        # 配置模板
├── docker-compose.yml  # Docker Compose 配置
├── deploy.sh           # Docker 一键部署脚本
└── start.sh            # 直接运行脚本（开发者）
```

## 文档

| 文档 | 内容 |
|---|---|
| [wiki/01-architecture.md](wiki/01-architecture.md) | 整体架构与设计决策 |
| [wiki/02-frontend.md](wiki/02-frontend.md) | 前端结构、主题系统 |
| [wiki/03-backend.md](wiki/03-backend.md) | API 设计、数据库 Schema |
| [wiki/04-admin.md](wiki/04-admin.md) | 后台使用指南 |
| [wiki/05-deployment.md](wiki/05-deployment.md) | 完整部署文档 |

## License

MIT
