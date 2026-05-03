# 03 后端

## 项目结构

```
backend/
├── cmd/server/main.go      # 程序入口
└── internal/
    ├── config/config.go    # 环境变量
    ├── db/db.go            # 数据库连接
    ├── model/              # 数据模型
    ├── service/            # 业务逻辑
    ├── handler/            # HTTP 处理器
    ├── middleware/         # 中间件
    └── router/router.go    # 路由
```

### 分层职责

```
请求 → Handler（解析参数）→ Service（业务逻辑）→ Model（数据库操作）
                                    ↓
                               返回 JSON
```

- **Handler**：只负责解析 HTTP 请求参数、调用 Service、返回 JSON，不包含业务逻辑
- **Service**：所有业务逻辑在这里，比如"字数统计"、"24小时内同一访客只计一次 UV"
- **Model**：定义数据库表结构，不包含逻辑

## 数据库表

### posts — 文章

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | INTEGER | 主键 |
| `title` | TEXT | 标题 |
| `slug` | TEXT | URL 路径，唯一 |
| `content` | TEXT | Markdown 原文（不序列化到公开 JSON） |
| `summary` | TEXT | 摘要 |
| `word_count` | INTEGER | 字数，保存时自动计算 |
| `tags` | TEXT | JSON 数组，如 `["Go","并发"]` |
| `is_pinned` | BOOLEAN | 是否置顶 |
| `is_published` | BOOLEAN | 是否已发布（草稿不出现在前台） |

### columns / column_posts — 专栏

专栏是有序的文章集合。`column_posts` 是关联表，记录文章在专栏内的顺序：

```
columns: id, name, slug, description, is_published
column_posts: column_id, post_id, sort_order  ← sort_order 决定顺序
```

### garden_sections / garden_items — 知识花园

```
garden_sections: id, name, slug, page, sort_order, is_visible
  page 字段区分: 'garden'（知识花园）或 'works'（作品集自定义板块）

garden_items: id, section_id, type, title, url, source, summary, tags, embed_url
  type: article / video / book / website / account
  tags: JSON 数组，支持多标签筛选
```

### profile — 主页个人信息

只有一行数据（id=1），保存时用 upsert（有则更新，无则插入）：

```
profile: id=1, name, subtitle, bio, skills(JSON), social(JSON), timeline(JSON)
```

### post_views / page_views / daily_stats — 统计

```
post_views: post_id, visitor_hash, created_at
  visitor_hash = SHA256(IP + UserAgent) 前32位，不存原始 IP

page_views: path, visitor_hash, referrer, created_at

daily_stats: date, path, pv, uv  ← 聚合缓存，加速查询
```

## API 路由

### 公开接口（无需鉴权）

```
GET  /api/profile              # 主页个人信息
GET  /api/posts                # 文章列表（?page=&tag=&q=）
GET  /api/posts/:slug          # 单篇文章（含 prev/next 相邻篇）
GET  /api/tags                 # 所有标签
GET  /api/columns              # 专栏列表（仅已发布）
GET  /api/columns/:slug        # 专栏详情 + 有序文章列表
GET  /api/projects             # 项目列表
GET  /api/garden/sections      # 知识花园板块
GET  /api/works/sections       # 作品集自定义板块
GET  /api/garden/items         # 板块条目（?section_id=&type=&tag=）
POST /api/track/page           # 页面浏览埋点
POST /api/track/post/:id       # 文章浏览埋点
```

### 管理接口（需要 JWT）

所有 `/api/admin/` 开头的路由都需要在 Header 里携带：

```
Authorization: Bearer <token>
```

token 通过 `POST /api/auth/login` 获取，有效期 7 天。

## Markdown 处理

文章保存时，后端用 `goldmark` 库做两件事：

1. **字数统计**：去除 Markdown 语法符号后，汉字逐个计数，英文按单词计数
2. **HTML 渲染**：公开接口返回 `content_html` 字段，前端直接插入 DOM

```go
// 字数统计：中文字符 + 英文单词
func WordCount(source string) int {
    plain := markdownSymbols.ReplaceAllString(source, " ")
    // 遍历字符，汉字 +1，英文单词 +1
}

// 预计阅读时间
func ReadingMinutes(wordCount int) int {
    return max(1, wordCount/300)  // 中文阅读速度约 300字/分钟
}
```

## JWT 鉴权流程

```
1. POST /api/auth/login { username, password }
        ↓ 验证通过
2. 返回 { token: "eyJ..." }  (有效期7天)
        ↓
3. 后续管理请求 Header 带 Authorization: Bearer eyJ...
        ↓ JWTAuth 中间件验证
4. 通过 → 执行 Handler；失败 → 返回 401
```

前端把 token 存在 `localStorage('admin-token')` 里，`api.js` 的 `request()` 函数自动从 localStorage 读取并附加到请求头。

## 如何新增一个 API

以新增"友情链接"功能为例：

1. `model/link.go` — 定义 `Link` 结构体（含 GORM 标签）
2. `db/db.go` — 在 `AutoMigrate` 里加 `&model.Link{}`
3. `service/link.go` — 写 `List()`、`Create()`、`Delete()` 方法
4. `handler/link.go` — 写 `List`、`AdminCreate`、`AdminDelete` Handler
5. `router/router.go` — 注册路由
6. `frontend/js/api.js` — 加对应的前端调用方法

每层职责清晰，照着现有代码依葫芦画瓢即可。
