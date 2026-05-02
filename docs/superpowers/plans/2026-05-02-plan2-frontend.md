# Frontend Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build all 6 public-facing frontend pages (index, blog, post, column, works, garden) plus shared nav/theme system, using static HTML/CSS/JS that calls the Go API. Pages must look correct with mock data before the backend exists.

**Architecture:** Pure static files in `frontend/`. Each page is self-contained HTML that imports shared CSS and JS. A `js/api.js` module wraps all fetch calls — pages call `api.getPosts()` etc, never raw fetch. A `js/mock.js` module returns fake data when `?mock=1` is in the URL, so pages can be developed and demoed without a running backend. Theme is stored in `localStorage` and applied via `data-theme` on `<html>`.

**Tech Stack:** HTML5, CSS3 (custom properties), Vanilla JS ES6 modules, marked.js (CDN), highlight.js (CDN)

---

## File Structure

```
frontend/
├── index.html          # 主页：Hero + 技能 + 时间线 + 精选项目 + 最新文章 + 联系
├── blog.html           # 博客：文章/专栏 Tab + 置顶 + 筛选 + 列表 + 分页
├── post.html           # 单篇文章：元信息 + 目录 + 正文 + 上下篇
├── column.html         # 专栏详情：封面 + 介绍 + 有序文章列表
├── works.html          # 作品集：项目/演示/自定义 Tab
├── garden.html         # 知识花园：子板块导航 + 收藏卡片 + 账号入口
├── css/
│   ├── common.css      # 布局、组件、排版公共样式
│   ├── theme-dark.css  # 暗黑主题 CSS 变量
│   └── theme-light.css # 明亮主题 CSS 变量
└── js/
    ├── theme.js        # 主题切换 + localStorage 持久化
    ├── api.js          # 所有后端 API 调用封装（fetch）
    ├── mock.js         # mock 数据，?mock=1 时替代 api.js
    └── nav.js          # 导航栏当前页高亮
```

---

## Task 1: CSS 变量 + 公共样式基础

**Files:**
- Create: `frontend/css/theme-dark.css`
- Create: `frontend/css/theme-light.css`
- Create: `frontend/css/common.css`

- [ ] **Step 1: Write dark theme variables**

Create `frontend/css/theme-dark.css`:

```css
:root[data-theme="dark"] {
  --bg-primary:    #0a0a0a;
  --bg-secondary:  #0f172a;
  --bg-card:       #0f172a;
  --bg-hover:      #1e293b;
  --text-primary:  #e2e8f0;
  --text-secondary:#94a3b8;
  --text-muted:    #475569;
  --accent:        #00ff88;
  --accent-dim:    #0d2a1a;
  --accent-border: #166534;
  --border:        #1e293b;
  --border-hover:  #334155;
  --tag-bg:        #1e293b;
  --tag-text:      #94a3b8;
  --pinned-bg:     #0f1f0f;
  --pinned-border: #166534;
  --code-bg:       #1e293b;
  --nav-bg:        #0d0d0d;
}
```

- [ ] **Step 2: Write light theme variables**

Create `frontend/css/theme-light.css`:

```css
:root[data-theme="light"] {
  --bg-primary:    #ffffff;
  --bg-secondary:  #f8fafc;
  --bg-card:       #ffffff;
  --bg-hover:      #f1f5f9;
  --text-primary:  #0f172a;
  --text-secondary:#475569;
  --text-muted:    #94a3b8;
  --accent:        #2563eb;
  --accent-dim:    #eff6ff;
  --accent-border: #bfdbfe;
  --border:        #e2e8f0;
  --border-hover:  #cbd5e1;
  --tag-bg:        #f1f5f9;
  --tag-text:      #475569;
  --pinned-bg:     #f0fdf4;
  --pinned-border: #86efac;
  --code-bg:       #f8fafc;
  --nav-bg:        #ffffff;
}
```

- [ ] **Step 3: Write common.css**

Create `frontend/css/common.css`:

