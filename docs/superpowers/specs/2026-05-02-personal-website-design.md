# 个人网站设计文档 — Phase 1 核心展示站

## 概览

个人品牌展示网站，面向所有访客（招聘者、技术社区、普通读者）。支持博客文章、项目展示、系统演示、知识花园四大内容模块，配套自建 CMS 后台进行内容管理。双主题（极客暗黑 / 现代明亮）切换，Docker Compose 容器化部署。

Phase 2（可视化页面构建器）单独立项，本文档不涉及。

---

## 技术栈

### 前端
- 原生 HTML5 / CSS3 / Vanilla JS (ES6+)
- marked.js — Markdown 客户端渲染
- highlight.js — 代码块高亮
- CSS 自定义变量 — 双主题切换

### 后端
- Go 1.22+
- Gin — HTTP 框架
- GORM — ORM
- SQLite — 数据库（单文件，零配置）
- goldmark — 服务端 Markdown 解析（计算字数等元信息）
- golang-jwt/jwt — JWT 鉴权

### 基础设施
- Docker + Docker Compose
- Nginx — 静态文件托管 + 反向代理 + SSL 终止
- Let's Encrypt — HTTPS 证书
- Google Analytics 4 — 第三方统计
- GitHub Actions — CI（build & push Docker image）

---

## 项目目录结构

```
blog/
├── frontend/
│   ├── index.html          # 主页
│   ├── blog.html           # 博客列表（文章+专栏两视图）
│   ├── post.html           # 单篇文章
│   ├── column.html         # 专栏详情页
│   ├── works.html          # 作品集
│   ├── garden.html         # 知识花园
│   ├── admin/              # CMS 后台（受 JWT 保护）
│   │   ├── index.html      # 登录页
│   │   ├── dashboard.html  # 统计面板
│   │   ├── posts.html      # 文章管理
│   │   └── editor.html     # Markdown 在线编辑器
│   ├── css/
│   │   ├── common.css      # 公共样式
│   │   ├── theme-dark.css  # 暗黑主题变量
│   │   └── theme-light.css # 明亮主题变量
│   └── js/
│       ├── theme.js        # 主题切换 + 持久化
│       ├── api.js          # fetch 封装，统一调用后端
│       └── markdown.js     # marked.js 初始化配置
├── backend/
│   ├── cmd/server/main.go
│   └── internal/
│       ├── handler/        # Gin 路由处理器
│       ├── model/          # GORM 模型
│       ├── service/        # 业务逻辑
│       └── middleware/     # JWT 中间件、埋点中间件
├── nginx/
│   ├── nginx.conf
│   └── ssl/
├── data/                   # SQLite 文件（.gitignore）
├── docker-compose.yml
├── docker-compose.dev.yml  # 本地开发覆盖配置
└── .github/workflows/ci.yml
```

---

## 页面设计

### 导航栏（所有页面公用）
- Logo：`dyt`（monospace 字体，accent 色）
- 链接：主页 / 博客 / 作品集 / 知识花园
- 右侧：主题切换按钮（🌙 / ☀️）

### index.html — 主页
内容区块从上到下：
1. **Hero**：名字、职业标签、一句话介绍、GitHub/邮件链接
2. **技能标签**：Go / Java / Python / Docker 等，彩色徽章
3. **经历时间线**：教育 + 工作经历，简洁时间轴样式
4. **精选项目卡片**：2-3 个精选项目，从 works 数据动态加载
5. **最新文章**：最新 3-5 篇博客，从 blog API 动态加载
6. **联系方式**：GitHub / 邮件 / 其他社交入口

### blog.html — 博客列表
页面分两个视图 Tab 切换：**文章** / **专栏**

**文章视图：**
- 置顶文章区（绿色边框，📌 标识，置顶文章排在列表最前）
- 标签筛选栏（全部 + 各 tag）
- 关键词搜索框（前端过滤 or 调 API）
- 文章卡片列表：标题、摘要、日期、👁 浏览人次、📖 预计阅读时长、字数、所属专栏（若有）
- 分页（每页 10 条）