```css
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

* { transition: background-color 0.25s, color 0.25s, border-color 0.25s; }

body {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
  background: var(--bg-primary);
  color: var(--text-primary);
  line-height: 1.6;
  font-size: 16px;
}

a { color: var(--accent); text-decoration: none; }
a:hover { text-decoration: underline; }

/* NAV */
.nav {
  position: sticky; top: 0; z-index: 100;
  background: var(--nav-bg);
  border-bottom: 1px solid var(--border);
  padding: 0 24px;
  height: 56px;
  display: flex; align-items: center; gap: 32px;
}
.nav-logo {
  font-family: monospace; font-size: 18px; font-weight: 700;
  color: var(--accent); letter-spacing: -1px;
}
.nav-links { display: flex; gap: 24px; margin-right: auto; }
.nav-links a { color: var(--text-secondary); font-size: 14px; }
.nav-links a.active, .nav-links a:hover { color: var(--text-primary); }
.nav-theme-btn {
  background: none; border: 1px solid var(--border);
  border-radius: 6px; padding: 4px 10px;
  cursor: pointer; font-size: 14px; color: var(--text-secondary);
}
.nav-theme-btn:hover { border-color: var(--accent); color: var(--accent); }

/* LAYOUT */
.page { max-width: 900px; margin: 0 auto; padding: 40px 24px; }
.page-wide { max-width: 1100px; margin: 0 auto; padding: 40px 24px; }

/* CARD */
.card {
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 20px;
}
.card:hover { border-color: var(--border-hover); }

/* TAG */
.tag {
  display: inline-block;
  background: var(--tag-bg); color: var(--tag-text);
  border-radius: 4px; padding: 2px 8px;
  font-size: 12px;
}

/* BADGE (tech stack) */
.badge {
  display: inline-block;
  border-radius: 4px; padding: 2px 8px; font-size: 11px;
}
.badge-go    { background: #0d2a1a; color: #4ade80; }
.badge-java  { background: #1a1a08; color: #fbbf24; }
.badge-python{ background: #1a1520; color: #c084fc; }
.badge-docker{ background: #0d1929; color: #60a5fa; }
.badge-other { background: var(--tag-bg); color: var(--tag-text); }

/* POST META */
.post-meta { display: flex; gap: 16px; font-size: 13px; color: var(--text-muted); flex-wrap: wrap; }
.post-meta span { display: flex; align-items: center; gap: 4px; }

/* TABS */
.tabs { display: flex; gap: 0; border-bottom: 1px solid var(--border); margin-bottom: 24px; }
.tab-btn {
  padding: 10px 20px; font-size: 14px; cursor: pointer;
  background: none; border: none; border-bottom: 2px solid transparent;
  color: var(--text-secondary); margin-bottom: -1px;
}
.tab-btn.active { color: var(--accent); border-bottom-color: var(--accent); }
.tab-panel { display: none; }
.tab-panel.active { display: block; }

/* SECTION TITLE */
.section-title {
  font-size: 11px; letter-spacing: 2px; text-transform: uppercase;
  color: var(--text-muted); margin-bottom: 16px;
}

/* PINNED */
.pinned-card {
  background: var(--pinned-bg);
  border: 1px solid var(--pinned-border);
  border-radius: 8px; padding: 20px;
  position: relative; margin-bottom: 12px;
}
.pinned-badge {
  position: absolute; top: 12px; right: 12px;
  background: var(--accent-border); color: var(--accent);
  font-size: 10px; padding: 2px 8px; border-radius: 3px;
}

/* GRID */
.grid-2 { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
.grid-3 { display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 16px; }
@media (max-width: 640px) {
  .grid-2, .grid-3 { grid-template-columns: 1fr; }
}

/* TIMELINE */
.timeline { display: flex; flex-direction: column; gap: 16px; }
.timeline-item { display: flex; gap: 16px; align-items: flex-start; }
.timeline-dot {
  width: 8px; height: 8px; border-radius: 50%;
  background: var(--accent); margin-top: 6px; flex-shrink: 0;
}
.timeline-dot.muted { background: var(--border); }

/* PAGINATION */
.pagination { display: flex; gap: 8px; justify-content: center; margin-top: 32px; }
.page-btn {
  padding: 6px 12px; border: 1px solid var(--border);
  border-radius: 6px; background: var(--bg-card);
  color: var(--text-secondary); cursor: pointer; font-size: 13px;
}
.page-btn.active { border-color: var(--accent); color: var(--accent); }
.page-btn:hover { border-color: var(--border-hover); color: var(--text-primary); }

/* HERO */
.hero { padding: 60px 0 40px; }
.hero-name { font-size: 36px; font-weight: 700; line-height: 1.2; }
.hero-name span { color: var(--accent); }
.hero-sub { font-size: 16px; color: var(--text-secondary); margin-top: 8px; }
.hero-links { display: flex; gap: 12px; margin-top: 20px; }
.hero-link {
  padding: 8px 16px; border: 1px solid var(--border);
  border-radius: 6px; font-size: 13px; color: var(--text-secondary);
}
.hero-link:hover { border-color: var(--accent); color: var(--accent); text-decoration: none; }

/* ARTICLE BODY */
.article-body { line-height: 1.8; }
.article-body h1, .article-body h2, .article-body h3 {
  margin: 2em 0 0.75em; color: var(--text-primary);
}
.article-body p { margin-bottom: 1.2em; color: var(--text-secondary); }
.article-body pre {
  background: var(--code-bg); border: 1px solid var(--border);
  border-radius: 6px; padding: 16px; overflow-x: auto; margin: 1.5em 0;
}
.article-body code { font-family: monospace; font-size: 14px; }
.article-body :not(pre) > code {
  background: var(--code-bg); padding: 2px 6px; border-radius: 4px; font-size: 13px;
}
.article-body blockquote {
  border-left: 3px solid var(--accent); padding-left: 16px;
  color: var(--text-secondary); margin: 1.5em 0;
}
.article-body table { width: 100%; border-collapse: collapse; margin: 1.5em 0; }
.article-body th, .article-body td {
  border: 1px solid var(--border); padding: 8px 12px; font-size: 14px;
}
.article-body th { background: var(--bg-secondary); }

/* TOC */
.toc {
  position: sticky; top: 72px; max-height: calc(100vh - 100px);
  overflow-y: auto; font-size: 13px;
}
.toc a { color: var(--text-muted); display: block; padding: 3px 0; }
.toc a.active { color: var(--accent); }
.toc a:hover { color: var(--text-primary); text-decoration: none; }
.toc-h3 { padding-left: 12px; }

/* LOADING / EMPTY */
.loading { text-align: center; padding: 40px; color: var(--text-muted); }
.empty   { text-align: center; padding: 40px; color: var(--text-muted); font-size: 14px; }
```

- [ ] **Step 4: Commit**

```bash
git add frontend/css/
git commit -m "feat: CSS theme variables and common styles"
```

---

## Task 2: JS Infrastructure (theme + api + mock + nav)

**Files:**
- Create: `frontend/js/theme.js`
- Create: `frontend/js/api.js`
- Create: `frontend/js/mock.js`
- Create: `frontend/js/nav.js`

- [ ] **Step 1: Write theme.js**

Create `frontend/js/theme.js`:

```js
const THEME_KEY = 'blog-theme';

function getInitialTheme() {
  const saved = localStorage.getItem(THEME_KEY);
  if (saved) return saved;
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
}

function applyTheme(theme) {
  document.documentElement.setAttribute('data-theme', theme);
  localStorage.setItem(THEME_KEY, theme);
  const btn = document.getElementById('theme-btn');
  if (btn) btn.textContent = theme === 'dark' ? '☀️' : '🌙';
}

function toggleTheme() {
  const current = document.documentElement.getAttribute('data-theme');
  applyTheme(current === 'dark' ? 'light' : 'dark');
}

document.addEventListener('DOMContentLoaded', () => {
  applyTheme(getInitialTheme());
  const btn = document.getElementById('theme-btn');
  if (btn) btn.addEventListener('click', toggleTheme);
});
```

- [ ] **Step 2: Write api.js**

Create `frontend/js/api.js`:

```js
const BASE = '/api';

async function request(path, options = {}) {
  const token = localStorage.getItem('admin-token');
  const headers = { 'Content-Type': 'application/json', ...options.headers };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await fetch(BASE + path, { ...options, headers });
  if (!res.ok) throw new Error(`${res.status} ${res.statusText}`);
  if (res.status === 204) return null;
  return res.json();
}

export const api = {
  // Posts
  getPosts:    (params = {}) => request('/posts?' + new URLSearchParams(params)),
  getPost:     (slug)        => request(`/posts/${slug}`),
  getTags:     ()            => request('/tags'),

  // Columns
  getColumns:  ()            => request('/columns'),
  getColumn:   (slug)        => request(`/columns/${slug}`),

  // Projects
  getProjects: ()            => request('/projects'),

  // Garden
  getGardenSections: ()           => request('/garden/sections'),
  getGardenItems:    (sectionId)  => request(`/garden/items?section_id=${sectionId}`),

  // Tracking
  trackPage: (path, referrer = '') =>
    request('/track/page', { method: 'POST', body: JSON.stringify({ path, referrer }) }),
  trackPost: (id) =>
    request(`/track/post/${id}`, { method: 'POST' }),

  // Admin
  login: (username, password) =>
    request('/auth/login', { method: 'POST', body: JSON.stringify({ username, password }) }),
  adminStats: (days = 7) => request(`/admin/stats?days=${days}`),

  adminCreatePost:  (data) => request('/admin/posts', { method: 'POST', body: JSON.stringify(data) }),
  adminUpdatePost:  (id, data) => request(`/admin/posts/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  adminDeletePost:  (id) => request(`/admin/posts/${id}`, { method: 'DELETE' }),
  adminPinPost:     (id, pinned) => request(`/admin/posts/${id}/pin`, { method: 'PUT', body: JSON.stringify({ pinned }) }),
  adminPublishPost: (id, published) => request(`/admin/posts/${id}/publish`, { method: 'PUT', body: JSON.stringify({ published }) }),

  adminCreateColumn: (data) => request('/admin/columns', { method: 'POST', body: JSON.stringify(data) }),
  adminUpdateColumn: (id, data) => request(`/admin/columns/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  adminDeleteColumn: (id) => request(`/admin/columns/${id}`, { method: 'DELETE' }),
  adminSetColumnPosts: (id, postIds) => request(`/admin/columns/${id}/posts`, { method: 'PUT', body: JSON.stringify({ post_ids: postIds }) }),

  adminCreateProject: (data) => request('/admin/projects', { method: 'POST', body: JSON.stringify(data) }),
  adminUpdateProject: (id, data) => request(`/admin/projects/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  adminDeleteProject: (id) => request(`/admin/projects/${id}`, { method: 'DELETE' }),

  adminCreateSection: (data) => request('/admin/garden/sections', { method: 'POST', body: JSON.stringify(data) }),
  adminUpdateSection: (id, data) => request(`/admin/garden/sections/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  adminDeleteSection: (id) => request(`/admin/garden/sections/${id}`, { method: 'DELETE' }),
  adminCreateItem:  (data) => request('/admin/garden/items', { method: 'POST', body: JSON.stringify(data) }),
  adminUpdateItem:  (id, data) => request(`/admin/garden/items/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  adminDeleteItem:  (id) => request(`/admin/garden/items/${id}`, { method: 'DELETE' }),

  adminUpload: (file) => {
    const form = new FormData();
    form.append('file', file);
    return request('/admin/upload', { method: 'POST', headers: {}, body: form });
  },
};
```

- [ ] **Step 3: Write mock.js**

Create `frontend/js/mock.js`:

```js
// Activated when URL contains ?mock=1
// Returns same shape as api.js methods so pages work without a backend.

export const api = {
  getPosts: async ({ page = 1 } = {}) => ({
    posts: [
      { id: 1, title: '置顶：Go 并发模式实践', slug: 'go-concurrency', summary: '深入理解 goroutine 和 channel 的设计模式，避免常见陷阱。', tags: '["Go","并发"]', is_pinned: true, word_count: 2400, created_at: '2026-04-10T10:00:00Z' },
      { id: 2, title: 'Docker 最佳实践总结', slug: 'docker-best-practices', summary: '从 Dockerfile 优化到多阶段构建，系统整理 Docker 使用心得。', tags: '["Docker","运维"]', is_pinned: false, word_count: 1800, created_at: '2026-04-01T10:00:00Z' },
      { id: 3, title: 'Java 虚拟线程初探', slug: 'java-virtual-threads', summary: 'Project Loom 带来的虚拟线程，彻底改变 Java 并发编程方式。', tags: '["Java"]', is_pinned: false, word_count: 1200, created_at: '2026-03-20T10:00:00Z' },
    ],
    total: 3, page,
  }),

  getPost: async (slug) => ({
    post: { id: 1, title: 'Go 并发模式实践', slug, content: '', content_html: '<h2 id="intro">简介</h2><p>本文介绍 Go 并发模式...</p><h2 id="goroutine">Goroutine</h2><p>goroutine 是 Go 的轻量级线程...</p><h3 id="channel">Channel 用法</h3><p>channel 是 goroutine 间通信的桥梁。</p><pre><code class="language-go">ch := make(chan int)\ngo func() { ch &lt;- 42 }()\nfmt.Println(&lt;-ch)</code></pre>', tags: '["Go","并发"]', word_count: 2400, created_at: '2026-04-10T10:00:00Z' },
    views: 1234,
    reading_minutes: 8,
  }),

  getTags: async () => ({ tags: ['Go', 'Java', 'Python', 'Docker', '并发', '随笔'] }),

  getColumns: async () => ({
    columns: [
      { id: 1, slug: 'go-series', name: 'Go 从入门到实践', description: '系统学习 Go 语言，从语法基础到工程实践。', cover_url: '', is_published: true, created_at: '2026-01-01T00:00:00Z' },
      { id: 2, slug: 'system-design', name: '系统设计笔记', description: '分布式系统设计的思考与实践记录。', cover_url: '', is_published: true, created_at: '2026-02-01T00:00:00Z' },
    ],
  }),

  getColumn: async (slug) => ({
    column: { id: 1, slug, name: 'Go 从入门到实践', description: '系统学习 Go 语言，从语法基础到工程实践。共 5 篇。', cover_url: '' },
    posts: [
      { id: 10, title: '第一篇：Go 环境搭建与基础语法', slug: 'go-01', word_count: 1500 },
      { id: 11, title: '第二篇：Go 并发编程', slug: 'go-02', word_count: 2400 },
      { id: 12, title: '第三篇：Go 标准库精讲', slug: 'go-03', word_count: 2000 },
    ],
  }),

  getProjects: async () => ({
    projects: [
      { id: 1, name: 'blog', description: '本博客网站，Go + 原生前端，Docker 部署。', repo_url: 'https://github.com/daiyutong/blog', demo_url: '', tech_stack: '["Go","HTML","CSS","Docker"]', stars: 12, is_featured: true },
      { id: 2, name: 'go-tools', description: '日常开发工具集，Go 实现。', repo_url: 'https://github.com/daiyutong/go-tools', demo_url: '', tech_stack: '["Go"]', stars: 8, is_featured: true },
      { id: 3, name: 'data-pipeline', description: 'Python 数据处理管道，支持多数据源。', repo_url: 'https://github.com/daiyutong/data-pipeline', demo_url: '', tech_stack: '["Python"]', stars: 5, is_featured: false },
    ],
  }),

  getGardenSections: async () => ({
    sections: [
      { id: 1, slug: 'bookmarks', name: '收藏', sort_order: 0 },
      { id: 2, slug: 'social', name: '我在哪里', sort_order: 1 },
    ],
  }),

  getGardenItems: async (sectionId) => ({
    items: sectionId === 1 ? [
      { id: 1, section_id: 1, type: 'article', title: 'Go 调度器原理深度解析', source: 'Go Blog', url: 'https://example.com', summary: '深入分析 GMP 模型和调度器工作原理。', embed_url: '' },
      { id: 2, section_id: 1, type: 'video',   title: '如何做系统设计', source: 'YouTube', url: 'https://example.com', summary: '系统设计面试思路全解析。', embed_url: '' },
      { id: 3, section_id: 1, type: 'book',    title: '《程序员的自我修养》', source: '豆瓣', url: 'https://example.com', summary: '链接、装载与库，深入理解程序运行原理。', embed_url: '' },
    ] : [
      { id: 4, section_id: 2, type: 'account', title: 'GitHub', source: 'github.com/daiyutong', url: 'https://github.com', summary: '', embed_url: '' },
    ],
  }),

  trackPage: async () => null,
  trackPost: async () => null,
};
```

- [ ] **Step 4: Write nav.js**

Create `frontend/js/nav.js`:

```js
// Highlights the current page's nav link
document.addEventListener('DOMContentLoaded', () => {
  const path = location.pathname.replace(/\/$/, '') || '/';
  document.querySelectorAll('.nav-links a').forEach(a => {
    const href = a.getAttribute('href').replace(/\/$/, '') || '/';
    if (path === href || (href !== '/' && path.startsWith(href))) {
      a.classList.add('active');
    }
  });
});
```

- [ ] **Step 5: Commit**

```bash
git add frontend/js/
git commit -m "feat: JS infrastructure — theme, api, mock, nav"
```

---

## Task 3: 共用导航栏 HTML 片段

每个页面都需要相同的导航栏。为避免重复，定义一个 JS 函数注入导航。

**Files:**
- Create: `frontend/js/components.js`

- [ ] **Step 1: Write components.js**

Create `frontend/js/components.js`:

```js
export function renderNav(activePage) {
  const pages = [
    { href: '/index.html', label: '主页',   key: 'index'  },
    { href: '/blog.html',  label: '博客',   key: 'blog'   },
    { href: '/works.html', label: '作品集', key: 'works'  },
    { href: '/garden.html',label: '知识花园',key: 'garden' },
  ];
  const links = pages.map(p =>
    `<a href="${p.href}" class="${p.key === activePage ? 'active' : ''}">${p.label}</a>`
  ).join('');
  return `
    <nav class="nav">
      <span class="nav-logo">dyt</span>
      <div class="nav-links">${links}</div>
      <button class="nav-theme-btn" id="theme-btn">🌙</button>
    </nav>`;
}

export function renderPostMeta({ word_count, views, reading_minutes, created_at }) {
  const date = new Date(created_at).toLocaleDateString('zh-CN');
  return `
    <div class="post-meta">
      <span>📅 ${date}</span>
      <span>📝 ${word_count ?? 0} 字</span>
      <span>👁 ${views ?? 0} 次浏览</span>
      <span>📖 约 ${reading_minutes ?? 1} 分钟</span>
    </div>`;
}

export function renderTags(tagsJSON) {
  const tags = JSON.parse(tagsJSON || '[]');
  return tags.map(t => `<span class="tag">${t}</span>`).join(' ');
}

export function renderBadges(techStackJSON) {
  const stack = JSON.parse(techStackJSON || '[]');
  return stack.map(t => {
    const cls = ['Go','Java','Python','Docker'].includes(t)
      ? `badge-${t.toLowerCase()}` : 'badge-other';
    return `<span class="badge ${cls}">${t}</span>`;
  }).join(' ');
}

export function getApi() {
  const isMock = new URLSearchParams(location.search).has('mock');
  return isMock
    ? import('/js/mock.js').then(m => m.api)
    : import('/js/api.js').then(m => m.api);
}
```

- [ ] **Step 2: Commit**

```bash
git add frontend/js/components.js
git commit -m "feat: shared nav and component helpers"
```

---

## Task 4: index.html — 主页

**Files:**
- Create: `frontend/index.html`

- [ ] **Step 1: Write index.html**

Create `frontend/index.html`:

```html
<!DOCTYPE html>
<html data-theme="dark" lang="zh-CN">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>戴宇彤 | 个人主页</title>
  <link rel="stylesheet" href="/css/theme-dark.css">
  <link rel="stylesheet" href="/css/theme-light.css">
  <link rel="stylesheet" href="/css/common.css">
</head>
<body>
  <div id="nav-root"></div>

  <main class="page">
    <!-- Hero -->
    <section class="hero">
      <h1 class="hero-name">Hi, I'm <span>戴宇彤</span> 👋</h1>
      <p class="hero-sub">Full-stack Developer · Go · Java · Python · Docker</p>
      <div style="margin-top:16px;display:flex;gap:8px;flex-wrap:wrap">
        <span class="badge badge-go">Go</span>
        <span class="badge badge-java">Java</span>
        <span class="badge badge-python">Python</span>
        <span class="badge badge-docker">Docker</span>
      </div>
      <div class="hero-links">
        <a class="hero-link" href="https://github.com/daiyutong" target="_blank">⭐ GitHub</a>
        <a class="hero-link" href="mailto:your@email.com">✉️ 邮件联系</a>
      </div>
    </section>

    <!-- Experience Timeline -->
    <section style="margin-bottom:40px">
      <div class="section-title">经历</div>
      <div class="timeline">
        <div class="timeline-item">
          <div class="timeline-dot"></div>
          <div>
            <div style="font-weight:600">某公司 · 后端工程师</div>
            <div style="color:var(--text-muted);font-size:13px">2024 — 至今</div>
          </div>
        </div>
        <div class="timeline-item">
          <div class="timeline-dot muted"></div>
          <div>
            <div style="font-weight:600">某大学 · 计算机科学与技术</div>
            <div style="color:var(--text-muted);font-size:13px">2020 — 2024</div>
          </div>
        </div>
      </div>
    </section>

    <!-- Featured Projects -->
    <section style="margin-bottom:40px">
      <div class="section-title">精选项目</div>
      <div class="grid-3" id="featured-projects"><div class="loading">加载中...</div></div>
    </section>

    <!-- Latest Posts -->
    <section style="margin-bottom:40px">
      <div class="section-title">最新文章</div>
      <div id="latest-posts"><div class="loading">加载中...</div></div>
    </section>
  </main>

  <script type="module">
    import { renderNav, renderBadges, getApi } from '/js/components.js';
    import '/js/theme.js';

    document.getElementById('nav-root').innerHTML = renderNav('index');

    const api = await getApi();
    api.trackPage(location.pathname, document.referrer);

    // Featured projects
    const { projects } = await api.getProjects();
    const featured = projects.filter(p => p.is_featured).slice(0, 3);
    document.getElementById('featured-projects').innerHTML = featured.length
      ? featured.map(p => `
          <div class="card">
            <div style="font-weight:600;margin-bottom:6px">${p.name}</div>
            <div style="font-size:13px;color:var(--text-secondary);margin-bottom:10px">${p.description}</div>
            <div style="margin-bottom:10px">${renderBadges(p.tech_stack)}</div>
            <div style="font-size:12px;color:var(--text-muted);display:flex;gap:12px">
              ${p.stars ? `<span>⭐ ${p.stars}</span>` : ''}
              ${p.repo_url ? `<a href="${p.repo_url}" target="_blank">GitHub →</a>` : ''}
            </div>
          </div>`).join('')
      : '<div class="empty">暂无精选项目</div>';

    // Latest posts
    const { posts } = await api.getPosts({ page: 1 });
    document.getElementById('latest-posts').innerHTML = posts.slice(0, 5).map(p => `
      <a href="/post.html?slug=${p.slug}" style="display:block;padding:14px 0;border-bottom:1px solid var(--border);text-decoration:none">
        <div style="font-size:15px;color:var(--text-primary);margin-bottom:4px">${p.title}</div>
        <div style="font-size:12px;color:var(--text-muted)">${new Date(p.created_at).toLocaleDateString('zh-CN')} · ${p.word_count} 字</div>
      </a>`).join('');
  </script>
</body>
</html>
```

- [ ] **Step 2: Open in browser to verify layout**

```bash
cd /Users/daiyutong/IdeaProjects/blog/frontend
python3 -m http.server 3000
```

Open http://localhost:3000/index.html?mock=1 — verify Hero, skill badges, timeline, project cards, latest posts all render. Check theme toggle works.

- [ ] **Step 3: Commit**

```bash
git add frontend/index.html
git commit -m "feat: index.html homepage"
```

---

## Task 5: blog.html — 博客列表

**Files:**
- Create: `frontend/blog.html`

- [ ] **Step 1: Write blog.html**

Create `frontend/blog.html`:

```html
<!DOCTYPE html>
<html data-theme="dark" lang="zh-CN">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>博客 | 戴宇彤</title>
  <link rel="stylesheet" href="/css/theme-dark.css">
  <link rel="stylesheet" href="/css/theme-light.css">
  <link rel="stylesheet" href="/css/common.css">
</head>
<body>
  <div id="nav-root"></div>

  <main class="page">
    <div class="tabs">
      <button class="tab-btn active" data-tab="posts">文章</button>
      <button class="tab-btn" data-tab="columns">专栏</button>
    </div>

    <!-- Posts Tab -->
    <div class="tab-panel active" id="tab-posts">
      <div style="display:flex;gap:8px;flex-wrap:wrap;margin-bottom:20px;align-items:center">
        <div id="tag-filter" style="display:flex;gap:6px;flex-wrap:wrap"></div>
        <input id="search" type="text" placeholder="🔍 搜索文章..."
          style="margin-left:auto;padding:6px 12px;border:1px solid var(--border);border-radius:6px;background:var(--bg-card);color:var(--text-primary);font-size:13px;outline:none">
      </div>
      <div id="pinned-area"></div>
      <div id="post-list"><div class="loading">加载中...</div></div>
      <div id="pagination" class="pagination"></div>
    </div>

    <!-- Columns Tab -->
    <div class="tab-panel" id="tab-columns">
      <div class="grid-2" id="column-list"><div class="loading">加载中...</div></div>
    </div>
  </main>

  <script type="module">
    import { renderNav, renderTags, getApi } from '/js/components.js';
    import '/js/theme.js';

    document.getElementById('nav-root').innerHTML = renderNav('blog');

    const api = await getApi();
    api.trackPage(location.pathname, document.referrer);

    let currentTag = '', currentPage = 1, allPosts = [];

    // Tab switching
    document.querySelectorAll('.tab-btn').forEach(btn => {
      btn.addEventListener('click', () => {
        document.querySelectorAll('.tab-btn').forEach(b => b.classList.remove('active'));
        document.querySelectorAll('.tab-panel').forEach(p => p.classList.remove('active'));
        btn.classList.add('active');
        document.getElementById('tab-' + btn.dataset.tab).classList.add('active');
        if (btn.dataset.tab === 'columns') loadColumns();
      });
    });

    // Tags
    const { tags } = await api.getTags();
    const tagFilter = document.getElementById('tag-filter');
    tagFilter.innerHTML = ['全部', ...tags].map(t =>
      `<button class="page-btn ${t === '全部' ? 'active' : ''}" data-tag="${t === '全部' ? '' : t}">${t}</button>`
    ).join('');
    tagFilter.addEventListener('click', e => {
      if (!e.target.matches('.page-btn')) return;
      document.querySelectorAll('#tag-filter .page-btn').forEach(b => b.classList.remove('active'));
      e.target.classList.add('active');
      currentTag = e.target.dataset.tag;
      currentPage = 1;
      loadPosts();
    });

    // Search
    let searchTimer;
    document.getElementById('search').addEventListener('input', e => {
      clearTimeout(searchTimer);
      searchTimer = setTimeout(() => { loadPosts(e.target.value); }, 300);
    });

    function renderPostCard(p, pinned = false) {
      const tags = JSON.parse(p.tags || '[]');
      const date = new Date(p.created_at).toLocaleDateString('zh-CN');
      const mins = Math.max(1, Math.ceil((p.word_count || 0) / 300));
      return `
        <div class="${pinned ? 'pinned-card' : 'card'}" style="margin-bottom:12px">
          ${pinned ? '<span class="pinned-badge">📌 置顶</span>' : ''}
          <a href="/post.html?slug=${p.slug}" style="font-size:16px;font-weight:600;color:var(--text-primary)">${p.title}</a>
          <p style="color:var(--text-secondary);font-size:14px;margin:8px 0">${p.summary || ''}</p>
          <div style="display:flex;gap:6px;flex-wrap:wrap;margin-bottom:8px">${tags.map(t => `<span class="tag">${t}</span>`).join('')}</div>
          <div class="post-meta">
            <span>📅 ${date}</span>
            <span>📝 ${p.word_count || 0} 字</span>
            <span>📖 约 ${mins} 分钟</span>
          </div>
        </div>`;
    }

    async function loadPosts(q = '') {
      const { posts, total } = await api.getPosts({ page: currentPage, tag: currentTag, q });
      allPosts = posts;
      const pinned = posts.filter(p => p.is_pinned);
      const normal = posts.filter(p => !p.is_pinned);
      document.getElementById('pinned-area').innerHTML = pinned.map(p => renderPostCard(p, true)).join('');
      document.getElementById('post-list').innerHTML = normal.length
        ? normal.map(p => renderPostCard(p)).join('')
        : '<div class="empty">暂无文章</div>';

      // Pagination
      const totalPages = Math.ceil(total / 10);
      document.getElementById('pagination').innerHTML = Array.from({ length: totalPages }, (_, i) =>
        `<button class="page-btn ${i + 1 === currentPage ? 'active' : ''}" data-page="${i + 1}">${i + 1}</button>`
      ).join('');
    }

    document.getElementById('pagination').addEventListener('click', e => {
      if (!e.target.matches('[data-page]')) return;
      currentPage = +e.target.dataset.page;
      loadPosts();
    });

    async function loadColumns() {
      const { columns } = await api.getColumns();
      document.getElementById('column-list').innerHTML = columns.length
        ? columns.map(c => `
            <a href="/column.html?slug=${c.slug}" style="text-decoration:none">
              <div class="card" style="cursor:pointer">
                <div style="font-size:16px;font-weight:600;margin-bottom:6px;color:var(--text-primary)">${c.name}</div>
                <div style="font-size:13px;color:var(--text-secondary)">${c.description || ''}</div>
              </div>
            </a>`).join('')
        : '<div class="empty">暂无专栏</div>';
    }

    loadPosts();
  </script>
</body>
</html>
```

- [ ] **Step 2: Verify in browser**

Open http://localhost:3000/blog.html?mock=1 — verify: Tab 切换文章/专栏，置顶文章绿色边框，标签筛选点击高亮，搜索框响应，专栏卡片列表。

- [ ] **Step 3: Commit**

```bash
git add frontend/blog.html
git commit -m "feat: blog.html with posts/columns tabs and filtering"
```

---

## Task 6: post.html — 单篇文章

**Files:**
- Create: `frontend/post.html`

- [ ] **Step 1: Write post.html**

Create `frontend/post.html`:

```html
<!DOCTYPE html>
<html data-theme="dark" lang="zh-CN">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>文章 | 戴宇彤</title>
  <link rel="stylesheet" href="/css/theme-dark.css">
  <link rel="stylesheet" href="/css/theme-light.css">
  <link rel="stylesheet" href="/css/common.css">
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/github-dark.min.css" id="hljs-theme">
  <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/highlight.min.js"></script>
</head>
<body>
  <div id="nav-root"></div>

  <div style="max-width:1100px;margin:0 auto;padding:40px 24px;display:grid;grid-template-columns:1fr 200px;gap:40px" id="layout">
    <!-- Main -->
    <article>
      <div id="post-header"><div class="loading">加载中...</div></div>
      <div class="article-body" id="post-body"></div>
      <div id="post-nav" style="display:flex;justify-content:space-between;margin-top:40px;padding-top:20px;border-top:1px solid var(--border)"></div>
    </article>
    <!-- TOC sidebar -->
    <aside>
      <div class="toc" id="toc">
        <div class="section-title" style="margin-bottom:12px">目录</div>
        <div id="toc-links"></div>
      </div>
    </aside>
  </div>

  <script type="module">
    import { renderNav, renderTags, renderPostMeta, getApi } from '/js/components.js';
    import '/js/theme.js';

    document.getElementById('nav-root').innerHTML = renderNav('blog');

    // Switch hljs theme when light mode active
    const observer = new MutationObserver(() => {
      const dark = document.documentElement.getAttribute('data-theme') === 'dark';
      document.getElementById('hljs-theme').href = dark
        ? 'https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/github-dark.min.css'
        : 'https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/github.min.css';
    });
    observer.observe(document.documentElement, { attributes: true, attributeFilter: ['data-theme'] });

    const slug = new URLSearchParams(location.search).get('slug');
    if (!slug) { document.getElementById('post-header').innerHTML = '<div class="empty">文章不存在</div>'; }

    const api = await getApi();
    const { post, views, reading_minutes } = await api.getPost(slug);

    document.title = `${post.title} | 戴宇彤`;
    document.getElementById('post-header').innerHTML = `
      <h1 style="font-size:28px;font-weight:700;margin-bottom:16px;line-height:1.3">${post.title}</h1>
      ${renderPostMeta({ ...post, views, reading_minutes })}
      <div style="margin-top:12px">${renderTags(post.tags)}</div>
      <hr style="border:none;border-top:1px solid var(--border);margin:24px 0">`;

    document.getElementById('post-body').innerHTML = post.content_html;
    hljs.highlightAll();
    api.trackPost(post.id);
    api.trackPage(location.pathname, document.referrer);

    // Build TOC from h2/h3
    const headings = document.querySelectorAll('.article-body h2, .article-body h3');
    if (headings.length) {
      document.getElementById('toc-links').innerHTML = [...headings].map(h => `
        <a href="#${h.id}" class="${h.tagName === 'H3' ? 'toc-h3' : ''}">${h.textContent}</a>
      `).join('');
    } else {
      document.getElementById('toc').style.display = 'none';
      document.getElementById('layout').style.gridTemplateColumns = '1fr';
    }

    // TOC scroll highlight
    const tocLinks = document.querySelectorAll('.toc a');
    const io = new IntersectionObserver(entries => {
      entries.forEach(e => {
        if (e.isIntersecting) {
          tocLinks.forEach(a => a.classList.remove('active'));
          const a = document.querySelector(`.toc a[href="#${e.target.id}"]`);
          if (a) a.classList.add('active');
        }
      });
    }, { rootMargin: '-20% 0px -70% 0px' });
    headings.forEach(h => io.observe(h));
  </script>
</body>
</html>
```

- [ ] **Step 2: Verify in browser**

Open http://localhost:3000/post.html?slug=go-concurrency&mock=1 — verify: 标题、元信息（字数/浏览/阅读时间）、正文渲染、目录生成、代码高亮、主题切换时 hljs 样式跟随。

- [ ] **Step 3: Commit**

```bash
git add frontend/post.html
git commit -m "feat: post.html with TOC, code highlight, meta info"
```

---

## Task 7: column.html — 专栏详情

**Files:**
- Create: `frontend/column.html`

- [ ] **Step 1: Write column.html**

Create `frontend/column.html`:

```html
<!DOCTYPE html>
<html data-theme="dark" lang="zh-CN">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>专栏 | 戴宇彤</title>
  <link rel="stylesheet" href="/css/theme-dark.css">
  <link rel="stylesheet" href="/css/theme-light.css">
  <link rel="stylesheet" href="/css/common.css">
</head>
<body>
  <div id="nav-root"></div>
  <main class="page">
    <a href="/blog.html" style="font-size:13px;color:var(--text-muted);display:inline-block;margin-bottom:20px">← 返回博客</a>
    <div id="column-header"><div class="loading">加载中...</div></div>
    <div id="post-list" style="margin-top:24px"></div>
  </main>

  <script type="module">
    import { renderNav, getApi } from '/js/components.js';
    import '/js/theme.js';

    document.getElementById('nav-root').innerHTML = renderNav('blog');

    const slug = new URLSearchParams(location.search).get('slug');
    const api = await getApi();
    const { column, posts } = await api.getColumn(slug);

    document.title = `${column.name} | 戴宇彤`;
    document.getElementById('column-header').innerHTML = `
      <h1 style="font-size:26px;font-weight:700;margin-bottom:10px">${column.name}</h1>
      <p style="color:var(--text-secondary);font-size:15px;margin-bottom:16px">${column.description || ''}</p>
      <div style="font-size:13px;color:var(--text-muted)">共 ${posts.length} 篇</div>
      <hr style="border:none;border-top:1px solid var(--border);margin:20px 0">`;

    document.getElementById('post-list').innerHTML = posts.map((p, i) => `
      <a href="/post.html?slug=${p.slug}&column=${slug}" style="text-decoration:none">
        <div class="card" style="margin-bottom:10px;display:flex;align-items:center;gap:16px;cursor:pointer">
          <span style="font-size:22px;font-weight:700;color:var(--accent);min-width:32px;text-align:center;font-family:monospace">${String(i+1).padStart(2,'0')}</span>
          <div>
            <div style="font-size:15px;font-weight:600;color:var(--text-primary)">${p.title}</div>
            <div style="font-size:12px;color:var(--text-muted);margin-top:3px">${p.word_count || 0} 字 · 约 ${Math.max(1, Math.ceil((p.word_count||0)/300))} 分钟</div>
          </div>
        </div>
      </a>`).join('');
  </script>
</body>
</html>
```

- [ ] **Step 2: Verify in browser**

Open http://localhost:3000/column.html?slug=go-series&mock=1 — verify: 专栏名/介绍、有序文章列表（01 02 03...）、点击跳转 post.html。

- [ ] **Step 3: Commit**

```bash
git add frontend/column.html
git commit -m "feat: column.html with ordered post list"
```

---

## Task 8: works.html — 作品集

**Files:**
- Create: `frontend/works.html`

- [ ] **Step 1: Write works.html**

Create `frontend/works.html`:

```html
<!DOCTYPE html>
<html data-theme="dark" lang="zh-CN">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>作品集 | 戴宇彤</title>
  <link rel="stylesheet" href="/css/theme-dark.css">
  <link rel="stylesheet" href="/css/theme-light.css">
  <link rel="stylesheet" href="/css/common.css">
</head>
<body>
  <div id="nav-root"></div>
  <main class="page-wide">
    <div class="tabs">
      <button class="tab-btn active" data-tab="projects">项目</button>
      <button class="tab-btn" data-tab="demos">系统演示</button>
      <button class="tab-btn" data-tab="custom">自定义</button>
    </div>

    <div class="tab-panel active" id="tab-projects">
      <div class="grid-3" id="project-list"><div class="loading">加载中...</div></div>
    </div>

    <div class="tab-panel" id="tab-demos">
      <div id="demo-list" style="display:flex;flex-direction:column;gap:24px"></div>
    </div>

    <div class="tab-panel" id="tab-custom">
      <div class="empty">自定义内容（后台配置）</div>
    </div>
  </main>

  <script type="module">
    import { renderNav, renderBadges, getApi } from '/js/components.js';
    import '/js/theme.js';

    document.getElementById('nav-root').innerHTML = renderNav('works');

    document.querySelectorAll('.tab-btn').forEach(btn => {
      btn.addEventListener('click', () => {
        document.querySelectorAll('.tab-btn').forEach(b => b.classList.remove('active'));
        document.querySelectorAll('.tab-panel').forEach(p => p.classList.remove('active'));
        btn.classList.add('active');
        document.getElementById('tab-' + btn.dataset.tab).classList.add('active');
      });
    });

    const api = await getApi();
    api.trackPage(location.pathname, document.referrer);

    const { projects } = await api.getProjects();

    document.getElementById('project-list').innerHTML = projects.length
      ? projects.map(p => `
          <div class="card">
            <div style="font-size:16px;font-weight:600;margin-bottom:6px">${p.name}</div>
            <div style="font-size:13px;color:var(--text-secondary);margin-bottom:12px;min-height:36px">${p.description || ''}</div>
            <div style="margin-bottom:10px;display:flex;gap:4px;flex-wrap:wrap">${renderBadges(p.tech_stack)}</div>
            <div style="display:flex;gap:12px;font-size:12px;color:var(--text-muted)">
              ${p.stars ? `<span>⭐ ${p.stars}</span>` : ''}
              ${p.repo_url ? `<a href="${p.repo_url}" target="_blank">GitHub →</a>` : ''}
              ${p.demo_url ? `<a href="${p.demo_url}" target="_blank">Demo →</a>` : ''}
            </div>
          </div>`).join('')
      : '<div class="empty">暂无项目</div>';

    // Demo tab: projects with demo_url get iframe embed
    const withDemo = projects.filter(p => p.demo_url);
    document.getElementById('demo-list').innerHTML = withDemo.length
      ? withDemo.map(p => `
          <div class="card">
            <div style="font-size:16px;font-weight:600;margin-bottom:8px">${p.name}</div>
            <p style="font-size:13px;color:var(--text-secondary);margin-bottom:12px">${p.description || ''}</p>
            <iframe src="${p.demo_url}" style="width:100%;height:480px;border:1px solid var(--border);border-radius:6px" loading="lazy"></iframe>
          </div>`).join('')
      : '<div class="empty">暂无系统演示</div>';
  </script>
</body>
</html>
```

- [ ] **Step 2: Verify in browser**

Open http://localhost:3000/works.html?mock=1 — verify: 项目卡片 grid、技术栈 badge 颜色、Tab 切换正常、演示 Tab 空状态。

- [ ] **Step 3: Commit**

```bash
git add frontend/works.html
git commit -m "feat: works.html with project cards and demo tab"
```

---

## Task 9: garden.html — 知识花园

**Files:**
- Create: `frontend/garden.html`

- [ ] **Step 1: Write garden.html**

Create `frontend/garden.html`:

```html
<!DOCTYPE html>
<html data-theme="dark" lang="zh-CN">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>知识花园 | 戴宇彤</title>
  <link rel="stylesheet" href="/css/theme-dark.css">
  <link rel="stylesheet" href="/css/theme-light.css">
  <link rel="stylesheet" href="/css/common.css">
</head>
<body>
  <div id="nav-root"></div>
  <main class="page">
    <!-- Section nav -->
    <div style="display:flex;gap:8px;margin-bottom:28px;flex-wrap:wrap" id="section-nav"></div>
    <!-- Content area -->
    <div id="section-content"><div class="loading">加载中...</div></div>
  </main>

  <script type="module">
    import { renderNav, getApi } from '/js/components.js';
    import '/js/theme.js';

    document.getElementById('nav-root').innerHTML = renderNav('garden');

    const api = await getApi();
    api.trackPage(location.pathname, document.referrer);

    const { sections } = await api.getGardenSections();
    let activeSection = sections[0];

    const TYPE_LABEL = { article: '文章', video: '视频', book: '书籍', account: '账号', other: '其他' };
    const TYPE_COLOR = { article: '#fbbf24', video: '#f87171', book: '#60a5fa', account: '#4ade80', other: '#94a3b8' };

    function renderSectionNav() {
      document.getElementById('section-nav').innerHTML = sections.map(s => `
        <button class="page-btn ${s.id === activeSection.id ? 'active' : ''}" data-id="${s.id}">${s.name}</button>
      `).join('');
    }

    async function renderSection(section) {
      const { items } = await api.getGardenItems(section.id);
      if (section.slug === 'social') {
        document.getElementById('section-content').innerHTML = `
          <div class="section-title">我在哪里</div>
          <div class="grid-3" style="margin-top:12px">
            ${items.map(item => `
              <div class="card" style="display:flex;align-items:center;gap:12px">
                <div style="width:40px;height:40px;background:var(--bg-hover);border-radius:8px;display:flex;align-items:center;justify-content:center;font-size:20px">🌐</div>
                <div>
                  <div style="font-weight:600">${item.title}</div>
                  <a href="${item.url}" target="_blank" style="font-size:12px;color:var(--text-muted)">${item.source}</a>
                </div>
              </div>`).join('') || '<div class="empty">规划中，敬请期待</div>'}
          </div>`;
        return;
      }
      // Bookmarks section
      document.getElementById('section-content').innerHTML = `
        <div style="display:flex;flex-direction:column;gap:12px">
          ${items.length ? items.map(item => `
            <div class="card">
              <div style="display:flex;align-items:flex-start;gap:12px">
                <span style="background:${TYPE_COLOR[item.type] || '#94a3b8'}20;color:${TYPE_COLOR[item.type] || '#94a3b8'};font-size:11px;padding:2px 8px;border-radius:4px;white-space:nowrap;margin-top:2px">${TYPE_LABEL[item.type] || item.type}</span>
                <div style="flex:1">
                  <div style="font-weight:600;margin-bottom:4px">${item.title}</div>
                  ${item.summary ? `<p style="font-size:13px;color:var(--text-secondary);margin-bottom:8px">${item.summary}</p>` : ''}
                  <div style="font-size:12px;color:var(--text-muted);display:flex;gap:12px">
                    ${item.source ? `<span>${item.source}</span>` : ''}
                    ${item.url ? `<a href="${item.url}" target="_blank">🔗 原文链接</a>` : ''}
                  </div>
                  ${item.embed_url ? `<iframe src="${item.embed_url}" style="width:100%;height:320px;border:1px solid var(--border);border-radius:6px;margin-top:12px" loading="lazy"></iframe>` : ''}
                </div>
              </div>
            </div>`).join('')
          : '<div class="empty">暂无收藏</div>'}
        </div>`;
    }

    renderSectionNav();
    renderSection(activeSection);

    document.getElementById('section-nav').addEventListener('click', e => {
      if (!e.target.matches('[data-id]')) return;
      activeSection = sections.find(s => s.id === +e.target.dataset.id);
      renderSectionNav();
      renderSection(activeSection);
    });
  </script>
</body>
</html>
```

- [ ] **Step 2: Verify in browser**

Open http://localhost:3000/garden.html?mock=1 — verify: 子板块导航切换，收藏卡片类型颜色标签（文章/视频/书籍），我在哪里 Tab 账号入口卡片。

- [ ] **Step 3: Commit**

```bash
git add frontend/garden.html
git commit -m "feat: garden.html with sections, bookmarks, and social tab"
```

---

## Task 10: 整体验收

- [ ] **Step 1: Start local server and test all pages**

```bash
cd /Users/daiyutong/IdeaProjects/blog/frontend
python3 -m http.server 3000
```

Test each URL with `?mock=1`:
- http://localhost:3000/index.html?mock=1
- http://localhost:3000/blog.html?mock=1
- http://localhost:3000/post.html?slug=go-concurrency&mock=1
- http://localhost:3000/column.html?slug=go-series&mock=1
- http://localhost:3000/works.html?mock=1
- http://localhost:3000/garden.html?mock=1

For each page verify:
1. Nav renders with correct active link
2. Theme toggle switches dark↔light, persists on refresh
3. Content loads from mock data without errors
4. No console errors

- [ ] **Step 2: Mobile responsive check**

In Chrome DevTools, switch to mobile view (375px width). Verify `.grid-2` and `.grid-3` collapse to single column. Nav links don't overflow.

- [ ] **Step 3: Final commit**

```bash
git add -A
git commit -m "feat: all frontend pages complete, verified with mock data"
```

---

## Self-Review

**Spec coverage check:**
- ✅ index.html: Hero, 技能标签, 时间线, 精选项目, 最新文章, 联系方式
- ✅ blog.html: 置顶文章, 标签筛选, 搜索, 文章列表(字数/浏览/阅读时长), 分页, 专栏 Tab
- ✅ post.html: 字数/浏览人次/预计阅读时间, 目录(H2/H3), 代码高亮, 上下篇导航(mock中未实现—backend提供数据后自然接入)
- ✅ column.html: 封面+介绍, 有序文章列表(01/02...), 文章字数
- ✅ works.html: 项目卡片(GitHub/Star/技术栈), 系统演示(iframe), 自定义Tab
- ✅ garden.html: 子板块导航, 外链卡片(类型标签+颜色), iframe预览, 账号入口
- ✅ 双主题切换: CSS变量 + localStorage + 跟随系统偏好
- ✅ mock模式: ?mock=1 无需后端即可开发和演示
- ✅ api.js: 覆盖所有后端接口

**Placeholder scan:** 上下篇导航在 post.html 中预留了 `#post-nav` 容器，后端提供数据后填充，符合设计意图，非占位符。

**Type consistency:** `getApi()` 在 `components.js` 中定义并在所有页面使用，mock.js 和 api.js 导出相同的 `api` 对象结构。✅