**专栏视图：**
- 专栏卡片列表：封面图、专栏名、简介、文章数、最近更新时间
- 点击专栏卡片 → 进入 `column.html`（专栏详情页）

### column.html — 专栏详情页
- 专栏封面 + 名称 + 介绍
- 有序文章列表（第1篇→第N篇），含每篇的标题、字数、浏览人次
- 当前阅读进度提示（若已读过某篇，可用 localStorage 标记）
- 上一篇 / 下一篇在专栏内导航（`post.html` 内也有）

### post.html — 单篇文章
- 文章元信息：标题、日期、标签、**字数 / 浏览人次 / 预计阅读时间**
- 右侧悬浮目录（从 H2/H3 自动生成，高亮当前章节）
- Markdown 正文（marked.js 渲染 + highlight.js 代码高亮）
- 上一篇 / 下一篇导航（全局时间序；若属于专栏，优先显示专栏内上下篇）
- 所属专栏入口（若有），点击可跳转至 `column.html`
- 页面加载时自动向 `POST /api/track/post/:id` 上报浏览

### works.html — 作品集
Tab 切换三个子板块：
1. **项目**：卡片网格，每张卡片含项目名、描述、技术栈徽章、GitHub 链接、Star 数
2. **系统演示**：iframe 嵌入或视频播放，附项目说明
3. **自定义分区**：后台可配置的自由内容区

### garden.html — 知识花园
- **子板块导航**：收藏 / 我在哪里 / 自定义（由 `garden_sections` 表驱动，支持增删和排序，Phase 2 实现拖拽，Phase 1 后台手动排序）
- **收藏子板块**：外链卡片，按类型打标签（文章 / 视频 / 书籍 / 其他），卡片含标题、来源、摘要、原文链接、可选 iframe 预览
- **我在哪里子板块**：第三方账号入口卡片，含平台图标、昵称、链接、二维码（规划中状态可标灰）
- 感想文章复用博客系统，在 garden 标签下聚合显示

### admin/ — CMS 后台
- `index.html`：用户名/密码登录，获取 JWT，存 localStorage
- `dashboard.html`：总 PV/UV（今日/7日/30日）折线图 + 热门文章 Top 10 + GA4 入口链接
- `posts.html`：文章列表，支持新建/编辑/删除/置顶切换/发布切换
- `editor.html`：左右分栏 Markdown 编辑器（左：原文，右：实时预览），支持图片上传

---

## API 设计

### 公开接口
```
GET  /api/posts              # 文章列表（?page=&tag=&q=&pinned=&column_id=）
GET  /api/posts/:slug        # 单篇文章详情（含渲染 HTML，含所属专栏上下篇）
GET  /api/tags               # 标签列表
GET  /api/columns            # 专栏列表
GET  /api/columns/:slug      # 专栏详情（含有序文章列表）
GET  /api/projects           # 项目列表
GET  /api/garden/sections    # 知识花园子板块列表
GET  /api/garden/items       # 收藏列表（?section_id=&type=）
POST /api/track/page         # 页面浏览埋点
POST /api/track/post/:id     # 文章浏览埋点
```

### 鉴权接口（JWT）
```
POST /api/auth/login

GET    /api/admin/stats          # 统计数据（PV/UV/热门文章）
POST   /api/admin/posts          # 新建文章
PUT    /api/admin/posts/:id      # 编辑文章
DELETE /api/admin/posts/:id      # 删除文章
PUT    /api/admin/posts/:id/pin  # 切换置顶
POST   /api/admin/columns          # 新建专栏
PUT    /api/admin/columns/:id      # 编辑专栏
DELETE /api/admin/columns/:id      # 删除专栏
PUT    /api/admin/columns/:id/posts # 更新专栏文章列表及顺序
POST   /api/admin/projects       # 新建项目
PUT    /api/admin/projects/:id   # 编辑项目
DELETE /api/admin/projects/:id   # 删除项目
POST   /api/admin/garden/sections        # 新建子板块
PUT    /api/admin/garden/sections/:id    # 编辑/排序子板块
DELETE /api/admin/garden/sections/:id   # 删除子板块
POST   /api/admin/garden/items          # 新建收藏
PUT    /api/admin/garden/items/:id      # 编辑收藏
DELETE /api/admin/garden/items/:id      # 删除收藏
POST   /api/admin/upload         # 图片上传
```

---

## 数据库 Schema（SQLite）

```sql
-- 文章
CREATE TABLE posts (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  title       TEXT NOT NULL,
  slug        TEXT NOT NULL UNIQUE,
  content     TEXT NOT NULL,          -- 原始 Markdown
  summary     TEXT,
  word_count  INTEGER DEFAULT 0,      -- 写入时由 goldmark 计算
  tags        TEXT DEFAULT '[]',      -- JSON 数组
  is_pinned   BOOLEAN DEFAULT FALSE,
  is_published BOOLEAN DEFAULT FALSE,
  created_at  DATETIME,
  updated_at  DATETIME
);

-- 专栏
CREATE TABLE columns (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  name        TEXT NOT NULL,
  slug        TEXT NOT NULL UNIQUE,
  cover_url   TEXT,
  description TEXT,
  is_published BOOLEAN DEFAULT FALSE,
  created_at  DATETIME,
  updated_at  DATETIME
);

-- 专栏文章关联（有序）
CREATE TABLE column_posts (
  column_id   INTEGER REFERENCES columns(id),
  post_id     INTEGER REFERENCES posts(id),
  sort_order  INTEGER NOT NULL,           -- 文章在专栏内的顺序
  PRIMARY KEY (column_id, post_id)
);

-- 项目
CREATE TABLE projects (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  name        TEXT NOT NULL,
  description TEXT,
  repo_url    TEXT,
  demo_url    TEXT,
  tech_stack  TEXT DEFAULT '[]',      -- JSON 数组
  stars       INTEGER DEFAULT 0,
  is_featured BOOLEAN DEFAULT FALSE,
  sort_order  INTEGER DEFAULT 0,      -- Phase 2 拖拽排序预留
  created_at  DATETIME
);

-- 知识花园子板块
CREATE TABLE garden_sections (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  name        TEXT NOT NULL,
  slug        TEXT NOT NULL UNIQUE,
  sort_order  INTEGER DEFAULT 0,      -- Phase 2 拖拽排序预留
  is_visible  BOOLEAN DEFAULT TRUE,
  created_at  DATETIME
);

-- 知识花园收藏条目
CREATE TABLE garden_items (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  section_id  INTEGER REFERENCES garden_sections(id),
  type        TEXT NOT NULL,          -- article / video / book / account / other
  title       TEXT NOT NULL,
  url         TEXT,
  source      TEXT,                   -- 来源名称
  summary     TEXT,
  embed_url   TEXT,                   -- iframe 嵌入 URL
  sort_order  INTEGER DEFAULT 0,      -- Phase 2 拖拽排序预留
  created_at  DATETIME
);

-- 文章浏览记录
CREATE TABLE post_views (
  id           INTEGER PRIMARY KEY AUTOINCREMENT,
  post_id      INTEGER REFERENCES posts(id),
  visitor_hash TEXT NOT NULL,         -- SHA256(IP + UA)，不存原始 IP
  created_at   DATETIME
);

-- 页面浏览记录
CREATE TABLE page_views (
  id           INTEGER PRIMARY KEY AUTOINCREMENT,
  path         TEXT NOT NULL,
  visitor_hash TEXT NOT NULL,
  referrer     TEXT,
  created_at   DATETIME
);

-- 每日统计聚合缓存（定时任务每日凌晨聚合，加速查询）
CREATE TABLE daily_stats (
  id    INTEGER PRIMARY KEY AUTOINCREMENT,
  date  TEXT NOT NULL,                -- YYYY-MM-DD
  path  TEXT NOT NULL,
  pv    INTEGER DEFAULT 0,
  uv    INTEGER DEFAULT 0,
  UNIQUE(date, path)
);
```

---

## 统计系统详细设计

### 采集架构
双轨并行，互为补充：

| 维度 | 自建埋点 | Google Analytics 4 |
|------|---------|-------------------|
| 数据归属 | 自有，存 SQLite | Google 平台 |
| 文章浏览人次 | ✓ 实时展示在文章页 | — |
| PV/UV 趋势 | ✓ CMS 后台可见 | ✓ 更丰富维度 |
| 地域/设备 | — | ✓ |
| 来源/搜索词 | 基础 referrer | ✓ |
| 隐私 | visitor_hash 匿名化 | 受 GA4 隐私政策约束 |

### 埋点实现
- 每个页面底部引入 `js/tracker.js`，页面加载后 `POST /api/track/page`（path + referrer）
- 文章页额外 `POST /api/track/post/:id`
- `visitor_hash = SHA256(ClientIP + UserAgent)`，不存储原始 IP
- 同一 visitor_hash 对同一文章，24 小时内只计 1 次 UV

### 文章页三指标计算
- **字数**：文章保存时由 goldmark 解析 Markdown，统计中文字符 + 英文单词数，存入 `posts.word_count`
- **浏览人次**：`SELECT COUNT(DISTINCT visitor_hash) FROM post_views WHERE post_id = ?`
- **预计阅读时间**：`ceil(word_count / 300)` 分钟（中文阅读速度约 300字/分钟）

### CMS 统计面板数据接口
`GET /api/admin/stats?range=7d` 返回：
```json
{
  "total_pv": 12345,
  "total_uv": 3456,
  "daily": [{"date": "2026-05-01", "pv": 234, "uv": 89}, ...],
  "top_posts": [{"title": "...", "views": 1234}, ...]
}
```

---

## 双主题切换

```css
/* theme-dark.css */
:root[data-theme="dark"] {
  --bg-primary: #0a0a0a;
  --bg-secondary: #0f172a;
  --text-primary: #e2e8f0;
  --text-secondary: #94a3b8;
  --accent: #00ff88;
  --border: #1e293b;
}

/* theme-light.css */
:root[data-theme="light"] {
  --bg-primary: #ffffff;
  --bg-secondary: #f8fafc;
  --text-primary: #0f172a;
  --text-secondary: #475569;
  --accent: #2563eb;
  --border: #e2e8f0;
}
```

切换逻辑（`theme.js`）：
1. 初始化：读 `localStorage.getItem('theme')`，若无则取 `prefers-color-scheme`
2. 设置 `document.documentElement.setAttribute('data-theme', theme)`
3. 切换时更新 localStorage + 重设属性
4. 过渡动画：`* { transition: background-color 0.3s, color 0.3s, border-color 0.3s }`

---

## Docker 部署

### docker-compose.yml
```yaml
services:
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./frontend:/usr/share/nginx/html:ro
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
      - ./nginx/ssl:/etc/nginx/ssl:ro
    depends_on:
      - api

  api:
    build: ./backend
    expose:
      - "8080"
    volumes:
      - ./data:/data
    environment:
      - DB_PATH=/data/blog.db
      - JWT_SECRET=${JWT_SECRET}
      - GA_MEASUREMENT_ID=${GA_MEASUREMENT_ID}
    restart: unless-stopped

volumes:
  # data/ 目录挂载，SQLite 文件持久化
```

### Nginx 关键配置
```nginx
location /api/ {
    proxy_pass http://api:8080;
    proxy_set_header X-Real-IP $remote_addr;
}

location / {
    root /usr/share/nginx/html;
    try_files $uri $uri.html $uri/ =404;
}
```

---

## Phase 2 预留设计

以下字段和约定在 Phase 1 写入，为 Phase 2 可视化构建器打基础：

- `projects.sort_order`、`garden_sections.sort_order`、`garden_items.sort_order`：拖拽排序直接更新此字段
- 前端组件 CSS class 统一命名规范：`.card`、`.section`、`.grid-2`、`.grid-3` 等，便于后续抽象为可配置组件
- 后端 API 全 RESTful，构建器直接复用，无需新增接口

---

## 不在 Phase 1 范围内

- 评论系统
- 邮件订阅
- 全文搜索（Elasticsearch）
- 可视化页面构建器（Phase 2 单独立项）
- 多用户 / 权限系统
